// src/tests.rs
use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
use cosmwasm_std::{coins, OwnedDeps, MemoryStorage};
use cosmwasm_std::testing::{MockApi, MockQuerier};

use crate::contract::{instantiate, execute, query};
use crate::msg::{InstantiateMsg, ExecuteMsg, QueryMsg};
use crate::state::{
    StudentProgress, SchedulePreferences, StudyIntensity, AcademicStanding, 
    SubjectScheduleInfo, ClassSchedule, DayOfWeek, TimeSlot,
    OptimizationCriteria, RecommendationAlgorithm
};
use crate::error::ContractError;

#[test]
fn test_instantiate() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    let msg = InstantiateMsg {
        owner: Some("admin".to_string()),
        ipfs_gateway: "https://ipfs.io/ipfs/".to_string(),
        max_subjects_per_semester: 6,
        recommendation_algorithm: RecommendationAlgorithm::Balanced,
    };

    let res = instantiate(deps.as_mut(), env, info, msg).unwrap();
    assert_eq!(0, res.messages.len());
    
    // Check that state was properly initialized
    let state_query = QueryMsg::GetState {};
    let state_res = query(deps.as_ref(), mock_env(), state_query).unwrap();
    assert!(!state_res.is_empty());
}

#[test] 
fn test_instantiate_with_invalid_owner() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    let msg = InstantiateMsg {
        owner: Some("invalid-address".to_string()),
        ipfs_gateway: "https://ipfs.io/ipfs/".to_string(),
        max_subjects_per_semester: 6,
        recommendation_algorithm: RecommendationAlgorithm::Balanced,
    };

    // Should use sender as owner if invalid address provided
    let res = instantiate(deps.as_mut(), env, info, msg);
    assert!(res.is_ok());
}

#[test]
fn test_register_student_progress() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    // First instantiate the contract
    let instantiate_msg = InstantiateMsg {
        owner: None,
        ipfs_gateway: "https://ipfs.io/ipfs/".to_string(),
        max_subjects_per_semester: 6,
        recommendation_algorithm: RecommendationAlgorithm::Balanced,
    };
    instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

    // Create student progress
    let student_progress = StudentProgress {
        student_id: "student1".to_string(),
        institution_id: "inst1".to_string(),
        course_id: "cs101".to_string(),
        curriculum_id: "cs2024".to_string(),
        current_semester: 1,
        completed_subjects: vec![],
        current_subjects: vec![],
        total_credits: 0,
        total_credits_required: 180,
        gpa: "0.0".to_string(),
        preferences: SchedulePreferences {
            max_subjects_per_semester: 5,
            preferred_study_intensity: StudyIntensity::Moderate,
            priority_subjects: vec![],
            avoid_subjects: vec![],
            preferred_days: vec![DayOfWeek::Monday, DayOfWeek::Wednesday],
            preferred_times: vec![TimeSlot::Morning],
            balance_theory_practice: true,
            prefer_prerequisites_early: true,
            graduation_target: Some("2028".to_string()),
        },
        academic_standing: AcademicStanding::Good,
        expected_graduation: Some("2028".to_string()),
    };

    let msg = ExecuteMsg::RegisterStudentProgress { student_progress };
    let res = execute(deps.as_mut(), env, info, msg).unwrap();
    
    assert_eq!(res.attributes.len(), 3);
    assert_eq!(res.attributes[0].key, "method");
    assert_eq!(res.attributes[0].value, "register_student_progress");
}

