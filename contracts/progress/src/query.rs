use cosmwasm_std::{Deps, StdResult, Order, Storage, Env};
use crate::msg::*;
use crate::state::*;
use std::collections::HashMap;

pub fn get_state(deps: Deps) -> StdResult<StateResponse> {
    let state = STATE.load(deps.storage)?;
    Ok(StateResponse { state })
}

pub fn get_config(deps: Deps) -> StdResult<ProgressConfigResponse> {
    let config = PROGRESS_CONFIG.load(deps.storage)?;
    Ok(ProgressConfigResponse { config })
}

pub fn get_student_progress(deps: Deps, student_id: String) -> StdResult<StudentProgressResponse> {
    let progress = STUDENT_PROGRESS.load(deps.storage, &student_id)?;
    Ok(StudentProgressResponse { progress })
}

pub fn get_student_dashboard(
    deps: Deps, 
    student_id: String, 
    include_peer_comparison: Option<bool>
) -> StdResult<StudentDashboardResponse> {
    let mut dashboard = STUDENT_DASHBOARDS.load(deps.storage, &student_id)?;
    
    // If peer comparison not included in stored dashboard but requested, generate it
    if include_peer_comparison == Some(true) && dashboard.peer_comparisons.is_none() {
        let student_progress = STUDENT_PROGRESS.load(deps.storage, &student_id)?;
        dashboard.peer_comparisons = Some(generate_peer_comparison(deps.storage, &student_progress)?);
    }
    
    Ok(StudentDashboardResponse { dashboard })
}

pub fn get_institution_analytics(
    deps: Deps, 
    institution_id: String, 
    _period: Option<AnalyticsPeriod>
) -> StdResult<InstitutionAnalyticsResponse> {
    let analytics = INSTITUTION_ANALYTICS.load(deps.storage, &institution_id)?;
    Ok(InstitutionAnalyticsResponse { analytics })
}

pub fn get_students_by_institution(
    deps: Deps,
    institution_id: String,
    status_filter: Option<AcademicStatus>,
    limit: Option<u32>,
    start_after: Option<String>,
) -> StdResult<StudentsListResponse> {
    let all_students = STUDENTS_BY_INSTITUTION.load(deps.storage, &institution_id).unwrap_or_default();
    
    let mut filtered_students = Vec::new();
    let mut total_count = 0u32;
    let limit = limit.unwrap_or(50).min(100); // Cap at 100 for gas management
    let mut skip_count = 0u32;
    
    // Determine skip count if start_after is provided
    if let Some(start_after_id) = &start_after {
        skip_count = all_students.iter()
            .position(|id| id == start_after_id)
            .map(|pos| pos as u32 + 1)
            .unwrap_or(0);
    }
    
    for student_id in all_students.iter().skip(skip_count as usize) {
        if let Ok(progress) = STUDENT_PROGRESS.load(deps.storage, student_id) {
            // Apply status filter if provided
            if let Some(ref filter_status) = status_filter {
                if progress.academic_status != *filter_status {
                    continue;
                }
            }
            
            total_count += 1;
            
            if filtered_students.len() < limit as usize {
                let completion_percentage = if progress.total_credits_required > 0 {
                    (progress.total_credits_completed * 100) / progress.total_credits_required
                } else {
                    0
                };
                
                filtered_students.push(StudentSummary {
                    student_id: student_id.clone(),
                    institution_id: progress.institution_id.clone(),
                    course_id: progress.course_id.clone(),
                    current_semester: progress.current_semester,
                    gpa: progress.gpa.clone(),
                    academic_status: progress.academic_status.clone(),
                    completion_percentage,
                });
            }
        }
    }
    
    let has_more = (skip_count + filtered_students.len() as u32) < total_count;
    
    Ok(StudentsListResponse {
        students: filtered_students,
        total_count,
        has_more,
    })
}

pub fn get_students_by_course(
    deps: Deps,
    course_id: String,
    status_filter: Option<AcademicStatus>,
    limit: Option<u32>,
    start_after: Option<String>,
) -> StdResult<StudentsListResponse> {
    let all_students = STUDENTS_BY_COURSE.load(deps.storage, &course_id).unwrap_or_default();
    
    let mut filtered_students = Vec::new();
    let mut total_count = 0u32;
    let limit = limit.unwrap_or(50).min(100);
    let mut skip_count = 0u32;
    
    if let Some(start_after_id) = &start_after {
        skip_count = all_students.iter()
            .position(|id| id == start_after_id)
            .map(|pos| pos as u32 + 1)
            .unwrap_or(0);
    }
    
    for student_id in all_students.iter().skip(skip_count as usize) {
        if let Ok(progress) = STUDENT_PROGRESS.load(deps.storage, student_id) {
            if let Some(ref filter_status) = status_filter {
                if progress.academic_status != *filter_status {
                    continue;
                }
            }
            
            total_count += 1;
            
            if filtered_students.len() < limit as usize {
                let completion_percentage = if progress.total_credits_required > 0 {
                    (progress.total_credits_completed * 100) / progress.total_credits_required
                } else {
                    0
                };
                
                filtered_students.push(StudentSummary {
                    student_id: student_id.clone(),
                    institution_id: progress.institution_id.clone(),
                    course_id: progress.course_id.clone(),
                    current_semester: progress.current_semester,
                    gpa: progress.gpa.clone(),
                    academic_status: progress.academic_status.clone(),
                    completion_percentage,
                });
            }
        }
    }
    
    let has_more = (skip_count + filtered_students.len() as u32) < total_count;
    
    Ok(StudentsListResponse {
        students: filtered_students,
        total_count,
        has_more,
    })
}

pub fn get_at_risk_students(
    deps: Deps,
    institution_id: Option<String>,
    risk_level_filter: Option<RiskSeverity>,
    limit: Option<u32>,
) -> StdResult<AtRiskStudentsResponse> {
    let limit = limit.unwrap_or(50).min(100);
    let mut at_risk_students = Vec::new();
    let mut risk_counts = HashMap::new();
    let mut total_interventions = 0u32;
    let mut successful_interventions = 0u32;
    
    // Get students to check
    let students_to_check = if let Some(ref inst_id) = institution_id {
        STUDENTS_BY_INSTITUTION.load(deps.storage, &inst_id).unwrap_or_default()
    } else {
        // Get all students (limited for gas management)
        STUDENT_PROGRESS
            .range(deps.storage, None, None, Order::Ascending)
            .take(500)
            .map(|item| item.map(|(id, _)| id))
            .collect::<StdResult<Vec<String>>>()?
    };
    
    for student_id in students_to_check.iter().take(limit as usize) {
        if let Ok(progress) = STUDENT_PROGRESS.load(deps.storage, student_id) {
            // Determine current risk level
            let risk_level = determine_overall_risk_level(&progress);
            
            // Count all risk levels for summary
            *risk_counts.entry(risk_level.clone()).or_insert(0u32) += 1;
            
            // Filter by risk level if specified
            if let Some(ref filter_level) = risk_level_filter {
                if risk_level != *filter_level {
                    continue;
                }
            }
            
            // Only include students with moderate to critical risk
            if matches!(risk_level, RiskSeverity::Moderate | RiskSeverity::High | RiskSeverity::Critical) {
                let risk_factors: Vec<RiskFactorType> = progress.risk_factors
                    .iter()
                    .map(|rf| rf.risk_type.clone())
                    .collect();
                
                let intervention_priority = match risk_level {
                    RiskSeverity::Critical => 100,
                    RiskSeverity::High => 75,
                    RiskSeverity::Moderate => 50,
                    RiskSeverity::Low => 25,
                };
                
                at_risk_students.push(AtRiskStudent {
                    student_id: student_id.clone(),
                    risk_level: risk_level.clone(),
                    risk_factors,
                    intervention_priority,
                    last_updated: progress.last_updated,
                });
                
                // Count interventions (simplified)
                total_interventions += progress.risk_factors.len() as u32;
                // Assume 70% success rate for interventions
                successful_interventions += (progress.risk_factors.len() as u32 * 70) / 100;
            }
        }
    }
    
    // Sort by intervention priority (highest first)
    at_risk_students.sort_by(|a, b| b.intervention_priority.cmp(&a.intervention_priority));
    
    let summary = RiskSummary {
        total_at_risk: at_risk_students.len() as u32,
        critical_risk_count: *risk_counts.get(&RiskSeverity::Critical).unwrap_or(&0),
        high_risk_count: *risk_counts.get(&RiskSeverity::High).unwrap_or(&0),
        medium_risk_count: *risk_counts.get(&RiskSeverity::Moderate).unwrap_or(&0),
        intervention_success_rate: if total_interventions > 0 {
            (successful_interventions * 100) / total_interventions
        } else {
            0
        },
    };
    
    Ok(AtRiskStudentsResponse {
        at_risk_students,
        summary,
    })
}

