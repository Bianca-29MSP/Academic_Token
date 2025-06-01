use cosmwasm_std::{
    entry_point, to_json_binary, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult,
    StdError, Uint128,
};
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg, ValidateDegreeRequirementsResponse};
use crate::state::{
    Config, CurriculumRequirements, DegreeValidationResult,
    CONFIG, CURRICULUM_REQUIREMENTS, VALIDATION_CACHE,
};
use sha2::{Sha256, Digest};

#[entry_point]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> StdResult<Response> {
    let config = Config {
        admin: info.sender.clone(), // Clone here
        module_address: deps.api.addr_validate(&msg.module_address)?,
    };
    
    CONFIG.save(deps.storage, &config)?;
    
    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("admin", info.sender) // Now we can use it here
        .add_attribute("module_address", msg.module_address))
}

#[entry_point]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> StdResult<Response> {
    match msg {
        ExecuteMsg::SetCurriculumRequirements {
            curriculum_id,
            min_credits,
            required_subjects,
            min_gpa,
            additional_requirements,
        } => set_curriculum_requirements(
            deps,
            info,
            curriculum_id,
            min_credits,
            required_subjects,
            min_gpa,
            additional_requirements,
        ),
        ExecuteMsg::ValidateDegreeRequirements {
            student_id,
            curriculum_id,
            institution_id,
            final_gpa,
            total_credits,
            completed_subjects,
            signatures,
            requested_date,
        } => validate_degree_requirements(
            deps,
            env,
            info,
            student_id,
            curriculum_id,
            institution_id,
            final_gpa,
            total_credits,
            completed_subjects,
            signatures,
            requested_date,
        ),
        ExecuteMsg::ClearValidationCache {
            student_id,
            curriculum_id,
        } => clear_validation_cache(deps, info, student_id, curriculum_id),
    }
}

fn set_curriculum_requirements(
    deps: DepsMut,
    info: MessageInfo,
    curriculum_id: String,
    min_credits: Uint128,
    required_subjects: Vec<String>,
    min_gpa: Option<String>,
    additional_requirements: Vec<String>,
) -> StdResult<Response> {
    let config = CONFIG.load(deps.storage)?;
    if info.sender != config.admin {
        return Err(StdError::generic_err("Unauthorized"));
    }
    
    let requirements = CurriculumRequirements {
        curriculum_id: curriculum_id.clone(),
        min_credits,
        required_subjects,
        min_gpa,
        additional_requirements,
    };
    
    CURRICULUM_REQUIREMENTS.save(deps.storage, &curriculum_id, &requirements)?;
    
    Ok(Response::new()
        .add_attribute("method", "set_curriculum_requirements")
        .add_attribute("curriculum_id", curriculum_id))
}

