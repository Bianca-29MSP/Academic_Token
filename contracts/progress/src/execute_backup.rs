use cosmwasm_std::{DepsMut, Env, MessageInfo, Response, StdResult, Order, Storage};
use crate::error::ContractError;
use crate::state::*;
use crate::msg::AnalyticsRefreshMode;
use crate::state::TrendDirection;
// use crate::analytics::AnalyticsEngine; // Commented out as it doesn't exist yet

pub fn update_config(
    deps: DepsMut,
    info: MessageInfo,
    new_owner: Option<String>,
    analytics_enabled: Option<bool>,
    update_frequency: Option<UpdateFrequency>,
    analytics_depth: Option<AnalyticsDepth>,
    retention_period_days: Option<u32>,
) -> Result<Response, ContractError> {
    let mut state = STATE.load(deps.storage)?;
    
    // Check if sender is owner
    if info.sender != state.owner {
        return Err(ContractError::Unauthorized {
            owner: state.owner.to_string(),
        });
    }
    
    // Update owner if provided
    if let Some(new_owner_addr) = new_owner {
        let validated_addr = deps.api.addr_validate(&new_owner_addr)?;
        state.owner = validated_addr;
    }
    
    // Update analytics enabled flag
    if let Some(enabled) = analytics_enabled {
        state.analytics_enabled = enabled;
    }
    
    STATE.save(deps.storage, &state)?;
    
    // Update config
    let mut config = PROGRESS_CONFIG.load(deps.storage)?;
    
    if let Some(frequency) = update_frequency {
        config.update_frequency = frequency;
    }
    
    if let Some(depth) = analytics_depth {
        config.analytics_depth = depth;
    }
    
    if let Some(retention) = retention_period_days {
        config.retention_period_days = retention;
    }
    
    PROGRESS_CONFIG.save(deps.storage, &config)?;
    
    Ok(Response::new()
        .add_attribute("method", "update_config")
        .add_attribute("updated_by", info.sender))
}

pub fn update_student_progress(
    deps: DepsMut,
    env: Env,
    student_progress: StudentProgress,
    force_analytics_refresh: Option<bool>,
) -> Result<Response, ContractError> {
    // Validate student progress
    if student_progress.student_id.is_empty() {
        return Err(ContractError::StudentNotFound {
            student_id: "empty".to_string(),
        });
    }
    
    if student_progress.institution_id.is_empty() {
        return Err(ContractError::InstitutionNotFound {
            institution_id: "empty".to_string(),
        });
    }
    
    // Check if student exists
    let existing = STUDENT_PROGRESS.may_load(deps.storage, &student_progress.student_id)?;
    let is_new_student = existing.is_none();
    
    // Update indices for new students
    if is_new_student {
        // Add to institution index
        STUDENTS_BY_INSTITUTION.update(
            deps.storage,
            &student_progress.institution_id,
            |students| -> StdResult<Vec<String>> {
                let mut list = students.unwrap_or_default();
                if !list.contains(&student_progress.student_id) {
                    list.push(student_progress.student_id.clone());
                }
                Ok(list)
            },
        )?;
        
        // Add to course index
        STUDENTS_BY_COURSE.update(
            deps.storage,
            &student_progress.course_id,
            |students| -> StdResult<Vec<String>> {
                let mut list = students.unwrap_or_default();
                if !list.contains(&student_progress.student_id) {
                    list.push(student_progress.student_id.clone());
                }
                Ok(list)
            },
        )?;
        
        // Update total count
        STATE.update(deps.storage, |mut state| -> StdResult<_> {
            state.total_students += 1;
            if state.total_institutions == 0 || 
               !STUDENTS_BY_INSTITUTION.has(deps.storage, &student_progress.institution_id) {
                state.total_institutions += 1;
            }
            state.total_progress_records += 1;
            Ok(state)
        })?;
    }
    
    // Remove from old status index if updating existing student
    if let Some(old_progress) = existing {
        let old_status = format!("{:?}", old_progress.academic_status);
        STUDENTS_BY_STATUS.update(
            deps.storage,
            &old_status,
            |students| -> StdResult<Vec<String>> {
                let mut list = students.unwrap_or_default();
                list.retain(|id| id != &student_progress.student_id);
                Ok(list)
            },
        )?;
        
        // Update risk level index
        let old_risk_level = determine_risk_level(&old_progress);
        STUDENTS_BY_RISK_LEVEL.update(
            deps.storage,
            &format!("{:?}", old_risk_level),
            |students| -> StdResult<Vec<String>> {
                let mut list = students.unwrap_or_default();
                list.retain(|id| id != &student_progress.student_id);
                Ok(list)
            },
        )?;
    }
    
    // Save updated progress with timestamp
    let mut updated_progress = student_progress;
    updated_progress.last_updated = env.block.time.seconds();
    
    // TODO: Calculate performance metrics
    // updated_progress.performance_metrics = AnalyticsEngine::calculate_performance_metrics(
    //     &updated_progress, 
    //     None
    // )?;
    
    // TODO: Generate graduation forecast
    // updated_progress.graduation_forecast = AnalyticsEngine::generate_graduation_forecast(&updated_progress)?;
    
    // TODO: Identify risk factors
    // updated_progress.risk_factors = AnalyticsEngine::identify_risk_factors(&updated_progress);
    
    STUDENT_PROGRESS.save(deps.storage, &updated_progress.student_id, &updated_progress)?;
    
    // Store historical snapshot if configured
    let config = PROGRESS_CONFIG.load(deps.storage)?;
    if matches!(config.analytics_depth, AnalyticsDepth::Advanced | AnalyticsDepth::Comprehensive) {
        PROGRESS_HISTORY.save(
            deps.storage,
            (&updated_progress.student_id, env.block.time.seconds()),
            &updated_progress,
        )?;
    }
    
    // Update status index
    let new_status = format!("{:?}", updated_progress.academic_status);
    STUDENTS_BY_STATUS.update(
        deps.storage,
        &new_status,
        |students| -> StdResult<Vec<String>> {
            let mut list = students.unwrap_or_default();
            if !list.contains(&updated_progress.student_id) {
                list.push(updated_progress.student_id.clone());
            }
            Ok(list)
        },
    )?;
    
    // Update risk level index
    let new_risk_level = determine_risk_level(&updated_progress);
    STUDENTS_BY_RISK_LEVEL.update(
        deps.storage,
        &format!("{:?}", new_risk_level),
        |students| -> StdResult<Vec<String>> {
            let mut list = students.unwrap_or_default();
            if !list.contains(&updated_progress.student_id) {
                list.push(updated_progress.student_id.clone());
            }
            Ok(list)
        },
    )?;
    
    // Generate analytics if enabled and forced
    let state = STATE.load(deps.storage)?;
    if state.analytics_enabled && force_analytics_refresh == Some(true) {
        // Trigger analytics update for this student's institution
        let period = AnalyticsPeriod {
            start_date: "current_year".to_string(),
            end_date: "current_year".to_string(),
            period_type: PeriodType::AcademicYear,
        };
        
        // Note: This would ideally be handled differently to avoid ownership issues
        // For now, we'll skip the analytics generation in this context
    }
    
    Ok(Response::new()
        .add_attribute("method", "update_student_progress")
        .add_attribute("student_id", updated_progress.student_id)
        .add_attribute("institution_id", updated_progress.institution_id)
        .add_attribute("is_new_student", is_new_student.to_string())
        .add_attribute("timestamp", env.block.time.seconds().to_string()))
}

