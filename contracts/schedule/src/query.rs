// src/query.rs
use cosmwasm_std::{Deps, StdResult, Order};
use crate::msg::{
    StateResponse, ScheduleConfigResponse, StudentProgressResponse,
    SubjectScheduleInfoResponse, ScheduleRecommendationResponse, AcademicPathResponse,
    StudentPathsResponse, AvailableSubjectsResponse, OptimalPathResponse,
    GraduationTimelineResponse, SubjectSequenceResponse, WorkloadAnalysisResponse,
    IpfsCacheStatusResponse, ScheduleStatisticsResponse, SemesterSubjects,
    SubjectWorkload, WorkloadRating, PopularSubject,
};
use crate::state::*;
use crate::ipfs::{SubjectContent, is_content_cached, fetch_ipfs_content, estimate_weekly_workload, analyze_subject_difficulty};

pub fn query_state(deps: Deps) -> StdResult<StateResponse> {
    let state = STATE.load(deps.storage)?;
    Ok(StateResponse { state })
}

pub fn query_config(deps: Deps) -> StdResult<ScheduleConfigResponse> {
    let config = SCHEDULE_CONFIG.load(deps.storage)?;
    Ok(ScheduleConfigResponse { config })
}

pub fn query_student_progress(deps: Deps, student_id: String) -> StdResult<StudentProgressResponse> {
    let progress = STUDENT_PROGRESS
        .may_load(deps.storage, &student_id)?
        .ok_or_else(|| cosmwasm_std::StdError::not_found("Student progress not found"))?;
    
    Ok(StudentProgressResponse { progress })
}

pub fn query_subject_schedule_info(deps: Deps, subject_id: String) -> StdResult<SubjectScheduleInfoResponse> {
    let subject_info = SUBJECT_SCHEDULE_INFO
        .may_load(deps.storage, &subject_id)?
        .ok_or_else(|| cosmwasm_std::StdError::not_found("Subject schedule info not found"))?;
    
    Ok(SubjectScheduleInfoResponse { subject_info })
}

pub fn query_schedule_recommendation(
    deps: Deps, 
    student_id: String, 
    semester: u32
) -> StdResult<ScheduleRecommendationResponse> {
    let recommendation = SCHEDULE_RECOMMENDATIONS
        .may_load(deps.storage, (&student_id, semester))?;
    
    Ok(ScheduleRecommendationResponse { recommendation })
}

pub fn query_academic_path(deps: Deps, path_id: String) -> StdResult<AcademicPathResponse> {
    let path = ACADEMIC_PATHS.may_load(deps.storage, &path_id)?;
    Ok(AcademicPathResponse { path })
}

pub fn query_student_paths(
    deps: Deps, 
    student_id: String, 
    include_inactive: Option<bool>
) -> StdResult<StudentPathsResponse> {
    let path_ids = STUDENT_PATHS
        .may_load(deps.storage, &student_id)?
        .unwrap_or_default();
    
    let mut paths = Vec::new();
    let mut active_path_id = None;
    let include_inactive = include_inactive.unwrap_or(false);
    
    for path_id in path_ids {
        if let Ok(path) = ACADEMIC_PATHS.load(deps.storage, &path_id) {
            let is_active = matches!(path.status, PathStatus::Active);
            let should_include = include_inactive || 
                               matches!(path.status, PathStatus::Active | PathStatus::Draft);
            
            if should_include {
                if is_active {
                    active_path_id = Some(path_id.clone());
                }
                paths.push(path);
            }
        }
    }
    
    Ok(StudentPathsResponse {
        paths,
        active_path_id,
    })
}

