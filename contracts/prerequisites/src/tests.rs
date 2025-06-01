use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
use cosmwasm_std::{coins, from_json};
use crate::contract::{execute, instantiate, query};
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg, CompletedSubjectMsg};
use crate::state::{PrerequisiteGroup, GroupType, LogicType};

const OWNER: &str = "owner";
const STUDENT: &str = "student123";

fn proper_instantiate() -> (cosmwasm_std::OwnedDeps<cosmwasm_std::MemoryStorage, cosmwasm_std::testing::MockApi, cosmwasm_std::testing::MockQuerier>, cosmwasm_std::Env, cosmwasm_std::MessageInfo) {
    let mut deps = mock_dependencies();
    let env = mock_env();
    let info = mock_info(OWNER, &coins(1000, "earth"));

    let msg = InstantiateMsg {
        owner: Some(OWNER.to_string()),
    };

    let res = instantiate(deps.as_mut(), env.clone(), info.clone(), msg).unwrap();
    assert_eq!(0, res.messages.len());

    (deps, env, info)
}

#[test]
fn test_instantiate() {
    let (deps, _env, _info) = proper_instantiate();

    // Query the state
    let res = query(deps.as_ref(), mock_env(), QueryMsg::GetState {}).unwrap();
    let state_response: crate::msg::StateResponse = from_json(&res).unwrap();
    
    assert_eq!(state_response.owner, OWNER);
    assert_eq!(state_response.total_subjects, 0);
    assert_eq!(state_response.total_verifications, 0);
}

#[test]
fn test_register_prerequisites() {
    let (mut deps, env, _info) = proper_instantiate();

    // Create prerequisite groups
    let prereq_group = PrerequisiteGroup {
        id: "calc1_prereq".to_string(),
        subject_id: "CALC2".to_string(),
        group_type: GroupType::All,
        minimum_credits: 0,
        minimum_completed_subjects: 1,
        subject_ids: vec!["CALC1".to_string()],
        logic: LogicType::And,
        priority: 1,
        confidence: 95,
    };

    let msg = ExecuteMsg::RegisterPrerequisites {
        subject_id: "CALC2".to_string(),
        prerequisites: vec![prereq_group],
    };

    let info = mock_info(OWNER, &[]);
    let res = execute(deps.as_mut(), env, info, msg).unwrap();

    assert_eq!(res.attributes[0].value, "register_prerequisites");
    assert_eq!(res.attributes[1].value, "CALC2");
    assert_eq!(res.attributes[2].value, "1");
}

#[test]
fn test_update_student_record() {
    let (mut deps, env, _info) = proper_instantiate();

    let completed_subject = CompletedSubjectMsg {
        subject_id: "CALC1".to_string(),
        credits: 4,
        completion_date: "2024-01-15".to_string(),
        grade: 8550, // 85.50 represented as 8550
        nft_token_id: "nft_calc1_123".to_string(),
    };

    let msg = ExecuteMsg::UpdateStudentRecord {
        student_id: STUDENT.to_string(),
        completed_subject,
    };

    let info = mock_info("institution", &[]);
    let res = execute(deps.as_mut(), env, info, msg).unwrap();

    assert_eq!(res.attributes[0].value, "update_student_record");
    assert_eq!(res.attributes[1].value, STUDENT);
    assert_eq!(res.attributes[2].value, "CALC1");
    assert_eq!(res.attributes[3].value, "4");
}

#[test]
fn test_verify_enrollment_no_prerequisites() {
    let (mut deps, env, _info) = proper_instantiate();

    // Try to enroll in a subject with no prerequisites
    let msg = ExecuteMsg::VerifyEnrollment {
        student_id: STUDENT.to_string(),
        subject_id: "INTRO_MATH".to_string(),
    };

    let info = mock_info("anyone", &[]);
    let res = execute(deps.as_mut(), env, info, msg).unwrap();

    assert_eq!(res.attributes[0].value, "verify_enrollment");
    assert_eq!(res.attributes[3].value, "true"); // can_enroll should be true
}

