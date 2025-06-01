// src/analytics.rs
use cosmwasm_std::{StdResult, StdError};
use crate::state::*;

/// Core analytics engine for progress tracking
pub struct AnalyticsEngine;

impl AnalyticsEngine {
    /// Calculate comprehensive performance metrics for a student
    pub fn calculate_performance_metrics(
        student_progress: &StudentProgress,
        institution_benchmarks: Option<&InstitutionAnalytics>,
    ) -> StdResult<PerformanceMetrics> {
        let total_subjects = student_progress.completed_subjects.len() + 
                           student_progress.failed_subjects.len();
        
        if total_subjects == 0 {
            return Ok(PerformanceMetrics {
                success_rate: 0,
                retry_rate: 0,
                grade_trend: GradeTrend::Stable,
                consistency_score: 0,
                attendance_average: 0,
                assignment_completion_rate: 0,
                study_hours_per_credit: 0,
                credits_per_semester_avg: 0,
                time_to_completion_ratio: 100,
                prerequisite_success_rate: 0,
                percentile_rank_gpa: 50,
                percentile_rank_progress: 50,
                graduation_probability: 50,
                time_extension_probability: 50,
            });
        }
        
        // Calculate success rate
        let success_count = student_progress.completed_subjects.len() as u32;
        let success_rate = (success_count * 100) / total_subjects as u32;
        
        // Calculate retry rate
        let retry_count = student_progress.failed_subjects.len() as u32;
        let retry_rate = (retry_count * 100) / total_subjects as u32;
        
        // Calculate grade trend
        let grade_trend = Self::calculate_grade_trend(&student_progress.completed_subjects)?;
        
        // Calculate consistency score (inverse of grade variance)
        let consistency_score = Self::calculate_consistency_score(&student_progress.completed_subjects)?;
        
        // Calculate average attendance from current subjects
        let attendance_average = Self::calculate_average_attendance(&student_progress.current_subjects);
        
        // Calculate assignment completion rate
        let assignment_completion_rate = Self::calculate_assignment_completion_rate(&student_progress.current_subjects);
        
        // Calculate study efficiency
        let study_hours_per_credit = Self::calculate_study_efficiency(&student_progress.completed_subjects);
        
        // Calculate credits per semester
        let credits_per_semester_avg = if student_progress.current_semester > 0 {
            student_progress.total_credits_completed / student_progress.current_semester
        } else {
            0
        };
        
        // Calculate time to completion ratio
        let time_to_completion_ratio = Self::calculate_time_completion_ratio(student_progress);
        
        // Calculate prerequisite success rate (simplified)
        let prerequisite_success_rate = if success_count > 0 { 
            (success_count * 100) / (success_count + retry_count) 
        } else { 
            0 
        };
        
        // Calculate percentile ranks (if benchmarks available)
        let (percentile_rank_gpa, percentile_rank_progress) = if let Some(benchmarks) = institution_benchmarks {
            Self::calculate_percentile_ranks(student_progress, benchmarks)?
        } else {
            (50, 50) // Default to median if no benchmarks
        };
        
        // Calculate graduation probability
        let graduation_probability = Self::calculate_graduation_probability(student_progress)?;
        
        // Calculate time extension probability
        let time_extension_probability = Self::calculate_time_extension_probability(student_progress)?;
        
        Ok(PerformanceMetrics {
            success_rate,
            retry_rate,
            grade_trend,
            consistency_score,
            attendance_average,
            assignment_completion_rate,
            study_hours_per_credit,
            credits_per_semester_avg,
            time_to_completion_ratio,
            prerequisite_success_rate,
            percentile_rank_gpa,
            percentile_rank_progress,
            graduation_probability,
            time_extension_probability,
        })
    }
    
