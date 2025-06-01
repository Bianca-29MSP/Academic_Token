use cosmwasm_std::{DepsMut, Env, MessageInfo, Response, StdResult};
use crate::error::ContractError;
use crate::state::{
    SubjectInfo, AnalysisMethod, EquivalenceType, EquivalenceStatus,
    Equivalence, STATE, EQUIVALENCES, EQUIVALENCE_INDEX, ANALYSIS_RESULTS, AnalysisResult,
    DetailedAnalysisResult, QualityMetrics, DETAILED_ANALYSIS_RESULTS,
    TransferRequest, TransferStatus, TRANSFER_REQUESTS, STUDENT_TRANSFERS,
};
use crate::msg::{EquivalenceRegistration, AnalysisOptions};
use crate::ipfs::{
    MultilingualSyllabusContent, SyllabusContent, compare_multilingual_content,
    cache_ipfs_content, fetch_ipfs_content, calculate_content_hash, verify_content_integrity,
    assess_content_quality_multilingual, calculate_enhanced_semantic_similarity,
    generate_detailed_recommendations, ContentQualityAssessment, Language,
    SimilarityAnalysis, WorkloadBreakdown,
    calculate_list_similarity, AnalysisWeights,
};


/// Register a new equivalence between subjects
pub fn execute_register_equivalence(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    source_subject: SubjectInfo,
    target_subject: SubjectInfo,
    analysis_method: AnalysisMethod,
    notes: Option<String>,
) -> Result<Response, ContractError> {
    // Check authorization (only owner for now)
    let state = STATE.load(deps.storage)?;
    if info.sender != state.owner {
        return Err(ContractError::Unauthorized {});
    }
    
    // Validate input
    if source_subject.subject_id == target_subject.subject_id {
        return Err(ContractError::SelfEquivalenceNotAllowed {
            subject_id: source_subject.subject_id,
        });
    }
    
    // Check if equivalence already exists (CLONANDO para evitar borrow issues)
    let source_id = source_subject.subject_id.clone();
    let target_id = target_subject.subject_id.clone();
    let equiv_key = (source_id.as_str(), target_id.as_str());
    
    if EQUIVALENCE_INDEX.may_load(deps.storage, equiv_key)?.is_some() {
        return Err(ContractError::EquivalenceAlreadyExists {
            source_id: source_id.clone(),
            target_id: target_id.clone(),
        });
    }
    
    // Generate equivalence ID
    let equivalence_id = format!(
        "{}-{}-{}",
        source_id,
        target_id,
        env.block.time.seconds()
    );
    
    // Calculate initial similarity based on basic metadata
    let initial_similarity = calculate_basic_similarity(&source_subject, &target_subject);
    
    // Determine initial status and other values ANTES de mover
    let status = match &analysis_method {
        AnalysisMethod::Manual => EquivalenceStatus::Pending,
        AnalysisMethod::Institutional => EquivalenceStatus::Approved,
        _ => EquivalenceStatus::Analyzing,
    };
    
    let equivalence_type = if status == EquivalenceStatus::Approved {
        EquivalenceType::Full
    } else {
        EquivalenceType::None
    };
    
    let approved_timestamp = if status == EquivalenceStatus::Approved {
        Some(env.block.time.seconds())
    } else {
        None
    };
    
    let approved_by = if status == EquivalenceStatus::Approved {
        Some(info.sender.clone())
    } else {
        None
    };
    
    let confidence_score = match &analysis_method {
        AnalysisMethod::Institutional => 100,
        AnalysisMethod::Manual => 0,
        _ => 50,
    };
    
    // Create equivalence record
    let equivalence = Equivalence {
        id: equivalence_id.clone(),
        source_subject,
        target_subject,
        equivalence_type,
        similarity_percentage: initial_similarity,
        status,
        analysis_method,
        created_timestamp: env.block.time.seconds(),
        approved_timestamp,
        approved_by,
        notes: notes.unwrap_or_default(),
        confidence_score,
    };
    
    // Save equivalence
    EQUIVALENCES.save(deps.storage, &equivalence_id, &equivalence)?;
    
    // Update index
    EQUIVALENCE_INDEX.save(deps.storage, equiv_key, &equivalence_id)?;
    
    // Update state
    STATE.update(deps.storage, |mut state| -> Result<_, ContractError> {
        state.total_equivalences += 1;
        Ok(state)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "register_equivalence")
        .add_attribute("equivalence_id", equivalence_id)
        .add_attribute("source_subject", equivalence.source_subject.subject_id)
        .add_attribute("target_subject", equivalence.target_subject.subject_id)
        .add_attribute("similarity_percentage", equivalence.similarity_percentage.to_string())
        .add_attribute("status", format!("{:?}", equivalence.status)))
}

/// Cache IPFS content in the contract storage
pub fn execute_cache_ipfs_content(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    ipfs_link: String,
    content: MultilingualSyllabusContent,
) -> Result<Response, ContractError> {
    // For now, allow anyone to cache content
    // TODO: Add oracle authorization check
    
    // Validate content structure
    if content.content.title.is_empty() {
        return Err(ContractError::InvalidContent {
            reason: "Title cannot be empty".to_string(),
        });
    }
    
    // Calculate content hash
    let calculated_hash = calculate_content_hash(&content.content)
        .map_err(|e| ContractError::ContentProcessingError { 
            reason: e.to_string() 
        })?;
    
    // Cache the content
    cache_ipfs_content(deps, &ipfs_link, &content)
        .map_err(|e| ContractError::StorageError { 
            reason: e.to_string() 
        })?;
    
    Ok(Response::new()
        .add_attribute("method", "cache_ipfs_content")
        .add_attribute("ipfs_link", ipfs_link)
        .add_attribute("content_title", content.content.title)
        .add_attribute("primary_language", format!("{:?}", content.primary_language))
        .add_attribute("content_hash", calculated_hash)
        .add_attribute("cached_by", info.sender))
}

