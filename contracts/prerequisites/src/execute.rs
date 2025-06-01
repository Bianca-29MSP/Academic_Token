use cosmwasm_std::{DepsMut, Env, MessageInfo, Response};
use crate::error::ContractError;
use crate::msg::CompletedSubjectMsg;
use crate::state::{
    PrerequisiteGroup, StudentRecord, CompletedSubject, STATE, PREREQUISITES, 
    STUDENT_RECORDS, VerificationResult, VERIFICATIONS, GroupType, LogicType
};
use crate::ipfs::{
    SubjectContent, cache_ipfs_content, fetch_ipfs_content, 
    analyze_prerequisite_relationship, verify_content_integrity
};

/// Register prerequisites for a subject
pub fn execute_register_prerequisites(
    deps: DepsMut,
    info: MessageInfo,
    subject_id: String,
    prerequisites: Vec<PrerequisiteGroup>,
) -> Result<Response, ContractError> {
    // Check if sender is owner
    let state = STATE.load(deps.storage)?;
    if info.sender != state.owner {
        return Err(ContractError::Unauthorized {});
    }
    
    // Validate prerequisites
    for prereq in &prerequisites {
        if prereq.subject_id != subject_id {
            return Err(ContractError::InvalidPrerequisite {
                reason: "Subject ID mismatch".to_string(),
            });
        }
    }
    
    // Save prerequisites
    PREREQUISITES.save(deps.storage, &subject_id, &prerequisites)?;
    
    // Update state
    STATE.update(deps.storage, |mut state| -> Result<_, ContractError> {
        state.total_subjects += 1;
        Ok(state)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "register_prerequisites")
        .add_attribute("subject_id", subject_id)
        .add_attribute("groups_count", prerequisites.len().to_string()))
}

/// Update student record with completed subjects
pub fn execute_update_student_record(
    deps: DepsMut,
    _info: MessageInfo,
    student_id: String,
    completed_subject: CompletedSubjectMsg,
) -> Result<Response, ContractError> {
    // In production, verify sender is authorized (institution or student)
    // For now, we'll allow any sender
    
    let completed = CompletedSubject {
        subject_id: completed_subject.subject_id.clone(),
        credits: completed_subject.credits,
        completion_date: completed_subject.completion_date,
        grade: completed_subject.grade, // Now both are u32, no conversion needed
        nft_token_id: completed_subject.nft_token_id,
        ipfs_link: completed_subject.ipfs_link, // Added IPFS link
    };
    
    // Update or create student record
    let mut record = STUDENT_RECORDS
        .may_load(deps.storage, &student_id)?
        .unwrap_or_else(|| StudentRecord {
            student_id: student_id.clone(),
            completed_subjects: vec![],
            total_credits: 0,
        });
    
    // Check if subject already completed
    if record.completed_subjects.iter().any(|s| s.subject_id == completed.subject_id) {
        return Err(ContractError::SubjectAlreadyCompleted {
            subject_id: completed.subject_id,
        });
    }
    
    // Add completed subject
    record.total_credits += completed.credits;
    record.completed_subjects.push(completed);
    
    // Save updated record
    STUDENT_RECORDS.save(deps.storage, &student_id, &record)?;
    
    Ok(Response::new()
        .add_attribute("method", "update_student_record")
        .add_attribute("student_id", student_id)
        .add_attribute("subject_id", completed_subject.subject_id)
        .add_attribute("total_credits", record.total_credits.to_string()))
}

/// Verify if a student can enroll in a subject
pub fn execute_verify_enrollment(
    deps: DepsMut,
    env: Env,
    _info: MessageInfo,
    student_id: String,
    subject_id: String,
) -> Result<Response, ContractError> {
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
    
    // Perform verification
    let mut result = verify_prerequisites(&student_record, &prerequisites)?;
    
    // Set the verification timestamp
    result.verification_timestamp = env.block.time.seconds();
    
    // Create verification ID
    let verification_id = format!("{}-{}-{}", student_id, subject_id, env.block.time.seconds());
    
    // Save verification result
    VERIFICATIONS.save(deps.storage, &verification_id, &result)?;
    
    // Update state
    STATE.update(deps.storage, |mut state| -> Result<_, ContractError> {
        state.total_verifications += 1;
        Ok(state)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "verify_enrollment")
        .add_attribute("student_id", student_id)
        .add_attribute("subject_id", subject_id)
        .add_attribute("can_enroll", result.can_enroll.to_string())
        .add_attribute("verification_id", verification_id))
}