pub fn get_top_performers(
    deps: Deps,
    institution_id: Option<String>,
    course_id: Option<String>,
    metric: PerformanceMetric,
    limit: Option<u32>,
) -> StdResult<TopPerformersResponse> {
    let limit = limit.unwrap_or(20).min(50);
    let mut performers = Vec::new();
    let mut metric_values = Vec::new();
    
    // Get students to evaluate
    let students_to_check = if let Some(ref inst_id) = institution_id {
        STUDENTS_BY_INSTITUTION.load(deps.storage, &inst_id).unwrap_or_default()
    } else if let Some(ref course_id_filter) = course_id {
        STUDENTS_BY_COURSE.load(deps.storage, &course_id_filter).unwrap_or_default()
    } else {
        STUDENT_PROGRESS
            .range(deps.storage, None, None, Order::Ascending)
            .take(200)
            .map(|item| item.map(|(id, _)| id))
            .collect::<StdResult<Vec<String>>>()?
    };
    
    for student_id in students_to_check {
        if let Ok(progress) = STUDENT_PROGRESS.load(deps.storage, &student_id) {
            // Filter by course if specified and not already filtered
            if let Some(ref course_filter) = course_id {
                if institution_id.is_some() && progress.course_id != *course_filter {
                    continue;
                }
            }
            
            let (metric_value, metric_string) = match metric {
                PerformanceMetric::GPA => {
                    if let Ok(gpa) = progress.gpa.parse::<f64>() {
                        (gpa as u32 * 100, progress.gpa.clone()) // Convert to comparable integer
                    } else {
                        continue;
                    }
                },
                PerformanceMetric::ProgressSpeed => {
                    let speed = if progress.current_semester > 0 {
                        (progress.total_credits_completed * 100) / progress.current_semester
                    } else {
                        0
                    };
                    (speed, speed.to_string())
                },
                PerformanceMetric::ConsistencyScore => {
                    (progress.performance_metrics.consistency_score, 
                     progress.performance_metrics.consistency_score.to_string())
                },
                PerformanceMetric::ImprovementRate => {
                    // Calculate improvement based on grade trend
                    let improvement = match progress.performance_metrics.grade_trend {
                        GradeTrend::StronglyIncreasing => 100,
                        GradeTrend::ModeratelyIncreasing => 75,
                        GradeTrend::Stable => 50,
                        GradeTrend::ModeratelyDecreasing => 25,
                        GradeTrend::StronglyDecreasing => 0,
                        GradeTrend::Volatile => 30,
                    };
                    (improvement, improvement.to_string())
                },
                PerformanceMetric::CompletionRate => {
                    (progress.performance_metrics.success_rate, 
                     progress.performance_metrics.success_rate.to_string())
                },
            };
            
            metric_values.push(metric_value);
            
            // Generate achievement highlights
            let mut achievement_highlights = Vec::new();
            
            if let Ok(gpa) = progress.gpa.parse::<f64>() {
                if gpa >= 9.0 {
                    achievement_highlights.push("Excellent GPA (9.0+)".to_string());
                }
            }
            
            if progress.performance_metrics.success_rate >= 95 {
                achievement_highlights.push("High success rate (95%+)".to_string());
            }
            
            if progress.milestones_achieved.len() >= 3 {
                achievement_highlights.push(format!("{} milestones achieved", progress.milestones_achieved.len()));
            }
            
            if progress.failed_subjects.is_empty() {
                achievement_highlights.push("No failed subjects".to_string());
            }
            
            performers.push((
                TopPerformer {
                    student_id: student_id.clone(),
                    metric_value: metric_string,
                    percentile_rank: 0, // Will be calculated after sorting
                    achievement_highlights,
                },
                metric_value,
            ));
        }
    }
    
    // Sort by metric value (highest first)
    performers.sort_by(|a, b| b.1.cmp(&a.1));
    
    // Calculate percentile ranks and take top performers
    let total_count = performers.len();
    let mut top_performers = Vec::new();
    
    for (i, (mut performer, _metric_value)) in performers.into_iter().enumerate() {
        if i >= limit as usize {
            break;
        }
        
        performer.percentile_rank = if total_count > 0 {
            ((total_count - i - 1) * 100) / total_count
        } else {
            100
        } as u32;
        
        top_performers.push(performer);
    }
    
    // Calculate benchmark value (median)
    let benchmark_value = if !metric_values.is_empty() {
        metric_values.sort();
        let median_index = metric_values.len() / 2;
        format!("{}", metric_values[median_index])
    } else {
        "0".to_string()
    };
    
    Ok(TopPerformersResponse {
        top_performers,
        metric_used: metric,
        benchmark_value,
    })
}