#[test]
fn test_register_student_with_empty_id() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    // First instantiate the contract
    let instantiate_msg = InstantiateMsg {
        owner: None,
        ipfs_gateway: "https://ipfs.io/ipfs/".to_string(),
        max_subjects_per_semester: 6,
        recommendation_algorithm: RecommendationAlgorithm::Balanced,
    };
    instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

    // Create student progress with empty ID
    let student_progress = StudentProgress {
        student_id: "".to_string(), // Empty ID should fail
        institution_id: "inst1".to_string(),
        course_id: "cs101".to_string(),
        curriculum_id: "cs2024".to_string(),
        current_semester: 1,
        completed_subjects: vec![],
        current_subjects: vec![],
        total_credits: 0,
        total_credits_required: 180,
        gpa: "0.0".to_string(),
        preferences: SchedulePreferences {
            max_subjects_per_semester: 5,
            preferred_study_intensity: StudyIntensity::Moderate,
            priority_subjects: vec![],
            avoid_subjects: vec![],
            preferred_days: vec![],
            preferred_times: vec![],
            balance_theory_practice: true,
            prefer_prerequisites_early: true,
            graduation_target: None,
        },
        academic_standing: AcademicStanding::Good,
        expected_graduation: None,
    };

    let msg = ExecuteMsg::RegisterStudentProgress { student_progress };
    let err = execute(deps.as_mut(), env, info, msg).unwrap_err();
    
    match err {
        ContractError::InvalidScheduleConfig { reason } => {
            assert_eq!(reason, "Student ID cannot be empty");
        }
        _ => panic!("Expected InvalidScheduleConfig error"),
    }
}

#[test]
fn test_register_subject_schedule_info() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    // First instantiate the contract
    let instantiate_msg = InstantiateMsg {
        owner: None,
        ipfs_gateway: "https://ipfs.io/ipfs/".to_string(),
        max_subjects_per_semester: 6,
        recommendation_algorithm: RecommendationAlgorithm::Balanced,
    };
    instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

    // Create subject schedule info
    let subject_info = SubjectScheduleInfo {
        subject_id: "CALC101".to_string(),
        title: "Calculus I".to_string(),
        credits: 4,
        prerequisites: vec![],
        corequisites: vec![],
        semester_offered: vec![1, 2], // Offered in both semesters
        max_students: Some(50),
        difficulty_level: 60,
        workload_hours: 12,
        is_elective: false,
        department: "Mathematics".to_string(),
        professor: Some("Dr. Smith".to_string()),
        schedule_info: ClassSchedule {
            days: vec![DayOfWeek::Monday, DayOfWeek::Wednesday, DayOfWeek::Friday],
            time_slots: vec![TimeSlot::Morning],
            location: Some("Room 101".to_string()),
            online_option: false,
        },
        ipfs_link: Some("QmExample123".to_string()),
    };

    let msg = ExecuteMsg::RegisterSubjectScheduleInfo { subject_info };
    let res = execute(deps.as_mut(), env, info, msg).unwrap();
    
    assert_eq!(res.attributes.len(), 2);
    assert_eq!(res.attributes[0].key, "method");
    assert_eq!(res.attributes[0].value, "register_subject_schedule_info");
    assert_eq!(res.attributes[1].key, "subject_id");
    assert_eq!(res.attributes[1].value, "CALC101");
}