/// Analyze equivalence using IPFS content
pub fn execute_analyze_equivalence(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    equivalence_id: String,
    force_reanalysis: Option<bool>,
) -> Result<Response, ContractError> {
    // Load equivalence record
    let mut equivalence = EQUIVALENCES.load(deps.storage, &equivalence_id)?;
    
    // Check if already analyzed and not forcing reanalysis
    if equivalence.status != EquivalenceStatus::Analyzing && 
       equivalence.status != EquivalenceStatus::Pending &&
       !force_reanalysis.unwrap_or(false) {
        return Err(ContractError::EquivalenceAlreadyAnalyzed { 
            equivalence_id: equivalence_id.clone() 
        });
    }
    
    // Update status to analyzing
    equivalence.status = EquivalenceStatus::Analyzing;
    EQUIVALENCES.save(deps.storage, &equivalence_id, &equivalence)?;
    
    // Try to fetch IPFS content for both subjects
    let source_content_result = fetch_ipfs_content(deps.as_ref(), &equivalence.source_subject.ipfs_link);
    let target_content_result = fetch_ipfs_content(deps.as_ref(), &equivalence.target_subject.ipfs_link);
    
    let (content_similarity, use_ipfs_content) = match (source_content_result, target_content_result) {
        (Ok(source_content), Ok(target_content)) => {
            // Verify content integrity
            if !verify_content_integrity(&source_content.content, &equivalence.source_subject.content_hash)
                .unwrap_or(false) {
                return Err(ContractError::ContentIntegrityError { 
                    subject_id: equivalence.source_subject.subject_id.clone() 
                });
            }
            
            if !verify_content_integrity(&target_content.content, &equivalence.target_subject.content_hash)
                .unwrap_or(false) {
                return Err(ContractError::ContentIntegrityError { 
                    subject_id: equivalence.target_subject.subject_id.clone() 
                });
            }
            
            // Perform multilingual content analysis
            let similarity_analysis = compare_multilingual_content(&source_content, &target_content)
                .map_err(|e| ContractError::ContentProcessingError { 
                    reason: e.to_string() 
                })?;
            
            // Store detailed analysis results
            let analysis_result = AnalysisResult {
                equivalence_id: equivalence_id.clone(),
                content_similarity: similarity_analysis.overall_similarity,
                structure_similarity: similarity_analysis.structural_similarity,
                credit_compatibility: calculate_credit_compatibility(&equivalence.source_subject, &equivalence.target_subject),
                level_compatibility: calculate_level_compatibility(&equivalence.source_subject, &equivalence.target_subject),
                overall_score: similarity_analysis.overall_similarity,
                recommendation: determine_equivalence_type(similarity_analysis.overall_similarity),
                analysis_details: format!(
                    "Title: {}%, Desc: {}%, Obj: {}%, Topics: {}%, Workload: {}%",
                    similarity_analysis.title_similarity,
                    similarity_analysis.description_similarity,
                    similarity_analysis.objectives_similarity,
                    similarity_analysis.topics_similarity,
                    similarity_analysis.workload_similarity
                ),
                analyzed_timestamp: env.block.time.seconds(),
            };
            
            // Save analysis result
            let analysis_id = format!("analysis_{}", equivalence_id);
            ANALYSIS_RESULTS.save(deps.storage, &analysis_id, &analysis_result)?;
            
            (similarity_analysis.overall_similarity, true)
        }
        _ => {
            // Fallback to metadata-only analysis
            let metadata_similarity = calculate_basic_similarity(
                &equivalence.source_subject, 
                &equivalence.target_subject
            );
            (metadata_similarity, false)
        }
    };
    
    // Combine with metadata similarity if we have IPFS content
    let final_similarity = if use_ipfs_content {
        // Weight: 80% content, 20% metadata
        let metadata_similarity = calculate_basic_similarity(
            &equivalence.source_subject, 
            &equivalence.target_subject
        );
        (content_similarity * 80 + metadata_similarity * 20) / 100
    } else {
        content_similarity
    };
    
    // Update equivalence with analysis results
    equivalence.similarity_percentage = final_similarity;
    equivalence.confidence_score = calculate_confidence_score(
        final_similarity, 
        use_ipfs_content,
        &equivalence.analysis_method
    );
    
    // Determine final status
    let state = STATE.load(deps.storage)?;
    if final_similarity >= state.auto_approval_threshold {
        equivalence.status = EquivalenceStatus::Approved;
        equivalence.equivalence_type = determine_equivalence_type(final_similarity);
        equivalence.approved_timestamp = Some(env.block.time.seconds());
        equivalence.approved_by = Some(info.sender.clone());
    } else if final_similarity >= 50 {
        equivalence.status = EquivalenceStatus::UnderReview;
        equivalence.equivalence_type = EquivalenceType::Conditional;
    } else {
        equivalence.status = EquivalenceStatus::Rejected;
        equivalence.equivalence_type = EquivalenceType::None;
    }
    
    // Save updated equivalence
    EQUIVALENCES.save(deps.storage, &equivalence_id, &equivalence)?;
    
    // Update total analyses count
    STATE.update(deps.storage, |mut state| -> Result<_, ContractError> {
        state.total_analyses += 1;
        Ok(state)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "analyze_equivalence")
        .add_attribute("equivalence_id", equivalence_id)
        .add_attribute("final_similarity", final_similarity.to_string())
        .add_attribute("confidence_score", equivalence.confidence_score.to_string())
        .add_attribute("status", format!("{:?}", equivalence.status))
        .add_attribute("used_ipfs_content", use_ipfs_content.to_string())
        .add_attribute("equivalence_type", format!("{:?}", equivalence.equivalence_type)))
}