fn validate_degree_requirements(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    student_id: String,
    curriculum_id: String,
    _institution_id: String,
    final_gpa: String,
    total_credits: u64,
    completed_subjects: Vec<String>,
    _signatures: Vec<String>,
    _requested_date: String,
) -> StdResult<Response> {
    let config = CONFIG.load(deps.storage)?;
    if info.sender != config.module_address && info.sender != config.admin {
        return Err(StdError::generic_err("Unauthorized"));
    }
    
    let requirements = CURRICULUM_REQUIREMENTS.load(deps.storage, &curriculum_id)?;
    
    let total_credits_u128 = Uint128::from(total_credits);
    let credits_met = total_credits_u128 >= requirements.min_credits;
    
    let gpa_met = if let Some(min_gpa) = requirements.min_gpa {
        final_gpa.parse::<f64>().unwrap_or(0.0) >= min_gpa.parse::<f64>().unwrap_or(0.0)
    } else {
        true
    };
    
    let missing_subjects: Vec<String> = requirements.required_subjects
        .iter()
        .filter(|subj| !completed_subjects.contains(subj))
        .cloned()
        .collect();
    
    let subjects_met = missing_subjects.is_empty();
    let is_valid = credits_met && gpa_met && subjects_met;
    
    let mut requirements_met = Vec::new();
    if credits_met {
        requirements_met.push(format!("Credits: {}/{}", total_credits, requirements.min_credits));
    }
    if gpa_met {
        requirements_met.push(format!("GPA: {}", final_gpa));
    }
    if subjects_met {
        requirements_met.push("All required subjects".to_string());
    }
    
    let mut missing_requirements = Vec::new();
    if !credits_met {
        missing_requirements.push(format!("Credits: need {}, have {}", requirements.min_credits, total_credits));
    }
    if !gpa_met {
        missing_requirements.push(format!("GPA below minimum"));
    }
    if !subjects_met {
        missing_requirements.push(format!("Missing subjects: {:?}", missing_subjects));
    }
    
    let validation_string = format!("{}-{}-{}-{}", student_id, curriculum_id, final_gpa, total_credits);
    let mut hasher = Sha256::new();
    hasher.update(validation_string.as_bytes());
    let validation_hash = format!("{:x}", hasher.finalize());
    
    let validation_result = DegreeValidationResult {
        is_valid,
        validation_score: if is_valid { "100".to_string() } else { "0".to_string() },
        message: if is_valid { "Requirements met".to_string() } else { "Requirements not met".to_string() },
        requirements_met: requirements_met.clone(),
        missing_requirements: missing_requirements.clone(),
        total_credits_completed: total_credits_u128,
        gpa: final_gpa.clone(),
    };
    
    VALIDATION_CACHE.save(deps.storage, (&student_id, &curriculum_id), &validation_result)?;
    
    let response_data = ValidateDegreeRequirementsResponse {
        is_valid,
        message: validation_result.message.clone(),
        degree_type: "Bachelor".to_string(),
        curriculum_version: "1.0".to_string(),
        validation_hash: validation_hash.clone(),
        requirements_met,
        missing_requirements,
    };
    
    Ok(Response::new()
        .add_attribute("method", "validate_degree_requirements")
        .add_attribute("is_valid", is_valid.to_string())
        .set_data(to_json_binary(&response_data)?))
}

fn clear_validation_cache(
    deps: DepsMut,
    info: MessageInfo,
    student_id: String,
    curriculum_id: String,
) -> StdResult<Response> {
    let config = CONFIG.load(deps.storage)?;
    if info.sender != config.admin {
        return Err(StdError::generic_err("Unauthorized"));
    }
    
    VALIDATION_CACHE.remove(deps.storage, (&student_id, &curriculum_id));
    
    Ok(Response::new()
        .add_attribute("method", "clear_validation_cache")
        .add_attribute("student_id", student_id)
        .add_attribute("curriculum_id", curriculum_id))
}

#[entry_point]
pub fn query(deps: Deps, _env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetConfig {} => to_json_binary(&CONFIG.load(deps.storage)?),
        QueryMsg::GetCurriculumRequirements { curriculum_id } => {
            to_json_binary(&CURRICULUM_REQUIREMENTS.load(deps.storage, &curriculum_id)?)
        }
        QueryMsg::GetValidationResult {
            student_id,
            curriculum_id,
        } => {
            let result = VALIDATION_CACHE.may_load(deps.storage, (&student_id, &curriculum_id))?
                .ok_or_else(|| StdError::generic_err("No validation result found"))?;
            to_json_binary(&result)
        }
        QueryMsg::GetMissingRequirements {
            student_id: _,
            curriculum_id,
            completed_subjects,
        } => {
            let requirements = CURRICULUM_REQUIREMENTS.load(deps.storage, &curriculum_id)?;
            let missing: Vec<String> = requirements.required_subjects
                .iter()
                .filter(|subj| !completed_subjects.contains(subj))
                .cloned()
                .collect();
            to_json_binary(&missing)
        }
    }
}
