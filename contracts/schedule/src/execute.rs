// src/execute.rs
use cosmwasm_std::{DepsMut, Env, MessageInfo, Response, StdResult, Order};
use crate::error::ContractError;
use crate::state::*;
use crate::ipfs::{SubjectContent, fetch_ipfs_content, estimate_weekly_workload, analyze_subject_difficulty};

pub fn update_config(
    deps: DepsMut,
    info: MessageInfo,
    ipfs_gateway: Option<String>,
    max_subjects_per_semester: Option<u32>,
    recommendation_algorithm: Option<RecommendationAlgorithm>,
    new_owner: Option<String>,
) -> Result<Response, ContractError> {
    let mut state = STATE.load(deps.storage)?;
    
    if info.sender != state.owner {
        return Err(ContractError::Unauthorized {
            owner: state.owner.to_string(),
        });
    }
    
    if let Some(new_owner_addr) = new_owner {
        let validated_addr = deps.api.addr_validate(&new_owner_addr)?;
        state.owner = validated_addr;
    }
    
    if let Some(gateway) = ipfs_gateway {
        state.ipfs_gateway = gateway;
    }
    
    STATE.save(deps.storage, &state)?;
    
    if max_subjects_per_semester.is_some() || recommendation_algorithm.is_some() {
        let mut config = SCHEDULE_CONFIG.load(deps.storage)?;
        
        if let Some(max_subjects) = max_subjects_per_semester {
            config.max_subjects_per_semester = max_subjects;
            config.default_preferences.max_subjects_per_semester = max_subjects;
        }
        
        if let Some(algorithm) = recommendation_algorithm {
            config.recommendation_algorithm = algorithm;
        }
        
        SCHEDULE_CONFIG.save(deps.storage, &config)?;
    }
    
    Ok(Response::new()
        .add_attribute("method", "update_config")
        .add_attribute("updated_by", info.sender))
}

pub fn register_student_progress(
    deps: DepsMut,
    env: Env,
    student_progress: StudentProgress,
) -> Result<Response, ContractError> {
    if student_progress.student_id.is_empty() {
        return Err(ContractError::InvalidScheduleConfig {
            reason: "Student ID cannot be empty".to_string(),
        });
    }
    
    let existing = STUDENT_PROGRESS.may_load(deps.storage, &student_progress.student_id)?;
    if existing.is_none() {
        STATE.update(deps.storage, |mut state| -> StdResult<_> {
            state.total_students += 1;
            Ok(state)
        })?;
    }
    
    STUDENT_PROGRESS.save(deps.storage, &student_progress.student_id, &student_progress)?;
    
    Ok(Response::new()
        .add_attribute("method", "register_student_progress")
        .add_attribute("student_id", student_progress.student_id)
        .add_attribute("timestamp", env.block.time.seconds().to_string()))
}

pub fn update_student_preferences(
    deps: DepsMut,
    student_id: String,
    preferences: SchedulePreferences,
) -> Result<Response, ContractError> {
    STUDENT_PROGRESS.update(deps.storage, &student_id, |progress| -> Result<_, ContractError> {
        let mut updated_progress = progress.ok_or(ContractError::StudentNotFound { 
            student_id: student_id.clone() 
        })?;
        updated_progress.preferences = preferences;
        Ok(updated_progress)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "update_student_preferences")
        .add_attribute("student_id", student_id))
}

pub fn register_subject_schedule_info(
    deps: DepsMut,
    subject_info: SubjectScheduleInfo,
) -> Result<Response, ContractError> {
    if subject_info.subject_id.is_empty() {
        return Err(ContractError::InvalidScheduleConfig {
            reason: "Subject ID cannot be empty".to_string(),
        });
    }
    
    SUBJECT_SCHEDULE_INFO.save(deps.storage, &subject_info.subject_id, &subject_info)?;
    
    Ok(Response::new()
        .add_attribute("method", "register_subject_schedule_info")
        .add_attribute("subject_id", subject_info.subject_id))
}