/// Core logic for prerequisite verification
fn verify_prerequisites(
    student_record: &StudentRecord,
    prerequisites: &[PrerequisiteGroup],
) -> Result<VerificationResult, ContractError> {
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
            verification_timestamp: 0, // Will be set by caller
            details: "No prerequisites required".to_string(),
            used_ipfs_content: false, // Added IPFS usage indicator
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
            used_ipfs_content: false, // Added IPFS usage indicator
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
        verification_timestamp: 0, // Will be set by caller
        details,
        used_ipfs_content: false, // Will be updated if IPFS content is used
    })
}

/// Cache IPFS content in the contract storage
pub fn execute_cache_ipfs_content(
    deps: DepsMut,
    info: MessageInfo,
    ipfs_link: String,
    content: SubjectContent,
) -> Result<Response, ContractError> {
    // Check if sender is owner (in production, might want more granular permissions)
    let state = STATE.load(deps.storage)?;
    if info.sender != state.owner {
        return Err(ContractError::Unauthorized {});
    }
    
    // Verify content integrity if hash is provided
    if !content.content_hash.is_empty() {
        if !verify_content_integrity(&content, &content.content_hash) {
            return Err(ContractError::InvalidData {
                reason: "Content hash verification failed".to_string(),
            });
        }
    }
    
    // Cache the content
    cache_ipfs_content(deps, &ipfs_link, &content)
        .map_err(|e| ContractError::StorageError { 
            reason: e.to_string() 
        })?;
    
    Ok(Response::new()
        .add_attribute("method", "cache_ipfs_content")
        .add_attribute("ipfs_link", ipfs_link)
        .add_attribute("content_title", content.title)
        .add_attribute("content_code", content.code))
}

/// Analyze prerequisite relationships using IPFS content
pub fn execute_analyze_prerequisite_relationship(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    source_subject_id: String,
    target_subject_id: String,
    source_ipfs_link: String,
    target_ipfs_link: String,
) -> Result<Response, ContractError> {
    // Check if sender is owner (in production, might want more granular permissions)
    let state = STATE.load(deps.storage)?;
    if info.sender != state.owner {
        return Err(ContractError::Unauthorized {});
    }
    
    // Fetch content from cache
    let source_content = fetch_ipfs_content(deps.as_ref(), &source_ipfs_link)
        .map_err(|_| ContractError::IpfsContentNotFound { 
            ipfs_link: source_ipfs_link.clone() 
        })?;
    
    let target_content = fetch_ipfs_content(deps.as_ref(), &target_ipfs_link)
        .map_err(|_| ContractError::IpfsContentNotFound { 
            ipfs_link: target_ipfs_link.clone() 
        })?;
    
    // Perform analysis
    let analysis = analyze_prerequisite_relationship(&source_content, &target_content);
    
    // If the analysis suggests a prerequisite relationship, we could optionally
    // auto-register the prerequisite (commented out for safety)
    /*
    if analysis.is_prerequisite && analysis.confidence_score >= 70 {
        // Auto-register the prerequisite
        let prerequisite_group = PrerequisiteGroup {
            id: format!("auto-{}-{}", source_subject_id, target_subject_id),
            subject_id: target_subject_id.clone(),
            group_type: GroupType::All,
            minimum_credits: 0,
            minimum_completed_subjects: 1,
            subject_ids: vec![source_subject_id.clone()],
            logic: LogicType::And,
            priority: 1,
            confidence: analysis.confidence_score,
            ipfs_link: Some(source_ipfs_link.clone()),
        };
        
        let mut existing_prereqs = PREREQUISITES
            .may_load(deps.storage, &target_subject_id)?
            .unwrap_or_default();
        
        existing_prereqs.push(prerequisite_group);
        PREREQUISITES.save(deps.storage, &target_subject_id, &existing_prereqs)?;
    }
    */
    
    Ok(Response::new()
        .add_attribute("method", "analyze_prerequisite_relationship")
        .add_attribute("source_subject_id", source_subject_id)
        .add_attribute("target_subject_id", target_subject_id)
        .add_attribute("is_prerequisite", analysis.is_prerequisite.to_string())
        .add_attribute("confidence_score", analysis.confidence_score.to_string())
        .add_attribute("knowledge_overlap", analysis.knowledge_overlap.to_string())
        .add_attribute("competency_overlap", analysis.competency_overlap.to_string()))
}