pub fn get_graduation_forecast(
    deps: Deps,
    student_id: String,
    _scenario: ForecastScenario,
) -> StdResult<GraduationForecastResponse> {
    let student_progress = STUDENT_PROGRESS.load(deps.storage, &student_id)?;
    
    // Base forecast
    let base_forecast = student_progress.graduation_forecast.clone();
    
    // Generate scenario comparisons
    let mut scenario_comparisons = Vec::new();
    
    // Current scenario
    scenario_comparisons.push(ScenarioComparison {
        scenario: ForecastScenario::Current,
        estimated_date: base_forecast.estimated_graduation_date.clone(),
        probability: base_forecast.confidence_level,
        key_assumptions: vec![
            "Current pace maintained".to_string(),
            "No additional delays".to_string(),
        ],
    });
    
    // Optimistic scenario
    let optimistic_semesters = base_forecast.remaining_semesters.saturating_sub(1);
    scenario_comparisons.push(ScenarioComparison {
        scenario: ForecastScenario::Optimistic,
        estimated_date: format!("Semester {}", student_progress.current_semester + optimistic_semesters),
        probability: (base_forecast.confidence_level + 20).min(100),
        key_assumptions: vec![
            "Accelerated course load".to_string(),
            "Summer sessions utilized".to_string(),
            "No failed subjects".to_string(),
        ],
    });
    
    // Pessimistic scenario
    let pessimistic_semesters = base_forecast.remaining_semesters + 2;
    scenario_comparisons.push(ScenarioComparison {
        scenario: ForecastScenario::Pessimistic,
        estimated_date: format!("Semester {}", student_progress.current_semester + pessimistic_semesters),
        probability: base_forecast.confidence_level.saturating_sub(30),
        key_assumptions: vec![
            "Some subjects need to be retaken".to_string(),
            "Reduced course load due to difficulties".to_string(),
            "Potential academic probation".to_string(),
        ],
    });
    
    // Intervention scenario
    let intervention_semesters = if student_progress.risk_factors.len() > 2 {
        base_forecast.remaining_semesters.saturating_sub(1)
    } else {
        base_forecast.remaining_semesters
    };
    scenario_comparisons.push(ScenarioComparison {
        scenario: ForecastScenario::Intervention,
        estimated_date: format!("Semester {}", student_progress.current_semester + intervention_semesters),
        probability: (base_forecast.confidence_level + 15).min(100),
        key_assumptions: vec![
            "Risk factors addressed".to_string(),
            "Academic support utilized".to_string(),
            "Study habits improved".to_string(),
        ],
    });
    
    // Generate confidence factors
    let mut confidence_factors = Vec::new();
    
    if let Ok(gpa) = student_progress.gpa.parse::<f64>() {
        if gpa >= 8.0 {
            confidence_factors.push("Strong academic performance".to_string());
        } else if gpa < 7.0 {
            confidence_factors.push("Academic performance concerns".to_string());
        }
    }
    
    if student_progress.performance_metrics.success_rate >= 90 {
        confidence_factors.push("High subject completion rate".to_string());
    }
    
    if student_progress.risk_factors.len() > 2 {
        confidence_factors.push("Multiple risk factors present".to_string());
    }
    
    if matches!(student_progress.performance_metrics.grade_trend, GradeTrend::StronglyIncreasing | GradeTrend::ModeratelyIncreasing) {
        confidence_factors.push("Improving grade trend".to_string());
    }
    
    Ok(GraduationForecastResponse {
        forecast: base_forecast,
        scenario_comparisons,
        confidence_factors,
    })
}

pub fn get_comparative_analytics(
    deps: Deps,
    institution_id: String,
    comparison_group: ComparisonGroup,
    metrics: Vec<String>,
) -> StdResult<ComparativeAnalyticsResponse> {
    let institution_analytics = INSTITUTION_ANALYTICS.load(deps.storage, &institution_id)?;
    
    let mut institution_metrics = Vec::new();
    let mut percentile_rankings = Vec::new();
    let mut improvement_opportunities = Vec::new();
    
    // Get comparison benchmarks
    let benchmarks = get_comparison_benchmarks(deps.storage, &comparison_group)?;
    
    for metric_name in metrics {
        let (institution_value, comparison_value, trend) = match metric_name.as_str() {
            "average_gpa" => {
                let inst_gpa = institution_analytics.average_gpa.parse::<f64>().unwrap_or(0.0);
                let comp_gpa = benchmarks.get("average_gpa").unwrap_or(&"7.5".to_string()).parse::<f64>().unwrap_or(7.5);
                let trend = if inst_gpa > comp_gpa { crate::state::TrendDirection::Improving } 
                          else if inst_gpa < comp_gpa { crate::state::TrendDirection::Declining } 
                          else { crate::state::TrendDirection::Stable };
                (
                    institution_analytics.average_gpa.clone(),
                    benchmarks.get("average_gpa").unwrap_or(&"7.5".to_string()).clone(),
                    trend
                )
            },
            "completion_rate" => {
                let inst_rate = institution_analytics.completion_rates.overall_completion_rate;
                let comp_rate = benchmarks.get("completion_rate").unwrap_or(&"80".to_string()).parse::<u32>().unwrap_or(80);
                let trend = if inst_rate > comp_rate { crate::state::TrendDirection::Improving } 
                          else if inst_rate < comp_rate { crate::state::TrendDirection::Declining } 
                          else { crate::state::TrendDirection::Stable };
                (
                    inst_rate.to_string(),
                    comp_rate.to_string(),
                    trend
                )
            },
            "dropout_rate" => {
                let inst_rate = institution_analytics.dropout_rate;
                let comp_rate = benchmarks.get("dropout_rate").unwrap_or(&"15".to_string()).parse::<u32>().unwrap_or(15);
                let trend = if inst_rate < comp_rate { crate::state::TrendDirection::Improving } 
                          else if inst_rate > comp_rate { crate::state::TrendDirection::Declining } 
                          else { crate::state::TrendDirection::Stable };
                (
                    inst_rate.to_string(),
                    comp_rate.to_string(),
                    trend
                )
            },
            "on_track_percentage" => {
                let inst_rate = institution_analytics.on_track_percentage;
                let comp_rate = benchmarks.get("on_track_percentage").unwrap_or(&"75".to_string()).parse::<u32>().unwrap_or(75);
                let trend = if inst_rate > comp_rate { crate::state::TrendDirection::Improving } 
                          else if inst_rate < comp_rate { crate::state::TrendDirection::Declining } 
                          else { crate::state::TrendDirection::Stable };
                (
                    inst_rate.to_string(),
                    comp_rate.to_string(),
                    trend
                )
            },
            _ => continue,
        };
        
        let performance_ratio = if let (Ok(inst_val), Ok(comp_val)) = (institution_value.parse::<f64>(), comparison_value.parse::<f64>()) {
            if comp_val != 0.0 {
                format!("{:.2}", inst_val / comp_val)
            } else {
                "N/A".to_string()
            }
        } else {
            "N/A".to_string()
        };
        
        institution_metrics.push(MetricComparison {
            metric_name: metric_name.clone(),
            institution_value: institution_value.clone(),
            comparison_value: comparison_value.clone(),
            performance_ratio,
            trend,
        });
        
        // Calculate percentile ranking
        let percentile = calculate_metric_percentile(&metric_name, &institution_value, &comparison_group);
        let rank_category = match percentile {
            90..=100 => RankCategory::TopTier,
            75..=89 => RankCategory::HighPerforming,
            60..=74 => RankCategory::AboveAverage,
            40..=59 => RankCategory::Average,
            25..=39 => RankCategory::BelowAverage,
            10..=24 => RankCategory::LowPerforming,
            _ => RankCategory::Bottom,
        };
        
        percentile_rankings.push(PercentileRanking {
            metric_name: metric_name.clone(),
            percentile,
            rank_category,
        });
        
        // Generate improvement opportunities
        if percentile < 60 {
            improvement_opportunities.push(format!("Focus on improving {}", metric_name.replace("_", " ")));
        }
    }
    
    // Add general improvement suggestions
    if institution_analytics.dropout_rate > 20 {
        improvement_opportunities.push("Implement early intervention programs for at-risk students".to_string());
    }
    
    if institution_analytics.average_gpa.parse::<f64>().unwrap_or(0.0) < 7.0 {
        improvement_opportunities.push("Enhance academic support services".to_string());
    }
    
    Ok(ComparativeAnalyticsResponse {
        institution_metrics,
        percentile_rankings,
        improvement_opportunities,
    })
}