pub fn batch_update_student_progress(
    deps: DepsMut,
    env: Env,
    student_progress_list: Vec<StudentProgress>,
    analytics_refresh_mode: AnalyticsRefreshMode,
) -> Result<Response, ContractError> {
    let mut updated_count = 0u32;
    let mut new_students_count = 0u32;
    let mut affected_institutions = std::collections::HashSet::new();
    
    // Process each student individually to avoid ownership issues
    for student_progress in student_progress_list {
        if student_progress.student_id.is_empty() {
            continue; // Skip invalid entries
        }
        
        let is_new = STUDENT_PROGRESS.may_load(deps.storage, &student_progress.student_id)?.is_none();
        if is_new {
            new_students_count += 1;
        }
        
        affected_institutions.insert(student_progress.institution_id.clone());
        
        // Process each student update separately
        match update_single_student_progress_internal(deps.storage, &env, student_progress) {
            Ok(_) => updated_count += 1,
            Err(_) => continue, // Continue with next student on error
        }
    }
    
    // Handle analytics refresh based on mode
    match analytics_refresh_mode {
        AnalyticsRefreshMode::None => {
            // Do nothing
        },
        AnalyticsRefreshMode::Individual => {
            // Already handled in individual updates
        },
        AnalyticsRefreshMode::Batch => {
            // Refresh analytics for all affected institutions
            let period = AnalyticsPeriod {
                start_date: "current_year".to_string(),
                end_date: "current_year".to_string(),
                period_type: PeriodType::AcademicYear,
            };
            
            for institution_id in affected_institutions.iter() {
                let _ = generate_institution_analytics_internal(
                    deps,
                    env,
                    institution_id.clone(),
                    period.clone(),
                    true,
                );
            }
        },
        AnalyticsRefreshMode::Background => {
            // In a real implementation, this would schedule background tasks
            // For now, just mark that background refresh is needed
        },
    }
    
    Ok(Response::new()
        .add_attribute("method", "batch_update_student_progress")
        .add_attribute("updated_count", updated_count.to_string())
        .add_attribute("new_students", new_students_count.to_string())
        .add_attribute("affected_institutions", affected_institutions.len().to_string()))
}

pub fn record_subject_completion(
    deps: DepsMut,
    env: Env,
    student_id: String,
    subject_id: String,
    final_grade: u32,
    completion_date: String,
    study_hours: Option<u32>,
    difficulty_rating: Option<u32>,
    satisfaction_rating: Option<u32>,
    nft_token_id: String,
) -> Result<Response, ContractError> {
    STUDENT_PROGRESS.update(deps.storage, &student_id, |progress| -> Result<_, ContractError> {
        let mut updated_progress = progress.ok_or(ContractError::StudentNotFound { 
            student_id: student_id.clone() 
        })?;
        
        // Remove from current subjects if present
        let was_current = updated_progress.current_subjects
            .iter()
            .find(|s| s.subject_id == subject_id)
            .map(|s| s.credits)
            .unwrap_or(0);
        
        updated_progress.current_subjects.retain(|s| s.subject_id != subject_id);
        
        // Determine credits (assume 4 if not found in current subjects)
        let credits = if was_current > 0 { was_current } else { 4 };
        
        // Create completed subject record
        let completed_subject = CompletedSubjectProgress {
            subject_id: subject_id.clone(),
            title: subject_id.clone(), // Would ideally fetch from subject registry
            credits,
            final_grade,
            letter_grade: grade_to_letter(final_grade),
            completion_date: completion_date.clone(),
            semester_taken: updated_progress.current_semester,
            attempt_number: calculate_attempt_number(&updated_progress, &subject_id),
            study_hours_logged: study_hours,
            difficulty_rating,
            satisfaction_rating,
            nft_token_id: nft_token_id.clone(),
            competencies_gained: vec![], // Would be derived from subject content
        };
        
        // Check if this is a retake of a failed subject
        let was_failed = updated_progress.failed_subjects
            .iter()
            .position(|f| f.subject_id == subject_id);
        
        if final_grade >= 60 {
            // Passing grade
            updated_progress.completed_subjects.push(completed_subject);
            updated_progress.total_credits_completed += credits;
            updated_progress.credits_in_progress = updated_progress.credits_in_progress.saturating_sub(credits);
            
            // Remove from failed subjects if it was a retake
            if let Some(failed_index) = was_failed {
                updated_progress.failed_subjects.remove(failed_index);
            }
            
            // Update grade distribution
            match final_grade {
                90..=100 => updated_progress.grade_distribution.excellent_count += 1,
                80..=89 => updated_progress.grade_distribution.good_count += 1,
                70..=79 => updated_progress.grade_distribution.satisfactory_count += 1,
                60..=69 => updated_progress.grade_distribution.poor_count += 1,
                _ => {} // This shouldn't happen for passing grades
            }
        } else {
            // Failing grade
            let failed_subject = FailedSubjectProgress {
                subject_id: subject_id.clone(),
                title: subject_id.clone(),
                credits,
                final_grade,
                completion_date: completion_date.clone(),
                semester_taken: updated_progress.current_semester,
                attempt_number: completed_subject.attempt_number,
                failure_reasons: vec!["Low grade".to_string()],
                remediation_required: true,
                retry_available: completed_subject.attempt_number < 3,
                retry_recommendations: vec![
                    "Review course material".to_string(),
                    "Seek tutoring support".to_string(),
                ],
            };
            
            // Remove previous failed attempt if exists
            if let Some(failed_index) = was_failed {
                updated_progress.failed_subjects[failed_index] = failed_subject;
            } else {
                updated_progress.failed_subjects.push(failed_subject);
            }
            
            updated_progress.grade_distribution.fail_count += 1;
            updated_progress.credits_in_progress = updated_progress.credits_in_progress.saturating_sub(credits);
        }
        
        // Recalculate GPA
        update_gpa(&mut updated_progress);
        
        // Update academic status based on performance
        update_academic_status_internal(&mut updated_progress);
        
        // Update last modified timestamp
        updated_progress.last_updated = env.block.time.seconds();
        
        Ok(updated_progress)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "record_subject_completion")
        .add_attribute("student_id", student_id)
        .add_attribute("subject_id", subject_id)
        .add_attribute("final_grade", final_grade.to_string())
        .add_attribute("passed", (final_grade >= 60).to_string())
        .add_attribute("nft_token_id", nft_token_id))
}

pub fn record_subject_enrollment(
    deps: DepsMut,
    student_id: String,
    subject_id: String,
    enrollment_date: String,
    _expected_completion: String,
) -> Result<Response, ContractError> {
    STUDENT_PROGRESS.update(deps.storage, &student_id, |progress| -> Result<_, ContractError> {
        let mut updated_progress = progress.ok_or(ContractError::StudentNotFound { 
            student_id: student_id.clone() 
        })?;
        
        // Check if already enrolled
        if updated_progress.current_subjects.iter().any(|s| s.subject_id == subject_id) {
            return Err(ContractError::InvalidAnalyticsConfig {
                reason: "Student already enrolled in this subject".to_string(),
            });
        }
        
        // Check if already completed
        if updated_progress.completed_subjects.iter().any(|s| s.subject_id == subject_id) {
            return Err(ContractError::InvalidAnalyticsConfig {
                reason: "Subject already completed".to_string(),
            });
        }
        
        // Create enrollment record
        let enrolled_subject = CurrentSubjectProgress {
            subject_id: subject_id.clone(),
            title: subject_id.clone(), // Would fetch from subject registry
            credits: 4, // Default, would fetch from subject info
            enrollment_date: enrollment_date.clone(),
            current_grade: None,
            attendance_rate: None,
            assignments_completed: 0,
            assignments_total: 10, // Default, would fetch from subject info
            study_hours_logged: 0,
            predicted_final_grade: None,
            risk_level: RiskLevel::Low, // Default
            last_activity_date: Some(enrollment_date.clone()),
        };
        
        updated_progress.current_subjects.push(enrolled_subject);
        updated_progress.credits_in_progress += 4; // Default credits
        
        Ok(updated_progress)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "record_subject_enrollment")
        .add_attribute("student_id", student_id)
        .add_attribute("subject_id", subject_id)
        .add_attribute("enrollment_date", enrollment_date))
}