    /// Calculate grade trend from completed subjects
    fn calculate_grade_trend(completed_subjects: &[CompletedSubjectProgress]) -> StdResult<GradeTrend> {
        if completed_subjects.len() < 3 {
            return Ok(GradeTrend::Stable);
        }
        
        // Take last 6 subjects for trend analysis
        let recent_subjects = if completed_subjects.len() > 6 {
            &completed_subjects[completed_subjects.len()-6..]
        } else {
            completed_subjects
        };
        
        let first_half_avg = recent_subjects.iter()
            .take(recent_subjects.len() / 2)
            .map(|s| s.final_grade)
            .sum::<u32>() as f64 / (recent_subjects.len() / 2) as f64;
            
        let second_half_avg = recent_subjects.iter()
            .skip(recent_subjects.len() / 2)
            .map(|s| s.final_grade)
            .sum::<u32>() as f64 / (recent_subjects.len() - recent_subjects.len() / 2) as f64;
        
        let difference = second_half_avg - first_half_avg;
        
        match difference {
            d if d > 5.0 => Ok(GradeTrend::StronglyIncreasing),
            d if d > 2.0 => Ok(GradeTrend::ModeratelyIncreasing),
            d if d > -2.0 => Ok(GradeTrend::Stable),
            d if d > -5.0 => Ok(GradeTrend::ModeratelyDecreasing),
            _ => Ok(GradeTrend::StronglyDecreasing),
        }
    }
    
    /// Calculate consistency score (100 - normalized variance)
    fn calculate_consistency_score(completed_subjects: &[CompletedSubjectProgress]) -> StdResult<u32> {
        if completed_subjects.is_empty() {
            return Ok(0);
        }
        
        let grades: Vec<f64> = completed_subjects.iter()
            .map(|s| s.final_grade as f64)
            .collect();
        
        let mean = grades.iter().sum::<f64>() / grades.len() as f64;
        let variance = grades.iter()
            .map(|g| (g - mean).powi(2))
            .sum::<f64>() / grades.len() as f64;
        
        let std_dev = variance.sqrt();
        
        // Normalize std_dev to 0-100 scale (assuming max reasonable std_dev is 20)
        let normalized_variance = (std_dev / 20.0 * 100.0).min(100.0);
        let consistency_score = (100.0 - normalized_variance).max(0.0) as u32;
        
        Ok(consistency_score)
    }
    
    /// Calculate average attendance from current subjects
    fn calculate_average_attendance(current_subjects: &[CurrentSubjectProgress]) -> u32 {
        if current_subjects.is_empty() {
            return 0;
        }
        
        let total_attendance: u32 = current_subjects.iter()
            .map(|s| s.attendance_rate.unwrap_or(0))
            .sum();
        
        total_attendance / current_subjects.len() as u32
    }
    
    /// Calculate assignment completion rate
    fn calculate_assignment_completion_rate(current_subjects: &[CurrentSubjectProgress]) -> u32 {
        if current_subjects.is_empty() {
            return 0;
        }
        
        let total_completed: u32 = current_subjects.iter()
            .map(|s| s.assignments_completed)
            .sum();
        
        let total_assignments: u32 = current_subjects.iter()
            .map(|s| s.assignments_total)
            .sum();
        
        if total_assignments > 0 {
            (total_completed * 100) / total_assignments
        } else {
            0
        }
    }
    
    /// Calculate study hours per credit
    fn calculate_study_efficiency(completed_subjects: &[CompletedSubjectProgress]) -> u32 {
        let total_hours: u32 = completed_subjects.iter()
            .map(|s| s.study_hours_logged.unwrap_or(0))
            .sum();
        
        let total_credits: u32 = completed_subjects.iter()
            .map(|s| s.credits)
            .sum();
        
        if total_credits > 0 {
            total_hours / total_credits
        } else {
            0
        }
    }
    
    /// Calculate time to completion ratio
    fn calculate_time_completion_ratio(student_progress: &StudentProgress) -> u32 {
        let completion_percentage = if student_progress.total_credits_required > 0 {
            (student_progress.total_credits_completed * 100) / student_progress.total_credits_required
        } else {
            0
        };
        
        let time_percentage = if student_progress.current_semester > 0 {
            // Assuming 8 semesters is standard completion time
            (student_progress.current_semester * 100) / 8
        } else {
            0
        };
        
        if time_percentage > 0 {
            (completion_percentage * 100) / time_percentage
        } else {
            100
        }
    }
    
