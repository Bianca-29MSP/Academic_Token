// src/contract.rs do EQUIVALENCE CONTRACT

#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{
    to_json_binary, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult,
};
use cw2::set_contract_version;

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg};
use crate::state::{State, STATE};
use crate::execute::{
    execute_register_equivalence, execute_analyze_equivalence, execute_analyze_equivalence_enhanced,
    execute_approve_equivalence, execute_submit_transfer_request,
    execute_process_transfer_request, execute_batch_register_equivalences,
    execute_update_config, execute_cache_ipfs_content
};
use crate::query::{
    query_state, query_equivalence, query_find_equivalence,
    query_list_equivalences_by_institution, query_analysis_result,
    query_detailed_analysis_result, query_transfer_request, 
    query_list_student_transfers, query_check_equivalence, 
    query_statistics, query_debug_equivalences
};

// Contract name and version
const CONTRACT_NAME: &str = "crates.io:equivalence-contract";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    let state = State {
        owner: msg.owner
            .map(|o| deps.api.addr_validate(&o))
            .transpose()?
            .unwrap_or(info.sender),
        total_equivalences: 0,
        total_analyses: 0,
        auto_approval_threshold: msg.auto_approval_threshold.unwrap_or(85),
    };
    
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;
    STATE.save(deps.storage, &state)?;
    
    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("owner", state.owner)
        .add_attribute("auto_approval_threshold", state.auto_approval_threshold.to_string()))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::RegisterEquivalence { source_subject, target_subject, analysis_method, notes } => {
            execute_register_equivalence(deps, env, info, source_subject, target_subject, analysis_method, notes)
        }
        ExecuteMsg::AnalyzeEquivalence { equivalence_id, force_reanalysis } => {
            execute_analyze_equivalence(deps, env, info, equivalence_id, force_reanalysis)
        }
        ExecuteMsg::AnalyzeEquivalenceEnhanced { equivalence_id, analysis_options } => {
            execute_analyze_equivalence_enhanced(deps, env, info, equivalence_id, analysis_options)
        }
        ExecuteMsg::ApproveEquivalence { equivalence_id, equivalence_type, similarity_percentage, notes } => {
            execute_approve_equivalence(deps, env, info, equivalence_id, equivalence_type, similarity_percentage, notes)
        }
        ExecuteMsg::SubmitTransferRequest { student_id, source_institution, target_institution, completed_subjects, requested_equivalences } => {
            execute_submit_transfer_request(deps, env, info, student_id, source_institution, target_institution, completed_subjects, requested_equivalences)
        }
        ExecuteMsg::ProcessTransferRequest { transfer_id, approved_equivalences, notes } => {
            execute_process_transfer_request(deps, env, info, transfer_id, approved_equivalences, notes)
        }
        ExecuteMsg::BatchRegisterEquivalences { equivalences } => {
            execute_batch_register_equivalences(deps, env, info, equivalences)
        }
        ExecuteMsg::UpdateConfig { auto_approval_threshold, new_owner } => {
            execute_update_config(deps, info, auto_approval_threshold, new_owner)
        }
        ExecuteMsg::CacheIpfsContent { ipfs_link, content } => {
            execute_cache_ipfs_content(deps, env, info, ipfs_link, content)
        }
    }
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetState {} => {
            to_json_binary(&query_state(deps)?)
        }
        QueryMsg::GetEquivalence { equivalence_id } => {
            to_json_binary(&query_equivalence(deps, equivalence_id)?)
        }
        QueryMsg::FindEquivalence { source_subject_id, target_subject_id } => {
            to_json_binary(&query_find_equivalence(deps, source_subject_id, target_subject_id)?)
        }
        QueryMsg::ListEquivalencesByInstitution { institution_id, limit, start_after } => {
            to_json_binary(&query_list_equivalences_by_institution(deps, institution_id, limit, start_after)?)
        }
        QueryMsg::GetAnalysisResult { analysis_id } => {
            to_json_binary(&query_analysis_result(deps, analysis_id)?)
        }
        QueryMsg::GetDetailedAnalysisResult { analysis_id } => {
            to_json_binary(&query_detailed_analysis_result(deps, analysis_id)?)
        }
        QueryMsg::GetTransferRequest { transfer_id } => {
            to_json_binary(&query_transfer_request(deps, transfer_id)?)
        }
        QueryMsg::ListStudentTransfers { student_id, limit } => {
            to_json_binary(&query_list_student_transfers(deps, student_id, limit)?)
        }
        QueryMsg::CheckEquivalence { source_subject_id, target_subject_id, minimum_similarity } => {
            to_json_binary(&query_check_equivalence(deps, source_subject_id, target_subject_id, minimum_similarity)?)
        }
        QueryMsg::GetStatistics { institution_id } => {
            to_json_binary(&query_statistics(deps, institution_id)?)
        }
        QueryMsg::DebugEquivalences {} => {
            to_json_binary(&query_debug_equivalences(deps)?)
        }
    }
}