pub fn get_progress_trends(
    deps: Deps,
    entity_id: String,
    entity_type: EntityType,
    _period: AnalyticsPeriod,
    metrics: Vec<String>,
) -> StdResult<ProgressTrendsResponse> {
    let mut trend_data = Vec::new();
    let mut forecasted_values = Vec::new();
    
    match entity_type {
        EntityType::Student => {
            // Get historical data for student
            let _current_progress = STUDENT_PROGRESS.load(deps.storage, &entity_id)?;
            
            // Get historical snapshots
            let history: Vec<_> = PROGRESS_HISTORY
                .prefix(&entity_id)
                .range(deps.storage, None, None, Order::Ascending)
                .take(20)
                .collect::<StdResult<Vec<_>>>()?;
            
            for (timestamp, historical_progress) in history {
                for metric_name in &metrics {
                    let value = match metric_name.as_str() {
                        "gpa" => historical_progress.gpa.clone(),
                        "credits_completed" => historical_progress.total_credits_completed.to_string(),
                        "completion_percentage" => {
                            if historical_progress.total_credits_required > 0 {
                                ((historical_progress.total_credits_completed * 100) / historical_progress.total_credits_required).to_string()
                            } else {
                                "0".to_string()
                            }
                        },
                        _ => continue,
                    };
                    
                    trend_data.push(TrendDataPoint {
                        period: timestamp.to_string(),
                        value,
                        period_change: "0".to_string(), // Would calculate actual change
                        cumulative_change: "0".to_string(), // Would calculate from baseline
                    });
                }
            }
        },
        EntityType::Institution => {
            // Get historical analytics for institution
            let analytics_history: Vec<_> = ANALYTICS_HISTORY
                .prefix(&entity_id)
                .range(deps.storage, None, None, Order::Ascending)
                .take(12) // Last 12 periods
                .collect::<StdResult<Vec<_>>>()?;
            
            for (timestamp, historical_analytics) in analytics_history {
                for metric_name in &metrics {
                    let value = match metric_name.as_str() {
                        "average_gpa" => historical_analytics.average_gpa.clone(),
                        "completion_rate" => historical_analytics.completion_rates.overall_completion_rate.to_string(),
                        "dropout_rate" => historical_analytics.dropout_rate.to_string(),
                        "total_students" => historical_analytics.total_students.to_string(),
                        _ => continue,
                    };
                    
                    trend_data.push(TrendDataPoint {
                        period: timestamp.to_string(),
                        value,
                        period_change: "0".to_string(),
                        cumulative_change: "0".to_string(),
                    });
                }
            }
        },
        _ => {
            // Other entity types would be implemented similarly
        }
    }
    
    // Generate trend analysis
    let trend_analysis = TrendAnalysis {
        overall_direction: if trend_data.len() >= 2 {
            // Simple trend calculation
            let first_value = trend_data.first().unwrap().value.parse::<f64>().unwrap_or(0.0);
            let last_value = trend_data.last().unwrap().value.parse::<f64>().unwrap_or(0.0);
            
            if last_value > first_value * 1.05 {
                crate::state::TrendDirection::Improving
            } else if last_value < first_value * 0.95 {
                crate::state::TrendDirection::Declining
            } else {
                crate::state::TrendDirection::Stable
            }
        } else {
            crate::state::TrendDirection::Stable
        },
        volatility: VolatilityLevel::Low, // Would calculate actual volatility
        seasonal_patterns: vec!["Higher performance in fall semesters".to_string()],
        significant_events: vec![
            SignificantEvent {
                date: "2024-09-01".to_string(),
                event_type: "Academic Year Start".to_string(),
                impact_description: "Enrollment spike".to_string(),
                impact_magnitude: 15,
            }
        ],
    };
    
    // Generate forecasts (simplified)
    for i in 1..=3 {
        forecasted_values.push(ForecastPoint {
            period: format!("Future+{}", i),
            predicted_value: "Forecast".to_string(),
            confidence_interval_low: "Lower".to_string(),
            confidence_interval_high: "Upper".to_string(),
            confidence_level: 70,
        });
    }
    
    Ok(ProgressTrendsResponse {
        trend_data,
        trend_analysis,
        forecasted_values,
    })
}

// Continue with remaining query functions...

pub fn get_subject_performance(
    deps: Deps,
    subject_id: String,
    institution_id: Option<String>,
    _period: Option<AnalyticsPeriod>,
) -> StdResult<SubjectPerformanceResponse> {
    // Analyze subject performance across students
    let mut enrollment_count = 0u32;
    let mut total_grades = 0u32;
    let mut grade_sum = 0u32;
    let mut completion_count = 0u32;
    let mut grade_distribution = HashMap::new();
    
    // Get all students (filtered by institution if specified)
    let students_to_check: Vec<String> = if let Some(inst_id) = institution_id {
        STUDENTS_BY_INSTITUTION.load(deps.storage, &inst_id).unwrap_or_default()
    } else {
        STUDENT_PROGRESS
            .range(deps.storage, None, None, Order::Ascending)
            .take(500)
            .map(|item| item.map(|(id, _)| id))
            .collect::<StdResult<Vec<String>>>()?
    };
    
    for student_id in students_to_check {
        if let Ok(progress) = STUDENT_PROGRESS.load(deps.storage, &student_id) {
            // Check completed subjects
            for completed in &progress.completed_subjects {
                if completed.subject_id == subject_id {
                    enrollment_count += 1;
                    completion_count += 1;
                    total_grades += 1;
                    grade_sum += completed.final_grade;
                    
                    let grade_range = match completed.final_grade {
                        90..=100 => "A (90-100)",
                        80..=89 => "B (80-89)",
                        70..=79 => "C (70-79)",
                        60..=69 => "D (60-69)",
                        _ => "F (0-59)",
                    };
                    *grade_distribution.entry(grade_range).or_insert(0u32) += 1;
                }
            }
            
            // Check failed subjects
            for failed in &progress.failed_subjects {
                if failed.subject_id == subject_id {
                    enrollment_count += 1;
                    total_grades += 1;
                    grade_sum += failed.final_grade;
                    
                    *grade_distribution.entry("F (0-59)").or_insert(0u32) += 1;
                }
            }
            
            // Check current subjects
            for current in &progress.current_subjects {
                if current.subject_id == subject_id {
                    enrollment_count += 1;
                }
            }
        }
    }
    
    let average_grade = if total_grades > 0 {
        grade_sum / total_grades
    } else {
        0
    };
    
    let completion_rate = if enrollment_count > 0 {
        (completion_count * 100) / enrollment_count
    } else {
        0
    };
    
    let subject_analytics = SubjectAnalytics {
        subject_id: subject_id.clone(),
        subject_title: subject_id.clone(), // Would fetch from subject registry
        enrollment_count,
        completion_rate,
        average_grade,
        difficulty_rating: 50, // Would calculate from student feedback
        student_satisfaction: 75, // Would calculate from ratings
        prerequisite_success_correlation: 80, // Would calculate from data
        retry_rate: if enrollment_count > 0 {
            ((enrollment_count - completion_count) * 100) / enrollment_count
        } else {
            0
        },
        common_failure_reasons: vec![
            "Insufficient prerequisite preparation".to_string(),
            "Heavy workload".to_string(),
        ],
    };
    
    let grade_distribution_points: Vec<GradeDistributionPoint> = grade_distribution
        .into_iter()
        .map(|(range, count)| {
            let percentage = if total_grades > 0 {
                (count * 100) / total_grades
            } else {
                0
            };
            GradeDistributionPoint {
                grade_range: range.to_string(),
                student_count: count,
                percentage,
            }
        })
        .collect();
    
    let success_predictors = vec![
        SuccessPredictor {
            predictor_name: "Strong prerequisite performance".to_string(),
            correlation_strength: 85,
            description: "Students with high grades in prerequisites tend to succeed".to_string(),
        },
        SuccessPredictor {
            predictor_name: "Regular class attendance".to_string(),
            correlation_strength: 75,
            description: "Attendance rate strongly correlates with final grades".to_string(),
        },
    ];
    
    let improvement_recommendations = vec![
        "Strengthen prerequisite enforcement".to_string(),
        "Provide additional tutoring resources".to_string(),
        "Implement early warning system for at-risk students".to_string(),
    ];
    
    Ok(SubjectPerformanceResponse {
        subject_analytics,
        grade_distribution: grade_distribution_points,
        success_predictors,
        improvement_recommendations,
    })
}