pub fn update_current_subject_progress(
    deps: DepsMut,
    env: Env,
    student_id: String,
    subject_id: String,
    current_grade: Option<u32>,
    attendance_rate: Option<u32>,
    assignments_completed: u32,
    study_hours_logged: u32,
) -> Result<Response, ContractError> {
    STUDENT_PROGRESS.update(deps.storage, &student_id, |progress| -> Result<_, ContractError> {
        let mut updated_progress = progress.ok_or(ContractError::StudentNotFound { 
            student_id: student_id.clone() 
        })?;
        
        // Find and update the current subject
        let subject_found = updated_progress.current_subjects
            .iter_mut()
            .find(|s| s.subject_id == subject_id);
        
        if let Some(subject) = subject_found {
            if let Some(grade) = current_grade {
                subject.current_grade = Some(grade);
                // Update predicted final grade based on current performance
                subject.predicted_final_grade = Some(predict_final_grade(grade, assignments_completed, subject.assignments_total));
            }
            
            if let Some(attendance) = attendance_rate {
                subject.attendance_rate = Some(attendance);
            }
            
            subject.assignments_completed = assignments_completed;
            subject.study_hours_logged = study_hours_logged;
            subject.last_activity_date = Some(env.block.time.seconds().to_string());
            
            // Update risk level based on performance indicators
            subject.risk_level = calculate_subject_risk_level(
                subject.current_grade,
                subject.attendance_rate,
                subject.assignments_completed,
                subject.assignments_total,
            );
            
            updated_progress.last_updated = env.block.time.seconds();
            
            Ok(updated_progress)
        } else {
            Err(ContractError::InvalidAnalyticsConfig {
                reason: "Subject not found in current enrollments".to_string(),
            })
        }
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "update_current_subject_progress")
        .add_attribute("student_id", student_id)
        .add_attribute("subject_id", subject_id)
        .add_attribute("assignments_completed", assignments_completed.to_string()))
}

pub fn record_milestone_achievement(
    deps: DepsMut,
    env: Env,
    student_id: String,
    milestone_id: String,
    achievement_date: String,
    notes: Option<String>,
) -> Result<Response, ContractError> {
    STUDENT_PROGRESS.update(deps.storage, &student_id, |progress| -> Result<_, ContractError> {
        let mut updated_progress = progress.ok_or(ContractError::StudentNotFound { 
            student_id: student_id.clone() 
        })?;
        
        // Check if milestone already achieved
        if updated_progress.milestones_achieved.iter().any(|m| m.milestone_id == milestone_id) {
            return Err(ContractError::InvalidAnalyticsConfig {
                reason: "Milestone already achieved".to_string(),
            });
        }
        
        // Determine milestone type and significance
        let (milestone_type, significance) = determine_milestone_type_and_significance(
            &milestone_id,
            &updated_progress,
        );
        
        let milestone = AcademicMilestone {
            milestone_id: milestone_id.clone(),
            milestone_type,
            title: format!("Milestone: {}", milestone_id),
            description: notes.unwrap_or_else(|| "Academic milestone achieved".to_string()),
            achieved_date: achievement_date.clone(),
            semester_achieved: updated_progress.current_semester,
            credits_at_achievement: updated_progress.total_credits_completed,
            gpa_at_achievement: updated_progress.gpa.clone(),
            significance,
        };
        
        updated_progress.milestones_achieved.push(milestone);
        
        // Remove from upcoming milestones if present
        updated_progress.upcoming_milestones.retain(|m| m.milestone_id != milestone_id);
        
        // Generate new upcoming milestones
        generate_upcoming_milestones(&mut updated_progress);
        
        updated_progress.last_updated = env.block.time.seconds();
        
        Ok(updated_progress)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "record_milestone_achievement")
        .add_attribute("student_id", student_id)
        .add_attribute("milestone_id", milestone_id)
        .add_attribute("achievement_date", achievement_date))
}

pub fn add_risk_factor(
    deps: DepsMut,
    env: Env,
    student_id: String,
    risk_type: RiskFactorType,
    severity: RiskSeverity,
    description: String,
    intervention_recommendations: Vec<String>,
) -> Result<Response, ContractError> {
    STUDENT_PROGRESS.update(deps.storage, &student_id, |progress| -> Result<_, ContractError> {
        let mut updated_progress = progress.ok_or(ContractError::StudentNotFound { 
            student_id: student_id.clone() 
        })?;
        
        // Remove existing risk factor of the same type if present
        updated_progress.risk_factors.retain(|rf| rf.risk_type != risk_type);
        
        let monitoring_frequency = match severity {
            RiskSeverity::Critical => MonitoringFrequency::Weekly,
            RiskSeverity::High => MonitoringFrequency::Biweekly,
            RiskSeverity::Moderate => MonitoringFrequency::Monthly,
            RiskSeverity::Low => MonitoringFrequency::Quarterly,
        };
        
        let risk_factor = RiskFactor {
            risk_type: risk_type.clone(),
            severity: severity.clone(),
            description: description.clone(),
            early_warning_indicators: generate_early_warning_indicators(&risk_type),
            intervention_recommendations,
            monitoring_frequency,
        };
        
        updated_progress.risk_factors.push(risk_factor);
        updated_progress.last_updated = env.block.time.seconds();
        
        Ok(updated_progress)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "add_risk_factor")
        .add_attribute("student_id", student_id)
        .add_attribute("risk_type", format!("{:?}", risk_type))
        .add_attribute("severity", format!("{:?}", severity)))
}

pub fn remove_risk_factor(
    deps: DepsMut,
    student_id: String,
    risk_type: RiskFactorType,
) -> Result<Response, ContractError> {
    STUDENT_PROGRESS.update(deps.storage, &student_id, |progress| -> Result<_, ContractError> {
        let mut updated_progress = progress.ok_or(ContractError::StudentNotFound { 
            student_id: student_id.clone() 
        })?;
        
        let initial_count = updated_progress.risk_factors.len();
        updated_progress.risk_factors.retain(|rf| rf.risk_type != risk_type);
        let final_count = updated_progress.risk_factors.len();
        
        if initial_count == final_count {
            return Err(ContractError::InvalidAnalyticsConfig {
                reason: "Risk factor not found".to_string(),
            });
        }
        
        Ok(updated_progress)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "remove_risk_factor")
        .add_attribute("student_id", student_id)
        .add_attribute("risk_type", format!("{:?}", risk_type)))
}

pub fn generate_institution_analytics(
    deps: DepsMut,
    env: Env,
    institution_id: String,
    period: AnalyticsPeriod,
    force_refresh: Option<bool>,
) -> Result<Response, ContractError> {
    let analytics = generate_institution_analytics_internal(
        &mut deps,
        &env,
        institution_id.clone(),
        period,
        force_refresh.unwrap_or(false),
    )?;
    
    Ok(Response::new()
        .add_attribute("method", "generate_institution_analytics")
        .add_attribute("institution_id", institution_id)
        .add_attribute("total_students", analytics.total_students.to_string())
        .add_attribute("average_gpa", analytics.average_gpa))
}