/// Enhanced analysis with multilingual support and quality metrics
pub fn execute_analyze_equivalence_enhanced(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    equivalence_id: String,
    analysis_options: AnalysisOptions,
) -> Result<Response, ContractError> {
    // Load equivalence record
    let mut equivalence = EQUIVALENCES.load(deps.storage, &equivalence_id)?;
    
    // Check if already analyzed
    if equivalence.status != EquivalenceStatus::Analyzing && 
       equivalence.status != EquivalenceStatus::Pending {
        return Err(ContractError::EquivalenceAlreadyAnalyzed { 
            equivalence_id: equivalence_id.clone() 
        });
    }
    
    // Update status
    equivalence.status = EquivalenceStatus::Analyzing;
    EQUIVALENCES.save(deps.storage, &equivalence_id, &equivalence)?;
    
    // Try to fetch multilingual content
    let source_content_result = fetch_ipfs_content(deps.as_ref(), &equivalence.source_subject.ipfs_link);
    let target_content_result = fetch_ipfs_content(deps.as_ref(), &equivalence.target_subject.ipfs_link);
    
    let (analysis_result, used_content) = match (source_content_result, target_content_result) {
        (Ok(source_content), Ok(target_content)) => {
            // Verify content integrity
            if !verify_content_integrity(&source_content.content, &equivalence.source_subject.content_hash)
                .unwrap_or(false) {
                return Err(ContractError::ContentIntegrityError { 
                    subject_id: equivalence.source_subject.subject_id.clone() 
                });
            }
            
            if !verify_content_integrity(&target_content.content, &equivalence.target_subject.content_hash)
                .unwrap_or(false) {
                return Err(ContractError::ContentIntegrityError { 
                    subject_id: equivalence.target_subject.subject_id.clone() 
                });
            }
            
            // Perform comprehensive enhanced analysis
            let comprehensive_result = perform_comprehensive_analysis_enhanced(
                &source_content,
                &target_content,
                &equivalence,
                &analysis_options,
            )?;
            
            (comprehensive_result, true)
        }
        _ => {
            // Fallback to metadata-only analysis
            let metadata_similarity = calculate_basic_similarity(
                &equivalence.source_subject, 
                &equivalence.target_subject
            );
            
            let basic_result = create_basic_analysis_result(
                &equivalence,
                metadata_similarity,
            );
            
            (basic_result, false)
        }
    };
    
    // Save detailed analysis result with timestamp
    let mut final_analysis_result = analysis_result;
    final_analysis_result.base_result.analyzed_timestamp = env.block.time.seconds();
    
    let analysis_id = format!("detailed_analysis_{}", equivalence_id);
    DETAILED_ANALYSIS_RESULTS.save(deps.storage, &analysis_id, &final_analysis_result)?;
    
    // Update equivalence with results
    equivalence.similarity_percentage = final_analysis_result.base_result.overall_score;
    equivalence.confidence_score = final_analysis_result.quality_metrics.analysis_confidence;
    
    // Add quality information to notes
    let quality_note = format!(
        "\n[Analysis Quality] Overall: {}%, Confidence: {}%, Language Compatibility: {}%",
        final_analysis_result.quality_metrics.overall_quality,
        final_analysis_result.quality_metrics.analysis_confidence,
        final_analysis_result.language_compatibility
    );
    equivalence.notes.push_str(&quality_note);
    
    // Determine final status
    let state = STATE.load(deps.storage)?;
    if final_analysis_result.base_result.overall_score >= state.auto_approval_threshold &&
       final_analysis_result.quality_metrics.analysis_confidence >= 70 {
        equivalence.status = EquivalenceStatus::Approved;
        equivalence.equivalence_type = determine_equivalence_type(final_analysis_result.base_result.overall_score);
        equivalence.approved_timestamp = Some(env.block.time.seconds());
        equivalence.approved_by = Some(info.sender.clone());
    } else if final_analysis_result.base_result.overall_score >= 50 {
        equivalence.status = EquivalenceStatus::UnderReview;
        equivalence.equivalence_type = EquivalenceType::Conditional;
    } else {
        equivalence.status = EquivalenceStatus::Rejected;
        equivalence.equivalence_type = EquivalenceType::None;
    }
    
    // Save updated equivalence
    EQUIVALENCES.save(deps.storage, &equivalence_id, &equivalence)?;
    
    // Update total analyses count
    STATE.update(deps.storage, |mut state| -> Result<_, ContractError> {
        state.total_analyses += 1;
        Ok(state)
    })?;
    
    // Log recommendations as events
    let mut response = Response::new()
        .add_attribute("method", "analyze_equivalence_enhanced")
        .add_attribute("equivalence_id", equivalence_id)
        .add_attribute("final_similarity", final_analysis_result.base_result.overall_score.to_string())
        .add_attribute("quality_score", final_analysis_result.quality_metrics.overall_quality.to_string())
        .add_attribute("confidence_score", final_analysis_result.quality_metrics.analysis_confidence.to_string())
        .add_attribute("language_compatibility", final_analysis_result.language_compatibility.to_string())
        .add_attribute("status", format!("{:?}", equivalence.status))
        .add_attribute("used_content", used_content.to_string())
        .add_attribute("equivalence_type", format!("{:?}", equivalence.equivalence_type));
    
    // Add first 3 recommendations as attributes
    for (i, recommendation) in final_analysis_result.recommendations.iter().take(3).enumerate() {
        response = response.add_attribute(format!("recommendation_{}", i + 1), recommendation);
    }
    
    Ok(response)
}