pub fn get_cohort_analysis(
    deps: Deps,
    institution_id: String,
    course_id: String,
    enrollment_year: u32,
    _analysis_type: CohortAnalysisType,
) -> StdResult<CohortAnalysisResponse> {
    // Get students from the specific cohort
    let all_students = STUDENTS_BY_INSTITUTION.load(deps.storage, &institution_id).unwrap_or_default();
    let mut cohort_students = Vec::new();
    
    for student_id in all_students {
        if let Ok(progress) = STUDENT_PROGRESS.load(deps.storage, &student_id) {
            if progress.course_id == course_id {
                // Simple enrollment year check (would parse enrollment_date properly)
                if progress.enrollment_date.contains(&enrollment_year.to_string()) {
                    cohort_students.push(progress);
                }
            }
        }
    }
    
    let cohort_size = cohort_students.len() as u32;
    let current_active = cohort_students.iter()
        .filter(|p| matches!(p.academic_status, AcademicStatus::Active))
        .count() as u32;
    let graduated = cohort_students.iter()
        .filter(|p| matches!(p.academic_status, AcademicStatus::Graduated))
        .count() as u32;
    let withdrawn = cohort_students.iter()
        .filter(|p| matches!(p.academic_status, AcademicStatus::Withdrawn | AcademicStatus::Expelled))
        .count() as u32;
    
    let average_time_to_completion = if graduated > 0 {
        let total_semesters: u32 = cohort_students.iter()
            .filter(|p| matches!(p.academic_status, AcademicStatus::Graduated))
            .map(|p| p.total_semesters_enrolled)
            .sum();
        total_semesters / graduated
    } else {
        0
    };
    
    let cohort_metrics = CohortMetrics {
        cohort_size,
        current_active,
        graduated,
        withdrawn,
        average_time_to_completion,
    };
    
    // Generate semester retention rates
    let mut semester_retention_rates = Vec::new();
    for semester in 1..=8 {
        let students_at_semester = cohort_students.iter()
            .filter(|p| p.current_semester >= semester || matches!(p.academic_status, AcademicStatus::Graduated))
            .count() as u32;
        
        let retention_rate = if cohort_size > 0 {
            (students_at_semester * 100) / cohort_size
        } else {
            0
        };
        
        semester_retention_rates.push(SemesterRetention {
            semester,
            retention_rate,
            dropout_count: cohort_size - students_at_semester,
            key_factors: vec![
                "Academic performance".to_string(),
                "Financial constraints".to_string(),
            ],
        });
    }
    
    let retention_analysis = RetentionAnalysis {
        semester_retention_rates,
        dropout_patterns: vec![
            DropoutPattern {
                semester_range: "1-2".to_string(),
                dropout_percentage: 15,
                common_reasons: vec![
                    "Academic adjustment difficulties".to_string(),
                    "Financial constraints".to_string(),
                ],
            },
            DropoutPattern {
                semester_range: "3-4".to_string(),
                dropout_percentage: 8,
                common_reasons: vec![
                    "Career change".to_string(),
                    "Transfer to other institution".to_string(),
                ],
            },
        ],
        retention_predictors: vec![
            "First semester GPA".to_string(),
            "High school performance".to_string(),
            "Financial aid status".to_string(),
        ],
    };
    
    // Calculate GPA quartiles
    let mut gpas: Vec<f64> = cohort_students.iter()
        .filter_map(|p| p.gpa.parse::<f64>().ok())
        .collect();
    gpas.sort_by(|a, b| a.partial_cmp(b).unwrap());
    
    let gpa_quartiles = if !gpas.is_empty() {
        let q1_index = gpas.len() / 4;
        let q2_index = gpas.len() / 2;
        let q3_index = (gpas.len() * 3) / 4;
        
        let mean = gpas.iter().sum::<f64>() / gpas.len() as f64;
        let variance = gpas.iter().map(|x| (x - mean).powi(2)).sum::<f64>() / gpas.len() as f64;
        let std_dev = variance.sqrt();
        
        GpaQuartiles {
            q1: format!("{:.2}", gpas[q1_index]),
            q2_median: format!("{:.2}", gpas[q2_index]),
            q3: format!("{:.2}", gpas[q3_index]),
            mean: format!("{:.2}", mean),
            standard_deviation: format!("{:.2}", std_dev),
        }
    } else {
        GpaQuartiles {
            q1: "0.00".to_string(),
            q2_median: "0.00".to_string(),
            q3: "0.00".to_string(),
            mean: "0.00".to_string(),
            standard_deviation: "0.00".to_string(),
        }
    };
    
    let performance_distribution = PerformanceDistribution {
        gpa_quartiles,
        credit_completion_rates: vec![
            CompletionRateDistribution {
                rate_range: "90-100%".to_string(),
                student_count: cohort_students.iter()
                    .filter(|p| p.performance_metrics.success_rate >= 90)
                    .count() as u32,
                percentage: 0, // Would calculate
            },
        ],
        time_to_graduation_distribution: vec![
            TimeDistribution {
                time_range_semesters: "8-10".to_string(),
                student_count: graduated,
                percentage: 0, // Would calculate
            },
        ],
    };
    
    let comparative_insights = vec![
        "Cohort performing above institutional average".to_string(),
        "Retention rate higher than previous cohorts".to_string(),
    ];
    
    Ok(CohortAnalysisResponse {
        cohort_metrics,
        retention_analysis,
        performance_distribution,
        comparative_insights,
    })
}