#[test]
fn test_verify_enrollment_with_prerequisites_success() {
    let (mut deps, env, _info) = proper_instantiate();

    // 1. Register prerequisites for CALC2 (requires CALC1)
    let prereq_group = PrerequisiteGroup {
        id: "calc1_prereq".to_string(),
        subject_id: "CALC2".to_string(),
        group_type: GroupType::All,
        minimum_credits: 0,
        minimum_completed_subjects: 1,
        subject_ids: vec!["CALC1".to_string()],
        logic: LogicType::And,
        priority: 1,
        confidence: 95,
    };

    let msg = ExecuteMsg::RegisterPrerequisites {
        subject_id: "CALC2".to_string(),
        prerequisites: vec![prereq_group],
    };

    let info = mock_info(OWNER, &[]);
    execute(deps.as_mut(), env.clone(), info, msg).unwrap();

    // 2. Add CALC1 completion for student
    let completed_subject = CompletedSubjectMsg {
        subject_id: "CALC1".to_string(),
        credits: 4,
        completion_date: "2024-01-15".to_string(),
        grade: 8550, // 85.50 represented as 8550
        nft_token_id: "nft_calc1_123".to_string(),
    };

    let msg = ExecuteMsg::UpdateStudentRecord {
        student_id: STUDENT.to_string(),
        completed_subject,
    };

    let info = mock_info("institution", &[]);
    execute(deps.as_mut(), env.clone(), info, msg).unwrap();

    // 3. Try to enroll in CALC2 (should succeed)
    let msg = ExecuteMsg::VerifyEnrollment {
        student_id: STUDENT.to_string(),
        subject_id: "CALC2".to_string(),
    };

    let info = mock_info("anyone", &[]);
    let res = execute(deps.as_mut(), env, info, msg).unwrap();

    assert_eq!(res.attributes[0].value, "verify_enrollment");
    assert_eq!(res.attributes[3].value, "true"); // can_enroll should be true
}

#[test]
fn test_verify_enrollment_with_prerequisites_failure() {
    let (mut deps, env, _info) = proper_instantiate();

    // 1. Register prerequisites for CALC2 (requires CALC1)
    let prereq_group = PrerequisiteGroup {
        id: "calc1_prereq".to_string(),
        subject_id: "CALC2".to_string(),
        group_type: GroupType::All,
        minimum_credits: 0,
        minimum_completed_subjects: 1,
        subject_ids: vec!["CALC1".to_string()],
        logic: LogicType::And,
        priority: 1,
        confidence: 95,
    };

    let msg = ExecuteMsg::RegisterPrerequisites {
        subject_id: "CALC2".to_string(),
        prerequisites: vec![prereq_group],
    };

    let info = mock_info(OWNER, &[]);
    execute(deps.as_mut(), env.clone(), info, msg).unwrap();

    // 2. Try to enroll in CALC2 WITHOUT completing CALC1 (should fail)
    let msg = ExecuteMsg::VerifyEnrollment {
        student_id: STUDENT.to_string(),
        subject_id: "CALC2".to_string(),
    };

    let info = mock_info("anyone", &[]);
    let res = execute(deps.as_mut(), env, info, msg).unwrap();

    assert_eq!(res.attributes[0].value, "verify_enrollment");
    assert_eq!(res.attributes[3].value, "false"); // can_enroll should be false
}

#[test]
fn test_unauthorized_register_prerequisites() {
    let (mut deps, env, _info) = proper_instantiate();

    let prereq_group = PrerequisiteGroup {
        id: "test".to_string(),
        subject_id: "TEST".to_string(),
        group_type: GroupType::None,
        minimum_credits: 0,
        minimum_completed_subjects: 0,
        subject_ids: vec![],
        logic: LogicType::None,
        priority: 1,
        confidence: 100,
    };

    let msg = ExecuteMsg::RegisterPrerequisites {
        subject_id: "TEST".to_string(),
        prerequisites: vec![prereq_group],
    };

    // Try with non-owner
    let info = mock_info("not_owner", &[]);
    let err = execute(deps.as_mut(), env, info, msg).unwrap_err();
    
    match err {
        crate::error::ContractError::Unauthorized {} => {}
        _ => panic!("Expected Unauthorized error"),
    }
}