    /// Calculate percentile ranks against institution benchmarks
    fn calculate_percentile_ranks(
        student_progress: &StudentProgress,
        benchmarks: &InstitutionAnalytics,
    ) -> StdResult<(u32, u32)> {
        // Parse student GPA
        let student_gpa = student_progress.gpa.parse::<f64>()
            .map_err(|_| StdError::generic_err("Invalid GPA format"))?;
        
        // Parse benchmark average GPA
        let benchmark_gpa = benchmarks.average_gpa.parse::<f64>()
            .map_err(|_| StdError::generic_err("Invalid benchmark GPA format"))?;
        
        // Simple percentile calculation (in real implementation, would use actual distribution)
        let gpa_percentile = if student_gpa >= benchmark_gpa {
            50 + ((student_gpa - benchmark_gpa) / benchmark_gpa * 40.0) as u32
        } else {
            50 - ((benchmark_gpa - student_gpa) / benchmark_gpa * 40.0) as u32
        }.min(100);
        
        // Calculate progress percentile based on completion rate vs time
        let progress_percentile = if student_progress.total_credits_required > 0 {
            let completion_rate = (student_progress.total_credits_completed * 100) / 
                                 student_progress.total_credits_required;
            // Compare against on-track percentage
            if completion_rate >= benchmarks.on_track_percentage {
                (50 + (completion_rate - benchmarks.on_track_percentage)).min(100)
            } else {
                (50 - (benchmarks.on_track_percentage - completion_rate) / 2).max(0)
            }
        } else {
            50
        };
        
        Ok((gpa_percentile, progress_percentile))
    }
    
    /// Calculate graduation probability based on current trajectory
    fn calculate_graduation_probability(student_progress: &StudentProgress) -> StdResult<u32> {
        let mut probability = 50u32; // Base probability
        
        // Parse GPA
        if let Ok(gpa) = student_progress.gpa.parse::<f64>() {
            // Higher GPA increases probability
            if gpa >= 9.0 {
                probability += 30;
            } else if gpa >= 8.0 {
                probability += 20;
            } else if gpa >= 7.0 {
                probability += 10;
            } else if gpa < 6.0 {
                probability = probability.saturating_sub(20);
            }
        }
        
        // Progress rate affects probability
        if student_progress.total_credits_required > 0 {
            let completion_rate = (student_progress.total_credits_completed * 100) / 
                                 student_progress.total_credits_required;
            
            if completion_rate >= 75 {
                probability += 20;
            } else if completion_rate >= 50 {
                probability += 10;
            } else if completion_rate < 25 {
                probability = probability.saturating_sub(15);
            }
        }
        
        // Academic status affects probability
        match student_progress.academic_status {
            AcademicStatus::Active => {} // No change
            AcademicStatus::Probation => probability = probability.saturating_sub(25),
            AcademicStatus::Suspended => probability = probability.saturating_sub(40),
            AcademicStatus::Graduated => probability = 100,
            AcademicStatus::Withdrawn | AcademicStatus::Expelled => probability = 0,
            AcademicStatus::Transfer => probability = 90, // High probability at new institution
        }
        
        // Success rate in completed subjects
        let total_subjects = student_progress.completed_subjects.len() + 
                           student_progress.failed_subjects.len();
        if total_subjects > 0 {
            let success_rate = (student_progress.completed_subjects.len() * 100) / total_subjects;
            if success_rate >= 90 {
                probability += 15;
            } else if success_rate >= 80 {
                probability += 10;
            } else if success_rate < 70 {
                probability = probability.saturating_sub(10);
            }
        }
        
        Ok(probability.min(100))
    }
    