/// Approve equivalence manually
pub fn execute_approve_equivalence(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    equivalence_id: String,
    equivalence_type: EquivalenceType,
    similarity_percentage: u32,
    notes: Option<String>,
) -> Result<Response, ContractError> {
    // Check authorization
    let state = STATE.load(deps.storage)?;
    if info.sender != state.owner {
        return Err(ContractError::Unauthorized {});
    }
    
    // Validate similarity percentage
    if similarity_percentage > 100 {
        return Err(ContractError::InvalidSimilarityPercentage { 
            percentage: similarity_percentage 
        });
    }
    
    // Load and update equivalence
    let mut equivalence = EQUIVALENCES.load(deps.storage, &equivalence_id)?;
    
    equivalence.equivalence_type = equivalence_type.clone();
    equivalence.similarity_percentage = similarity_percentage;
    equivalence.status = EquivalenceStatus::Approved;
    equivalence.approved_timestamp = Some(env.block.time.seconds());
    equivalence.approved_by = Some(info.sender.clone());
    
    if let Some(note) = notes {
        equivalence.notes = if equivalence.notes.is_empty() {
            note
        } else {
            format!("{}\n[Manual Approval]: {}", equivalence.notes, note)
        };
    }
    
    // Update confidence score for manual approval
    equivalence.confidence_score = 95; // High confidence for manual approval
    
    EQUIVALENCES.save(deps.storage, &equivalence_id, &equivalence)?;
    
    Ok(Response::new()
        .add_attribute("method", "approve_equivalence")
        .add_attribute("equivalence_id", equivalence_id)
        .add_attribute("equivalence_type", format!("{:?}", equivalence_type))
        .add_attribute("similarity_percentage", similarity_percentage.to_string())
        .add_attribute("approved_by", info.sender))
}

/// Submit transfer request
pub fn execute_submit_transfer_request(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    student_id: String,
    source_institution: String,
    target_institution: String,
    completed_subjects: Vec<String>,
    requested_equivalences: Vec<String>,
) -> Result<Response, ContractError> {
    // Generate transfer ID
    let transfer_id = format!(
        "transfer_{}_{}",
        student_id,
        env.block.time.seconds()
    );
    
    // Create transfer request
    let transfer_request = TransferRequest {
        id: transfer_id.clone(),
        student_id: student_id.clone(),
        source_institution,
        target_institution,
        completed_subjects,
        requested_equivalences,
        approved_equivalences: Vec::new(),
        status: TransferStatus::Pending,
        submitted_timestamp: env.block.time.seconds(),
        processed_timestamp: None,
        processed_by: None,
        notes: String::new(),
    };
    
    // Save transfer request
    TRANSFER_REQUESTS.save(deps.storage, &transfer_id, &transfer_request)?;
    
    // Update student transfers index
    STUDENT_TRANSFERS.update(deps.storage, &student_id, |transfers| -> Result<_, ContractError> {
        let mut transfers = transfers.unwrap_or_default();
        transfers.push(transfer_id.clone());
        Ok(transfers)
    })?;
    
    Ok(Response::new()
        .add_attribute("method", "submit_transfer_request")
        .add_attribute("transfer_id", transfer_id)
        .add_attribute("student_id", student_id)
        .add_attribute("submitted_by", info.sender))
}

/// Process transfer request
pub fn execute_process_transfer_request(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    transfer_id: String,
    approved_equivalences: Vec<String>,
    notes: Option<String>,
) -> Result<Response, ContractError> {
    // Check authorization
    let state = STATE.load(deps.storage)?;
    if info.sender != state.owner {
        return Err(ContractError::Unauthorized {});
    }
    
    // Load and update transfer request
    let mut transfer = TRANSFER_REQUESTS.load(deps.storage, &transfer_id)?;
    
    transfer.approved_equivalences = approved_equivalences.clone();
    transfer.processed_timestamp = Some(env.block.time.seconds());
    transfer.processed_by = Some(info.sender.clone());
    
    if let Some(note) = notes {
        transfer.notes = note;
    }
    
    // Determine status based on approvals
    if approved_equivalences.is_empty() {
        transfer.status = TransferStatus::Rejected;
    } else if approved_equivalences.len() == transfer.requested_equivalences.len() {
        transfer.status = TransferStatus::FullyApproved;
    } else {
        transfer.status = TransferStatus::PartiallyApproved;
    }
    
    TRANSFER_REQUESTS.save(deps.storage, &transfer_id, &transfer)?;
    
    Ok(Response::new()
        .add_attribute("method", "process_transfer_request")
        .add_attribute("transfer_id", transfer_id)
        .add_attribute("status", format!("{:?}", transfer.status))
        .add_attribute("approved_count", approved_equivalences.len().to_string())
        .add_attribute("processed_by", info.sender))
}

/// Batch register equivalences
pub fn execute_batch_register_equivalences(
    mut deps: DepsMut,
    env: Env,
    info: MessageInfo,
    equivalences: Vec<EquivalenceRegistration>,
) -> Result<Response, ContractError> {
    let mut total_registered = 0u32;
    let mut total_failed = 0u32;
    
    for equiv_reg in equivalences {
        match execute_register_equivalence(
            deps.branch(),
            env.clone(),
            info.clone(),
            equiv_reg.source_subject,
            equiv_reg.target_subject,
            equiv_reg.analysis_method,
            equiv_reg.notes,
        ) {
            Ok(_) => total_registered += 1,
            Err(_) => {
                total_failed += 1;
                continue; // Skip failed registrations in batch
            }
        }
    }
    
    Ok(Response::new()
        .add_attribute("method", "batch_register_equivalences")
        .add_attribute("total_registered", total_registered.to_string())
        .add_attribute("total_failed", total_failed.to_string()))
}

