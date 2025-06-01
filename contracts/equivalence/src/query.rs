use cosmwasm_std::{Deps, StdResult, StdError, Order};
use cw_storage_plus::Bound;
use crate::state::{
    STATE, EQUIVALENCES, EQUIVALENCE_INDEX, ANALYSIS_RESULTS,
    DETAILED_ANALYSIS_RESULTS, TRANSFER_REQUESTS, STUDENT_TRANSFERS,
    Equivalence, DetailedAnalysisResult, TransferRequest,
};
use crate::msg::{
    StateResponse, EquivalenceResponse, EquivalencesResponse,
    AnalysisResponse, TransferResponse, TransfersResponse,
    EquivalenceCheckResponse, StatisticsResponse,
};
use crate::ipfs::{fetch_ipfs_content, assess_content_quality_multilingual};

pub fn query_state(deps: Deps) -> StdResult<StateResponse> {
    let state = STATE.load(deps.storage)?;
    Ok(StateResponse { state })
}

pub fn query_equivalence(deps: Deps, equivalence_id: String) -> StdResult<EquivalenceResponse> {
    let equivalence = EQUIVALENCES.may_load(deps.storage, &equivalence_id)?;
    Ok(EquivalenceResponse { equivalence })
}

pub fn query_find_equivalence(
    deps: Deps,
    source_subject_id: String,
    target_subject_id: String,
) -> StdResult<EquivalenceResponse> {
    let equiv_key = (source_subject_id.as_str(), target_subject_id.as_str());
    let equivalence_id = EQUIVALENCE_INDEX.may_load(deps.storage, equiv_key)?;
    
    let equivalence = match equivalence_id {
        Some(id) => EQUIVALENCES.may_load(deps.storage, &id)?,
        None => None,
    };
    
    Ok(EquivalenceResponse { equivalence })
}

pub fn query_list_equivalences_by_institution(
    deps: Deps,
    institution_id: String,
    limit: Option<u32>,
    start_after: Option<String>,
) -> StdResult<EquivalencesResponse> {
    let limit = limit.unwrap_or(30).min(100) as usize;
    
    // Fix the Bound type issue
    let start_bound = start_after.map(|s| Bound::ExclusiveRaw(s.into_bytes()));
    
    let equivalences: Vec<Equivalence> = EQUIVALENCES
        .range(deps.storage, start_bound, None, Order::Ascending)
        .take(limit)
        .filter_map(|item| {
            item.ok().and_then(|(_, equiv)| {
                if equiv.source_subject.institution_id == institution_id ||
                   equiv.target_subject.institution_id == institution_id {
                    Some(equiv)
                } else {
                    None
                }
            })
        })
        .collect();
    
    let total = equivalences.len() as u64;
    
    Ok(EquivalencesResponse {
        equivalences,
        total,
    })
}

pub fn query_analysis_result(deps: Deps, analysis_id: String) -> StdResult<AnalysisResponse> {
    let analysis = ANALYSIS_RESULTS.may_load(deps.storage, &analysis_id)?;
    Ok(AnalysisResponse { analysis })
}

pub fn query_detailed_analysis_result(
    deps: Deps,
    analysis_id: String,
) -> StdResult<DetailedAnalysisResult> {
    DETAILED_ANALYSIS_RESULTS
        .load(deps.storage, &analysis_id)
}

pub fn query_transfer_request(deps: Deps, transfer_id: String) -> StdResult<TransferResponse> {
    let transfer = TRANSFER_REQUESTS.may_load(deps.storage, &transfer_id)?;
    Ok(TransferResponse { transfer })
}

pub fn query_list_student_transfers(
    deps: Deps,
    student_id: String,
    limit: Option<u32>,
) -> StdResult<TransfersResponse> {
    let transfer_ids = STUDENT_TRANSFERS
        .may_load(deps.storage, &student_id)?
        .unwrap_or_default();
    
    let limit = limit.unwrap_or(30).min(100) as usize;
    
    let transfers: Vec<TransferRequest> = transfer_ids
        .into_iter()
        .take(limit)
        .filter_map(|id| TRANSFER_REQUESTS.may_load(deps.storage, &id).ok().flatten())
        .collect();
    
    let total = transfers.len() as u64;
    
    Ok(TransfersResponse {
        transfers,
        total,
    })
}