    /// Calculate probability of needing time extension
    fn calculate_time_extension_probability(student_progress: &StudentProgress) -> StdResult<u32> {
        let mut probability = 20u32; // Base probability
        
        // Current semester vs expected progress
        if student_progress.total_credits_required > 0 {
            let expected_semester = (student_progress.total_credits_completed * 8) / 
                                   student_progress.total_credits_required;
            
            if student_progress.current_semester > expected_semester + 2 {
                probability += 40;
            } else if student_progress.current_semester > expected_semester + 1 {
                probability += 25;
            } else if student_progress.current_semester < expected_semester {
                probability = probability.saturating_sub(15);
            }
        }
        
        // Failed subjects increase probability
        let failure_count = student_progress.failed_subjects.len() as u32;
        probability += failure_count * 10;
        
        // Academic status
        match student_progress.academic_status {
            AcademicStatus::Probation => probability += 30,
            AcademicStatus::Suspended => probability += 50,
            AcademicStatus::Active => {} // No change
            _ => {} // Other statuses don't affect this calculation
        }
        
        // Parse GPA impact
        if let Ok(gpa) = student_progress.gpa.parse::<f64>() {
            if gpa < 6.0 {
                probability += 25;
            } else if gpa < 7.0 {
                probability += 15;
            } else if gpa >= 9.0 {
                probability = probability.saturating_sub(10);
            }
        }
        
        Ok(probability.min(100))
    }
    
    /// Generate risk factors based on student data
    pub fn identify_risk_factors(student_progress: &StudentProgress) -> Vec<RiskFactor> {
        let mut risk_factors = Vec::new();
        
        // Academic performance risk
        if let Ok(gpa) = student_progress.gpa.parse::<f64>() {
            if gpa < 6.0 {
                risk_factors.push(RiskFactor {
                    risk_type: RiskFactorType::AcademicPerformance,
                    severity: RiskSeverity::Critical,
                    description: "GPA below minimum threshold".to_string(),
                    early_warning_indicators: vec![
                        "Multiple failing grades".to_string(),
                        "Declining grade trend".to_string(),
                    ],
                    intervention_recommendations: vec![
                        "Academic counseling".to_string(),
                        "Tutoring support".to_string(),
                        "Study skills workshop".to_string(),
                    ],
                    monitoring_frequency: MonitoringFrequency::Weekly,
                });
            } else if gpa < 7.0 {
                risk_factors.push(RiskFactor {
                    risk_type: RiskFactorType::AcademicPerformance,
                    severity: RiskSeverity::High,
                    description: "GPA approaching critical threshold".to_string(),
                    early_warning_indicators: vec![
                        "Consistent low grades".to_string(),
                        "Struggling with core subjects".to_string(),
                    ],
                    intervention_recommendations: vec![
                        "Academic advisor meeting".to_string(),
                        "Study group participation".to_string(),
                    ],
                    monitoring_frequency: MonitoringFrequency::Biweekly,
                });
            }
        }
        
        // Attendance risk from current subjects
        let low_attendance_count = student_progress.current_subjects.iter()
            .filter(|s| s.attendance_rate.unwrap_or(100) < 70)
            .count();
        
        if low_attendance_count > 0 {
            risk_factors.push(RiskFactor {
                risk_type: RiskFactorType::Attendance,
                severity: if low_attendance_count > 2 { 
                    RiskSeverity::High 
                } else { 
                    RiskSeverity::Moderate 
                },
                description: format!("Low attendance in {} subjects", low_attendance_count),
                early_warning_indicators: vec![
                    "Missing multiple classes".to_string(),
                    "Irregular attendance pattern".to_string(),
                ],
                intervention_recommendations: vec![
                    "Attendance counseling".to_string(),
                    "Flexible scheduling options".to_string(),
                ],
                monitoring_frequency: MonitoringFrequency::Weekly,
            });
        }
        
        // Workload risk
        let current_subjects_count = student_progress.current_subjects.len();
        if current_subjects_count > 8 {
            risk_factors.push(RiskFactor {
                risk_type: RiskFactorType::Workload,
                severity: RiskSeverity::Moderate,
                description: "Heavy course load".to_string(),
                early_warning_indicators: vec![
                    "Taking excessive subjects".to_string(),
                    "Declining performance".to_string(),
                ],
                intervention_recommendations: vec![
                    "Course load assessment".to_string(),
                    "Time management coaching".to_string(),
                ],
                monitoring_frequency: MonitoringFrequency::Monthly,
            });
        }
        
        // Prerequisite chain risk
        if student_progress.failed_subjects.len() > 2 {
            risk_factors.push(RiskFactor {
                risk_type: RiskFactorType::Prerequisites,
                severity: RiskSeverity::High,
                description: "Multiple failed subjects affecting prerequisite chains".to_string(),
                early_warning_indicators: vec![
                    "Failed prerequisite subjects".to_string(),
                    "Unable to enroll in advanced subjects".to_string(),
                ],
                intervention_recommendations: vec![
                    "Academic path replanning".to_string(),
                    "Prerequisite remediation".to_string(),
                ],
                monitoring_frequency: MonitoringFrequency::Monthly,
            });
        }
        
        risk_factors
    }
    