// ============ NOVOS TESTES MAIS RIGOROSOS ============

#[test]
fn test_complex_prerequisites_or_logic() {
    let (mut deps, env, _info) = proper_instantiate();

    // Subject requires EITHER MATH1 OR MATH2
    let prereq_group = PrerequisiteGroup {
        id: "math_prereq".to_string(),
        subject_id: "ADVANCED_MATH".to_string(),
        group_type: GroupType::Any, // ANY means OR logic
        minimum_credits: 0,
        minimum_completed_subjects: 1,
        subject_ids: vec!["MATH1".to_string(), "MATH2".to_string()],
        logic: LogicType::Or,
        priority: 1,
        confidence: 95,
    };

    let msg = ExecuteMsg::RegisterPrerequisites {
        subject_id: "ADVANCED_MATH".to_string(),
        prerequisites: vec![prereq_group],
    };

    let info = mock_info(OWNER, &[]);
    execute(deps.as_mut(), env.clone(), info, msg).unwrap();

    // Student completed only MATH2 (not MATH1)
    let completed_subject = CompletedSubjectMsg {
        subject_id: "MATH2".to_string(),
        credits: 4,
        completion_date: "2024-01-15".to_string(),
        grade: 8000,
        nft_token_id: "nft_math2_123".to_string(),
    };

    let msg = ExecuteMsg::UpdateStudentRecord {
        student_id: STUDENT.to_string(),
        completed_subject,
    };

    let info = mock_info("institution", &[]);
    execute(deps.as_mut(), env.clone(), info, msg).unwrap();

    // Should be able to enroll because MATH2 satisfies the OR requirement
    let msg = ExecuteMsg::VerifyEnrollment {
        student_id: STUDENT.to_string(),
        subject_id: "ADVANCED_MATH".to_string(),
    };

    let info = mock_info("anyone", &[]);
    let res = execute(deps.as_mut(), env, info, msg).unwrap();

    assert_eq!(res.attributes[3].value, "true");
}

#[test]
fn test_minimum_credits_requirement() {
    let (mut deps, env, _info) = proper_instantiate();

    // Subject requires minimum 8 credits from elective subjects
    let prereq_group = PrerequisiteGroup {
        id: "elective_credits".to_string(),
        subject_id: "CAPSTONE".to_string(),
        group_type: GroupType::Minimum,
        minimum_credits: 8, // Requires at least 8 credits
        minimum_completed_subjects: 2, // And at least 2 subjects
        subject_ids: vec!["ELECTIVE1".to_string(), "ELECTIVE2".to_string(), "ELECTIVE3".to_string()],
        logic: LogicType::And,
        priority: 1,
        confidence: 95,
    };

    let msg = ExecuteMsg::RegisterPrerequisites {
        subject_id: "CAPSTONE".to_string(),
        prerequisites: vec![prereq_group],
    };

    let info = mock_info(OWNER, &[]);
    execute(deps.as_mut(), env.clone(), info, msg).unwrap();

    // Add first elective (4 credits)
    let completed_subject1 = CompletedSubjectMsg {
        subject_id: "ELECTIVE1".to_string(),
        credits: 4,
        completion_date: "2024-01-15".to_string(),
        grade: 8000,
        nft_token_id: "nft_elec1_123".to_string(),
    };

    let msg = ExecuteMsg::UpdateStudentRecord {
        student_id: STUDENT.to_string(),
        completed_subject: completed_subject1,
    };

    let info = mock_info("institution", &[]);
    execute(deps.as_mut(), env.clone(), info, msg).unwrap();

    // Try to enroll - should fail (only 4 credits, need 8)
    let msg = ExecuteMsg::VerifyEnrollment {
        student_id: STUDENT.to_string(),
        subject_id: "CAPSTONE".to_string(),
    };

    let info = mock_info("anyone", &[]);
    let res = execute(deps.as_mut(), env.clone(), info, msg).unwrap();
    assert_eq!(res.attributes[3].value, "false");

    // Add second elective (4 credits) - now should have 8 total
    let completed_subject2 = CompletedSubjectMsg {
        subject_id: "ELECTIVE2".to_string(),
        credits: 4,
        completion_date: "2024-01-16".to_string(),
        grade: 8500,
        nft_token_id: "nft_elec2_123".to_string(),
    };

    let msg = ExecuteMsg::UpdateStudentRecord {
        student_id: STUDENT.to_string(),
        completed_subject: completed_subject2,
    };

    let info = mock_info("institution", &[]);
    execute(deps.as_mut(), env.clone(), info, msg).unwrap();

    // Now should succeed (8 credits, 2 subjects)
    let msg = ExecuteMsg::VerifyEnrollment {
        student_id: STUDENT.to_string(),
        subject_id: "CAPSTONE".to_string(),
    };

    let info = mock_info("anyone", &[]);
    let res = execute(deps.as_mut(), env, info, msg).unwrap();
    assert_eq!(res.attributes[3].value, "true");
}