pub fn batch_register_subjects(
    deps: DepsMut,
    subjects: Vec<SubjectScheduleInfo>,
) -> Result<Response, ContractError> {
    let mut registered_count = 0u32;
    
    for subject in subjects {
        if !subject.subject_id.is_empty() {
            SUBJECT_SCHEDULE_INFO.save(deps.storage, &subject.subject_id, &subject)?;
            registered_count += 1;
        }
    }
    
    Ok(Response::new()
        .add_attribute("method", "batch_register_subjects")
        .add_attribute("registered_count", registered_count.to_string()))
}

pub fn generate_schedule_recommendation(
    deps: DepsMut,
    env: Env,
    student_id: String,
    target_semester: u32,
    force_refresh: Option<bool>,
    custom_preferences: Option<SchedulePreferences>,
) -> Result<Response, ContractError> {
    let student_progress = STUDENT_PROGRESS.load(deps.storage, &student_id)
        .map_err(|_| ContractError::StudentNotFound { student_id: student_id.clone() })?;
    
    if force_refresh != Some(true) {
        if let Ok(_existing) = SCHEDULE_RECOMMENDATIONS.load(deps.storage, (&student_id, target_semester)) {
            return Ok(Response::new()
                .add_attribute("method", "generate_schedule_recommendation")
                .add_attribute("student_id", student_id)
                .add_attribute("semester", target_semester.to_string())
                .add_attribute("cached", "true"));
        }
    }
    
    let config = SCHEDULE_CONFIG.load(deps.storage)?;
    let preferences = custom_preferences.unwrap_or(student_progress.preferences.clone());
    
    let completed_subject_ids: Vec<String> = student_progress.completed_subjects
        .iter()
        .map(|s| s.subject_id.clone())
        .collect();
    
    let enrolled_subject_ids: Vec<String> = student_progress.current_subjects
        .iter()
        .map(|s| s.subject_id.clone())
        .collect();
    
    let mut recommended_subjects = Vec::new();
    let mut alternative_subjects = Vec::new();
    let mut total_credits = 0u32;
    let mut estimated_workload = 0u32;
    
    // Load all subjects and filter eligible ones
    let all_subjects: Result<Vec<_>, _> = SUBJECT_SCHEDULE_INFO
        .range(deps.storage, None, None, Order::Ascending)
        .collect();
    
    if let Ok(subjects) = all_subjects {
        for (subject_id, subject_info) in subjects {
            // Skip completed or enrolled subjects
            if completed_subject_ids.contains(&subject_id) || enrolled_subject_ids.contains(&subject_id) {
                continue;
            }
            
            // Check if offered in target semester
            if !subject_info.semester_offered.contains(&target_semester) {
                continue;
            }
            
            // Check prerequisites
            let prerequisites_met = subject_info.prerequisites.iter()
                .all(|prereq| completed_subject_ids.contains(prereq));
            
            if !prerequisites_met {
                continue;
            }
            
            // Determine priority
            let priority = if preferences.priority_subjects.contains(&subject_id) {
                Priority::High
            } else if subject_info.is_elective {
                Priority::Low
            } else {
                Priority::Medium
            };
            
            let reason = if !subject_info.is_elective {
                RecommendationReason::MandatoryForGraduation
            } else {
                RecommendationReason::OptimalSequencing
            };
            
            // Enhanced difficulty and workload calculation using IPFS if available
            let (enhanced_difficulty, enhanced_workload) = if let Some(ipfs_link) = &subject_info.ipfs_link {
                if let Ok(content) = fetch_ipfs_content(deps.as_ref(), ipfs_link) {
                    (analyze_subject_difficulty(&content), estimate_weekly_workload(&content))
                } else {
                    (subject_info.difficulty_level, subject_info.workload_hours)
                }
            } else {
                (subject_info.difficulty_level, subject_info.workload_hours)
            };
            
            let recommended_subject = RecommendedSubject {
                subject_id: subject_id.clone(),
                title: subject_info.title.clone(),
                credits: subject_info.credits,
                priority: priority.clone(),
                recommendation_reason: reason,
                prerequisites_met: true,
                estimated_difficulty: enhanced_difficulty,
                estimated_workload: enhanced_workload,
                schedule_conflicts: vec![], // Would implement real conflict detection
                alternative_semesters: subject_info.semester_offered.clone(),
            };
            
            // Add to appropriate list based on space and priority
            if recommended_subjects.len() < preferences.max_subjects_per_semester as usize && 
               total_credits + subject_info.credits <= 30 { // Max 30 credits per semester
                
                if matches!(priority, Priority::High | Priority::Medium) {
                    recommended_subjects.push(recommended_subject);
                    total_credits += subject_info.credits;
                    estimated_workload += enhanced_workload;
                } else {
                    alternative_subjects.push(recommended_subject);
                }
            } else {
                alternative_subjects.push(recommended_subject);
            }
        }
    }
    
    // Sort recommendations by priority and difficulty
    recommended_subjects.sort_by(|a, b| {
        match (&a.priority, &b.priority) {
            (Priority::Critical, _) => std::cmp::Ordering::Less,
            (_, Priority::Critical) => std::cmp::Ordering::Greater,
            (Priority::High, Priority::Medium | Priority::Low | Priority::Optional) => std::cmp::Ordering::Less,
            (Priority::Medium | Priority::Low | Priority::Optional, Priority::High) => std::cmp::Ordering::Greater,
            _ => a.estimated_difficulty.cmp(&b.estimated_difficulty),
        }
    });
    
    // Calculate difficulty score
    let difficulty_score = if !recommended_subjects.is_empty() {
        recommended_subjects.iter()
            .map(|s| s.estimated_difficulty)
            .sum::<u32>() / recommended_subjects.len() as u32
    } else {
        0
    };
    
    // Calculate completion percentage
    let total_required_credits = student_progress.total_credits_required;
    let completion_percentage = if total_required_credits > 0 {
        format!("{:.1}", (student_progress.total_credits as f64 / total_required_credits as f64) * 100.0)
    } else {
        "0.0".to_string()
    };
    
    // Generate notes based on analysis
    let mut notes = Vec::new();
    notes.push(format!("Generated using {} algorithm", format!("{:?}", config.recommendation_algorithm)));
    
    if total_credits > 24 {
        notes.push("Heavy credit load - consider reducing if struggling".to_string());
    }
    
    if difficulty_score > 80 {
        notes.push("High difficulty semester - plan extra study time".to_string());
    }
    
    if recommended_subjects.is_empty() {
        notes.push("No eligible subjects found for this semester".to_string());
    }
    
    // Calculate confidence based on available data
    let confidence_score = calculate_confidence_score(&recommended_subjects, &student_progress);
    
    let recommendation = ScheduleRecommendation {
        student_id: student_id.clone(),
        target_semester,
        recommended_subjects,
        alternative_subjects,
        total_credits,
        estimated_workload,
        difficulty_score,
        completion_percentage,
        notes,
        confidence_score,
        generated_timestamp: env.block.time.seconds(),
        algorithm_used: config.recommendation_algorithm,
    };
    
    SCHEDULE_RECOMMENDATIONS.save(deps.storage, (&student_id, target_semester), &recommendation)?;
    
    STATE.update(deps.storage, |mut state| -> StdResult<_> {
        state.total_recommendations += 1;
        Ok(state)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "generate_schedule_recommendation")
        .add_attribute("student_id", student_id)
        .add_attribute("semester", target_semester.to_string())
        .add_attribute("subjects_recommended", recommendation.recommended_subjects.len().to_string())
        .add_attribute("total_credits", total_credits.to_string())
        .add_attribute("confidence", confidence_score.to_string()))
}