pub fn query_available_subjects(
    deps: Deps, 
    student_id: String, 
    semester: u32, 
    include_electives: Option<bool>
) -> StdResult<AvailableSubjectsResponse> {
    let student_progress = STUDENT_PROGRESS
        .may_load(deps.storage, &student_id)?
        .ok_or_else(|| cosmwasm_std::StdError::not_found("Student progress not found"))?;
    
    let completed_subject_ids: Vec<String> = student_progress.completed_subjects
        .iter()
        .map(|s| s.subject_id.clone())
        .collect();
    
    let enrolled_subject_ids: Vec<String> = student_progress.current_subjects
        .iter()
        .map(|s| s.subject_id.clone())
        .collect();
    
    let mut available_subjects = Vec::new();
    let include_electives = include_electives.unwrap_or(true);
    
    let all_subjects: Result<Vec<_>, _> = SUBJECT_SCHEDULE_INFO
        .range(deps.storage, None, None, Order::Ascending)
        .collect();
    
    if let Ok(subjects) = all_subjects {
        for (subject_id, subject_info) in subjects {
            // Skip completed or enrolled
            if completed_subject_ids.contains(&subject_id) || 
               enrolled_subject_ids.contains(&subject_id) {
                continue;
            }
            
            // Skip electives if not requested
            if subject_info.is_elective && !include_electives {
                continue;
            }
            
            // Check semester availability
            if !subject_info.semester_offered.contains(&semester) {
                continue;
            }
            
            // Check prerequisites
            let missing_prerequisites: Vec<String> = subject_info.prerequisites
                .iter()
                .filter(|prereq_id| !completed_subject_ids.contains(prereq_id))
                .cloned()
                .collect();
            
            let can_enroll = missing_prerequisites.is_empty();
            
            // Determine conflicts (simplified - would implement real conflict detection)
            let conflicts = Vec::new();
            
            let recommendation_priority = if student_progress.preferences.priority_subjects.contains(&subject_id) {
                Some(Priority::High)
            } else if subject_info.is_elective {
                Some(Priority::Low)
            } else if can_enroll {
                Some(Priority::Medium)
            } else {
                None
            };
            
            available_subjects.push(AvailableSubject {
                subject_id: subject_id.clone(),
                title: subject_info.title.clone(),
                credits: subject_info.credits,
                can_enroll,
                missing_prerequisites,
                conflicts,
                recommendation_priority,
            });
        }
    }
    
    Ok(AvailableSubjectsResponse {
        available_subjects: AvailableSubjects {
            semester,
            subjects: available_subjects.clone(),
            total_count: available_subjects.len() as u32,
        },
    })
}

pub fn query_optimal_path(
    deps: Deps, 
    student_id: String, 
    criteria: OptimizationCriteria, 
    max_paths: Option<u32>
) -> StdResult<OptimalPathResponse> {
    let max_paths = max_paths.unwrap_or(3).min(10);
    
    let path_ids = STUDENT_PATHS
        .may_load(deps.storage, &student_id)?
        .unwrap_or_default();
    
    let mut ranked_paths = Vec::new();
    
    for path_id in path_ids.iter().take(max_paths as usize) {
        if let Ok(path) = ACADEMIC_PATHS.load(deps.storage, path_id) {
            let score = calculate_path_optimization_score(&path, &criteria);
            ranked_paths.push((path, score));
        }
    }
    
    // Sort by optimization score (higher is better)
    ranked_paths.sort_by(|a, b| b.1.cmp(&a.1));
    
    let paths: Vec<AcademicPath> = ranked_paths.into_iter().map(|(path, _)| path).collect();
    
    let recommendation_notes = generate_optimization_notes(&criteria, &paths);
    
    Ok(OptimalPathResponse {
        paths,
        recommendation_notes,
    })
}

