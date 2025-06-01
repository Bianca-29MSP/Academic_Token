use cosmwasm_std::{Deps, StdResult, Order};
use crate::msg::{
    PrerequisitesResponse, StudentRecordResponse, VerificationHistoryResponse, 
    StateResponse, IpfsCacheStatusResponse
};
use crate::state::{
    STATE, PREREQUISITES, STUDENT_RECORDS, VERIFICATIONS, 
    PrerequisiteGroup, StudentRecord, VerificationResult
};
use crate::ipfs::{SubjectContent, is_content_cached, fetch_ipfs_content};

/// Query prerequisites for a subject
pub fn query_prerequisites(deps: Deps, subject_id: String) -> StdResult<PrerequisitesResponse> {
    let prerequisites = PREREQUISITES
        .may_load(deps.storage, &subject_id)?
        .unwrap_or_default();
    
    Ok(PrerequisitesResponse {
        subject_id,
        prerequisites,
    })
}

/// Query student record
pub fn query_student_record(deps: Deps, student_id: String) -> StdResult<StudentRecordResponse> {
    let record = STUDENT_RECORDS
        .may_load(deps.storage, &student_id)?
        .ok_or_else(|| cosmwasm_std::StdError::not_found("Student record not found"))?;
    
    Ok(StudentRecordResponse { record })
}

/// Check enrollment eligibility
pub fn query_check_eligibility(
    deps: Deps, 
    student_id: String, 
    subject_id: String
) -> StdResult<VerificationResult> {
    // Load student record
    let student_record = STUDENT_RECORDS
        .may_load(deps.storage, &student_id)?
        .unwrap_or_else(|| StudentRecord {
            student_id: student_id.clone(),
            completed_subjects: vec![],
            total_credits: 0,
        });
    
    // Load prerequisites
    let prerequisites = PREREQUISITES
        .may_load(deps.storage, &subject_id)?
        .unwrap_or_default();
    
    // Perform verification (reuse logic from execute.rs)
    verify_prerequisites_query(&student_record, &prerequisites)
}

/// Query verification history for a student
pub fn query_verification_history(
    deps: Deps, 
    student_id: String, 
    limit: Option<u32>
) -> StdResult<VerificationHistoryResponse> {
    let limit = limit.unwrap_or(50).min(100) as usize;
    
    let verifications: Result<Vec<_>, _> = VERIFICATIONS
        .range(deps.storage, None, None, Order::Descending)
        .take(limit)
        .filter(|item| {
            if let Ok((key, _)) = item {
                key.starts_with(&student_id)
            } else {
                false
            }
        })
        .collect();
    
    match verifications {
        Ok(results) => Ok(VerificationHistoryResponse {
            verifications: results,
        }),
        Err(e) => Err(e),
    }
}

/// Query contract state
pub fn query_state(deps: Deps) -> StdResult<StateResponse> {
    let state = STATE.load(deps.storage)?;
    Ok(StateResponse {
        owner: state.owner.to_string(),
        total_subjects: state.total_subjects,
        total_verifications: state.total_verifications,
    })
}

/// Check if IPFS content is cached
pub fn query_ipfs_cache_status(deps: Deps, ipfs_link: String) -> StdResult<IpfsCacheStatusResponse> {
    let is_cached = is_content_cached(deps, &ipfs_link)?;
    Ok(IpfsCacheStatusResponse {
        is_cached,
        ipfs_link,
    })
}

/// Get cached IPFS content
pub fn query_cached_content(deps: Deps, ipfs_link: String) -> StdResult<SubjectContent> {
    fetch_ipfs_content(deps, &ipfs_link)
}

/// Core logic for prerequisite verification (query version)
fn verify_prerequisites_query(
    student_record: &StudentRecord,
    prerequisites: &[PrerequisiteGroup],
) -> StdResult<VerificationResult> {
    use crate::state::{GroupType, LogicType};
    
    let mut satisfied_groups = vec![];
    let mut unsatisfied_groups = vec![];
    let mut missing_prerequisites = vec![];
    
    // If no prerequisites, student can enroll
    if prerequisites.is_empty() {
        return Ok(VerificationResult {
            can_enroll: true,
            missing_prerequisites: vec![],
            satisfied_groups: vec![],
            unsatisfied_groups: vec![],
            verification_timestamp: 0,
            details: "No prerequisites required".to_string(),
            used_ipfs_content: false,
        });
    }
    
    // Check if subject has "NONE" type prerequisites
    if prerequisites.len() == 1 && matches!(prerequisites[0].group_type, GroupType::None) {
        return Ok(VerificationResult {
            can_enroll: true,
            missing_prerequisites: vec![],
            satisfied_groups: vec!["none".to_string()],
            unsatisfied_groups: vec![],
            verification_timestamp: 0,
            details: "No prerequisites required (explicit)".to_string(),
            used_ipfs_content: false,
        });
    }
    
    // Get list of completed subject IDs
    let completed_ids: Vec<String> = student_record
        .completed_subjects
        .iter()
        .map(|s| s.subject_id.clone())
        .collect();
    
    // Check each prerequisite group
    for group in prerequisites {
        let is_satisfied = match &group.group_type {
            GroupType::All => {
                // All subjects in the group must be completed
                group.subject_ids.iter().all(|id| completed_ids.contains(id))
            }
            GroupType::Any => {
                // At least one subject must be completed
                group.subject_ids.iter().any(|id| completed_ids.contains(id))
            }
            GroupType::Minimum => {
                // Check minimum credits and/or subjects
                let completed_in_group = group.subject_ids
                    .iter()
                    .filter(|id| completed_ids.contains(id))
                    .count() as u64;
                
                let credits_in_group: u64 = student_record
                    .completed_subjects
                    .iter()
                    .filter(|s| group.subject_ids.contains(&s.subject_id))
                    .map(|s| s.credits)
                    .sum();
                
                completed_in_group >= group.minimum_completed_subjects &&
                credits_in_group >= group.minimum_credits
            }
            GroupType::None => true,
        };
        
        if is_satisfied {
            satisfied_groups.push(group.id.clone());
        } else {
            unsatisfied_groups.push(group.id.clone());
            
            // Add missing subjects to the list
            for subject_id in &group.subject_ids {
                if !completed_ids.contains(subject_id) && !missing_prerequisites.contains(subject_id) {
                    missing_prerequisites.push(subject_id.clone());
                }
            }
        }
    }
    
    // Determine if student can enroll based on logic type
    let can_enroll = if prerequisites.is_empty() {
        true
    } else {
        match prerequisites[0].logic {
            LogicType::And => unsatisfied_groups.is_empty(),
            LogicType::Or => !satisfied_groups.is_empty(),
            LogicType::Xor => satisfied_groups.len() == 1,
            LogicType::Threshold => {
                // For threshold, we need custom logic
                // For now, require at least 50% of groups to be satisfied
                satisfied_groups.len() >= prerequisites.len() / 2
            }
            LogicType::None => true,
        }
    };
    
    let details = if can_enroll {
        "All prerequisites satisfied".to_string()
    } else {
        format!(
            "Missing prerequisites: {}. Unsatisfied groups: {}",
            missing_prerequisites.join(", "),
            unsatisfied_groups.join(", ")
        )
    };
    
    Ok(VerificationResult {
        can_enroll,
        missing_prerequisites,
        satisfied_groups,
        unsatisfied_groups,
        verification_timestamp: 0,
        details,
        used_ipfs_content: false, // Could be enhanced to use IPFS content analysis
    })
}
