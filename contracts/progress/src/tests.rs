#[cfg(test)]
mod tests {
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::from_json;
    use crate::msg::{InstantiateMsg, ExecuteMsg, QueryMsg};
    use crate::state::{UpdateFrequency, AnalyticsDepth, AcademicStatus, StudentProgress, PerformanceMetrics, GradeTrend, GraduationForecast};
    use crate::contract::{instantiate, execute, query};

    // Helper function to create a default student progress
    fn create_test_student_progress() -> StudentProgress {
        StudentProgress {
            student_id: "student1".to_string(),
            institution_id: "institution1".to_string(),
            course_id: "cs101".to_string(),
            curriculum_id: "cs_curriculum_2024".to_string(),
            current_semester: 3,
            total_semesters_enrolled: 3,
            academic_status: AcademicStatus::Active,
            total_credits_completed: 45,
            total_credits_required: 120,
            credits_in_progress: 15,
            gpa: "8.5".to_string(),
            cgpa: "8.2".to_string(),
            grade_distribution: crate::state::GradeDistribution {
                excellent_count: 8,
                good_count: 5,
                satisfactory_count: 2,
                poor_count: 0,
                fail_count: 0,
            },
            completed_subjects: vec![],
            current_subjects: vec![],
            failed_subjects: vec![],
            milestones_achieved: vec![],
            upcoming_milestones: vec![],
            performance_metrics: PerformanceMetrics {
                success_rate: 100,
                retry_rate: 0,
                grade_trend: GradeTrend::StronglyIncreasing,
                consistency_score: 85,
                attendance_average: 95,
                assignment_completion_rate: 98,
                study_hours_per_credit: 3,
                credits_per_semester_avg: 15,
                time_to_completion_ratio: 100,
                prerequisite_success_rate: 100,
                percentile_rank_gpa: 85,
                percentile_rank_progress: 90,
                graduation_probability: 95,
                time_extension_probability: 5,
            },
            graduation_forecast: GraduationForecast {
                estimated_graduation_date: "2026-05-15".to_string(),
                confidence_level: 85,
                remaining_semesters: 5,
                remaining_credits: 75,
                on_track: true,
                potential_delays: vec![],
                acceleration_opportunities: vec!["Summer courses".to_string()],
            },
            risk_factors: vec![],
            improvement_recommendations: vec![],
            enrollment_date: "2023-09-01".to_string(),
            last_updated: 1672531200, // Jan 1, 2023
            next_milestone_date: Some("2024-12-15".to_string()),
        }
    }

    #[test]
    fn test_instantiate() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let info = mock_info("creator", &[]);

        let msg = InstantiateMsg {
            owner: Some("admin".to_string()),
            analytics_enabled: true,
            update_frequency: UpdateFrequency::Daily,
            analytics_depth: AnalyticsDepth::Standard,
        };

        let res = instantiate(deps.as_mut(), env, info, msg).unwrap();
        