pub fn query_graduation_timeline(
    deps: Deps, 
    student_id: String, 
    path_id: Option<String>
) -> StdResult<GraduationTimelineResponse> {
    let student_progress = STUDENT_PROGRESS
        .may_load(deps.storage, &student_id)?
        .ok_or_else(|| cosmwasm_std::StdError::not_found("Student progress not found"))?;
    
    let total_required_credits = student_progress.total_credits_required;
    let completed_credits = student_progress.total_credits;
    let current_progress_percentage = if total_required_credits > 0 {
        format!("{:.1}%", (completed_credits as f64 / total_required_credits as f64) * 100.0)
    } else {
        "0.0%".to_string()
    };
    
    let remaining_credits = total_required_credits.saturating_sub(completed_credits);
    
    // Load specific path or find active path
    let academic_path = if let Some(path_id) = path_id {
        ACADEMIC_PATHS.may_load(deps.storage, &path_id)?
    } else {
        let path_ids = STUDENT_PATHS.may_load(deps.storage, &student_id)?
            .unwrap_or_default();
        
        let mut active_path = None;
        for pid in path_ids {
            if let Ok(path) = ACADEMIC_PATHS.load(deps.storage, &pid) {
                if matches!(path.status, PathStatus::Active) {
                    active_path = Some(path);
                    break;
                }
            }
        }
        active_path
    };
    
    // Calculate timeline
    let (estimated_graduation_semester, estimated_graduation_date) = if let Some(path) = &academic_path {
        (
            student_progress.current_semester + path.total_duration_semesters,
            path.expected_graduation_date.clone(),
        )
    } else {
        let credits_per_semester = match student_progress.preferences.preferred_study_intensity {
            StudyIntensity::Light => 12,
            StudyIntensity::Moderate => 18,
            StudyIntensity::Intensive => 24,
            StudyIntensity::Maximum => 30,
        };
        
        let remaining_semesters = (remaining_credits + credits_per_semester - 1) / credits_per_semester;
        let graduation_semester = student_progress.current_semester + remaining_semesters;
        let graduation_date = format!("Semester {}", graduation_semester);
        
        (graduation_semester, graduation_date)
    };
    
    // Find remaining subjects
    let completed_subject_ids: Vec<String> = student_progress.completed_subjects
        .iter()
        .map(|s| s.subject_id.clone())
        .collect();
    
    let mut remaining_subjects = Vec::new();
    let mut critical_path_subjects = Vec::new();
    
    if let Some(path) = &academic_path {
        for semester in &path.semesters {
            for subject in &semester.subjects {
                if !completed_subject_ids.contains(&subject.subject_id) {
                    remaining_subjects.push(subject.subject_id.clone());
                    
                    if subject.is_mandatory {
                        critical_path_subjects.push(subject.subject_id.clone());
                    }
                }
            }
        }
    }
    
    // Calculate risk factors
    let mut risk_factors = Vec::new();
    
    let credits_per_semester = match student_progress.preferences.preferred_study_intensity {
        StudyIntensity::Light => 12,
        StudyIntensity::Moderate => 18,
        StudyIntensity::Intensive => 24,
        StudyIntensity::Maximum => 30,
    };
    
    if remaining_credits > credits_per_semester * 8 {
        risk_factors.push("High remaining credit load may extend graduation".to_string());
    }
    
    if matches!(student_progress.preferences.preferred_study_intensity, StudyIntensity::Light) {
        risk_factors.push("Light study intensity may extend graduation timeline".to_string());
    }
    
    if critical_path_subjects.len() > 15 {
        risk_factors.push("Many critical path subjects remaining".to_string());
    }
    
    // Check academic standing
    if matches!(student_progress.academic_standing, AcademicStanding::Probation) {
        risk_factors.push("Academic probation may require reduced course load".to_string());
    }
    
    Ok(GraduationTimelineResponse {
        current_progress_percentage,
        estimated_graduation_semester,
        estimated_graduation_date,
        remaining_credits,
        remaining_subjects,
        critical_path_subjects,
        risk_factors,
    })
}

pub fn query_subject_sequence(
    deps: Deps, 
    student_id: String, 
    target_subjects: Vec<String>
) -> StdResult<SubjectSequenceResponse> {
    let student_progress = STUDENT_PROGRESS
        .may_load(deps.storage, &student_id)?
        .ok_or_else(|| cosmwasm_std::StdError::not_found("Student progress not found"))?;
    
    let completed_subject_ids: Vec<String> = student_progress.completed_subjects
        .iter()
        .map(|s| s.subject_id.clone())
        .collect();
    
    // Build prerequisite dependency graph
    let mut subject_dependencies = std::collections::HashMap::new();
    
    for subject_id in &target_subjects {
        if let Ok(subject_info) = SUBJECT_SCHEDULE_INFO.load(deps.storage, subject_id) {
            subject_dependencies.insert(subject_id.clone(), subject_info.prerequisites);
        }
    }
    
    // Topological sort to determine optimal sequence
    let mut sequence = Vec::new();
    let mut current_semester = student_progress.current_semester + 1;
    let mut remaining_targets = target_subjects.clone();
    let max_subjects_per_semester = student_progress.preferences.max_subjects_per_semester;
    
    while !remaining_targets.is_empty() && sequence.len() < 16 {
        let mut semester_subjects = Vec::new();
        let mut subjects_to_remove = Vec::new();
        
        for (index, subject_id) in remaining_targets.iter().enumerate() {
            if semester_subjects.len() >= max_subjects_per_semester as usize {
                break;
            }
            
            if let Some(prerequisites) = subject_dependencies.get(subject_id) {
                let prerequisites_met = prerequisites.iter().all(|prereq| {
                    completed_subject_ids.contains(prereq) ||
                    sequence.iter().any(|sem: &SemesterSubjects| sem.subjects.contains(prereq))
                });
                
                if prerequisites_met {
                    semester_subjects.push(subject_id.clone());
                    subjects_to_remove.push(index);
                }
            }
        }
        
        // Remove processed subjects (reverse order to maintain indices)
        for &index in subjects_to_remove.iter().rev() {
            remaining_targets.remove(index);
        }
        
        if !semester_subjects.is_empty() {
            sequence.push(SemesterSubjects {
                semester: current_semester,
                subjects: semester_subjects,
            });
            current_semester += 1;
        } else if !remaining_targets.is_empty() {
            break; // Circular dependency or missing prerequisites
        }
    }
    
    let mut notes = Vec::new();
    
    if !remaining_targets.is_empty() {
        notes.push(format!("Could not sequence {} subjects due to prerequisite constraints", 
                           remaining_targets.len()));
    }
    
    if sequence.len() > 8 {
        notes.push("Long graduation timeline - consider increasing course load".to_string());
    }
    
    Ok(SubjectSequenceResponse {
        recommended_sequence: sequence.clone(),
        total_semesters: sequence.len() as u32,
        notes,
    })
}