#[test]
fn test_batch_register_subjects() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    // First instantiate the contract
    let instantiate_msg = InstantiateMsg {
        owner: None,
        ipfs_gateway: "https://ipfs.io/ipfs/".to_string(),
        max_subjects_per_semester: 6,
        recommendation_algorithm: RecommendationAlgorithm::Balanced,
    };
    instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

    // Create multiple subjects
    let subjects = vec![
        SubjectScheduleInfo {
            subject_id: "CALC101".to_string(),
            title: "Calculus I".to_string(),
            credits: 4,
            prerequisites: vec![],
            corequisites: vec![],
            semester_offered: vec![1, 2],
            max_students: Some(50),
            difficulty_level: 60,
            workload_hours: 12,
            is_elective: false,
            department: "Mathematics".to_string(),
            professor: Some("Dr. Smith".to_string()),
            schedule_info: ClassSchedule {
                days: vec![DayOfWeek::Monday, DayOfWeek::Wednesday],
                time_slots: vec![TimeSlot::Morning],
                location: Some("Room 101".to_string()),
                online_option: false,
            },
            ipfs_link: None,
        },
        SubjectScheduleInfo {
            subject_id: "CALC102".to_string(),
            title: "Calculus II".to_string(),
            credits: 4,
            prerequisites: vec!["CALC101".to_string()],
            corequisites: vec![],
            semester_offered: vec![2],
            max_students: Some(45),
            difficulty_level: 75,
            workload_hours: 14,
            is_elective: false,
            department: "Mathematics".to_string(),
            professor: Some("Dr. Johnson".to_string()),
            schedule_info: ClassSchedule {
                days: vec![DayOfWeek::Tuesday, DayOfWeek::Thursday],
                time_slots: vec![TimeSlot::Afternoon],
                location: Some("Room 102".to_string()),
                online_option: true,
            },
            ipfs_link: None,
        },
    ];

    let msg = ExecuteMsg::BatchRegisterSubjects { subjects };
    let res = execute(deps.as_mut(), env, info, msg).unwrap();
    
    assert_eq!(res.attributes.len(), 2);
    assert_eq!(res.attributes[0].key, "method");
    assert_eq!(res.attributes[0].value, "batch_register_subjects");
    assert_eq!(res.attributes[1].key, "registered_count");
    assert_eq!(res.attributes[1].value, "2");
}

#[test]
fn test_generate_schedule_recommendation() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    setup_test_environment(&mut deps, env.clone(), info.clone());

    // Generate schedule recommendation
    let msg = ExecuteMsg::GenerateScheduleRecommendation {
        student_id: "student1".to_string(),
        target_semester: 1,
        force_refresh: Some(true),
        custom_preferences: None,
    };

    let res = execute(deps.as_mut(), env, info, msg).unwrap();
    
    assert_eq!(res.attributes[0].key, "method");
    assert_eq!(res.attributes[0].value, "generate_schedule_recommendation");
    assert_eq!(res.attributes[1].key, "student_id");
    assert_eq!(res.attributes[1].value, "student1");
}

#[test]
fn test_create_academic_path() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    setup_test_environment(&mut deps, env.clone(), info.clone());

    // Create academic path
    let msg = ExecuteMsg::CreateAcademicPath {
        student_id: "student1".to_string(),
        path_name: "Fast Track to Graduation".to_string(),
        optimization_criteria: OptimizationCriteria::Fastest,
        target_graduation_semester: Some(8),
    };

    let res = execute(deps.as_mut(), env, info, msg).unwrap();
    
    assert_eq!(res.attributes[0].key, "method");
    assert_eq!(res.attributes[0].value, "create_academic_path");
    assert_eq!(res.attributes[1].key, "student_id");
    assert_eq!(res.attributes[1].value, "student1");
}

#[test]
fn test_complete_subject() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    setup_test_environment(&mut deps, env.clone(), info.clone());

    // Complete a subject
    let msg = ExecuteMsg::CompleteSubject {
        student_id: "student1".to_string(),
        subject_id: "CALC101".to_string(),
        grade: 85,
        completion_date: "2024-12-01".to_string(),
        difficulty_rating: Some(70),
        workload_rating: Some(60),
        nft_token_id: "nft_123".to_string(),
    };

    let res = execute(deps.as_mut(), env, info, msg).unwrap();
    
    assert_eq!(res.attributes[0].key, "method");
    assert_eq!(res.attributes[0].value, "complete_subject");
    assert_eq!(res.attributes[2].key, "subject_id");
    assert_eq!(res.attributes[2].value, "CALC101");
}