pub fn create_academic_path(
    deps: DepsMut,
    env: Env,
    student_id: String,
    path_name: String,
    optimization_criteria: OptimizationCriteria,
    target_graduation_semester: Option<u32>,
) -> Result<Response, ContractError> {
    let student_progress = STUDENT_PROGRESS.load(deps.storage, &student_id)
        .map_err(|_| ContractError::StudentNotFound { student_id: student_id.clone() })?;
    
    let path_id = format!("{}_{}", student_id, env.block.time.seconds());
    
    let remaining_credits = student_progress.total_credits_required
        .saturating_sub(student_progress.total_credits);
    
    if remaining_credits == 0 {
        return Err(ContractError::InvalidScheduleConfig {
            reason: "Student has completed all required credits".to_string(),
        });
    }
    
    // Calculate duration based on optimization criteria
    let credits_per_semester = match optimization_criteria {
        OptimizationCriteria::Fastest => 24,
        OptimizationCriteria::Balanced => 18,
        OptimizationCriteria::EasiestFirst => 15,
        OptimizationCriteria::PrerequisiteOptimal => 20,
        OptimizationCriteria::Custom(_) => 18,
    };
    
    let estimated_duration = (remaining_credits + credits_per_semester - 1) / credits_per_semester;
    
    let target_graduation = target_graduation_semester
        .unwrap_or(student_progress.current_semester + estimated_duration);
    
    // Generate real semester plans
    let semesters = generate_realistic_semester_plans(
        deps.as_ref(),
        &student_progress,
        remaining_credits,
        target_graduation,
        &optimization_criteria,
    )?;
    
    let path_metrics = calculate_real_path_metrics(&semesters);
    
    let academic_path = AcademicPath {
        path_id: path_id.clone(),
        student_id: student_id.clone(),
        path_name: path_name.clone(),
        semesters,
        total_duration_semesters: estimated_duration,
        total_credits: student_progress.total_credits_required,
        expected_graduation_date: format!("Semester {}", target_graduation),
        optimization_criteria: optimization_criteria.clone(),
        path_metrics,
        created_timestamp: env.block.time.seconds(),
        last_updated: env.block.time.seconds(),
        is_optimized: false,
        status: PathStatus::Draft,
    };
    
    ACADEMIC_PATHS.save(deps.storage, &path_id, &academic_path)?;
    
    STUDENT_PATHS.update(deps.storage, &student_id, |paths| -> StdResult<Vec<String>> {
        let mut list = paths.unwrap_or_default();
        list.push(path_id.clone());
        Ok(list)
    })?;
    
    STATE.update(deps.storage, |mut state| -> StdResult<_> {
        state.total_paths += 1;
        Ok(state)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "create_academic_path")
        .add_attribute("student_id", student_id)
        .add_attribute("path_id", path_id)
        .add_attribute("estimated_duration", estimated_duration.to_string()))
}