/// Update contract configuration
pub fn execute_update_config(
    deps: DepsMut,
    info: MessageInfo,
    auto_approval_threshold: Option<u32>,
    new_owner: Option<String>,
) -> Result<Response, ContractError> {
    let mut state = STATE.load(deps.storage)?;
    
    // Only current owner can update
    if info.sender != state.owner {
        return Err(ContractError::Unauthorized {});
    }
    
    if let Some(threshold) = auto_approval_threshold {
        if threshold > 100 {
            return Err(ContractError::InvalidSimilarityPercentage { percentage: threshold });
        }
        state.auto_approval_threshold = threshold;
    }
    
    if let Some(owner_str) = new_owner {
        state.owner = deps.api.addr_validate(&owner_str)?;
    }
    
    STATE.save(deps.storage, &state)?;
    
    Ok(Response::new()
        .add_attribute("method", "update_config")
        .add_attribute("owner", state.owner)
        .add_attribute("auto_approval_threshold", state.auto_approval_threshold.to_string()))
}

// ================================
// HELPER FUNCTIONS
// ================================

/// Calculate basic similarity between subjects based on metadata
fn calculate_basic_similarity(source: &SubjectInfo, target: &SubjectInfo) -> u32 {
    let mut score = 0u32;
    let mut total_factors = 0u32;
    
    // Credits similarity (40% weight)
    let credit_diff = if source.credits > target.credits {
        source.credits - target.credits
    } else {
        target.credits - source.credits
    };
    let credit_similarity = if credit_diff == 0 {
        100
    } else if credit_diff <= 2 {
        80
    } else if credit_diff <= 4 {
        60
    } else {
        40
    };
    score += credit_similarity * 40;
    total_factors += 40;
    
    // Academic level similarity (30% weight)
    let level_similarity = if source.metadata.level == target.metadata.level {
        100
    } else {
        50
    };
    score += level_similarity * 30;
    total_factors += 30;
    
    // Department similarity (20% weight) - simple string comparison
    let dept_similarity = if source.metadata.department == target.metadata.department {
        100
    } else if source.metadata.department.to_lowercase().contains(&target.metadata.department.to_lowercase()) ||
              target.metadata.department.to_lowercase().contains(&source.metadata.department.to_lowercase()) {
        70
    } else {
        30
    };
    score += dept_similarity * 20;
    total_factors += 20;
    
    // Workload similarity (10% weight)
    let workload_diff = if source.metadata.workload_hours > target.metadata.workload_hours {
        source.metadata.workload_hours - target.metadata.workload_hours
    } else {
        target.metadata.workload_hours - source.metadata.workload_hours
    };
    let workload_similarity = if workload_diff <= 10 {
        100
    } else if workload_diff <= 30 {
        80
    } else {
        60
    };
    score += workload_similarity * 10;
    total_factors += 10;
    
    // Calculate weighted average
    if total_factors > 0 {
        score / total_factors
    } else {
        50 // Default similarity
    }
}

/// Calculate credit compatibility
fn calculate_credit_compatibility(source: &SubjectInfo, target: &SubjectInfo) -> u32 {
    let credit_diff = if source.credits > target.credits {
        source.credits - target.credits
    } else {
        target.credits - source.credits
    };
    
    match credit_diff {
        0 => 100,
        1..=2 => 90,
        3..=4 => 70,
        5..=8 => 50,
        _ => 20,
    }
}

/// Calculate academic level compatibility
fn calculate_level_compatibility(source: &SubjectInfo, target: &SubjectInfo) -> u32 {
    if source.metadata.level == target.metadata.level {
        100
    } else {
        match (&source.metadata.level, &target.metadata.level) {
            (crate::state::AcademicLevel::Undergraduate, crate::state::AcademicLevel::Graduate) |
            (crate::state::AcademicLevel::Graduate, crate::state::AcademicLevel::Undergraduate) => 60,
            (crate::state::AcademicLevel::Graduate, crate::state::AcademicLevel::Postgraduate) |
            (crate::state::AcademicLevel::Postgraduate, crate::state::AcademicLevel::Graduate) => 70,
            _ => 30,
        }
    }
}

/// Determine equivalence type based on similarity score
fn determine_equivalence_type(similarity: u32) -> EquivalenceType {
    match similarity {
        90..=100 => EquivalenceType::Full,
        70..=89 => EquivalenceType::Partial,
        50..=69 => EquivalenceType::Conditional,
        _ => EquivalenceType::None,
    }
}

/// Calculate confidence score based on various factors
fn calculate_confidence_score(
    similarity: u32,
    used_ipfs_content: bool,
    analysis_method: &AnalysisMethod,
) -> u32 {
    let mut base_confidence = similarity;
    
    // Boost confidence if we used IPFS content
    if used_ipfs_content {
        base_confidence = (base_confidence * 110).min(10000) / 100; // 10% boost, max 100
    }
    
    // Adjust based on analysis method
    match analysis_method {
        AnalysisMethod::Institutional => 100,
        AnalysisMethod::Manual => base_confidence.max(80), // Manual review adds confidence
        AnalysisMethod::Hybrid => (base_confidence * 105).min(10000) / 100, // 5% boost
        AnalysisMethod::Automatic => base_confidence,
    }
}

// ===================== ENHANCED ANALYSIS FUNCTIONS =====================