pub fn query_workload_analysis(
    deps: Deps, 
    _student_id: String, 
    _semester: u32, 
    proposed_subjects: Vec<String>
) -> StdResult<WorkloadAnalysisResponse> {
    let mut total_credits = 0u32;
    let mut estimated_hours_per_week = 0u32;
    let mut difficulty_scores = Vec::new();
    let mut subject_breakdown = Vec::new();
    let mut recommendations = Vec::new();
    
    for subject_id in &proposed_subjects {
        if let Ok(subject_info) = SUBJECT_SCHEDULE_INFO.load(deps.storage, subject_id) {
            total_credits += subject_info.credits;
            
            // Enhanced analysis using IPFS content if available
            let (weekly_hours, difficulty) = if let Some(ipfs_link) = &subject_info.ipfs_link {
                if let Ok(content) = fetch_ipfs_content(deps, ipfs_link) {
                    (estimate_weekly_workload(&content), analyze_subject_difficulty(&content))
                } else {
                    (subject_info.workload_hours, subject_info.difficulty_level)
                }
            } else {
                (subject_info.workload_hours, subject_info.difficulty_level)
            };
            
            estimated_hours_per_week += weekly_hours;
            difficulty_scores.push(difficulty);
            
            subject_breakdown.push(SubjectWorkload {
                subject_id: subject_id.clone(),
                credits: subject_info.credits,
                estimated_hours: weekly_hours,
                difficulty_contribution: difficulty,
            });
        }
    }
    
    let difficulty_score = if !difficulty_scores.is_empty() {
        difficulty_scores.iter().sum::<u32>() / difficulty_scores.len() as u32
    } else {
        0
    };
    
    let workload_rating = match estimated_hours_per_week {
        0..=15 => WorkloadRating::Light,
        16..=25 => WorkloadRating::Moderate,
        26..=35 => WorkloadRating::Heavy,
        _ => WorkloadRating::Excessive,
    };
    
    // Generate specific recommendations
    match workload_rating {
        WorkloadRating::Light => {
            recommendations.push("Light workload detected".to_string());
            if total_credits < 12 {
                recommendations.push("Consider adding more subjects to maintain good progress".to_string());
            }
        }
        WorkloadRating::Moderate => {
            recommendations.push("Well-balanced workload".to_string());
        }
        WorkloadRating::Heavy => {
            recommendations.push("Heavy workload - ensure good time management".to_string());
        }
        WorkloadRating::Excessive => {
            recommendations.push("Excessive workload - strongly recommend reducing subjects".to_string());
        }
    }
    
    if difficulty_score > 80 {
        recommendations.push("High difficulty semester - allocate extra study time".to_string());
    }
    
    if total_credits > 24 {
        recommendations.push("Very high credit load - verify institutional limits".to_string());
    }
    
    Ok(WorkloadAnalysisResponse {
        total_credits,
        estimated_hours_per_week,
        difficulty_score,
        workload_rating,
        recommendations,
        subject_breakdown,
    })
}