#[test]
fn test_duplicate_subject_completion_error() {
    let (mut deps, env, _info) = proper_instantiate();

    let completed_subject = CompletedSubjectMsg {
        subject_id: "CALC1".to_string(),
        credits: 4,
        completion_date: "2024-01-15".to_string(),
        grade: 8550,
        nft_token_id: "nft_calc1_123".to_string(),
    };

    // Add subject first time - should succeed
    let msg = ExecuteMsg::UpdateStudentRecord {
        student_id: STUDENT.to_string(),
        completed_subject: completed_subject.clone(),
    };

    let info = mock_info("institution", &[]);
    let res = execute(deps.as_mut(), env.clone(), info, msg);
    assert!(res.is_ok());

    // Try to add same subject again - should fail
    let msg = ExecuteMsg::UpdateStudentRecord {
        student_id: STUDENT.to_string(),
        completed_subject,
    };

    let info = mock_info("institution", &[]);
    let err = execute(deps.as_mut(), env, info, msg).unwrap_err();
    
    match err {
        crate::error::ContractError::SubjectAlreadyCompleted { .. } => {}
        _ => panic!("Expected SubjectAlreadyCompleted error"),
    }
}

#[test]
fn test_empty_string_ids() {
    let (mut deps, env, _info) = proper_instantiate();

    // Try to register prerequisites with empty subject_id
    let prereq_group = PrerequisiteGroup {
        id: "test".to_string(),
        subject_id: "".to_string(), // Empty string
        group_type: GroupType::None,
        minimum_credits: 0,
        minimum_completed_subjects: 0,
        subject_ids: vec![],
        logic: LogicType::None,
        priority: 1,
        confidence: 100,
    };

    let msg = ExecuteMsg::RegisterPrerequisites {
        subject_id: "VALID_SUBJECT".to_string(),
        prerequisites: vec![prereq_group],
    };

    let info = mock_info(OWNER, &[]);
    let err = execute(deps.as_mut(), env, info, msg).unwrap_err();
    
    match err {
        crate::error::ContractError::InvalidPrerequisite { .. } => {}
        _ => panic!("Expected InvalidPrerequisite error"),
    }
}