pub fn get_predictive_insights(
    deps: Deps,
    student_id: String,
    prediction_horizon_semesters: u32,
) -> StdResult<PredictiveInsightsResponse> {
    let student_progress = STUDENT_PROGRESS.load(deps.storage, &student_id)?;
    
    // Generate predictions
    let mut predictions = Vec::new();
    
    // Graduation probability prediction
    predictions.push(Prediction {
        prediction_type: PredictionType::GraduationProbability,
        predicted_outcome: format!("{}% probability", student_progress.performance_metrics.graduation_probability),
        confidence_level: 75,
        time_horizon: format!("{} semesters", prediction_horizon_semesters),
        key_factors: vec![
            "Current GPA trend".to_string(),
            "Credit completion rate".to_string(),
            "Historical performance".to_string(),
        ],
    });
    
    // Next semester GPA prediction
    let current_gpa = student_progress.gpa.parse::<f64>().unwrap_or(0.0);
    let predicted_gpa = match student_progress.performance_metrics.grade_trend {
        GradeTrend::StronglyIncreasing => current_gpa + 0.3,
        GradeTrend::ModeratelyIncreasing => current_gpa + 0.1,
        GradeTrend::Stable => current_gpa,
        GradeTrend::ModeratelyDecreasing => current_gpa - 0.1,
        GradeTrend::StronglyDecreasing => current_gpa - 0.3,
        GradeTrend::Volatile => current_gpa,
    }.min(10.0).max(0.0);
    
    predictions.push(Prediction {
        prediction_type: PredictionType::NextSemesterGPA,
        predicted_outcome: format!("{:.2}", predicted_gpa),
        confidence_level: 70,
        time_horizon: "1 semester".to_string(),
        key_factors: vec![
            "Grade trend analysis".to_string(),
            "Current subject performance".to_string(),
        ],
    });
    
    // Risk assessments
    let risk_assessments: Vec<RiskAssessment> = student_progress.risk_factors.iter().map(|rf| {
        let probability = match rf.severity {
            RiskSeverity::Critical => 85,
            RiskSeverity::High => 70,
            RiskSeverity::Moderate => 45,
            RiskSeverity::Low => 20,
        };
        
        RiskAssessment {
            risk_factor: rf.risk_type.clone(),
            probability,
            impact_severity: rf.severity.clone(),
            early_indicators: rf.early_warning_indicators.clone(),
            mitigation_strategies: rf.intervention_recommendations.clone(),
        }
    }).collect();
    
    // Opportunity identifications
    let mut opportunities = Vec::new();
    
    if student_progress.performance_metrics.graduation_probability > 80 {
        opportunities.push(Opportunity {
            opportunity_type: OpportunityType::AcceleratedCompletion,
            description: "Student shows potential for early graduation".to_string(),
            potential_benefit: "Graduate 1-2 semesters early".to_string(),
            action_required: vec![
                "Increase course load".to_string(),
                "Take summer courses".to_string(),
            ],
            success_probability: 70,
        });
    }
    
    if let Ok(gpa) = student_progress.gpa.parse::<f64>() {
        if gpa >= 8.5 {
            opportunities.push(Opportunity {
                opportunity_type: OpportunityType::AcademicImprovement,
                description: "Potential for academic honors".to_string(),
                potential_benefit: "Graduate with honors".to_string(),
                action_required: vec![
                    "Maintain current performance".to_string(),
                    "Consider research opportunities".to_string(),
                ],
                success_probability: 80,
            });
        }
    }
    
    // Recommended interventions
    let mut interventions = Vec::new();
    
    for risk_factor in &student_progress.risk_factors {
        let intervention_type = match risk_factor.risk_type {
            RiskFactorType::AcademicPerformance => InterventionType::AcademicSupport,
            RiskFactorType::Attendance => InterventionType::Counseling,
            RiskFactorType::TimeManagement => InterventionType::TimeManagementCoaching,
            _ => InterventionType::AcademicSupport,
        };
        
        interventions.push(Intervention {
            intervention_id: format!("{}_{:?}", student_id, risk_factor.risk_type),
            intervention_type,
            target_issue: risk_factor.description.clone(),
            recommended_actions: risk_factor.intervention_recommendations.clone(),
            expected_outcome: "Improved academic performance".to_string(),
            success_probability: 75,
            effort_required: match risk_factor.severity {
                RiskSeverity::Critical => crate::msg::EffortLevel::Intensive,
                RiskSeverity::High => crate::msg::EffortLevel::High,
                RiskSeverity::Moderate => crate::msg::EffortLevel::Medium,
                RiskSeverity::Low => crate::msg::EffortLevel::Low,
            },
            timeline: "1-2 semesters".to_string(),
        });
    }
    
    Ok(PredictiveInsightsResponse {
        predictions,
        risk_assessments,
        opportunity_identifications: opportunities,
        recommended_interventions: interventions,
    })
}

pub fn search_students(
    deps: Deps,
    query: StudentSearchQuery,
    limit: Option<u32>,
    start_after: Option<String>,
) -> StdResult<StudentSearchResponse> {
    let limit = limit.unwrap_or(50).min(100);
    let mut matching_students = Vec::new();
    let mut total_matches = 0u32;
    
    // Get students to search through
    let students_to_search: Vec<String> = if let Some(inst_id) = &query.institution_id {
        STUDENTS_BY_INSTITUTION.load(deps.storage, inst_id).unwrap_or_default()
    } else if let Some(course_id) = &query.course_id {
        STUDENTS_BY_COURSE.load(deps.storage, course_id).unwrap_or_default()
    } else {
        STUDENT_PROGRESS
            .range(deps.storage, None, None, Order::Ascending)
            .take(500)
            .map(|item| item.map(|(id, _)| id))
            .collect::<StdResult<Vec<String>>>()?
    };
    
    let mut skip_count = 0usize;
    if let Some(start_after_id) = &start_after {
        skip_count = students_to_search.iter()
            .position(|id| id == start_after_id)
            .map(|pos| pos + 1)
            .unwrap_or(0);
    }
    
    for student_id in students_to_search.iter().skip(skip_count) {
        if let Ok(progress) = STUDENT_PROGRESS.load(deps.storage, student_id) {
            let mut matches = true;
            let mut relevance_score = 0u32;
            let mut highlighted_attributes = Vec::new();
            
            // Apply filters
            if let Some(ref name_filter) = query.name_contains {
                // In a real implementation, would have student names
                if !student_id.to_lowercase().contains(&name_filter.to_lowercase()) {
                    matches = false;
                } else {
                    relevance_score += 10;
                    highlighted_attributes.push("Student ID match".to_string());
                }
            }
            
            if let Some(ref status_filter) = query.status {
                if progress.academic_status != *status_filter {
                    matches = false;
                } else {
                    relevance_score += 5;
                }
            }
            
            if let Some(ref gpa_min) = query.gpa_min {
                if let Ok(student_gpa) = progress.gpa.parse::<f64>() {
                    if let Ok(min_gpa) = gpa_min.parse::<f64>() {
                        if student_gpa < min_gpa {
                            matches = false;
                        } else {
                            relevance_score += 3;
                        }
                    }
                }
            }
            
            if let Some(ref gpa_max) = query.gpa_max {
                if let Ok(student_gpa) = progress.gpa.parse::<f64>() {
                    if let Ok(max_gpa) = gpa_max.parse::<f64>() {
                        if student_gpa > max_gpa {
                            matches = false;
                        } else {
                            relevance_score += 3;
                        }
                    }
                }
            }
            
            if let Some(credits_min) = query.credits_min {
                if progress.total_credits_completed < credits_min {
                    matches = false;
                } else {
                    relevance_score += 2;
                }
            }
            
            if let Some(credits_max) = query.credits_max {
                if progress.total_credits_completed > credits_max {
                    matches = false;
                } else {
                    relevance_score += 2;
                }
            }
            
            if let Some(ref risk_filter) = query.risk_level {
                let student_risk = determine_overall_risk_level(&progress);
                if student_risk != *risk_filter {
                    matches = false;
                } else {
                    relevance_score += 8;
                    highlighted_attributes.push("Risk level match".to_string());
                }
            }
            
            if matches {
                total_matches += 1;
                
                if matching_students.len() < limit as usize {
                    let completion_percentage = if progress.total_credits_required > 0 {
                        (progress.total_credits_completed * 100) / progress.total_credits_required
                    } else {
                        0
                    };
                    
                    let student_summary = StudentSummary {
                        student_id: student_id.clone(),
                        institution_id: progress.institution_id.clone(),
                        course_id: progress.course_id.clone(),
                        current_semester: progress.current_semester,
                        gpa: progress.gpa.clone(),
                        academic_status: progress.academic_status.clone(),
                        completion_percentage,
                    };
                    
                    matching_students.push(StudentSearchResult {
                        student_id: student_id.clone(),
                        basic_info: student_summary,
                        match_relevance: relevance_score,
                        highlighted_attributes,
                    });
                }
            }
        }
    }
    
    // Sort by relevance score
    matching_students.sort_by(|a, b| b.match_relevance.cmp(&a.match_relevance));
    
    let search_quality = match total_matches {
        0 => SearchQuality::Limited,
        1..=10 => SearchQuality::Good,
        11..=50 => SearchQuality::Excellent,
        _ => SearchQuality::Fair,
    };
    
    let search_metadata = SearchMetadata {
        query_time_ms: 50, // Would measure actual query time
        filters_applied: vec![
            if query.status.is_some() { "status".to_string() } else { "".to_string() },
            if query.gpa_min.is_some() { "gpa_min".to_string() } else { "".to_string() },
            if query.risk_level.is_some() { "risk_level".to_string() } else { "".to_string() },
        ].into_iter().filter(|s| !s.is_empty()).collect(),
        result_quality: search_quality,
    };
    
    Ok(StudentSearchResponse {
        students: matching_students,
        total_matches,
        search_metadata,
    })
}