pub fn query_ipfs_cache_status(deps: Deps, ipfs_link: String) -> StdResult<IpfsCacheStatusResponse> {
    let is_cached = is_content_cached(deps, &ipfs_link)?;
    
    let cache_timestamp = if is_cached {
        Some(1640995200u64)
    } else {
        None
    };
    
    Ok(IpfsCacheStatusResponse {
        is_cached,
        ipfs_link,
        cache_timestamp,
    })
}

pub fn query_cached_content(deps: Deps, ipfs_link: String) -> StdResult<SubjectContent> {
    fetch_ipfs_content(deps, &ipfs_link)
}

pub fn query_schedule_statistics(
    deps: Deps, 
    _student_id: Option<String>, 
    _institution_id: Option<String>
) -> StdResult<ScheduleStatisticsResponse> {
    let state = STATE.load(deps.storage)?;
    
    let total_students = state.total_students;
    let total_recommendations_generated = state.total_recommendations;
    let total_paths_created = state.total_paths;
    
    let average_time_to_graduation = if total_paths_created > 0 {
        Some(8u32) // 8 semesters average
    } else {
        None
    };
    
    let most_recommended_subjects = vec![
        PopularSubject {
            subject_id: "CALC101".to_string(),
            recommendation_count: 150,
            completion_rate: 85,
        },
        PopularSubject {
            subject_id: "PROG101".to_string(),
            recommendation_count: 120,
            completion_rate: 90,
        },
    ];
    
    let optimization_success_rate = if total_paths_created > 0 {
        75u32 // 75% success rate
    } else {
        0u32
    };
    
    Ok(ScheduleStatisticsResponse {
        total_students,
        total_recommendations_generated,
        total_paths_created,
        average_time_to_graduation,
        most_recommended_subjects,
        optimization_success_rate,
    })
}

// Helper functions
fn calculate_path_optimization_score(path: &AcademicPath, criteria: &OptimizationCriteria) -> u32 {
    let mut score = 0u32;
    
    match criteria {
        OptimizationCriteria::Fastest => {
            let max_duration = 16u32;
            score += (max_duration - path.total_duration_semesters.min(max_duration)) * 10;
            score += path.path_metrics.average_credits_per_semester.min(30);
        }
        OptimizationCriteria::Balanced => {
            let variance_penalty = path.path_metrics.workload_variance.min(50);
            score += 50 - variance_penalty;
            
            let avg_credits = path.path_metrics.average_credits_per_semester;
            if avg_credits >= 18 && avg_credits <= 20 {
                score += 30;
            }
        }
        OptimizationCriteria::EasiestFirst => {
            score += path.path_metrics.difficulty_progression_score;
        }
        OptimizationCriteria::PrerequisiteOptimal => {
            score += path.path_metrics.prerequisite_efficiency;
            score += path.path_metrics.schedule_flexibility;
        }
        OptimizationCriteria::Custom(_) => {
            score += path.path_metrics.prerequisite_efficiency / 2;
            score += (50 - path.path_metrics.workload_variance.min(50)) / 2;
        }
    }
    
    let risk_penalty = path.path_metrics.risk_factors.len() as u32 * 10;
    score = score.saturating_sub(risk_penalty);
    
    score
}

fn generate_optimization_notes(criteria: &OptimizationCriteria, paths: &[AcademicPath]) -> Vec<String> {
    let mut notes = Vec::new();
    
    if paths.is_empty() {
        notes.push("No existing paths found. Consider creating a new academic path.".to_string());
        return notes;
    }
    
    match criteria {
        OptimizationCriteria::Fastest => {
            notes.push("Paths ranked by graduation speed".to_string());
            let fastest_duration = paths.iter().map(|p| p.total_duration_semesters).min().unwrap_or(0);
            notes.push(format!("Fastest path completes in {} semesters", fastest_duration));
        }
        OptimizationCriteria::Balanced => {
            notes.push("Paths ranked by workload balance".to_string());
        }
        OptimizationCriteria::EasiestFirst => {
            notes.push("Paths ranked by difficulty progression".to_string());
        }
        OptimizationCriteria::PrerequisiteOptimal => {
            notes.push("Paths ranked by prerequisite efficiency".to_string());
        }
        OptimizationCriteria::Custom(_) => {
            notes.push("Paths ranked using custom optimization weights".to_string());
        }
    }
    
    notes
}