    /// Generate graduation forecast
    pub fn generate_graduation_forecast(student_progress: &StudentProgress) -> StdResult<GraduationForecast> {
        let remaining_credits = student_progress.total_credits_required
            .saturating_sub(student_progress.total_credits_completed);
        
        // Calculate remaining semesters based on average pace
        let avg_credits_per_semester = if student_progress.current_semester > 0 {
            student_progress.total_credits_completed / student_progress.current_semester
        } else {
            20 // Default assumption
        };
        
        let remaining_semesters = if avg_credits_per_semester > 0 {
            (remaining_credits + avg_credits_per_semester - 1) / avg_credits_per_semester
        } else {
            8 // Conservative estimate
        };
        
        // Estimate graduation date (simplified - would use actual calendar in real implementation)
        let estimated_graduation_date = format!(
            "Semester {}", 
            student_progress.current_semester + remaining_semesters
        );
        
        // Calculate confidence based on performance consistency
        let confidence_level = Self::calculate_graduation_probability(student_progress)?;
        
        // Check if on track (within 10% of expected timeline)
        let expected_timeline = 8; // Standard 8-semester program
        let current_timeline = student_progress.current_semester + remaining_semesters;
        let on_track = current_timeline <= expected_timeline + 1;
        
        // Identify potential delays
        let mut potential_delays = Vec::new();
        
        if student_progress.failed_subjects.len() > 0 {
            potential_delays.push(DelayFactor {
                factor_type: DelayFactorType::AcademicPerformance,
                description: "Need to retake failed subjects".to_string(),
                impact_semesters: (student_progress.failed_subjects.len() as u32 + 1) / 2,
                mitigation_strategies: vec![
                    "Summer courses".to_string(),
                    "Intensive tutoring".to_string(),
                ],
            });
        }
        
        if let Ok(gpa) = student_progress.gpa.parse::<f64>() {
            if gpa < 7.0 {
                potential_delays.push(DelayFactor {
                    factor_type: DelayFactorType::AcademicPerformance,
                    description: "Low GPA requiring performance improvement".to_string(),
                    impact_semesters: 1,
                    mitigation_strategies: vec![
                        "Academic support services".to_string(),
                        "Reduced course load".to_string(),
                    ],
                });
            }
        }
        
        // Acceleration opportunities
        let acceleration_opportunities = if student_progress.gpa.parse::<f64>().unwrap_or(0.0) >= 8.5 {
            vec![
                "Summer session enrollment".to_string(),
                "Credit by examination".to_string(),
                "Advanced placement".to_string(),
            ]
        } else {
            vec![]
        };
        
        Ok(GraduationForecast {
            estimated_graduation_date,
            confidence_level,
            remaining_semesters,
            remaining_credits,
            on_track,
            potential_delays,
            acceleration_opportunities,
        })
    }
}