pub fn get_analytics_summary(
    deps: Deps,
    _env: Env,
    scope: AnalyticsScope,
    quick_stats_only: Option<bool>,
) -> StdResult<AnalyticsSummaryResponse> {
    let quick_stats_only = quick_stats_only.unwrap_or(false);
    
    let (quick_stats, _scope_description) = match scope {
        AnalyticsScope::Global => {
            let state = STATE.load(deps.storage)?;
            (
                GlobalQuickStats {
                    total_students: state.total_students as u32,
                    active_students: 0, // Would calculate from status index
                    graduation_rate: 75, // Would calculate from data
                    average_gpa: "7.8".to_string(), // Would calculate global average
                    at_risk_percentage: 15, // Would calculate from risk indices
                },
                "Global system overview".to_string(),
            )
        },
        AnalyticsScope::Institution(ref inst_id) => {
            let students = STUDENTS_BY_INSTITUTION.load(deps.storage, &inst_id).unwrap_or_default();
            let total_students = students.len() as u32;
            
            // Calculate active students
            let active_students = students.iter()
                .filter_map(|id| STUDENT_PROGRESS.load(deps.storage, id).ok())
                .filter(|p| matches!(p.academic_status, AcademicStatus::Active))
                .count() as u32;
            
            (
                GlobalQuickStats {
                    total_students,
                    active_students,
                    graduation_rate: 80, // Would calculate from institution data
                    average_gpa: "7.9".to_string(), // Would calculate from students
                    at_risk_percentage: 12, // Would calculate from students
                },
                format!("Institution {} overview", inst_id),
            )
        },
        AnalyticsScope::Course(ref course_id) => {
            let students = STUDENTS_BY_COURSE.load(deps.storage, &course_id).unwrap_or_default();
            (
                GlobalQuickStats {
                    total_students: students.len() as u32,
                    active_students: students.len() as u32, // Simplified
                    graduation_rate: 85,
                    average_gpa: "8.1".to_string(),
                    at_risk_percentage: 8,
                },
                format!("Course {} overview", course_id),
            )
        },
        AnalyticsScope::Student(ref student_id) => {
            let progress = STUDENT_PROGRESS.load(deps.storage, &student_id)?;
            let completion_percentage = if progress.total_credits_required > 0 {
                (progress.total_credits_completed * 100) / progress.total_credits_required
            } else {
                0
            };
            
            (
                GlobalQuickStats {
                    total_students: 1,
                    active_students: if matches!(progress.academic_status, AcademicStatus::Active) { 1 } else { 0 },
                    graduation_rate: completion_percentage,
                    average_gpa: progress.gpa.clone(),
                    at_risk_percentage: if !progress.risk_factors.is_empty() { 100 } else { 0 },
                },
                format!("Student {} overview", student_id),
            )
        },
    };
    
    let summary = AnalyticsSummary {
        scope,
        period: AnalyticsPeriod {
            start_date: "current_period".to_string(),
            end_date: "current_period".to_string(),
            period_type: PeriodType::AcademicYear,
        },
        quick_stats,
        performance_overview: if quick_stats_only {
            PerformanceOverview {
                top_performing_institutions: vec![],
                improvement_leaders: vec![],
                areas_needing_attention: vec![],
            }
        } else {
            PerformanceOverview {
                top_performing_institutions: vec!["Institution A".to_string()],
                improvement_leaders: vec!["Institution B".to_string()],
                areas_needing_attention: vec!["Student retention".to_string()],
            }
        },
        trend_summary: if quick_stats_only {
            TrendSummary {
                enrollment_trend: crate::state::TrendDirection::Stable,
                performance_trend: crate::state::TrendDirection::Stable,
                completion_trend: crate::state::TrendDirection::Stable,
                key_trend_drivers: vec![],
            }
        } else {
            TrendSummary {
                enrollment_trend: crate::state::TrendDirection::Improving,
                performance_trend: crate::state::TrendDirection::Improving,
                completion_trend: crate::state::TrendDirection::Stable,
                key_trend_drivers: vec![
                    "Improved academic support".to_string(),
                    "Enhanced student services".to_string(),
                ],
            }
        },
    };
    
    let key_insights = if quick_stats_only {
        vec![]
    } else {
        vec![
            KeyInsight {
                insight_category: InsightCategory::PerformancePattern,
                title: "Strong academic performance trend".to_string(),
                description: "Overall GPA has improved over the last year".to_string(),
                supporting_metrics: vec!["Average GPA: 7.8".to_string()],
                confidence_level: 85,
            },
            KeyInsight {
                insight_category: InsightCategory::RiskIdentification,
                title: "At-risk student population stable".to_string(),
                description: "15% of students identified as at-risk, consistent with previous periods".to_string(),
                supporting_metrics: vec!["At-risk percentage: 15%".to_string()],
                confidence_level: 90,
            },
        ]
    };
    
    let actionable_recommendations = if quick_stats_only {
        vec![]
    } else {
        vec![
            "Implement early intervention programs for at-risk students".to_string(),
            "Expand academic support services".to_string(),
            "Monitor graduation rate trends closely".to_string(),
        ]
    };
    
    Ok(AnalyticsSummaryResponse {
        summary,
        key_insights,
        actionable_recommendations,
    })
}

// Helper functions for query operations

fn determine_overall_risk_level(progress: &StudentProgress) -> RiskSeverity {
    if progress.risk_factors.is_empty() {
        return RiskSeverity::Low;
    }
    
    let max_severity = progress.risk_factors.iter()
        .map(|rf| match rf.severity {
            RiskSeverity::Critical => 4,
            RiskSeverity::High => 3,
            RiskSeverity::Moderate => 2,
            RiskSeverity::Low => 1,
        })
        .max()
        .unwrap_or(1);
    
    match max_severity {
        4 => RiskSeverity::Critical,
        3 => RiskSeverity::High,
        2 => RiskSeverity::Moderate,
        _ => RiskSeverity::Low,
    }
}