#[test]
fn test_enroll_in_subject() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    setup_test_environment(&mut deps, env.clone(), info.clone());

    // Enroll in a subject
    let msg = ExecuteMsg::EnrollInSubject {
        student_id: "student1".to_string(),
        subject_id: "CALC101".to_string(),
        enrollment_date: "2024-01-15".to_string(),
        expected_completion: "2024-05-15".to_string(),
    };

    let res = execute(deps.as_mut(), env, info, msg).unwrap();
    
    assert_eq!(res.attributes[0].key, "method");
    assert_eq!(res.attributes[0].value, "enroll_in_subject");
    assert_eq!(res.attributes[2].key, "subject_id");
    assert_eq!(res.attributes[2].value, "CALC101");
}

#[test]
fn test_query_student_progress() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    setup_test_environment(&mut deps, env.clone(), info.clone());

    // Query student progress
    let query_msg = QueryMsg::GetStudentProgress {
        student_id: "student1".to_string(),
    };

    let res = query(deps.as_ref(), env, query_msg).unwrap();
    assert!(!res.is_empty());
}

#[test]
fn test_query_subject_schedule_info() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    setup_test_environment(&mut deps, env.clone(), info.clone());

    // Query subject schedule info
    let query_msg = QueryMsg::GetSubjectScheduleInfo {
        subject_id: "CALC101".to_string(),
    };

    let res = query(deps.as_ref(), env, query_msg).unwrap();
    assert!(!res.is_empty());
}

#[test]
fn test_query_available_subjects() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    setup_test_environment(&mut deps, env.clone(), info.clone());

    // Query available subjects
    let query_msg = QueryMsg::GetAvailableSubjects {
        student_id: "student1".to_string(),
        semester: 1,
        include_electives: Some(true),
    };

    let res = query(deps.as_ref(), env, query_msg).unwrap();
    assert!(!res.is_empty());
}

#[test]
fn test_query_schedule_statistics() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    setup_test_environment(&mut deps, env.clone(), info.clone());

    // Query schedule statistics
    let query_msg = QueryMsg::GetScheduleStatistics {
        student_id: None,
        institution_id: None,
    };

    let res = query(deps.as_ref(), env, query_msg).unwrap();
    assert!(!res.is_empty());
}

#[test]
fn test_update_config_unauthorized() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    // First instantiate the contract
    let instantiate_msg = InstantiateMsg {
        owner: Some("admin".to_string()),
        ipfs_gateway: "https://ipfs.io/ipfs/".to_string(),
        max_subjects_per_semester: 6,
        recommendation_algorithm: RecommendationAlgorithm::Balanced,
    };
    instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

    // Try to update config with unauthorized user
    let unauthorized_info = mock_info("unauthorized", &coins(1000, "earth"));
    let msg = ExecuteMsg::UpdateConfig {
        ipfs_gateway: Some("https://new-gateway.io/ipfs/".to_string()),
        max_subjects_per_semester: Some(8),
        recommendation_algorithm: Some(RecommendationAlgorithm::OptimalPath),
        new_owner: None,
    };

    let err = execute(deps.as_mut(), env, unauthorized_info, msg).unwrap_err();
    
    match err {
        ContractError::Unauthorized { .. } => {
            // Expected error
        }
        _ => panic!("Expected Unauthorized error"),
    }
}

#[test]
fn test_student_not_found_error() {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info("creator", &coins(1000, "earth"));

    // First instantiate the contract
    let instantiate_msg = InstantiateMsg {
        owner: None,
        ipfs_gateway: "https://ipfs.io/ipfs/".to_string(),
        max_subjects_per_semester: 6,
        recommendation_algorithm: RecommendationAlgorithm::Balanced,
    };
    instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

    // Try to generate recommendation for non-existent student
    let msg = ExecuteMsg::GenerateScheduleRecommendation {
        student_id: "nonexistent".to_string(),
        target_semester: 1,
        force_refresh: Some(true),
        custom_preferences: None,
    };

    let err = execute(deps.as_mut(), env, info, msg).unwrap_err();
    
    match err {
        ContractError::StudentNotFound { .. } => {
            // Expected error
        }
        _ => panic!("Expected StudentNotFound error"),
    }
}