#[test]
fn test_multiple_prerequisite_groups_and_logic() {
    let (mut deps, env, _info) = proper_instantiate();

    // Complex scenario: Subject requires (MATH1 AND MATH2) AND (PHYSICS1 OR PHYSICS2)
    let math_group = PrerequisiteGroup {
        id: "math_requirements".to_string(),
        subject_id: "QUANTUM_PHYSICS".to_string(),
        group_type: GroupType::All, // Must complete ALL math subjects
        minimum_credits: 0,
        minimum_completed_subjects: 2,
        subject_ids: vec!["MATH1".to_string(), "MATH2".to_string()],
        logic: LogicType::And,
        priority: 1,
        confidence: 95,
    };

    let physics_group = PrerequisiteGroup {
        id: "physics_requirements".to_string(),
        subject_id: "QUANTUM_PHYSICS".to_string(),
        group_type: GroupType::Any, // Must complete ANY physics subject
        minimum_credits: 0,
        minimum_completed_subjects: 1,
        subject_ids: vec!["PHYSICS1".to_string(), "PHYSICS2".to_string()],
        logic: LogicType::And,
        priority: 2,
        confidence: 95,
    };

    let msg = ExecuteMsg::RegisterPrerequisites {
        subject_id: "QUANTUM_PHYSICS".to_string(),
        prerequisites: vec![math_group, physics_group],
    };

    let info = mock_info(OWNER, &[]);
    execute(deps.as_mut(), env.clone(), info, msg).unwrap();

    // Add MATH1, MATH2, and PHYSICS1
    for (subject, credits) in [("MATH1", 4), ("MATH2", 4), ("PHYSICS1", 3)] {
        let completed_subject = CompletedSubjectMsg {
            subject_id: subject.to_string(),
            credits,
            completion_date: "2024-01-15".to_string(),
            grade: 8000,
            nft_token_id: format!("nft_{}_123", subject),
        };

        let msg = ExecuteMsg::UpdateStudentRecord {
            student_id: STUDENT.to_string(),
            completed_subject,
        };

        let info = mock_info("institution", &[]);
        execute(deps.as_mut(), env.clone(), info, msg).unwrap();
    }

    // Should be able to enroll (has all math + one physics)
    let msg = ExecuteMsg::VerifyEnrollment {
        student_id: STUDENT.to_string(),
        subject_id: "QUANTUM_PHYSICS".to_string(),
    };

    let info = mock_info("anyone", &[]);
    let res = execute(deps.as_mut(), env, info, msg).unwrap();
    assert_eq!(res.attributes[3].value, "true");
}

#[test]
fn test_query_nonexistent_student() {
    let (deps, _env, _info) = proper_instantiate();

    // Query student that doesn't exist
    let res = query(
        deps.as_ref(), 
        mock_env(), 
        QueryMsg::GetStudentRecord { 
            student_id: "nonexistent_student".to_string() 
        }
    );
    
    // Should return an error
    assert!(res.is_err());
}

#[test]
fn test_large_number_of_prerequisites() {
    let (mut deps, env, _info) = proper_instantiate();

    // Create many prerequisite subjects (stress test)
    let mut subject_ids = Vec::new();
    for i in 0..50 {
        subject_ids.push(format!("PREREQ_{}", i));
    }

    let prereq_group = PrerequisiteGroup {
        id: "massive_prerequisites".to_string(),
        subject_id: "FINAL_COURSE".to_string(),
        group_type: GroupType::All, // Must complete ALL 50 subjects
        minimum_credits: 0,
        minimum_completed_subjects: 50,
        subject_ids,
        logic: LogicType::And,
        priority: 1,
        confidence: 95,
    };

    let msg = ExecuteMsg::RegisterPrerequisites {
        subject_id: "FINAL_COURSE".to_string(),
        prerequisites: vec![prereq_group],
    };

    let info = mock_info(OWNER, &[]);
    let res = execute(deps.as_mut(), env.clone(), info, msg);
    
    // Should handle large number of prerequisites without crashing
    assert!(res.is_ok());

    // Test verification with incomplete prerequisites
    let msg = ExecuteMsg::VerifyEnrollment {
        student_id: STUDENT.to_string(),
        subject_id: "FINAL_COURSE".to_string(),
    };

    let info = mock_info("anyone", &[]);
    let res = execute(deps.as_mut(), env, info, msg).unwrap();
    assert_eq!(res.attributes[3].value, "false"); // Should fail (no subjects completed)
}