pub fn generate_student_dashboard(
    mut deps: DepsMut,
    env: Env,
    student_id: String,
    include_peer_comparison: Option<bool>,
    force_refresh: Option<bool>,
) -> Result<Response, ContractError> {
    // Check if dashboard exists and not forcing refresh
    if force_refresh != Some(true) {
        if let Ok(_existing) = STUDENT_DASHBOARDS.load(deps.storage, &student_id) {
            return Ok(Response::new()
                .add_attribute("method", "generate_student_dashboard")
                .add_attribute("student_id", student_id)
                .add_attribute("cached", "true"));
        }
    }
    
    // Load student progress
    let student_progress = STUDENT_PROGRESS.load(deps.storage, &student_id)
        .map_err(|_| ContractError::StudentNotFound { student_id: student_id.clone() })?;
    
    // Generate dashboard components
    let quick_stats = QuickStats {
        current_gpa: student_progress.gpa.clone(),
        credits_completed: student_progress.total_credits_completed,
        credits_remaining: student_progress.total_credits_required
            .saturating_sub(student_progress.total_credits_completed),
        completion_percentage: if student_progress.total_credits_required > 0 {
            (student_progress.total_credits_completed * 100) / student_progress.total_credits_required
        } else {
            0
        },
        semesters_remaining: student_progress.graduation_forecast.remaining_semesters,
        current_semester_subjects: student_progress.current_subjects.len() as u32,
    };
    
    let current_status = CurrentStatus {
        academic_standing: student_progress.academic_status.clone(),
        current_semester: student_progress.current_semester,
        enrollment_status: if student_progress.current_subjects.is_empty() {
            EnrollmentStatus::OnLeave
        } else {
            EnrollmentStatus::FullTime
        },
        next_milestone: student_progress.upcoming_milestones
            .first()
            .map(|m| m.title.clone()),
        days_to_next_milestone: None, // Would require date parsing
    };
    
    // Generate grade history
    let grade_history = generate_grade_history(&student_progress);
    
    // Generate progress timeline
    let progress_timeline = generate_progress_timeline(&student_progress);
    
    // Generate alerts
    let urgent_alerts = generate_urgent_alerts(&student_progress, &env);
    
    // Generate upcoming deadlines
    let upcoming_deadlines = generate_upcoming_deadlines(&student_progress);
    
    // Generate recommendations
    let recommendations = generate_dashboard_recommendations(&student_progress);
    
    // Generate recent achievements
    let recent_achievements = generate_recent_achievements(&student_progress);
    
    // Generate goal progress
    let progress_towards_goals = generate_goal_progress(&student_progress);
    
    // Generate performance insights
    let performance_insights = generate_performance_insights(&student_progress);
    
    // Generate risk warnings
    let risk_warnings = student_progress.risk_factors.iter().map(|rf| RiskWarning {
        warning_id: format!("{}_{:?}", student_id, rf.risk_type),
        risk_type: rf.risk_type.clone(),
        severity: rf.severity.clone(),
        probability: 70, // Would be calculated based on historical data
        description: rf.description.clone(),
        early_indicators: rf.early_warning_indicators.clone(),
        prevention_strategies: rf.intervention_recommendations.clone(),
        support_resources: vec!["Academic Support Center".to_string()],
    }).collect();
    
    // Generate peer comparison if requested
    let peer_comparisons = if include_peer_comparison == Some(true) {
        Some(generate_peer_comparison(deps.storage, &student_progress)?)
    } else {
        None
    };
    
    let dashboard = StudentDashboard {
        student_id: student_id.clone(),
        last_updated: env.block.time.seconds(),
        quick_stats,
        current_status,
        grade_history,
        progress_timeline,
        urgent_alerts,
        upcoming_deadlines,
        recommendations,
        recent_achievements,
        progress_towards_goals,
        performance_insights,
        risk_warnings,
        peer_comparisons,
    };
    
    // Save dashboard
    STUDENT_DASHBOARDS.save(deps.storage, &student_id, &dashboard)?;
    
    Ok(Response::new()
        .add_attribute("method", "generate_student_dashboard")
        .add_attribute("student_id", student_id)
        .add_attribute("components_generated", "full")
        .add_attribute("peer_comparison", include_peer_comparison.unwrap_or(false).to_string()))
}

// Continue with remaining execute functions...

pub fn batch_generate_dashboards(
    deps: DepsMut,
    env: Env,
    institution_id: String,
    student_ids: Option<Vec<String>>,
    include_analytics: bool,
) -> Result<Response, ContractError> {
    let students_to_process = if let Some(ids) = student_ids {
        ids
    } else {
        // Get all students from institution
        STUDENTS_BY_INSTITUTION.load(deps.storage, &institution_id).unwrap_or_default()
    };
    
    let mut generated_count = 0u32;
    let mut errors_count = 0u32;
    
    for student_id in students_to_process {
        match generate_student_dashboard(&mut deps, &env, student_id, Some(include_analytics), Some(true)) {
            Ok(_) => generated_count += 1,
            Err(_) => errors_count += 1,
        }
    }
    
    Ok(Response::new()
        .add_attribute("method", "batch_generate_dashboards")
        .add_attribute("institution_id", institution_id)
        .add_attribute("generated_count", generated_count.to_string())
        .add_attribute("errors_count", errors_count.to_string()))
}

pub fn run_global_analytics(
    deps: DepsMut,
    env: Env,
    period: AnalyticsPeriod,
    institutions: Option<Vec<String>>,
) -> Result<Response, ContractError> {
    let institutions_to_process = if let Some(inst_list) = institutions {
        inst_list
    } else {
        // Get all institutions
        get_all_institutions(deps.storage)?
    };
    
    let mut processed_count = 0u32;
    
    for institution_id in institutions_to_process {
        let _ = generate_institution_analytics_internal(
            &mut deps,
            &env,
            institution_id,
            period.clone(),
            true,
        );
        processed_count += 1;
    }
    
    Ok(Response::new()
        .add_attribute("method", "run_global_analytics")
        .add_attribute("processed_institutions", processed_count.to_string())
        .add_attribute("period_type", format!("{:?}", period.period_type)))
}

pub fn update_academic_status(
    deps: DepsMut,
    env: Env,
    student_id: String,
    new_status: AcademicStatus,
    effective_date: String,
    _reason: Option<String>,
) -> Result<Response, ContractError> {
    let old_status = {
        let progress = STUDENT_PROGRESS.load(deps.storage, &student_id)
            .map_err(|_| ContractError::StudentNotFound { student_id: student_id.clone() })?;
        progress.academic_status.clone()
    };
    
    STUDENT_PROGRESS.update(deps.storage, &student_id, |progress| -> Result<_, ContractError> {
        let mut updated_progress = progress.ok_or(ContractError::StudentNotFound { 
            student_id: student_id.clone() 
        })?;
        
        updated_progress.academic_status = new_status.clone();
        updated_progress.last_updated = env.block.time.seconds();
        
        Ok(updated_progress)
    })?;
    
    // Update status indices
    let old_status_key = format!("{:?}", old_status);
    let new_status_key = format!("{:?}", new_status);
    
    // Remove from old status index
    STUDENTS_BY_STATUS.update(
        deps.storage,
        &old_status_key,
        |students| -> StdResult<Vec<String>> {
            let mut list = students.unwrap_or_default();
            list.retain(|id| id != &student_id);
            Ok(list)
        },
    )?;
    
    // Add to new status index
    STUDENTS_BY_STATUS.update(
        deps.storage,
        &new_status_key,
        |students| -> StdResult<Vec<String>> {
            let mut list = students.unwrap_or_default();
            if !list.contains(&student_id) {
                list.push(student_id.clone());
            }
            Ok(list)
        },
    )?;
    
    Ok(Response::new()
        .add_attribute("method", "update_academic_status")
        .add_attribute("student_id", student_id)
        .add_attribute("old_status", format!("{:?}", old_status))
        .add_attribute("new_status", format!("{:?}", new_status))
        .add_attribute("effective_date", effective_date))
}

pub fn archive_old_data(
    deps: DepsMut,
    cutoff_date: String,
    dry_run: Option<bool>,
) -> Result<Response, ContractError> {
    let is_dry_run = dry_run.unwrap_or(false);
    
    // Parse cutoff date (simplified - would use proper date parsing)
    let cutoff_timestamp = cutoff_date.parse::<u64>().unwrap_or(0);
    
    let mut archived_count = 0u32;
    let mut total_checked = 0u32;
    
    // Archive old progress history entries
    let history_entries: Vec<_> = PROGRESS_HISTORY
        .range(deps.storage, None, None, Order::Ascending)
        .take(1000) // Limit for gas management
        .collect::<StdResult<Vec<_>>>()?;
    
    for ((student_id, timestamp), _progress) in history_entries {
        total_checked += 1;
        
        if timestamp < cutoff_timestamp {
            if !is_dry_run {
                PROGRESS_HISTORY.remove(deps.storage, (&student_id, timestamp));
            }
            archived_count += 1;
        }
    }
    
    Ok(Response::new()
        .add_attribute("method", "archive_old_data")
        .add_attribute("cutoff_date", cutoff_date)
        .add_attribute("dry_run", is_dry_run.to_string())
        .add_attribute("total_checked", total_checked.to_string())
        .add_attribute("archived_count", archived_count.to_string()))
}