// Helper function to set up test environment with sample data
fn setup_test_environment(
    deps: &mut OwnedDeps<MemoryStorage, MockApi, MockQuerier>,
    env: cosmwasm_std::Env,
    info: cosmwasm_std::MessageInfo,
) {
    // Instantiate contract
    let instantiate_msg = InstantiateMsg {
        owner: None,
        ipfs_gateway: "https://ipfs.io/ipfs/".to_string(),
        max_subjects_per_semester: 6,
        recommendation_algorithm: RecommendationAlgorithm::Balanced,
    };
    instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg).unwrap();

    // Register a student
    let student_progress = StudentProgress {
        student_id: "student1".to_string(),
        institution_id: "inst1".to_string(),
        course_id: "cs101".to_string(),
        curriculum_id: "cs2024".to_string(),
        current_semester: 1,
        completed_subjects: vec![],
        current_subjects: vec![],
        total_credits: 0,
        total_credits_required: 180,
        gpa: "0.0".to_string(),
        preferences: SchedulePreferences {
            max_subjects_per_semester: 5,
            preferred_study_intensity: StudyIntensity::Moderate,
            priority_subjects: vec!["CALC101".to_string()],
            avoid_subjects: vec![],
            preferred_days: vec![DayOfWeek::Monday, DayOfWeek::Wednesday],
            preferred_times: vec![TimeSlot::Morning],
            balance_theory_practice: true,
            prefer_prerequisites_early: true,
            graduation_target: Some("2028".to_string()),
        },
        academic_standing: AcademicStanding::Good,
        expected_graduation: Some("2028".to_string()),
    };

    let msg = ExecuteMsg::RegisterStudentProgress { student_progress };
    execute(deps.as_mut(), env.clone(), info.clone(), msg).unwrap();

    // Register some subjects
    let subjects = vec![
        SubjectScheduleInfo {
            subject_id: "CALC101".to_string(),
            title: "Calculus I".to_string(),
            credits: 4,
            prerequisites: vec![],
            corequisites: vec![],
            semester_offered: vec![1, 2],
            max_students: Some(50),
            difficulty_level: 60,
            workload_hours: 12,
            is_elective: false,
            department: "Mathematics".to_string(),
            professor: Some("Dr. Smith".to_string()),
            schedule_info: ClassSchedule {
                days: vec![DayOfWeek::Monday, DayOfWeek::Wednesday],
                time_slots: vec![TimeSlot::Morning],
                location: Some("Room 101".to_string()),
                online_option: false,
            },
            ipfs_link: None,
        },
        SubjectScheduleInfo {
            subject_id: "PROG101".to_string(),
            title: "Programming I".to_string(),
            credits: 4,
            prerequisites: vec![],
            corequisites: vec![],
            semester_offered: vec![1, 2],
            max_students: Some(40),
            difficulty_level: 65,
            workload_hours: 15,
            is_elective: false,
            department: "Computer Science".to_string(),
            professor: Some("Dr. Brown".to_string()),
            schedule_info: ClassSchedule {
                days: vec![DayOfWeek::Tuesday, DayOfWeek::Thursday],
                time_slots: vec![TimeSlot::Afternoon],
                location: Some("Lab 201".to_string()),
                online_option: true,
            },
            ipfs_link: None,
        },
    ];

    let msg = ExecuteMsg::BatchRegisterSubjects { subjects };
    execute(deps.as_mut(), env, info, msg).unwrap();
}

#[cfg(test)]
mod integration_tests {
    use super::*;

    #[test]
    fn test_full_student_lifecycle() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let info = mock_info("creator", &coins(1000, "earth"));

        setup_test_environment(&mut deps, env.clone(), info.clone());

        // 1. Generate schedule recommendation
        let recommendation_msg = ExecuteMsg::GenerateScheduleRecommendation {
            student_id: "student1".to_string(),
            target_semester: 1,
            force_refresh: Some(true),
            custom_preferences: None,
        };
        let res = execute(deps.as_mut(), env.clone(), info.clone(), recommendation_msg).unwrap();
        assert_eq!(res.attributes[0].value, "generate_schedule_recommendation");