fn perform_comprehensive_analysis_enhanced(
    source_content: &MultilingualSyllabusContent,
    target_content: &MultilingualSyllabusContent,
    equivalence: &Equivalence,
    options: &AnalysisOptions,
) -> Result<DetailedAnalysisResult, ContractError> {
    // Perform quality assessment on both contents
    let source_quality = assess_content_quality_multilingual(source_content);
    let target_quality = assess_content_quality_multilingual(target_content);
    
    // Perform enhanced multilingual comparison
    let similarity_analysis = compare_multilingual_enhanced(source_content, target_content)?;
    
    // Calculate various compatibility scores
    let credit_compatibility = calculate_credit_compatibility(
        &equivalence.source_subject, 
        &equivalence.target_subject
    );
    
    let level_compatibility = calculate_level_compatibility(
        &equivalence.source_subject, 
        &equivalence.target_subject
    );
    
    let (source_content_data, target_content_data, _, _) = 
        find_best_language_pair(source_content, target_content);
    
    let workload_compatibility = calculate_workload_compatibility(
        &source_content_data.workload_breakdown,
        &target_content_data.workload_breakdown,
    );
    
    let prerequisite_alignment = if options.check_prerequisites {
        calculate_prerequisite_alignment(
            &source_content_data.prerequisites,
            &target_content_data.prerequisites,
        )
    } else {
        100
    };
    
    let learning_outcome_alignment = if options.analyze_learning_outcomes {
        similarity_analysis.objectives_similarity
    } else {
        similarity_analysis.overall_similarity
    };
    
    let bibliography_overlap = similarity_analysis.bibliography_similarity;
    
    // Calculate language compatibility
    let language_compatibility = if source_content.primary_language == target_content.primary_language {
        100
    } else if source_content.auto_translated || target_content.auto_translated {
        70
    } else {
        85
    };
    
    // Calculate overall score with enhanced weights
    let overall_score = calculate_overall_score_with_quality(
        &similarity_analysis,
        credit_compatibility,
        level_compatibility,
        workload_compatibility,
        prerequisite_alignment,
        &source_quality,
        &target_quality,
    );
    
    // Assess overall quality metrics
    let quality_metrics = QualityMetrics {
        overall_quality: (source_quality.overall_score + target_quality.overall_score) / 2,
        content_completeness: (source_quality.completeness_score + target_quality.completeness_score) / 2,
        analysis_confidence: calculate_analysis_confidence(&source_quality, &target_quality, options),
        data_availability: (source_quality.completeness_score + target_quality.completeness_score) / 2,
        language_quality: (source_quality.language_quality_score + target_quality.language_quality_score) / 2,
    };
    
    // Generate detailed recommendations
    let recommendations = generate_detailed_recommendations(
        overall_score,
        &source_quality,
        &target_quality,
        &similarity_analysis,
    );
    
    // Create base analysis result
    let base_result = AnalysisResult {
        equivalence_id: equivalence.id.clone(),
        content_similarity: similarity_analysis.overall_similarity,
        structure_similarity: similarity_analysis.structural_similarity,
        credit_compatibility,
        level_compatibility,
        overall_score,
        recommendation: determine_equivalence_type(overall_score),
        analysis_details: format!(
            "Enhanced analysis - Content: {}%, Structure: {}%, Credits: {}%, Level: {}%, Language: {}%",
            similarity_analysis.overall_similarity, 
            similarity_analysis.structural_similarity,
            credit_compatibility, 
            level_compatibility,
            language_compatibility
        ),
        analyzed_timestamp: 0, // Will be set by caller
    };
    
    Ok(DetailedAnalysisResult {
        base_result,
        quality_metrics,
        language_compatibility,
        content_depth_score: options.minimum_content_depth.max(
            (source_quality.depth_score + target_quality.depth_score) / 2
        ),
        prerequisite_alignment,
        learning_outcome_alignment,
        workload_compatibility,
        bibliography_overlap,
        recommendations,
    })
}

fn create_basic_analysis_result(
    equivalence: &Equivalence,
    similarity: u32,
) -> DetailedAnalysisResult {
    let base_result = AnalysisResult {
        equivalence_id: equivalence.id.clone(),
        content_similarity: similarity,
        structure_similarity: 50, // Default
        credit_compatibility: calculate_credit_compatibility(
            &equivalence.source_subject,
            &equivalence.target_subject,
        ),
        level_compatibility: calculate_level_compatibility(
            &equivalence.source_subject,
            &equivalence.target_subject,
        ),
        overall_score: similarity,
        recommendation: determine_equivalence_type(similarity),
        analysis_details: "Basic metadata analysis only".to_string(),
        analyzed_timestamp: 0,
    };
    
    DetailedAnalysisResult {
        base_result,
        quality_metrics: QualityMetrics {
            overall_quality: 50,
            content_completeness: 30,
            analysis_confidence: 40,
            data_availability: 20,
            language_quality: 50,
        },
        language_compatibility: 100,
        content_depth_score: 0,
        prerequisite_alignment: 50,
        learning_outcome_alignment: 50,
        workload_compatibility: 50,
        bibliography_overlap: 0,
        recommendations: vec!["Limited analysis - content not available".to_string()],
    }
}

fn calculate_workload_compatibility(
    source: &WorkloadBreakdown,
    target: &WorkloadBreakdown,
) -> u32 {
    let total_source = source.theoretical_hours + source.practical_hours + 
                      source.laboratory_hours + source.field_work_hours + 
                      source.seminar_hours + source.individual_study_hours;
    
    let total_target = target.theoretical_hours + target.practical_hours + 
                      target.laboratory_hours + target.field_work_hours + 
                      target.seminar_hours + target.individual_study_hours;
    
    if total_source == 0 || total_target == 0 {
        return 50; // Default compatibility if no workload info
    }
    
    let diff = (total_source as i32 - total_target as i32).abs() as u32;
    let max_total = total_source.max(total_target);
    
    if diff * 100 / max_total <= 10 {
        100
    } else if diff * 100 / max_total <= 20 {
        85
    } else if diff * 100 / max_total <= 30 {
        70
    } else {
        50
    }
}

fn calculate_prerequisite_alignment(
    source_prereqs: &[String],
    target_prereqs: &[String],
) -> u32 {
    if source_prereqs.is_empty() && target_prereqs.is_empty() {
        return 100;
    }
    
    if source_prereqs.is_empty() || target_prereqs.is_empty() {
        return 80; // One has no prerequisites
    }
    
    let mut matches = 0;
    for prereq in source_prereqs {
        if target_prereqs.iter().any(|t| {
            crate::ipfs::calculate_semantic_similarity(prereq, t) > 70
        }) {
            matches += 1;
        }
    }
    
    let similarity = (matches * 100) / source_prereqs.len().max(target_prereqs.len());
    similarity as u32
}