pub fn optimize_academic_path(
    deps: DepsMut,
    env: Env,
    path_id: String,
    optimization_criteria: OptimizationCriteria,
    preserve_current_semester: Option<bool>,
) -> Result<Response, ContractError> {
    // Load the path first to avoid borrow conflicts
    let mut updated_path = ACADEMIC_PATHS.load(deps.storage, &path_id)
        .map_err(|_| ContractError::AcademicPathNotFound { path_id: path_id.clone() })?;
    
    // Load all necessary subject info before updating
    let mut subject_info_cache = std::collections::HashMap::new();
    for semester in &updated_path.semesters {
        for subject in &semester.subjects {
            if let Ok(info) = SUBJECT_SCHEDULE_INFO.load(deps.storage, &subject.subject_id) {
                subject_info_cache.insert(subject.subject_id.clone(), info);
            }
        }
    }
    
    // Update the path with cached data
    updated_path.optimization_criteria = optimization_criteria.clone();
    updated_path.is_optimized = true;
    updated_path.last_updated = env.block.time.seconds();
    
    // Apply optimization using cached data
    apply_path_optimization_with_cache(&mut updated_path, &optimization_criteria, 
                                     preserve_current_semester.unwrap_or(true), &subject_info_cache)?;
    
    // Save the updated path
    ACADEMIC_PATHS.save(deps.storage, &path_id, &updated_path)?;
    
    Ok(Response::new()
        .add_attribute("method", "optimize_academic_path")
        .add_attribute("path_id", path_id)
        .add_attribute("criteria", format!("{:?}", optimization_criteria)))
}