pub fn reset_student_progress(
    deps: DepsMut,
    info: MessageInfo,
    student_id: String,
    reason: String,
    backup: bool,
) -> Result<Response, ContractError> {
    let state = STATE.load(deps.storage)?;
    
    // Check if sender is owner (admin function)
    if info.sender != state.owner {
        return Err(ContractError::Unauthorized {
            owner: state.owner.to_string(),
        });
    }
    
    // Load existing progress for backup if requested
    let existing_progress = STUDENT_PROGRESS.load(deps.storage, &student_id)
        .map_err(|_| ContractError::StudentNotFound { student_id: student_id.clone() })?;
    
    if backup {
        // Store backup with special timestamp
        let backup_timestamp = 999999999u64; // Special marker for backups
        PROGRESS_HISTORY.save(
            deps.storage,
            (&student_id, backup_timestamp),
            &existing_progress,
        )?;
    }
    
    // Remove from all indices
    let institution_id = existing_progress.institution_id.clone();
    let course_id = existing_progress.course_id.clone();
    let status_key = format!("{:?}", existing_progress.academic_status);
    
    // Remove from institution index
    STUDENTS_BY_INSTITUTION.update(
        deps.storage,
        &institution_id,
        |students| -> StdResult<Vec<String>> {
            let mut list = students.unwrap_or_default();
            list.retain(|id| id != &student_id);
            Ok(list)
        },
    )?;
    
    // Remove from course index
    STUDENTS_BY_COURSE.update(
        deps.storage,
        &course_id,
        |students| -> StdResult<Vec<String>> {
            let mut list = students.unwrap_or_default();
            list.retain(|id| id != &student_id);
            Ok(list)
        },
    )?;
    
    // Remove from status index
    STUDENTS_BY_STATUS.update(
        deps.storage,
        &status_key,
        |students| -> StdResult<Vec<String>> {
            let mut list = students.unwrap_or_default();
            list.retain(|id| id != &student_id);
            Ok(list)
        },
    )?;
    
    // Remove main progress record
    STUDENT_PROGRESS.remove(deps.storage, &student_id);
    
    // Remove dashboard
    STUDENT_DASHBOARDS.remove(deps.storage, &student_id);
    
    // Update state counters
    STATE.update(deps.storage, |mut state| -> StdResult<_> {
        state.total_students = state.total_students.saturating_sub(1);
        state.total_progress_records = state.total_progress_records.saturating_sub(1);
        Ok(state)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "reset_student_progress")
        .add_attribute("student_id", student_id)
        .add_attribute("reason", reason)
        .add_attribute("backup_created", backup.to_string())
        .add_attribute("reset_by", info.sender))
}

// Helper functions

fn update_single_student_progress_internal(
    storage: &mut dyn Storage,
    env: &Env,
    student_progress: StudentProgress,
) -> Result<(), ContractError> {
    // Simple implementation for batch processing
    let mut updated_progress = student_progress;
    updated_progress.last_updated = env.block.time.seconds();
    
    STUDENT_PROGRESS.save(storage, &updated_progress.student_id, &updated_progress)?;
    Ok(())
}

fn determine_risk_level(progress: &StudentProgress) -> RiskSeverity {
    let mut risk_score = 0u32;
    
    // Check GPA
    if let Ok(gpa) = progress.gpa.parse::<f64>() {
        if gpa < 6.0 {
            risk_score += 40;
        } else if gpa < 7.0 {
            risk_score += 20;
        }
    }
    
    // Check failed subjects
    if progress.failed_subjects.len() > 2 {
        risk_score += 30;
    } else if progress.failed_subjects.len() > 0 {
        risk_score += 15;
    }
    
    // Check academic status
    match progress.academic_status {
        AcademicStatus::Probation => risk_score += 25,
        AcademicStatus::Suspended => risk_score += 50,
        _ => {}
    }
    
    // Check current subject performance
    let at_risk_subjects = progress.current_subjects.iter()
        .filter(|s| matches!(s.risk_level, RiskLevel::High | RiskLevel::Critical))
        .count();
    
    if at_risk_subjects > 1 {
        risk_score += 20;
    } else if at_risk_subjects > 0 {
        risk_score += 10;
    }
    
    match risk_score {
        0..=20 => RiskSeverity::Low,
        21..=40 => RiskSeverity::Moderate,
        41..=70 => RiskSeverity::High,
        _ => RiskSeverity::Critical,
    }
}

// Helper functions for execute.rs (continuação)

fn grade_to_letter(grade: u32) -> String {
    match grade {
        90..=100 => "A".to_string(),
        80..=89 => "B".to_string(),
        70..=79 => "C".to_string(),
        60..=69 => "D".to_string(),
        _ => "F".to_string(),
    }
}

fn calculate_attempt_number(progress: &StudentProgress, subject_id: &str) -> u32 {
    let failed_attempts = progress.failed_subjects.iter()
        .filter(|f| f.subject_id == subject_id)
        .count() as u32;
    
    failed_attempts + 1
}

fn update_gpa(progress: &mut StudentProgress) {
    let total_grade_points: u32 = progress.completed_subjects.iter()
        .map(|s| s.final_grade * s.credits)
        .sum();
    
    let total_credits: u32 = progress.completed_subjects.iter()
        .map(|s| s.credits)
        .sum();
    
    if total_credits > 0 {
        let gpa = (total_grade_points as f64) / (total_credits as f64) / 10.0;
        progress.gpa = format!("{:.2}", gpa);
        
        // Update CGPA (same as GPA in this simplified version)
        progress.cgpa = progress.gpa.clone();
    }
}

fn update_academic_status_internal(progress: &mut StudentProgress) {
    if let Ok(gpa) = progress.gpa.parse::<f64>() {
        progress.academic_status = match gpa {
            _ if gpa >= 9.0 => AcademicStatus::Active, // Excellent standing
            _ if gpa >= 7.0 => AcademicStatus::Active, // Good standing
            _ if gpa >= 6.0 => AcademicStatus::Probation, // Academic probation
            _ => AcademicStatus::Suspended, // Academic suspension risk
        };
    }
    
    // Check if graduated
    if progress.total_credits_completed >= progress.total_credits_required && 
       progress.total_credits_required > 0 {
        progress.academic_status = AcademicStatus::Graduated;
    }
}

fn predict_final_grade(current_grade: u32, assignments_completed: u32, assignments_total: u32) -> u32 {
    if assignments_total == 0 {
        return current_grade;
    }
    
    let completion_rate = (assignments_completed as f64) / (assignments_total as f64);
    let base_prediction = current_grade as f64;
    
    // Adjust prediction based on completion rate
    let adjustment = match completion_rate {
        rate if rate >= 0.9 => 5.0,  // Bonus for high completion
        rate if rate >= 0.7 => 0.0,  // No adjustment
        rate if rate >= 0.5 => -5.0, // Penalty for low completion
        _ => -10.0, // Higher penalty for very low completion
    };
    
    ((base_prediction + adjustment).max(0.0).min(100.0)) as u32
}

fn calculate_subject_risk_level(
    current_grade: Option<u32>,
    attendance_rate: Option<u32>,
    assignments_completed: u32,
    assignments_total: u32,
) -> RiskLevel {
    let mut risk_score = 0u32;
    
    // Grade-based risk
    if let Some(grade) = current_grade {
        match grade {
            0..=50 => risk_score += 40,
            51..=60 => risk_score += 25,
            61..=70 => risk_score += 10,
            _ => {} // No risk for good grades
        }
    }
    
    // Attendance-based risk
    if let Some(attendance) = attendance_rate {
        match attendance {
            0..=50 => risk_score += 30,
            51..=70 => risk_score += 15,
            71..=80 => risk_score += 5,
            _ => {} // Good attendance
        }
    }
    
    // Assignment completion risk
    if assignments_total > 0 {
        let completion_rate = (assignments_completed * 100) / assignments_total;
        match completion_rate {
            0..=30 => risk_score += 20,
            31..=50 => risk_score += 10,
            51..=70 => risk_score += 5,
            _ => {} // Good completion rate
        }
    }
    
    match risk_score {
        0..=10 => RiskLevel::Low,
        11..=25 => RiskLevel::Medium,
        26..=50 => RiskLevel::High,
        _ => RiskLevel::Critical,
    }
}

