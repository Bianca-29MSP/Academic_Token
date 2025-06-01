// Schedule Contract - src/contract.rs
use cosmwasm_std::{
    entry_point, to_json_binary, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult,
};
use cw2::set_contract_version;

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg};
use crate::state::{State, ScheduleConfig, STATE, SCHEDULE_CONFIG};
use crate::{execute, query};

const CONTRACT_NAME: &str = "crates.io:schedule";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[entry_point]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;

    let owner = msg.owner
        .and_then(|s| deps.api.addr_validate(&s).ok())
        .unwrap_or(info.sender);

    let state = State {
        owner: owner.clone(),
        ipfs_gateway: msg.ipfs_gateway.clone(),
        total_students: 0,
        total_recommendations: 0,
        total_paths: 0,
    };

    let config = ScheduleConfig {
        max_subjects_per_semester: msg.max_subjects_per_semester,
        recommendation_algorithm: msg.recommendation_algorithm,
        optimization_weights: crate::state::OptimizationWeights {
            graduation_speed: 25,
            workload_balance: 25,
            difficulty_distribution: 20,
            subject_availability: 15,
            student_preferences: 15,
        },
        default_preferences: crate::state::DefaultSchedulePreferences {
            max_subjects_per_semester: msg.max_subjects_per_semester,
            study_intensity: crate::state::StudyIntensity::Moderate,
            balance_theory_practice: true,
            prefer_prerequisites_early: true,
        },
    };

    STATE.save(deps.storage, &state)?;
    SCHEDULE_CONFIG.save(deps.storage, &config)?;

    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("owner", owner.to_string())
        .add_attribute("ipfs_gateway", msg.ipfs_gateway))
}

#[entry_point]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::UpdateConfig { 
            ipfs_gateway, 
            max_subjects_per_semester, 
            recommendation_algorithm, 
            new_owner 
        } => execute::update_config(deps, info, ipfs_gateway, max_subjects_per_semester, recommendation_algorithm, new_owner),
        
        ExecuteMsg::RegisterStudentProgress { student_progress } => 
            execute::register_student_progress(deps, env, student_progress),
        
        ExecuteMsg::UpdateStudentPreferences { student_id, preferences } => 
            execute::update_student_preferences(deps, student_id, preferences),
        
        ExecuteMsg::RegisterSubjectScheduleInfo { subject_info } => 
            execute::register_subject_schedule_info(deps, subject_info),
        
        ExecuteMsg::BatchRegisterSubjects { subjects } => 
            execute::batch_register_subjects(deps, subjects),
        
        ExecuteMsg::GenerateScheduleRecommendation { 
            student_id, 
            target_semester, 
            force_refresh, 
            custom_preferences 
        } => execute::generate_schedule_recommendation(deps, env, student_id, target_semester, force_refresh, custom_preferences),
        
        ExecuteMsg::CreateAcademicPath { 
            student_id, 
            path_name, 
            optimization_criteria, 
            target_graduation_semester 
        } => execute::create_academic_path(deps, env, student_id, path_name, optimization_criteria, target_graduation_semester),
        
        ExecuteMsg::OptimizeAcademicPath { 
            path_id, 
            optimization_criteria, 
            preserve_current_semester 
        } => execute::optimize_academic_path(deps, env, path_id, optimization_criteria, preserve_current_semester),
        
        ExecuteMsg::UpdateAcademicPath { 
            path_id, 
            semester_number, 
            new_subjects, 
            notes 
        } => execute::update_academic_path(deps, env, path_id, semester_number, new_subjects, notes),
        
        ExecuteMsg::ActivateAcademicPath { path_id } => 
            execute::activate_academic_path(deps, path_id),
        
        ExecuteMsg::CompleteSubject { 
            student_id, 
            subject_id, 
            grade, 
            completion_date, 
            difficulty_rating, 
            workload_rating, 
            nft_token_id 
        } => execute::complete_subject(deps, env, student_id, subject_id, grade, completion_date, difficulty_rating, workload_rating, nft_token_id),
        
        ExecuteMsg::EnrollInSubject { 
            student_id, 
            subject_id, 
            enrollment_date, 
            expected_completion 
        } => execute::enroll_in_subject(deps, student_id, subject_id, enrollment_date, expected_completion),
        
        ExecuteMsg::CacheIpfsContent { ipfs_link, content } => 
            execute::cache_ipfs_content(deps, ipfs_link, content),
        
        ExecuteMsg::GenerateAlternativeRecommendations { 
            student_id, 
            target_semester, 
            excluded_subjects 
        } => execute::generate_alternative_recommendations(deps, env, student_id, target_semester, excluded_subjects),
        
        ExecuteMsg::SimulateSchedule { 
            student_id, 
            hypothetical_completions, 
            target_semester 
        } => execute::simulate_schedule(deps, env, student_id, hypothetical_completions, target_semester),
    }
}

#[entry_point]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetState {} => to_json_binary(&query::query_state(deps)?),
        QueryMsg::GetConfig {} => to_json_binary(&query::query_config(deps)?),
        QueryMsg::GetStudentProgress { student_id } => to_json_binary(&query::query_student_progress(deps, student_id)?),
        QueryMsg::GetSubjectScheduleInfo { subject_id } => to_json_binary(&query::query_subject_schedule_info(deps, subject_id)?),
        QueryMsg::GetScheduleRecommendation { student_id, semester } => 
           to_json_binary(&query::query_schedule_recommendation(deps, student_id, semester)?),
       QueryMsg::GetAcademicPath { path_id } => 
           to_json_binary(&query::query_academic_path(deps, path_id)?),
       QueryMsg::GetStudentPaths { student_id, include_inactive } => 
           to_json_binary(&query::query_student_paths(deps, student_id, include_inactive)?),
       QueryMsg::GetAvailableSubjects { student_id, semester, include_electives } => 
           to_json_binary(&query::query_available_subjects(deps, student_id, semester, include_electives)?),
       QueryMsg::GetOptimalPath { student_id, criteria, max_paths } => 
           to_json_binary(&query::query_optimal_path(deps, student_id, criteria, max_paths)?),
       QueryMsg::GetGraduationTimeline { student_id, path_id } => 
           to_json_binary(&query::query_graduation_timeline(deps, student_id, path_id)?),
       QueryMsg::GetSubjectSequence { student_id, target_subjects } => 
           to_json_binary(&query::query_subject_sequence(deps, student_id, target_subjects)?),
       QueryMsg::GetWorkloadAnalysis { student_id, semester, proposed_subjects } => 
           to_json_binary(&query::query_workload_analysis(deps, student_id, semester, proposed_subjects)?),
       QueryMsg::GetIpfsCacheStatus { ipfs_link } => 
           to_json_binary(&query::query_ipfs_cache_status(deps, ipfs_link)?),
       QueryMsg::GetCachedContent { ipfs_link } => 
           to_json_binary(&query::query_cached_content(deps, ipfs_link)?),
       QueryMsg::GetScheduleStatistics { student_id, institution_id } => 
           to_json_binary(&query::query_schedule_statistics(deps, student_id, institution_id)?),
   }
}