pub fn update_academic_path(
    deps: DepsMut,
    _env: Env,
    path_id: String,
    semester_number: u32,
    new_subjects: Vec<String>,
    notes: Option<String>,
) -> Result<Response, ContractError> {
    // Load the path first to avoid borrow conflicts
    let mut updated_path = ACADEMIC_PATHS.load(deps.storage, &path_id)
        .map_err(|_| ContractError::AcademicPathNotFound { path_id: path_id.clone() })?;
    
    // Pre-load all subject information we need
    let mut subject_info_cache = std::collections::HashMap::new();
    for subject_id in &new_subjects {
        if let Ok(subject_info) = SUBJECT_SCHEDULE_INFO.load(deps.storage, subject_id) {
            subject_info_cache.insert(subject_id.clone(), subject_info);
        } else {
            return Err(ContractError::SubjectNotFound { subject_id: subject_id.clone() });
        }
    }
    
    // Update the path using cached data
    if let Some(semester) = updated_path.semesters.iter_mut()
        .find(|s| s.semester_number == semester_number) {
        
        let mut total_credits = 0u32;
        let mut planned_subjects = Vec::new();
        
        for subject_id in new_subjects {
            if let Some(subject_info) = subject_info_cache.get(&subject_id) {
                total_credits += subject_info.credits;
                planned_subjects.push(PlannedSubject {
                    subject_id: subject_id.clone(),
                    title: subject_info.title.clone(),
                    credits: subject_info.credits,
                    is_mandatory: !subject_info.is_elective,
                    backup_options: vec![],
                    placement_reason: "Manual update".to_string(),
                });
            }
        }
        
        semester.subjects = planned_subjects;
        semester.total_credits = total_credits;
        semester.estimated_workload = total_credits * 3;
        
        if let Some(note) = notes {
            semester.notes.push(note);
        }
        
        updated_path.last_updated = 1640995200u64; // Fixed timestamp
    }
    
    // Save the updated path
    ACADEMIC_PATHS.save(deps.storage, &path_id, &updated_path)?;
    
    Ok(Response::new()
        .add_attribute("method", "update_academic_path")
        .add_attribute("path_id", path_id)
        .add_attribute("semester", semester_number.to_string()))
}

pub fn activate_academic_path(
    deps: DepsMut,
    path_id: String,
) -> Result<Response, ContractError> {
    let path = ACADEMIC_PATHS.load(deps.storage, &path_id)
        .map_err(|_| ContractError::AcademicPathNotFound { path_id: path_id.clone() })?;
    
    let student_id = path.student_id.clone();
    
    // Deactivate other paths for this student
    let student_path_ids = STUDENT_PATHS.may_load(deps.storage, &student_id)?.unwrap_or_default();
    
    for existing_path_id in &student_path_ids {
        if existing_path_id != &path_id {
            if let Ok(mut existing_path) = ACADEMIC_PATHS.load(deps.storage, existing_path_id) {
                if matches!(existing_path.status, PathStatus::Active) {
                    existing_path.status = PathStatus::Draft;
                    ACADEMIC_PATHS.save(deps.storage, existing_path_id, &existing_path)?;
                }
            }
        }
    }
    
    // Activate the target path
    let mut updated_path = path;
    updated_path.status = PathStatus::Active;
    ACADEMIC_PATHS.save(deps.storage, &path_id, &updated_path)?;
    
    Ok(Response::new()
        .add_attribute("method", "activate_academic_path")
        .add_attribute("path_id", path_id)
        .add_attribute("student_id", student_id))
}

pub fn complete_subject(
    deps: DepsMut,
    _env: Env,
    student_id: String,
    subject_id: String,
    grade: u32,
    completion_date: String,
    difficulty_rating: Option<u32>,
    workload_rating: Option<u32>,
    nft_token_id: String,
) -> Result<Response, ContractError> {
    let subject_info = SUBJECT_SCHEDULE_INFO.load(deps.storage, &subject_id)
        .map_err(|_| ContractError::SubjectNotFound { subject_id: subject_id.clone() })?;
    
    STUDENT_PROGRESS.update(deps.storage, &student_id, |progress| -> Result<_, ContractError> {
        let mut updated_progress = progress.ok_or(ContractError::StudentNotFound { 
            student_id: student_id.clone() 
        })?;
        
        // Remove from current subjects
        updated_progress.current_subjects.retain(|s| s.subject_id != subject_id);
        
        let completed_subject = CompletedSubject {
            subject_id: subject_id.clone(),
            credits: subject_info.credits,
            grade,
            completion_date: completion_date.clone(),
            semester_taken: updated_progress.current_semester,
            difficulty_rating,
            workload_rating,
            nft_token_id: nft_token_id.clone(),
        };
        
        updated_progress.completed_subjects.push(completed_subject);
        updated_progress.total_credits += subject_info.credits;
        
        // Recalculate GPA
        let total_grade_points: u32 = updated_progress.completed_subjects
            .iter()
            .map(|s| s.grade * s.credits)
            .sum();
        let total_credits: u32 = updated_progress.completed_subjects
            .iter()
            .map(|s| s.credits)
            .sum();
        
        if total_credits > 0 {
            updated_progress.gpa = format!("{:.2}", total_grade_points as f64 / total_credits as f64 / 10.0);
        }
        
        Ok(updated_progress)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "complete_subject")
        .add_attribute("student_id", student_id)
        .add_attribute("subject_id", subject_id)
        .add_attribute("grade", grade.to_string())
        .add_attribute("nft_token_id", nft_token_id))
}