fn calculate_learning_outcome_alignment(
    source_objectives: &[String],
    target_objectives: &[String],
) -> u32 {
    calculate_list_similarity(source_objectives, target_objectives)
}

fn calculate_bibliography_overlap(
    source_bib: &[String],
    target_bib: &[String],
) -> u32 {
    if source_bib.is_empty() && target_bib.is_empty() {
        return 50; // No bibliography info
    }
    
    if source_bib.is_empty() || target_bib.is_empty() {
        return 30;
    }
    
    let mut matches = 0;
    for bib1 in source_bib {
        if target_bib.iter().any(|bib2| {
            crate::ipfs::calculate_semantic_similarity(bib1, bib2) > 80
        }) {
            matches += 1;
        }
    }
    
    (matches * 100 / source_bib.len().max(target_bib.len())) as u32
}

fn calculate_structure_similarity(
    source: &SyllabusContent,
    target: &SyllabusContent,
) -> u32 {
    let mut score = 0u32;
    let mut factors = 0u32;
    
    // Compare number of objectives
    let obj_diff = (source.objectives.len() as i32 - target.objectives.len() as i32).abs();
    if obj_diff <= 2 {
        score += 30;
    } else if obj_diff <= 4 {
        score += 20;
    } else {
        score += 10;
    }
    factors += 30;
    
    // Compare number of topics
    let topic_diff = (source.topics.len() as i32 - target.topics.len() as i32).abs();
    if topic_diff <= 3 {
        score += 30;
    } else if topic_diff <= 6 {
        score += 20;
    } else {
        score += 10;
    }
    factors += 30;
    
    // Compare presence of prerequisites
    let prereq_score = match (source.prerequisites.is_empty(), target.prerequisites.is_empty()) {
        (true, true) | (false, false) => 20,
        _ => 10,
    };
    score += prereq_score;
    factors += 20;
    
    // Compare course duration
    match (&source.duration_weeks, &target.duration_weeks) {
        (Some(d1), Some(d2)) => {
            let diff = (*d1 as i32 - *d2 as i32).abs();
            if diff <= 2 {
                score += 20;
            } else if diff <= 4 {
                score += 15;
            } else {
                score += 5;
            }
        }
        (None, None) => score += 15,
        _ => score += 5,
    }
    factors += 20;
    
    (score * 100) / factors
}

// Enhanced comparison function
fn compare_multilingual_enhanced(
    source: &MultilingualSyllabusContent,
    target: &MultilingualSyllabusContent,
) -> StdResult<SimilarityAnalysis> {
    // Find best language pair for comparison
    let (source_content, target_content, source_lang, target_lang) = 
        find_best_language_pair(source, target);
    
    let weights = AnalysisWeights::default();
    
    // Calculate individual similarity scores with language awareness
    let title_sim = calculate_enhanced_semantic_similarity(
        &source_content.title,
        &target_content.title,
        &source_lang,
        &target_lang,
    );
    
    let desc_sim = calculate_enhanced_semantic_similarity(
        &source_content.description,
        &target_content.description,
        &source_lang,
        &target_lang,
    );
    
    // For lists, we need to adapt the function
    let obj_sim = calculate_list_similarity_enhanced(
        &source_content.objectives,
        &target_content.objectives,
        &source_lang,
        &target_lang,
    );
    
    let topics_sim = calculate_list_similarity_enhanced(
        &source_content.topics,
        &target_content.topics,
        &source_lang,
        &target_lang,
    );
    
    let bib_sim = calculate_list_similarity_enhanced(
        &source_content.bibliography,
        &target_content.bibliography,
        &source_lang,
        &target_lang,
    );
    
    let method_sim = calculate_enhanced_semantic_similarity(
        &source_content.methodology,
        &target_content.methodology,
        &source_lang,
        &target_lang,
    );
    
    let workload_sim = calculate_workload_similarity(
        &source_content.workload_breakdown,
        &target_content.workload_breakdown,
    );
    
    let keywords_sim = calculate_list_similarity_enhanced(
        &source_content.keywords,
        &target_content.keywords,
        &source_lang,
        &target_lang,
    );
    
    let structural_sim = calculate_structure_similarity(source_content, target_content);
    
    // Calculate weighted overall similarity
    let total_weight = weights.title_weight + weights.description_weight + 
                      weights.objectives_weight + weights.topics_weight +
                      weights.bibliography_weight + weights.methodology_weight +
                      weights.workload_weight + weights.keywords_weight;
    
    let weighted_score = (title_sim * weights.title_weight +
                         desc_sim * weights.description_weight +
                         obj_sim * weights.objectives_weight +
                         topics_sim * weights.topics_weight +
                         bib_sim * weights.bibliography_weight +
                         method_sim * weights.methodology_weight +
                         workload_sim * weights.workload_weight +
                         keywords_sim * weights.keywords_weight) / total_weight;
    
    Ok(SimilarityAnalysis {
        overall_similarity: weighted_score,
        title_similarity: title_sim,
        description_similarity: desc_sim,
        objectives_similarity: obj_sim,
        topics_similarity: topics_sim,
        bibliography_similarity: bib_sim,
        methodology_similarity: method_sim,
        workload_similarity: workload_sim,
        keywords_similarity: keywords_sim,
        structural_similarity: structural_sim,
    })
}