fn determine_milestone_type_and_significance(
    milestone_id: &str,
    _progress: &StudentProgress,
) -> (MilestoneType, MilestoneSignificance) {
    match milestone_id {
        id if id.contains("first_semester") => (MilestoneType::FirstSemester, MilestoneSignificance::Standard),
        id if id.contains("quarter") => (MilestoneType::QuarterProgress, MilestoneSignificance::Major),
        id if id.contains("halfway") => (MilestoneType::HalfwayPoint, MilestoneSignificance::Major),
        id if id.contains("three_quarter") => (MilestoneType::ThreeQuarterProgress, MilestoneSignificance::Major),
        id if id.contains("graduation") => (MilestoneType::Graduation, MilestoneSignificance::Critical),
        id if id.contains("perfect") => (MilestoneType::PerfectSemester, MilestoneSignificance::Major),
        id if id.contains("gpa") => (MilestoneType::HighGPA, MilestoneSignificance::Standard),
        _ => (MilestoneType::Custom(milestone_id.to_string()), MilestoneSignificance::Minor),
    }
}

fn generate_upcoming_milestones(progress: &mut StudentProgress) {
    progress.upcoming_milestones.clear();
    
    // Calculate current completion percentage
    let completion_percentage = if progress.total_credits_required > 0 {
        (progress.total_credits_completed * 100) / progress.total_credits_required
    } else {
        0
    };
    
    // Quarter progress milestone
    if completion_percentage < 25 {
        progress.upcoming_milestones.push(UpcomingMilestone {
            milestone_id: "quarter_progress".to_string(),
            milestone_type: MilestoneType::QuarterProgress,
            title: "Quarter Progress".to_string(),
            description: "Complete 25% of total credits".to_string(),
            estimated_date: format!("Semester {}", progress.current_semester + 2),
            requirements: vec!["Complete required credits".to_string()],
            progress_percentage: (completion_percentage * 4).min(100),
            probability: 85,
        });
    }
    
    // Halfway milestone
    if completion_percentage < 50 {
        progress.upcoming_milestones.push(UpcomingMilestone {
            milestone_id: "halfway_point".to_string(),
            milestone_type: MilestoneType::HalfwayPoint,
            title: "Halfway Point".to_string(),
            description: "Complete 50% of total credits".to_string(),
            estimated_date: format!("Semester {}", progress.current_semester + 4),
            requirements: vec!["Complete half of required credits".to_string()],
            progress_percentage: (completion_percentage * 2).min(100),
            probability: 80,
        });
    }
    
    // Graduation milestone
    if completion_percentage >= 75 {
        progress.upcoming_milestones.push(UpcomingMilestone {
            milestone_id: "graduation".to_string(),
            milestone_type: MilestoneType::Graduation,
            title: "Graduation".to_string(),
            description: "Complete all degree requirements".to_string(),
            estimated_date: progress.graduation_forecast.estimated_graduation_date.clone(),
            requirements: vec![
                "Complete all required credits".to_string(),
                "Maintain minimum GPA".to_string(),
            ],
            progress_percentage: completion_percentage,
            probability: progress.graduation_forecast.confidence_level,
        });
    }
}

fn generate_early_warning_indicators(risk_type: &RiskFactorType) -> Vec<String> {
    match risk_type {
        RiskFactorType::AcademicPerformance => vec![
            "Declining grades across multiple subjects".to_string(),
            "Failure to meet assignment deadlines".to_string(),
            "Poor exam performance".to_string(),
        ],
        RiskFactorType::Attendance => vec![
            "Frequent absences".to_string(),
            "Late arrivals to class".to_string(),
            "Missing important sessions".to_string(),
        ],
        RiskFactorType::Engagement => vec![
            "Low participation in class discussions".to_string(),
            "Lack of interaction with instructors".to_string(),
            "Missing extracurricular activities".to_string(),
        ],
        RiskFactorType::TimeManagement => vec![
            "Late submission of assignments".to_string(),
            "Overwhelming course load".to_string(),
            "Poor study schedule management".to_string(),
        ],
        RiskFactorType::Prerequisites => vec![
            "Struggling with foundational concepts".to_string(),
            "Gaps in required knowledge".to_string(),
            "Difficulty with advanced topics".to_string(),
        ],
        _ => vec![
            "General academic concerns".to_string(),
            "Need for additional support".to_string(),
        ],
    }
}

pub fn generate_institution_analytics_internal(
    mut deps: DepsMut,
    env: Env,
    institution_id: String,
    period: AnalyticsPeriod,
    force_refresh: bool,
) -> Result<InstitutionAnalytics, ContractError> {
    // Check if analytics already exist and not forcing refresh
    if !force_refresh {
        if let Ok(existing) = INSTITUTION_ANALYTICS.load(deps.storage, &institution_id) {
            return Ok(existing);
        }
    }
    
    // Get all students from institution
    let student_ids = STUDENTS_BY_INSTITUTION.load(deps.storage, &institution_id).unwrap_or_default();
    
    let mut total_students = 0u32;
    let mut active_students = 0u32;
    let mut graduated_students = 0u32;
    let mut withdrawn_students = 0u32;
    
    let mut gpa_sum = 0.0f64;
    let mut gpa_count = 0u32;
    
    let mut grade_distribution = InstitutionGradeDistribution {
        excellent_percentage: 0,
        good_percentage: 0,
        satisfactory_percentage: 0,
        poor_percentage: 0,
        fail_percentage: 0,
    };
    
    let course_completion_rates = Vec::new();
    let mut risk_distribution = RiskDistribution {
        low_risk_count: 0,
        medium_risk_count: 0,
        high_risk_count: 0,
        critical_risk_count: 0,
    };
    
    // Process each student
    for student_id in &student_ids {
        if let Ok(progress) = STUDENT_PROGRESS.load(deps.storage, student_id) {
            total_students += 1;
            
            match progress.academic_status {
                AcademicStatus::Active => active_students += 1,
                AcademicStatus::Graduated => graduated_students += 1,
                AcademicStatus::Withdrawn | AcademicStatus::Expelled => withdrawn_students += 1,
                _ => {}
            }
            
            // Add to GPA calculation
            if let Ok(gpa) = progress.gpa.parse::<f64>() {
                gpa_sum += gpa;
                gpa_count += 1;
            }
            
            // Update grade distribution
            let total_grades = progress.completed_subjects.len() + progress.failed_subjects.len();
            if total_grades > 0 {
                let excellent = progress.grade_distribution.excellent_count;
                let good = progress.grade_distribution.good_count;
                let satisfactory = progress.grade_distribution.satisfactory_count;
                let poor = progress.grade_distribution.poor_count;
                let fail = progress.grade_distribution.fail_count;
                
                grade_distribution.excellent_percentage += (excellent * 100) / total_grades as u32;
                grade_distribution.good_percentage += (good * 100) / total_grades as u32;
                grade_distribution.satisfactory_percentage += (satisfactory * 100) / total_grades as u32;
                grade_distribution.poor_percentage += (poor * 100) / total_grades as u32;
                grade_distribution.fail_percentage += (fail * 100) / total_grades as u32;
            }
            
            // Update risk distribution
            let risk_level = determine_risk_level(&progress);
            match risk_level {
                RiskSeverity::Low => risk_distribution.low_risk_count += 1,
                RiskSeverity::Moderate => risk_distribution.medium_risk_count += 1,
                RiskSeverity::High => risk_distribution.high_risk_count += 1,
                RiskSeverity::Critical => risk_distribution.critical_risk_count += 1,
            }
        }
    }
    
    // Calculate averages
    if total_students > 0 {
        grade_distribution.excellent_percentage /= total_students;
        grade_distribution.good_percentage /= total_students;
        grade_distribution.satisfactory_percentage /= total_students;
        grade_distribution.poor_percentage /= total_students;
        grade_distribution.fail_percentage /= total_students;
    }
    
    let average_gpa = if gpa_count > 0 {
        format!("{:.2}", gpa_sum / gpa_count as f64)
    } else {
        "0.00".to_string()
    };
    
    let dropout_rate = if total_students > 0 {
        (withdrawn_students * 100) / total_students
    } else {
        0
    };
    
    let completion_rates = CompletionRates {
        overall_completion_rate: if total_students > 0 {
            (graduated_students * 100) / total_students
        } else {
            0
        },
        completion_rate_by_course: course_completion_rates,
        completion_rate_by_semester: vec![], // Would be calculated with semester data
    };
    
    let analytics = InstitutionAnalytics {
        institution_id: institution_id.clone(),
        analytics_period: period,
        total_students,
        active_students,
        graduated_students,
        dropout_rate,
        average_gpa,
        grade_distribution,
        completion_rates,
        average_time_to_graduation: 8, // Default estimate
        on_track_percentage: if active_students > 0 {
            ((active_students - risk_distribution.high_risk_count - risk_distribution.critical_risk_count) * 100) / active_students
        } else {
            0
        },
        risk_distribution,
        benchmark_comparisons: vec![], // Would be populated with actual benchmarks
        enrollment_trends: vec![],      // Would be calculated from historical data
        performance_trends: vec![],     // Would be calculated from historical data
        subject_analytics: vec![],      // Would be populated with subject data
        generated_timestamp: env.block.time.seconds(),
        next_update_due: env.block.time.seconds() + 86400, // 24 hours
    };
    
    // Save analytics
    INSTITUTION_ANALYTICS.save(deps.storage, &institution_id, &analytics)?;
    
    // Store historical snapshot
    ANALYTICS_HISTORY.save(
        deps.storage,
        (&institution_id, env.block.time.seconds()),
        &analytics,
    )?;
    
    Ok(analytics)
}