pub fn enroll_in_subject(
    deps: DepsMut,
    student_id: String,
    subject_id: String,
    enrollment_date: String,
    expected_completion: String,
) -> Result<Response, ContractError> {
    let subject_info = SUBJECT_SCHEDULE_INFO.load(deps.storage, &subject_id)
        .map_err(|_| ContractError::SubjectNotFound { subject_id: subject_id.clone() })?;
    
    STUDENT_PROGRESS.update(deps.storage, &student_id, |progress| -> Result<_, ContractError> {
        let mut updated_progress = progress.ok_or(ContractError::StudentNotFound { 
            student_id: student_id.clone() 
        })?;
        
        if updated_progress.current_subjects.iter().any(|s| s.subject_id == subject_id) {
            return Err(ContractError::InvalidScheduleConfig {
                reason: "Student already enrolled in this subject".to_string(),
            });
        }
        
        let enrolled_subject = EnrolledSubject {
            subject_id: subject_id.clone(),
            credits: subject_info.credits,
            enrollment_date,
            expected_completion,
            current_grade: None,
        };
        
        updated_progress.current_subjects.push(enrolled_subject);
        
        Ok(updated_progress)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "enroll_in_subject")
        .add_attribute("student_id", student_id)
        .add_attribute("subject_id", subject_id))
}

pub fn cache_ipfs_content(
    deps: DepsMut,
    ipfs_link: String,
    content: SubjectContent,
) -> Result<Response, ContractError> {
    // Cache IPFS content implementation
    IPFS_CACHE.save(deps.storage, &ipfs_link, &content)?;
    
    Ok(Response::new()
        .add_attribute("method", "cache_ipfs_content")
        .add_attribute("ipfs_link", ipfs_link))
}

pub fn generate_alternative_recommendations(
    deps: DepsMut,
    env: Env,
    student_id: String,
    target_semester: u32,
    excluded_subjects: Vec<String>,
) -> Result<Response, ContractError> {
    let student_progress = STUDENT_PROGRESS.load(deps.storage, &student_id)
        .map_err(|_| ContractError::StudentNotFound { student_id: student_id.clone() })?;
    
    let completed_subject_ids: Vec<String> = student_progress.completed_subjects
        .iter()
        .map(|s| s.subject_id.clone())
        .collect();
    
    let mut alternative_subjects = Vec::new();
    
    let all_subjects: Result<Vec<_>, _> = SUBJECT_SCHEDULE_INFO
        .range(deps.storage, None, None, Order::Ascending)
        .collect();
    
    if let Ok(subjects) = all_subjects {
        for (subject_id, subject_info) in subjects {
            if !excluded_subjects.contains(&subject_id) && 
               !completed_subject_ids.contains(&subject_id) &&
               subject_info.semester_offered.contains(&target_semester) {
                
                let prerequisites_met = subject_info.prerequisites.iter()
                    .all(|prereq| completed_subject_ids.contains(prereq));
                
                if prerequisites_met {
                    alternative_subjects.push(RecommendedSubject {
                        subject_id: subject_id.clone(),
                        title: subject_info.title.clone(),
                        credits: subject_info.credits,
                        priority: Priority::Medium,
                        recommendation_reason: RecommendationReason::OptimalSequencing,
                        prerequisites_met: true,
                        estimated_difficulty: subject_info.difficulty_level,
                        estimated_workload: subject_info.workload_hours,
                        schedule_conflicts: vec![],
                        alternative_semesters: subject_info.semester_offered.clone(),
                    });
                }
            }
        }
    }
    
    let recommendation = ScheduleRecommendation {
        student_id: student_id.clone(),
        target_semester,
        recommended_subjects: vec![],
        alternative_subjects,
        total_credits: 0,
        estimated_workload: 0,
        difficulty_score: 0,
        completion_percentage: "0".to_string(),
        notes: vec!["Alternative recommendations generated".to_string()],
        confidence_score: 60,
        generated_timestamp: env.block.time.seconds(),
        algorithm_used: RecommendationAlgorithm::Basic,
    };
    
    let alt_key = format!("{}_alt", student_id);
    SCHEDULE_RECOMMENDATIONS.save(deps.storage, (&alt_key, target_semester), &recommendation)?;
    
    Ok(Response::new()
        .add_attribute("method", "generate_alternative_recommendations")
        .add_attribute("student_id", student_id)
        .add_attribute("alternatives_count", recommendation.alternative_subjects.len().to_string()))
}