pub fn query_check_equivalence(
    deps: Deps,
    source_subject_id: String,
    target_subject_id: String,
    minimum_similarity: Option<u32>,
) -> StdResult<EquivalenceCheckResponse> {
    let equiv_key = (source_subject_id.as_str(), target_subject_id.as_str());
    let equivalence_id = EQUIVALENCE_INDEX.may_load(deps.storage, equiv_key)?;
    
    match equivalence_id {
        Some(id) => {
            let equivalence = EQUIVALENCES.load(deps.storage, &id)?;
            let min_sim = minimum_similarity.unwrap_or(70);
            
            Ok(EquivalenceCheckResponse {
                is_equivalent: equivalence.similarity_percentage >= min_sim,
                equivalence_type: Some(equivalence.equivalence_type),
                similarity_percentage: equivalence.similarity_percentage,
                equivalence_id: Some(id),
            })
        }
        None => Ok(EquivalenceCheckResponse {
            is_equivalent: false,
            equivalence_type: None,
            similarity_percentage: 0,
            equivalence_id: None,
        }),
    }
}

pub fn query_statistics(
    deps: Deps,
    institution_id: Option<String>,
) -> StdResult<StatisticsResponse> {
    let _state = STATE.load(deps.storage)?;
    
    let mut total_equivalences = 0u64;
    let mut approved_equivalences = 0u64;
    let mut pending_equivalences = 0u64;
    let mut total_similarity = 0u64;
    let mut count_with_similarity = 0u64;
    
    // Count equivalences
    for result in EQUIVALENCES.range(deps.storage, None, None, Order::Ascending) {
        if let Ok((_, equiv)) = result {
            // Filter by institution if specified
            if let Some(ref inst_id) = institution_id {
                if equiv.source_subject.institution_id != *inst_id &&
                   equiv.target_subject.institution_id != *inst_id {
                    continue;
                }
            }
            
            total_equivalences += 1;
            
            match equiv.status {
                crate::state::EquivalenceStatus::Approved => approved_equivalences += 1,
                crate::state::EquivalenceStatus::Pending => pending_equivalences += 1,
                _ => {}
            }
            
            if equiv.similarity_percentage > 0 {
                total_similarity += equiv.similarity_percentage as u64;
                count_with_similarity += 1;
            }
        }
    }
    
    let average_similarity = if count_with_similarity > 0 {
        (total_similarity / count_with_similarity) as u32
    } else {
        0
    };
    
    Ok(StatisticsResponse {
        total_equivalences,
        approved_equivalences,
        pending_equivalences,
        average_similarity,
        total_transfers: 0, // TODO: Count actual transfers
        successful_transfers: 0, // TODO: Count successful transfers
    })
}

pub fn query_debug_equivalences(deps: Deps) -> StdResult<Vec<String>> {
    let mut ids = Vec::new();
    
    for result in EQUIVALENCES.range(deps.storage, None, None, Order::Ascending) {
        if let Ok((_, equiv)) = result {
            ids.push(equiv.id);
        }
    }
    
    Ok(ids)
}

// Add the new enhanced query function
pub fn query_enhanced_analysis_with_recommendations(
    deps: Deps,
    equivalence_id: String,
) -> StdResult<crate::msg::EnhancedAnalysisResponse> {
    let analysis_id = format!("detailed_analysis_{}", equivalence_id);
    
    let detailed_result = DETAILED_ANALYSIS_RESULTS
        .may_load(deps.storage, &analysis_id)?
        .ok_or_else(|| StdError::not_found("Analysis result"))?;
    
    let equivalence = EQUIVALENCES
        .may_load(deps.storage, &equivalence_id)?
        .ok_or_else(|| StdError::not_found("Equivalence"))?;
    
    // Try to load content quality assessments if available
    let (source_quality, target_quality) = match (
        fetch_ipfs_content(deps, &equivalence.source_subject.ipfs_link),
        fetch_ipfs_content(deps, &equivalence.target_subject.ipfs_link),
    ) {
        (Ok(source), Ok(target)) => {
            let source_assessment = assess_content_quality_multilingual(&source);
            let target_assessment = assess_content_quality_multilingual(&target);
            (Some(source_assessment), Some(target_assessment))
        }
        _ => (None, None),
    };
    
    // Store quality metrics locally before moving
    let quality_threshold_met = detailed_result.quality_metrics.overall_quality >= 70;
    let content_available = source_quality.is_some() && target_quality.is_some();
    
    Ok(crate::msg::EnhancedAnalysisResponse {
        equivalence_id,
        analysis_result: detailed_result.base_result,
        quality_metrics: detailed_result.quality_metrics,
        language_compatibility: detailed_result.language_compatibility,
        recommendations: detailed_result.recommendations,
        source_quality_assessment: source_quality,
        target_quality_assessment: target_quality,
        confidence_factors: crate::msg::ConfidenceFactors {
            content_available,
            languages_match: detailed_result.language_compatibility == 100,
            quality_threshold_met,
            auto_translated: false, // Would need to check actual content
            comprehensive_analysis: true,
        },
    })
}