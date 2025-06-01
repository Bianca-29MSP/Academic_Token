#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{
    to_json_binary, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult,
};
use cw2::set_contract_version;

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg};
use crate::state::{State, STATE, PREREQUISITES};
use crate::execute::{execute_register_prerequisites, execute_update_student_record, execute_verify_enrollment, execute_cache_ipfs_content, execute_analyze_prerequisite_relationship};
use crate::query::{query_prerequisites, query_student_record, query_check_eligibility, query_state, query_verification_history, query_ipfs_cache_status, query_cached_content};

// Contract name and version
const CONTRACT_NAME: &str = "crates.io:prerequisites-contract";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

/// Instantiate the contract
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
        total_subjects: 0,
        total_verifications: 0,
    };
    
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;
    STATE.save(deps.storage, &state)?;
    
    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("owner", state.owner))
}

/// Execute contract functions
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::RegisterPrerequisites { subject_id, prerequisites } => {
            execute_register_prerequisites(deps, info, subject_id, prerequisites)
        }
        ExecuteMsg::UpdateStudentRecord { student_id, completed_subject } => {
            execute_update_student_record(deps, info, student_id, completed_subject)
        }
        ExecuteMsg::VerifyEnrollment { student_id, subject_id } => {
            execute_verify_enrollment(deps, env, info, student_id, subject_id)
        }
        ExecuteMsg::BatchRegisterPrerequisites { items } => {
            execute_batch_register_prerequisites(deps, info, items)
        }
        ExecuteMsg::UpdateOwner { new_owner } => {
            execute_update_owner(deps, info, new_owner)
        }
        ExecuteMsg::CacheIpfsContent { ipfs_link, content } => {
            execute_cache_ipfs_content(deps, info, ipfs_link, content)
        }
        ExecuteMsg::AnalyzePrerequisiteRelationship { 
            source_subject_id, 
            target_subject_id, 
            source_ipfs_link, 
            target_ipfs_link 
        } => {
            execute_analyze_prerequisite_relationship(
                deps, env, info, 
                source_subject_id, 
                target_subject_id, 
                source_ipfs_link, 
                target_ipfs_link
            )
        }
    }
}

/// Query contract state
#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetPrerequisites { subject_id } => {
            to_json_binary(&query_prerequisites(deps, subject_id)?)
        }
        QueryMsg::GetStudentRecord { student_id } => {
            to_json_binary(&query_student_record(deps, student_id)?)
        }
        QueryMsg::CheckEligibility { student_id, subject_id } => {
            to_json_binary(&query_check_eligibility(deps, student_id, subject_id)?)
        }
        QueryMsg::GetVerificationHistory { student_id, limit } => {
            to_json_binary(&query_verification_history(deps, student_id, limit)?)
        }
        QueryMsg::GetState {} => {
            to_json_binary(&query_state(deps)?)
        }
        QueryMsg::GetIpfsCacheStatus { ipfs_link } => {
            to_json_binary(&query_ipfs_cache_status(deps, ipfs_link)?)
        }
        QueryMsg::GetCachedContent { ipfs_link } => {
            to_json_binary(&query_cached_content(deps, ipfs_link)?)
        }
    }
}

/// Execute batch registration
fn execute_batch_register_prerequisites(
    mut deps: DepsMut,
    info: MessageInfo,
    items: Vec<crate::msg::PrerequisiteRegistration>,
) -> Result<Response, ContractError> {
    // Check owner
    let state = STATE.load(deps.storage)?;
    if info.sender != state.owner {
        return Err(ContractError::Unauthorized {});
    }
    
    let mut total_registered = 0u64;
    
    for item in items {
        PREREQUISITES.save(deps.branch().storage, &item.subject_id, &item.prerequisites)?;
        total_registered += 1;
    }
    
    // Update state
    STATE.update(deps.storage, |mut state| -> Result<_, ContractError> {
        state.total_subjects += total_registered;
        Ok(state)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "batch_register_prerequisites")
        .add_attribute("total_registered", total_registered.to_string()))
}

/// Execute owner update
fn execute_update_owner(
    deps: DepsMut,
    info: MessageInfo,
    new_owner: String,
) -> Result<Response, ContractError> {
    let mut state = STATE.load(deps.storage)?;
    
    // Only current owner can update
    if info.sender != state.owner {
        return Err(ContractError::Unauthorized {});
    }
    
    state.owner = deps.api.addr_validate(&new_owner)?;
    STATE.save(deps.storage, &state)?;
    
    Ok(Response::new()
        .add_attribute("method", "update_owner")
        .add_attribute("new_owner", new_owner))
}