pub fn simulate_schedule(
    deps: DepsMut,
    _env: Env,
    student_id: String,
    hypothetical_completions: Vec<String>,
    target_semester: u32,
) -> Result<Response, ContractError> {
    let mut simulated_progress = STUDENT_PROGRESS.load(deps.storage, &student_id)
        .map_err(|_| ContractError::StudentNotFound { student_id: student_id.clone() })?;
    
    let hypothetical_count = hypothetical_completions.len();
    for subject_id in &hypothetical_completions {
        if let Ok(subject_info) = SUBJECT_SCHEDULE_INFO.load(deps.storage, &subject_id) {
            if !simulated_progress.completed_subjects.iter().any(|s| s.subject_id == *subject_id) {
                simulated_progress.completed_subjects.push(CompletedSubject {
                    subject_id: subject_id.clone(),
                    credits: subject_info.credits,
                    grade: 85,
                    completion_date: "simulated".to_string(),
                    semester_taken: target_semester.saturating_sub(1),
                    difficulty_rating: None,
                    workload_rating: None,
                    nft_token_id: "simulated".to_string(),
                });
                simulated_progress.total_credits += subject_info.credits;
            }
        }
    }
    
    Ok(Response::new()
        .add_attribute("method", "simulate_schedule")
        .add_attribute("student_id", student_id)
        .add_attribute("hypothetical_count", hypothetical_count.to_string()))
}

// Helper functions
fn calculate_confidence_score(subjects: &[RecommendedSubject], progress: &StudentProgress) -> u32 {
    if subjects.is_empty() {
        return 0;
    }
    
    let high_priority_count = subjects.iter().filter(|s| matches!(s.priority, Priority::High | Priority::Critical)).count();
    let base_score = (high_priority_count * 100 / subjects.len()).min(100) as u32;
    
    // Adjust based on student's academic standing
    let standing_bonus = match progress.academic_standing {
        crate::state::AcademicStanding::Excellent => 20,
        crate::state::AcademicStanding::Good => 10,
        crate::state::AcademicStanding::Satisfactory => 0,
        crate::state::AcademicStanding::Probation => -20,
    };
    
    (base_score as i32 + standing_bonus).max(0).min(100) as u32
}