fn get_all_institutions(storage: &dyn Storage) -> StdResult<Vec<String>> {
    // Get all unique institution IDs from student records
    let mut institutions = std::collections::HashSet::new();
    
    let students: Vec<(String, StudentProgress)> = STUDENT_PROGRESS
        .range(storage, None, None, Order::Ascending)
        .take(500) // Limit for gas management
        .collect::<StdResult<Vec<_>>>()?;
    
    for (_, progress) in students {
        institutions.insert(progress.institution_id);
    }
    
    Ok(institutions.into_iter().collect())
}

// Dashboard generation helper functions

fn generate_grade_history(progress: &StudentProgress) -> Vec<GradeHistoryPoint> {
    let mut grade_history = Vec::new();
    let mut semester_data: std::collections::HashMap<u32, (u32, u32, u32)> = std::collections::HashMap::new();
    
    // Group completed subjects by semester
    for subject in &progress.completed_subjects {
        let entry = semester_data.entry(subject.semester_taken).or_insert((0, 0, 0));
        entry.0 += subject.final_grade * subject.credits; // Grade points
        entry.1 += subject.credits; // Total credits
        entry.2 += 1; // Subject count
    }
    
    // Convert to grade history points
    for semester in 1..=progress.current_semester {
        if let Some((grade_points, credits, subject_count)) = semester_data.get(&semester) {
            let semester_gpa = if *credits > 0 {
                format!("{:.2}", (*grade_points as f64) / (*credits as f64) / 10.0)
            } else {
                "0.00".to_string()
            };
            
            grade_history.push(GradeHistoryPoint {
                semester,
                gpa: semester_gpa,
                credits_earned: *credits,
                subjects_completed: *subject_count,
                trend_indicator: TrendDirection::Stable, // Would calculate actual trend
            });
        }
    }
    
    grade_history
}

fn generate_progress_timeline(progress: &StudentProgress) -> Vec<ProgressTimelineEvent> {
    let mut timeline = Vec::new();
    
    // Add enrollment event
    timeline.push(ProgressTimelineEvent {
        date: progress.enrollment_date.clone(),
        event_type: TimelineEventType::Enrollment,
        title: "Program Enrollment".to_string(),
        description: format!("Enrolled in {}", progress.course_id),
        impact_level: ImpactLevel::High,
    });
    
    // Add major milestones
    for milestone in &progress.milestones_achieved {
        timeline.push(ProgressTimelineEvent {
            date: milestone.achieved_date.clone(),
            event_type: TimelineEventType::MilestoneAchievement,
            title: milestone.title.clone(),
            description: milestone.description.clone(),
            impact_level: match milestone.significance {
                MilestoneSignificance::Critical => ImpactLevel::Critical,
                MilestoneSignificance::Major => ImpactLevel::High,
                MilestoneSignificance::Standard => ImpactLevel::Medium,
                MilestoneSignificance::Minor => ImpactLevel::Low,
            },
        });
    }
    
    // Add significant grade improvements
    if matches!(progress.performance_metrics.grade_trend, GradeTrend::StronglyIncreasing) {
        timeline.push(ProgressTimelineEvent {
            date: "Recent".to_string(),
            event_type: TimelineEventType::GradeImprovement,
            title: "Grade Improvement".to_string(),
            description: "Significant improvement in academic performance".to_string(),
            impact_level: ImpactLevel::Medium,
        });
    }
    
    timeline
}

fn generate_urgent_alerts(progress: &StudentProgress, _env: &Env) -> Vec<Alert> {
    let mut alerts = Vec::new();
    
    // GPA warning
    if let Ok(gpa) = progress.gpa.parse::<f64>() {
        if gpa < 6.0 {
            alerts.push(Alert {
                alert_id: format!("{}_gpa_critical", progress.student_id),
                alert_type: AlertType::GradeWarning,
                severity: AlertSeverity::Critical,
                title: "Critical GPA Warning".to_string(),
                message: "Your GPA is below the minimum threshold. Immediate action required.".to_string(),
                action_required: true,
                deadline: Some("End of semester".to_string()),
                related_subject: None,
            });
        } else if gpa < 7.0 {
            alerts.push(Alert {
                alert_id: format!("{}_gpa_warning", progress.student_id),
                alert_type: AlertType::GradeWarning,
                severity: AlertSeverity::Warning,
                title: "GPA Warning".to_string(),
                message: "Your GPA is approaching critical levels. Consider academic support.".to_string(),
                action_required: true,
                deadline: None,
                related_subject: None,
            });
        }
    }
    
    // Academic probation alert
    if matches!(progress.academic_status, AcademicStatus::Probation) {
        alerts.push(Alert {
            alert_id: format!("{}_probation", progress.student_id),
            alert_type: AlertType::AcademicProbation,
            severity: AlertSeverity::Critical,
            title: "Academic Probation".to_string(),
            message: "You are currently on academic probation. Meet with your advisor immediately.".to_string(),
            action_required: true,
            deadline: Some("Within 1 week".to_string()),
            related_subject: None,
        });
    }
    
    // At-risk subject alerts
    for subject in &progress.current_subjects {
        if matches!(subject.risk_level, RiskLevel::High | RiskLevel::Critical) {
            alerts.push(Alert {
                alert_id: format!("{}_{}_risk", progress.student_id, subject.subject_id),
                alert_type: AlertType::GradeWarning,
                severity: if matches!(subject.risk_level, RiskLevel::Critical) {
                    AlertSeverity::Critical
                } else {
                    AlertSeverity::Warning
                },
                title: format!("At Risk: {}", subject.title),
                message: "Your performance in this subject indicates high risk of failure.".to_string(),
                action_required: true,
                deadline: None,
                related_subject: Some(subject.subject_id.clone()),
            });
        }
    }
    
    alerts
}

fn generate_upcoming_deadlines(progress: &StudentProgress) -> Vec<Deadline> {
    let mut deadlines = Vec::new();
    
    // Generate assignment deadlines from current subjects
    for subject in &progress.current_subjects {
        if subject.assignments_completed < subject.assignments_total {
            let remaining = subject.assignments_total - subject.assignments_completed;
            deadlines.push(Deadline {
                deadline_id: format!("{}_{}_assignment", progress.student_id, subject.subject_id),
                title: format!("{} - Assignment Due", subject.title),
                description: format!("{} assignments remaining", remaining),
                due_date: "Next week".to_string(), // Would calculate actual dates
                days_remaining: 7,
                subject_id: Some(subject.subject_id.clone()),
                deadline_type: DeadlineType::Assignment,
                completion_status: if remaining > 5 {
                    CompletionStatus::NotStarted
                } else {
                    CompletionStatus::InProgress
                },
            });
        }
    }
    
    // Add registration deadlines
    deadlines.push(Deadline {
        deadline_id: format!("{}_registration", progress.student_id),
        title: "Course Registration".to_string(),
        description: "Registration opens for next semester".to_string(),
        due_date: "Next month".to_string(),
        days_remaining: 30,
        subject_id: None,
        deadline_type: DeadlineType::Registration,
        completion_status: CompletionStatus::NotStarted,
    });
    
    deadlines
}