fn get_comparison_benchmarks(
    _storage: &dyn Storage,
    comparison_group: &ComparisonGroup,
) -> StdResult<HashMap<String, String>> {
    // In a real implementation, this would load actual benchmark data
    let mut benchmarks = HashMap::new();
    
    match comparison_group {
        ComparisonGroup::NationalAverage => {
            benchmarks.insert("average_gpa".to_string(), "7.5".to_string());
            benchmarks.insert("completion_rate".to_string(), "80".to_string());
            benchmarks.insert("dropout_rate".to_string(), "15".to_string());
            benchmarks.insert("on_track_percentage".to_string(), "75".to_string());
        },
        ComparisonGroup::InstitutionType => {
            benchmarks.insert("average_gpa".to_string(), "7.8".to_string());
            benchmarks.insert("completion_rate".to_string(), "85".to_string());
            benchmarks.insert("dropout_rate".to_string(), "12".to_string());
            benchmarks.insert("on_track_percentage".to_string(), "80".to_string());
        },
        _ => {
            // Default benchmarks
            benchmarks.insert("average_gpa".to_string(), "7.5".to_string());
            benchmarks.insert("completion_rate".to_string(), "80".to_string());
           benchmarks.insert("dropout_rate".to_string(), "15".to_string());
           benchmarks.insert("on_track_percentage".to_string(), "75".to_string());
       }
   }
   
   Ok(benchmarks)
}

fn calculate_metric_percentile(metric_name: &str, value: &str, _comparison_group: &ComparisonGroup) -> u32 {
   // Simplified percentile calculation
   // In a real implementation, this would compare against actual distribution data
   match metric_name {
       "average_gpa" => {
           if let Ok(gpa) = value.parse::<f64>() {
               match gpa {
                   x if x >= 9.0 => 95,
                   x if x >= 8.5 => 85,
                   x if x >= 8.0 => 75,
                   x if x >= 7.5 => 60,
                   x if x >= 7.0 => 45,
                   x if x >= 6.5 => 30,
                   _ => 15,
               }
           } else {
               50
           }
       },
       "completion_rate" => {
           if let Ok(rate) = value.parse::<u32>() {
               match rate {
                   x if x >= 95 => 95,
                   x if x >= 90 => 85,
                   x if x >= 85 => 75,
                   x if x >= 80 => 60,
                   x if x >= 75 => 45,
                   x if x >= 70 => 30,
                   _ => 15,
               }
           } else {
               50
           }
       },
       "dropout_rate" => {
           if let Ok(rate) = value.parse::<u32>() {
               // Lower dropout rate is better, so invert the percentile
               match rate {
                   x if x <= 5 => 95,
                   x if x <= 10 => 85,
                   x if x <= 15 => 75,
                   x if x <= 20 => 60,
                   x if x <= 25 => 45,
                   x if x <= 30 => 30,
                   _ => 15,
               }
           } else {
               50
           }
       },
       _ => 50, // Default median percentile
   }
}

pub fn generate_peer_comparison(
   storage: &dyn Storage,
   student_progress: &StudentProgress,
) -> StdResult<PeerComparison> {
   // Get peers (students in same course and similar semester)
   let course_students = STUDENTS_BY_COURSE.load(storage, &student_progress.course_id).unwrap_or_default();
   
   let mut peer_gpas = Vec::new();
   let mut peer_progress_rates = Vec::new();
   let mut peer_credit_rates = Vec::new();
   let mut peer_grade_trends = Vec::new();
   
   for peer_id in course_students {
       if peer_id == student_progress.student_id {
           continue; // Skip self
       }
       
       if let Ok(peer_progress) = STUDENT_PROGRESS.load(storage, &peer_id) {
           // Only compare with students in similar semester range
           if (peer_progress.current_semester as i32 - student_progress.current_semester as i32).abs() <= 1 {
               if let Ok(gpa) = peer_progress.gpa.parse::<f64>() {
                   peer_gpas.push(gpa);
               }
               
               if peer_progress.total_credits_required > 0 {
                   let progress_rate = (peer_progress.total_credits_completed * 100) / peer_progress.total_credits_required;
                   peer_progress_rates.push(progress_rate);
               }
               
               if peer_progress.current_semester > 0 {
                   let credit_rate = peer_progress.total_credits_completed / peer_progress.current_semester;
                   peer_credit_rates.push(credit_rate);
               }
               
               peer_grade_trends.push(peer_progress.performance_metrics.grade_trend.clone());
           }
       }
   }
   
   // Calculate student's GPA percentile
   let student_gpa = student_progress.gpa.parse::<f64>().unwrap_or(0.0);
   let gpa_percentile = if !peer_gpas.is_empty() {
       let better_count = peer_gpas.iter().filter(|&&gpa| gpa < student_gpa).count();
       ((better_count * 100) / peer_gpas.len()) as u32
   } else {
       50
   };
   
   // Calculate progress percentile
   let student_progress_rate = if student_progress.total_credits_required > 0 {
       (student_progress.total_credits_completed * 100) / student_progress.total_credits_required
   } else {
       0
   };
   
   let progress_percentile = if !peer_progress_rates.is_empty() {
       let better_count = peer_progress_rates.iter().filter(|&&rate| rate < student_progress_rate).count();
       ((better_count * 100) / peer_progress_rates.len()) as u32
   } else {
       50
   };
   
   // Calculate credit completion rate comparison
   let student_credit_rate = if student_progress.current_semester > 0 {
       student_progress.total_credits_completed / student_progress.current_semester
   } else {
       0
   };
   
   let credit_completion_rate_vs_peers = if !peer_credit_rates.is_empty() {
       let avg_peer_rate = peer_credit_rates.iter().sum::<u32>() / peer_credit_rates.len() as u32;
       match student_credit_rate {
           rate if rate > avg_peer_rate + 2 => ComparisonResult::TopTier,
           rate if rate > avg_peer_rate => ComparisonResult::AboveAverage,
           rate if rate >= avg_peer_rate.saturating_sub(2) => ComparisonResult::Average,
           _ => ComparisonResult::BelowAverage,
       }
   } else {
       ComparisonResult::Average
   };
   
   // Calculate grade trend comparison
   let improving_trends = peer_grade_trends.iter()
       .filter(|trend| matches!(trend, GradeTrend::StronglyIncreasing | GradeTrend::ModeratelyIncreasing))
       .count();
   
   let grade_trend_vs_peers = match student_progress.performance_metrics.grade_trend {
       GradeTrend::StronglyIncreasing => ComparisonResult::TopTier,
       GradeTrend::ModeratelyIncreasing => {
           if improving_trends < peer_grade_trends.len() / 2 {
               ComparisonResult::AboveAverage
           } else {
               ComparisonResult::Average
           }
       },
       GradeTrend::Stable => ComparisonResult::Average,
       _ => ComparisonResult::BelowAverage,
   };
   
   // Generate anonymous peer insights
   let mut anonymous_peer_insights = Vec::new();
   
   if gpa_percentile > 75 {
       anonymous_peer_insights.push("Your GPA is higher than most of your peers".to_string());
   } else if gpa_percentile < 25 {
       anonymous_peer_insights.push("Your GPA is below most of your peers - consider seeking academic support".to_string());
   } else {
       anonymous_peer_insights.push("Your GPA is similar to your peers".to_string());
   }
   
   if progress_percentile > 75 {
       anonymous_peer_insights.push("You're progressing faster than most peers".to_string());
   } else if progress_percentile < 25 {
       anonymous_peer_insights.push("Consider accelerating your progress to match peer pace".to_string());
   }
   
   match credit_completion_rate_vs_peers {
       ComparisonResult::TopTier => {
           anonymous_peer_insights.push("You're completing credits at an exceptional rate".to_string());
       },
       ComparisonResult::BelowAverage => {
           anonymous_peer_insights.push("Consider increasing your course load to match peer progress".to_string());
       },
       _ => {}
   }
   
   Ok(PeerComparison {
       gpa_percentile,
       progress_percentile,
       credit_completion_rate_vs_peers,
       grade_trend_vs_peers,
       anonymous_peer_insights,
   })
}