fn generate_realistic_semester_plans(
    deps: cosmwasm_std::Deps,
    student_progress: &StudentProgress,
    remaining_credits: u32,
    target_graduation: u32,
    criteria: &OptimizationCriteria,
) -> Result<Vec<SemesterPlan>, ContractError> {
    let mut semesters = Vec::new();
    let mut current_semester = student_progress.current_semester + 1;
    let mut credits_left = remaining_credits;
    
    let completed_ids: Vec<String> = student_progress.completed_subjects
        .iter()
        .map(|s| s.subject_id.clone())
        .collect();
    
    while credits_left > 0 && current_semester <= target_graduation {
        let credits_this_semester = match criteria {
            OptimizationCriteria::Fastest => credits_left.min(24),
            OptimizationCriteria::Balanced => credits_left.min(18),
            OptimizationCriteria::EasiestFirst => credits_left.min(15),
            _ => credits_left.min(20),
        };
        
        let mut planned_subjects = Vec::new();
        let mut semester_credits = 0u32;
        
        // Find subjects that can be taken this semester
        let all_subjects: Result<Vec<_>, _> = SUBJECT_SCHEDULE_INFO
            .range(deps.storage, None, None, Order::Ascending)
            .collect();
        
        if let Ok(subjects) = all_subjects {
            for (subject_id, subject_info) in subjects {
                if semester_credits >= credits_this_semester {
                    break;
                }
                
                if completed_ids.contains(&subject_id) ||
                   planned_subjects.iter().any(|p: &PlannedSubject| p.subject_id == subject_id) {
                    continue;
                }
                
                if subject_info.semester_offered.contains(&((current_semester % 2) + 1)) {
                    let prerequisites_met = subject_info.prerequisites.iter()
                        .all(|prereq| completed_ids.contains(prereq) ||
                             planned_subjects.iter().any(|p| p.subject_id == *prereq));
                    
                    if prerequisites_met {
                        planned_subjects.push(PlannedSubject {
                            subject_id: subject_id.clone(),
                            title: subject_info.title,
                            credits: subject_info.credits,
                            is_mandatory: !subject_info.is_elective,
                            backup_options: vec![],
                            placement_reason: format!("Planned for semester {}", current_semester),
                        });
                        semester_credits += subject_info.credits;
                    }
                }
            }
        }
        
        semesters.push(SemesterPlan {
            semester_number: current_semester,
            subjects: planned_subjects,
            total_credits: semester_credits,
            estimated_difficulty: 50,
            estimated_workload: semester_credits * 3,
            notes: vec![],
            flexibility_score: 70,
        });
        
        credits_left = credits_left.saturating_sub(semester_credits);
        current_semester += 1;
    }
    
    Ok(semesters)
}

fn calculate_real_path_metrics(semesters: &[SemesterPlan]) -> PathMetrics {
    if semesters.is_empty() {
        return PathMetrics {
            total_duration_semesters: 0,
            average_credits_per_semester: 0,
            workload_variance: 0,
            difficulty_progression_score: 0,
            prerequisite_efficiency: 0,
            schedule_flexibility: 0,
            risk_factors: vec![],
        };
    }
    
    let total_credits: u32 = semesters.iter().map(|s| s.total_credits).sum();
    let average_credits = total_credits / semesters.len() as u32;
    
    // Calculate workload variance
    let variance: f64 = semesters.iter()
        .map(|s| (s.total_credits as f64 - average_credits as f64).powi(2))
        .sum::<f64>() / semesters.len() as f64;
    
    let mut risk_factors = Vec::new();
    
    // Check for heavy semesters
    if semesters.iter().any(|s| s.total_credits > 24) {
        risk_factors.push("Some semesters have excessive credit load".to_string());
    }
    
    // Check for prerequisite bottlenecks
    let total_subjects: usize = semesters.iter().map(|s| s.subjects.len()).sum();
    if total_subjects > semesters.len() * 6 {
        risk_factors.push("High subject density may cause scheduling conflicts".to_string());
    }
    
    PathMetrics {
        total_duration_semesters: semesters.len() as u32,
        average_credits_per_semester: average_credits,
        workload_variance: variance.sqrt() as u32,
        difficulty_progression_score: 75, // Would calculate based on actual difficulty progression
        prerequisite_efficiency: 80,
        schedule_flexibility: 70,
        risk_factors,
    }
}

fn apply_path_optimization_with_cache(
    path: &mut AcademicPath,
    criteria: &OptimizationCriteria,
    preserve_current: bool,
    _subject_cache: &std::collections::HashMap<String, SubjectScheduleInfo>,
) -> Result<(), ContractError> {
    let start_index = if preserve_current { 1 } else { 0 };
    
    for (index, semester) in path.semesters.iter_mut().enumerate() {
        if index < start_index {
            continue;
        }
        
        match criteria {
            OptimizationCriteria::Fastest => {
                if semester.total_credits < 24 {
                    semester.total_credits = 24.min(semester.total_credits + 6);
                    semester.estimated_workload = semester.total_credits * 3;
                }
            }
            OptimizationCriteria::Balanced => {
                semester.total_credits = 18;
                semester.estimated_workload = 54;
                semester.estimated_difficulty = 50;
            }
            OptimizationCriteria::EasiestFirst => {
                if index < 4 {
                    semester.estimated_difficulty = semester.estimated_difficulty.saturating_sub(15);
                    semester.total_credits = semester.total_credits.min(15);
                }
            }
            _ => {}
        }
    }
    
    Ok(())
}