fn generate_dashboard_recommendations(progress: &StudentProgress) -> Vec<Recommendation> {
    let mut recommendations = Vec::new();
    
    // GPA improvement recommendations
    if let Ok(gpa) = progress.gpa.parse::<f64>() {
        if gpa < 7.5 {
            recommendations.push(Recommendation {
                recommendation_id: format!("{}_gpa_improvement", progress.student_id),
                category: RecommendationCategory::AcademicImprovement,
                priority: if gpa < 6.0 { Priority::Critical } else { Priority::High },
                title: "Improve Academic Performance".to_string(),
                description: "Your GPA needs improvement to maintain good academic standing".to_string(),
                action_steps: vec![
                    "Schedule meeting with academic advisor".to_string(),
                    "Seek tutoring for challenging subjects".to_string(),
                    "Develop better study habits".to_string(),
                ],
                expected_benefit: "Improved grades and academic standing".to_string(),
                effort_level: EffortLevel::High,
            });
        }
    }
    
    // Study strategy recommendations
    if progress.performance_metrics.study_hours_per_credit < 10 {
        recommendations.push(Recommendation {
            recommendation_id: format!("{}_study_time", progress.student_id),
            category: RecommendationCategory::StudyStrategy,
            priority: Priority::Medium,
            title: "Increase Study Time".to_string(),
            description: "You may benefit from dedicating more time to studying".to_string(),
            action_steps: vec![
                "Create a structured study schedule".to_string(),
                "Find a dedicated study space".to_string(),
                "Use active learning techniques".to_string(),
            ],
            expected_benefit: "Better retention and improved grades".to_string(),
            effort_level: EffortLevel::Medium,
        });
    }
    
    // Time management recommendations
    if progress.current_subjects.len() as u32 > 6 {
        recommendations.push(Recommendation {
            recommendation_id: format!("{}_workload", progress.student_id),
            category: RecommendationCategory::TimeManagement,
            priority: Priority::Medium,
            title: "Manage Course Load".to_string(),
            description: "Consider reducing course load to improve performance".to_string(),
            action_steps: vec![
                "Evaluate current course difficulty".to_string(),
                "Consider dropping one course if struggling".to_string(),
                "Focus on core requirements first".to_string(),
            ],
            expected_benefit: "Better focus and improved grades".to_string(),
            effort_level: EffortLevel::Low,
        });
    }
    
    recommendations
}

fn generate_recent_achievements(progress: &StudentProgress) -> Vec<Achievement> {
    let mut achievements = Vec::new();
    
    // Recent milestones as achievements
    for milestone in progress.milestones_achieved.iter().rev().take(3) {
        achievements.push(Achievement {
            achievement_id: milestone.milestone_id.clone(),
            title: milestone.title.clone(),
            description: milestone.description.clone(),
            earned_date: milestone.achieved_date.clone(),
            achievement_type: match milestone.milestone_type {
                MilestoneType::Graduation => AchievementType::Academic,
                MilestoneType::HighGPA => AchievementType::Academic,
                MilestoneType::PerfectSemester => AchievementType::Academic,
                _ => AchievementType::Participation,
            },
            rarity: match milestone.significance {
                MilestoneSignificance::Critical => AchievementRarity::Legendary,
                MilestoneSignificance::Major => AchievementRarity::Epic,
                MilestoneSignificance::Standard => AchievementRarity::Rare,
                MilestoneSignificance::Minor => AchievementRarity::Common,
            },
            points_awarded: match milestone.significance {
                MilestoneSignificance::Critical => 500,
                MilestoneSignificance::Major => 200,
                MilestoneSignificance::Standard => 100,
                MilestoneSignificance::Minor => 50,
            },
        });
    }
    
    // Performance-based achievements
    if let Ok(gpa) = progress.gpa.parse::<f64>() {
        if gpa >= 9.0 {
            achievements.push(Achievement {
                achievement_id: format!("{}_excellent_gpa", progress.student_id),
                title: "Academic Excellence".to_string(),
                description: "Maintained GPA above 9.0".to_string(),
                earned_date: "Current".to_string(),
                achievement_type: AchievementType::Academic,
                rarity: AchievementRarity::Epic,
                points_awarded: 300,
            });
        }
    }
    
    achievements
}

fn generate_goal_progress(progress: &StudentProgress) -> Vec<GoalProgress> {
    let mut goals = Vec::new();
    
    // Graduation goal
    let completion_percentage = if progress.total_credits_required > 0 {
        (progress.total_credits_completed * 100) / progress.total_credits_required
    } else {
        0
    };
    
    goals.push(GoalProgress {
        goal_id: format!("{}_graduation", progress.student_id),
        goal_title: "Graduate".to_string(),
        goal_description: "Complete all degree requirements".to_string(),
        current_progress: progress.total_credits_completed,
        target_value: progress.total_credits_required,
        progress_percentage: completion_percentage,
        estimated_completion_date: Some(progress.graduation_forecast.estimated_graduation_date.clone()),
        on_track: progress.graduation_forecast.on_track,
    });
    
    // GPA goal
    if let Ok(current_gpa) = progress.gpa.parse::<f64>() {
        let target_gpa = 8.0; // Target GPA
        let gpa_progress = ((current_gpa / target_gpa) * 100.0).min(100.0) as u32;
        
        goals.push(GoalProgress {
            goal_id: format!("{}_gpa", progress.student_id),
            goal_title: "Achieve High GPA".to_string(),
            goal_description: "Maintain GPA above 8.0".to_string(),
            current_progress: (current_gpa * 100.0) as u32,
            target_value: (target_gpa * 100.0) as u32,
            progress_percentage: gpa_progress,
            estimated_completion_date: None,
            on_track: current_gpa >= target_gpa,
        });
    }
    
    goals
}

fn generate_performance_insights(progress: &StudentProgress) -> Vec<PerformanceInsight> {
    let mut insights = Vec::new();
    
    // Grade trend insight
    match progress.performance_metrics.grade_trend {
        GradeTrend::StronglyIncreasing => {
            insights.push(PerformanceInsight {
                insight_id: format!("{}_grade_trend", progress.student_id),
                insight_type: InsightType::TrendAnalysis,
                title: "Excellent Grade Improvement".to_string(),
                description: "Your grades are showing strong upward trend".to_string(),
                supporting_data: vec!["Consistent improvement across semesters".to_string()],
                actionable_advice: vec!["Continue current study strategies".to_string()],
                confidence_level: 90,
            });
        },
        GradeTrend::StronglyDecreasing => {
            insights.push(PerformanceInsight {
                insight_id: format!("{}_grade_decline", progress.student_id),
                insight_type: InsightType::PredictiveWarning,
                title: "Grade Decline Warning".to_string(),
                description: "Your grades are declining and need immediate attention".to_string(),
                supporting_data: vec!["Declining performance across multiple subjects".to_string()],
                actionable_advice: vec![
                    "Seek academic counseling".to_string(),
                    "Review study methods".to_string(),
                ],
                confidence_level: 85,
            });
        },
        _ => {}
    }
    
    // Strength identification
    if progress.performance_metrics.success_rate >= 90 {
        insights.push(PerformanceInsight {
            insight_id: format!("{}_high_success", progress.student_id),
            insight_type: InsightType::StrengthIdentification,
            title: "High Success Rate".to_string(),
            description: "You demonstrate excellent subject completion ability".to_string(),
            supporting_data: vec![
                format!("{}% success rate", progress.performance_metrics.success_rate)
            ],
            actionable_advice: vec![
                "Consider taking more challenging subjects".to_string(),
                "Explore advanced opportunities".to_string(),
            ],
            confidence_level: 95,
        });
    }
    
    insights
}