fn find_best_language_pair<'a>(
    source: &'a MultilingualSyllabusContent,
    target: &'a MultilingualSyllabusContent,
) -> (&'a SyllabusContent, &'a SyllabusContent, Language, Language) {
    // Priority 1: Same primary language
    if source.primary_language == target.primary_language {
        return (&source.content, &target.content, 
                source.primary_language.clone(), target.primary_language.clone());
    }
    
    // Priority 2: Check if target has translation in source's primary language
    for trans in &target.translations {
        if trans.language == source.primary_language {
            return (&source.content, &trans.content,
                    source.primary_language.clone(), trans.language.clone());
        }
    }
    
    // Priority 3: Check if source has translation in target's primary language
    for trans in &source.translations {
        if trans.language == target.primary_language {
            return (&trans.content, &target.content,
                    trans.language.clone(), target.primary_language.clone());
        }
    }
    
    // Priority 4: Find any common language in translations
    for source_trans in &source.translations {
        for target_trans in &target.translations {
            if source_trans.language == target_trans.language {
                return (&source_trans.content, &target_trans.content,
                        source_trans.language.clone(), target_trans.language.clone());
            }
        }
    }
    
    // Fallback: Use primary languages
    (&source.content, &target.content,
     source.primary_language.clone(), target.primary_language.clone())
}


fn calculate_list_similarity_enhanced(
    list1: &[String],
    list2: &[String],
    lang1: &Language,
    lang2: &Language,
) -> u32 {
    if list1.is_empty() && list2.is_empty() {
        return 100;
    }
    
    if list1.is_empty() || list2.is_empty() {
        return 0;
    }
    
    let mut total_similarity = 0u32;
    let mut comparison_count = 0u32;
    
    // Compare each item in list1 with best match in list2
    for item1 in list1 {
        let mut best_match = 0u32;
        for item2 in list2 {
            let similarity = calculate_enhanced_semantic_similarity(item1, item2, lang1, lang2);
            best_match = best_match.max(similarity);
        }
        total_similarity += best_match;
        comparison_count += 1;
    }
    
    // Compare each item in list2 with best match in list1
    for item2 in list2 {
        let mut best_match = 0u32;
        for item1 in list1 {
            let similarity = calculate_enhanced_semantic_similarity(item2, item1, lang2, lang1);
            best_match = best_match.max(similarity);
        }
        total_similarity += best_match;
        comparison_count += 1;
    }
    
    if comparison_count > 0 {
        total_similarity / comparison_count
    } else {
        0
    }
}

fn calculate_workload_similarity(
    w1: &WorkloadBreakdown,
    w2: &WorkloadBreakdown,
) -> u32 {
    let total1 = w1.theoretical_hours + w1.practical_hours + w1.laboratory_hours + 
                w1.field_work_hours + w1.seminar_hours + w1.individual_study_hours;
    let total2 = w2.theoretical_hours + w2.practical_hours + w2.laboratory_hours + 
                w2.field_work_hours + w2.seminar_hours + w2.individual_study_hours;
    
    if total1 == 0 && total2 == 0 {
        return 100;
    }
    
    if total1 == 0 || total2 == 0 {
        return 0;
    }
    
    // Calculate similarity for each component
    let theo_sim = calculate_hour_similarity(w1.theoretical_hours, w2.theoretical_hours);
    let prac_sim = calculate_hour_similarity(w1.practical_hours, w2.practical_hours);
    let lab_sim = calculate_hour_similarity(w1.laboratory_hours, w2.laboratory_hours);
    let field_sim = calculate_hour_similarity(w1.field_work_hours, w2.field_work_hours);
    let sem_sim = calculate_hour_similarity(w1.seminar_hours, w2.seminar_hours);
    let study_sim = calculate_hour_similarity(w1.individual_study_hours, w2.individual_study_hours);
    
    // Weighted average (theoretical and practical hours are more important)
    (theo_sim * 30 + prac_sim * 25 + lab_sim * 15 + 
     field_sim * 10 + sem_sim * 10 + study_sim * 10) / 100
}

fn calculate_hour_similarity(h1: u32, h2: u32) -> u32 {
    if h1 == h2 {
        return 100;
    }
    
    if h1 == 0 && h2 == 0 {
        return 100;
    }
    
    let diff = if h1 > h2 { h1 - h2 } else { h2 - h1 };
    let max_hours = h1.max(h2);
    
    if max_hours == 0 {
        return 100;
    }
    
    // Similarity decreases with relative difference
    let similarity = if diff * 100 / max_hours <= 20 {
        100 - (diff * 100 / max_hours) * 2
    } else {
        100 - diff * 100 / max_hours
    };
    
    similarity.max(0)
}

fn calculate_overall_score_with_quality(
    similarity: &SimilarityAnalysis,
    credit_compat: u32,
    level_compat: u32,
    workload_compat: u32,
    prereq_align: u32,
    source_quality: &ContentQualityAssessment,
    target_quality: &ContentQualityAssessment,
) -> u32 {
    // Base score from similarity analysis
    let base_score = (
        similarity.overall_similarity * 40 +
        credit_compat * 15 +
        level_compat * 15 +
        workload_compat * 10 +
        prereq_align * 10 +
        similarity.structural_similarity * 10
    ) / 100;
    
    // Apply quality penalty if content quality is low
    let avg_quality = (source_quality.overall_score + target_quality.overall_score) / 2;
    let quality_factor = if avg_quality >= 80 {
        100
    } else if avg_quality >= 60 {
        95
    } else if avg_quality >= 40 {
        90
    } else {
        80
    };
    
    (base_score * quality_factor) / 100
}

fn calculate_analysis_confidence(
    source_quality: &ContentQualityAssessment,
    target_quality: &ContentQualityAssessment,
    options: &AnalysisOptions,
) -> u32 {
    let base_confidence = (source_quality.overall_score + target_quality.overall_score) / 2;
    
    let option_boost = 
        if options.use_semantic_analysis { 10 } else { 0 } +
        if options.check_prerequisites { 5 } else { 0 } +
        if options.analyze_learning_outcomes { 5 } else { 0 };
    
    let translation_penalty = 
        (100 - source_quality.translation_quality_score).min(20) +
        (100 - target_quality.translation_quality_score).min(20);
    
    (base_confidence + option_boost).saturating_sub(translation_penalty).min(100)
}