        // 2. Enroll in recommended subject
        let enroll_msg = ExecuteMsg::EnrollInSubject {
            student_id: "student1".to_string(),
            subject_id: "CALC101".to_string(),
            enrollment_date: "2024-01-15".to_string(),
            expected_completion: "2024-05-15".to_string(),
        };
        let res = execute(deps.as_mut(), env.clone(), info.clone(), enroll_msg).unwrap();
        assert_eq!(res.attributes[0].value, "enroll_in_subject");

        // 3. Complete the subject
        let complete_msg = ExecuteMsg::CompleteSubject {
            student_id: "student1".to_string(),
            subject_id: "CALC101".to_string(),
            grade: 85,
            completion_date: "2024-05-15".to_string(),
            difficulty_rating: Some(70),
            workload_rating: Some(60),
            nft_token_id: "nft_123".to_string(),
        };
        let res = execute(deps.as_mut(), env.clone(), info.clone(), complete_msg).unwrap();
        assert_eq!(res.attributes[0].value, "complete_subject");

        // 4. Create academic path
        let path_msg = ExecuteMsg::CreateAcademicPath {
            student_id: "student1".to_string(),
            path_name: "Standard Path".to_string(),
            optimization_criteria: OptimizationCriteria::Balanced,
            target_graduation_semester: Some(8),
        };
        let res = execute(deps.as_mut(), env.clone(), info.clone(), path_msg).unwrap();
        assert_eq!(res.attributes[0].value, "create_academic_path");

        // 5. Query student progress to verify updates
        let query_msg = QueryMsg::GetStudentProgress {
            student_id: "student1".to_string(),
        };
        let res = query(deps.as_ref(), env, query_msg).unwrap();
        assert!(!res.is_empty());
    }

    #[test]
    fn test_recommendation_algorithm_variations() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let info = mock_info("creator", &coins(1000, "earth"));

        // Test different recommendation algorithms
        let algorithms = vec![
            RecommendationAlgorithm::Basic,
            RecommendationAlgorithm::Balanced,
            RecommendationAlgorithm::OptimalPath,
            RecommendationAlgorithm::MachineLearning,
        ];

        for algorithm in algorithms {
            let instantiate_msg = InstantiateMsg {
                owner: None,
                ipfs_gateway: "https://ipfs.io/ipfs/".to_string(),
                max_subjects_per_semester: 6,
                recommendation_algorithm: algorithm,
            };

            let res = instantiate(deps.as_mut(), env.clone(), info.clone(), instantiate_msg);
            assert!(res.is_ok());

            // Query config to verify algorithm was set
            let query_msg = QueryMsg::GetConfig {};
            let res = query(deps.as_ref(), env.clone(), query_msg).unwrap();
            assert!(!res.is_empty());
        }
    }

    #[test]
    fn test_optimization_criteria_variations() {
        let mut deps = mock_dependencies();
        let env = mock_env();
        let info = mock_info("creator", &coins(1000, "earth"));

        setup_test_environment(&mut deps, env.clone(), info.clone());

        // Test different optimization criteria
        let criteria_list = vec![
            OptimizationCriteria::Fastest,
            OptimizationCriteria::Balanced,
            OptimizationCriteria::EasiestFirst,
            OptimizationCriteria::PrerequisiteOptimal,
        ];

        for (i, criteria) in criteria_list.iter().enumerate() {
            let path_msg = ExecuteMsg::CreateAcademicPath {
                student_id: "student1".to_string(),
                path_name: format!("Path {}", i + 1),
                optimization_criteria: criteria.clone(),
                target_graduation_semester: Some(8),
            };

            let res = execute(deps.as_mut(), env.clone(), info.clone(), path_msg);
            assert!(res.is_ok());
        }
    }
}