        assert_eq!(res.attributes.len(), 3);
        assert_eq!(res.attributes[0].key, "method");
        assert_eq!(res.attributes[0].value, "instantiate");
        assert_eq!(res.attributes[1].key, "owner");
        assert_eq!(res.attributes[1].value, "admin");
        assert_eq!(res.attributes[2].key, "analytics_enabled");
        assert_eq!(res.attributes[2].value, "true");
    }

    #[test]
    fn test_query_state() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let info = mock_info("creator", &[]);

        // First instantiate
        let msg = InstantiateMsg {
            owner: None,
            analytics_enabled: true,
            update_frequency: UpdateFrequency::Daily,
            analytics_depth: AnalyticsDepth::Standard,
        };

        instantiate(deps.as_mut(), env.clone(), info, msg).unwrap();

        // Query state
        let query_msg = QueryMsg::GetState {};
        let res = query(deps.as_ref(), env, query_msg).unwrap();
        let state_response: crate::msg::StateResponse = from_json(&res).unwrap();

        assert_eq!(state_response.state.total_students, 0);
        assert_eq!(state_response.state.total_institutions, 0);
        assert_eq!(state_response.state.analytics_enabled, true);
    }

    #[test]
    fn test_update_student_progress() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let info = mock_info("admin", &[]);

        // Instantiate first
        let instantiate_msg = InstantiateMsg {
            owner: None,
            analytics_enabled: true,
            update_frequency: UpdateFrequency::Daily,
            analytics_depth: AnalyticsDepth::Standard,
        };

        instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

        // Update student progress
        let student_progress = create_test_student_progress();
        let update_msg = ExecuteMsg::UpdateStudentProgress {
            student_progress: student_progress.clone(),
            force_analytics_refresh: Some(false),
        };

        let res = execute(deps.as_mut(), env.clone(), info, update_msg).unwrap();

        assert_eq!(res.attributes.len(), 5);
        assert_eq!(res.attributes[0].key, "method");
        assert_eq!(res.attributes[0].value, "update_student_progress");
        assert_eq!(res.attributes[1].key, "student_id");
        assert_eq!(res.attributes[1].value, "student1");
        assert_eq!(res.attributes[4].key, "timestamp");
        // Check if student was marked as new
        let is_new_student_attr = res.attributes.iter().find(|attr| attr.key == "is_new_student").unwrap();
        assert_eq!(is_new_student_attr.value, "true");

        // Query the student progress
        let query_msg = QueryMsg::GetStudentProgress {
            student_id: "student1".to_string(),
        };
        let res = query(deps.as_ref(), env, query_msg).unwrap();
        let progress_response: crate::msg::StudentProgressResponse = from_json(&res).unwrap();

        assert_eq!(progress_response.progress.student_id, "student1");
        assert_eq!(progress_response.progress.gpa, "8.5");
        assert_eq!(progress_response.progress.total_credits_completed, 45);
    }

    #[test]
    fn test_batch_update_student_progress() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let info = mock_info("admin", &[]);

        // Instantiate first
        let instantiate_msg = InstantiateMsg {
            owner: None,
            analytics_enabled: true,
            update_frequency: UpdateFrequency::Daily,
            analytics_depth: AnalyticsDepth::Standard,
        };

        instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

        // Create multiple students
        let mut student1 = create_test_student_progress();
        student1.student_id = "student1".to_string();
        
        let mut student2 = create_test_student_progress();
        student2.student_id = "student2".to_string();
        student2.gpa = "7.8".to_string();

        let batch_msg = ExecuteMsg::BatchUpdateStudentProgress {
            student_progress_list: vec![student1, student2],
            analytics_refresh_mode: crate::msg::AnalyticsRefreshMode::None,
        };

        let res = execute(deps.as_mut(), env, info, batch_msg).unwrap();

        assert_eq!(res.attributes[0].key, "method");
        assert_eq!(res.attributes[0].value, "batch_update_student_progress");
        assert_eq!(res.attributes[1].key, "updated_count");
        assert_eq!(res.attributes[1].value, "2");
    }

    #[test]
    fn test_record_subject_completion() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let info = mock_info("admin", &[]);

        // Instantiate and add a student first
        let instantiate_msg = InstantiateMsg {
            owner: None,
            analytics_enabled: true,
            update_frequency: UpdateFrequency::Daily,
            analytics_depth: AnalyticsDepth::Standard,
        };

        instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

        // Add student
        let student_progress = create_test_student_progress();
        let update_msg = ExecuteMsg::UpdateStudentProgress {
            student_progress,
            force_analytics_refresh: Some(false),
        };
        execute(deps.as_mut(), env.clone(), info.clone(), update_msg).unwrap();

        // Record subject completion
        let completion_msg = ExecuteMsg::RecordSubjectCompletion {
            student_id: "student1".to_string(),
            subject_id: "math101".to_string(),
            final_grade: 85,
            completion_date: "2024-05-15".to_string(),
            study_hours: Some(45),
            difficulty_rating: Some(7),
            satisfaction_rating: Some(8),
            nft_token_id: "nft_math101_student1".to_string(),
        };

        let res = execute(deps.as_mut(), env, info, completion_msg).unwrap();

        assert_eq!(res.attributes[0].key, "method");
        assert_eq!(res.attributes[0].value, "record_subject_completion");
        assert_eq!(res.attributes[3].key, "final_grade");
        assert_eq!(res.attributes[3].value, "85");
        assert_eq!(res.attributes[4].key, "passed");
        assert_eq!(res.attributes[4].value, "true");
    }

    #[test]
    fn test_update_config() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let admin_info = mock_info("admin", &[]);
        let user_info = mock_info("user", &[]);

        // Instantiate with admin
        let instantiate_msg = InstantiateMsg {
            owner: Some("admin".to_string()),
            analytics_enabled: true,
            update_frequency: UpdateFrequency::Daily,
            analytics_depth: AnalyticsDepth::Standard,
        };

        instantiate(deps.as_mut(), env.clone(), admin_info.clone(), instantiate_msg).unwrap();

        // Try to update config as non-admin (should fail)
        let update_msg = ExecuteMsg::UpdateConfig {
            new_owner: None,
            analytics_enabled: Some(false),
            update_frequency: None,
            analytics_depth: None,
            retention_period_days: None,
        };

        let res = execute(deps.as_mut(), env.clone(), user_info, update_msg.clone());
        assert!(res.is_err());

        // Update config as admin (should succeed)
        let res = execute(deps.as_mut(), env, admin_info, update_msg).unwrap();

        assert_eq!(res.attributes[0].key, "method");
        assert_eq!(res.attributes[0].value, "update_config");
    }

    #[test]
    fn test_generate_student_dashboard() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let info = mock_info("admin", &[]);

        // Setup
        let instantiate_msg = InstantiateMsg {
            owner: None,
            analytics_enabled: true,
            update_frequency: UpdateFrequency::Daily,
            analytics_depth: AnalyticsDepth::Standard,
        };

        instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

        // Add student
        let student_progress = create_test_student_progress();
        let update_msg = ExecuteMsg::UpdateStudentProgress {
            student_progress,
            force_analytics_refresh: Some(false),
        };
        execute(deps.as_mut(), env.clone(), info.clone(), update_msg).unwrap();

        // Generate dashboard
        let dashboard_msg = ExecuteMsg::GenerateStudentDashboard {
            student_id: "student1".to_string(),
            include_peer_comparison: Some(false),
            force_refresh: Some(true),
        };

        let res = execute(deps.as_mut(), env, info, dashboard_msg).unwrap();

        assert_eq!(res.attributes[0].key, "method");
        assert_eq!(res.attributes[0].value, "generate_student_dashboard");
        assert_eq!(res.attributes[1].key, "student_id");
        assert_eq!(res.attributes[1].value, "student1");
    }

    #[test]
    fn test_academic_status_update() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let info = mock_info("admin", &[]);

        // Setup
        let instantiate_msg = InstantiateMsg {
            owner: None,
            analytics_enabled: true,
            update_frequency: UpdateFrequency::Daily,
            analytics_depth: AnalyticsDepth::Standard,
        };

        instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

        // Add student
        let student_progress = create_test_student_progress();
        let update_msg = ExecuteMsg::UpdateStudentProgress {
            student_progress,
            force_analytics_refresh: Some(false),
        };
        execute(deps.as_mut(), env.clone(), info.clone(), update_msg).unwrap();

        // Update academic status
        let status_msg = ExecuteMsg::UpdateAcademicStatus {
            student_id: "student1".to_string(),
            new_status: AcademicStatus::Graduated,
            effective_date: "2024-05-15".to_string(),
            reason: Some("Completed all requirements".to_string()),
        };

        let res = execute(deps.as_mut(), env, info, status_msg).unwrap();

        assert_eq!(res.attributes[0].key, "method");
        assert_eq!(res.attributes[0].value, "update_academic_status");
        assert_eq!(res.attributes[3].key, "new_status");
        assert_eq!(res.attributes[3].value, "Graduated");
    }

    #[test]
    fn test_students_by_institution_query() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let info = mock_info("admin", &[]);

        // Setup
        let instantiate_msg = InstantiateMsg {
            owner: None,
            analytics_enabled: true,
            update_frequency: UpdateFrequency::Daily,
            analytics_depth: AnalyticsDepth::Standard,
        };

        instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

        // Add multiple students to same institution
        for i in 1..=3 {
            let mut student = create_test_student_progress();
            student.student_id = format!("student{}", i);
            
            let update_msg = ExecuteMsg::UpdateStudentProgress {
                student_progress: student,
                force_analytics_refresh: Some(false),
            };
            execute(deps.as_mut(), env.clone(), info.clone(), update_msg).unwrap();
        }

        // Query students by institution
        let query_msg = QueryMsg::GetStudentsByInstitution {
            institution_id: "institution1".to_string(),
            status_filter: None,
            limit: None,
            start_after: None,
        };

        let res = query(deps.as_ref(), env, query_msg).unwrap();
        let students_response: crate::msg::StudentsListResponse = from_json(&res).unwrap();

        assert_eq!(students_response.students.len(), 3);
        assert_eq!(students_response.total_count, 3);
        assert!(!students_response.has_more);
    }

    #[test]
    fn test_risk_factor_management() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let info = mock_info("admin", &[]);

        // Setup
        let instantiate_msg = InstantiateMsg {
            owner: None,
            analytics_enabled: true,
            update_frequency: UpdateFrequency::Daily,
            analytics_depth: AnalyticsDepth::Standard,
        };

        instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

        // Add student
        let student_progress = create_test_student_progress();
        let update_msg = ExecuteMsg::UpdateStudentProgress {
            student_progress,
            force_analytics_refresh: Some(false),
        };
        execute(deps.as_mut(), env.clone(), info.clone(), update_msg).unwrap();

        // Add risk factor
        let risk_msg = ExecuteMsg::AddRiskFactor {
            student_id: "student1".to_string(),
            risk_type: crate::state::RiskFactorType::AcademicPerformance,
            severity: crate::state::RiskSeverity::Moderate,
            description: "Declining grades in math courses".to_string(),
            intervention_recommendations: vec![
                "Additional tutoring".to_string(),
                "Study group participation".to_string(),
            ],
        };

        let res = execute(deps.as_mut(), env.clone(), info.clone(), risk_msg).unwrap();
        assert_eq!(res.attributes[0].value, "add_risk_factor");

        // Remove risk factor
        let remove_msg = ExecuteMsg::RemoveRiskFactor {
            student_id: "student1".to_string(),
            risk_type: crate::state::RiskFactorType::AcademicPerformance,
        };

        let res = execute(deps.as_mut(), env, info, remove_msg).unwrap();
        assert_eq!(res.attributes[0].value, "remove_risk_factor");
    }

    #[test]
    fn test_config_query() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let info = mock_info("creator", &[]);

        // Instantiate
        let msg = InstantiateMsg {
            owner: None,
            analytics_enabled: true,
            update_frequency: UpdateFrequency::Weekly,
            analytics_depth: AnalyticsDepth::Advanced,
        };

        instantiate(deps.as_mut(), env.clone(), info, msg).unwrap();

        // Query config
        let query_msg = QueryMsg::GetConfig {};
        let res = query(deps.as_ref(), env, query_msg).unwrap();
        let config_response: crate::msg::ProgressConfigResponse = from_json(&res).unwrap();

        assert!(matches!(config_response.config.update_frequency, UpdateFrequency::Weekly));
        assert!(matches!(config_response.config.analytics_depth, AnalyticsDepth::Advanced));
        assert_eq!(config_response.config.retention_period_days, 365);
    }
}
