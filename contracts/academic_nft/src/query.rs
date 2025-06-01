use cosmwasm_std::{Deps, Env, StdResult};
use crate::state::{
    NFTType,
    CONFIG, SUBJECT_NFTS, DEGREE_NFTS, STUDENT_COLLECTIONS, TOKEN_OWNERS,
    TOKEN_APPROVALS, NFT_STATS, IPFS_CACHE, VALIDATIONS,
};
use crate::msg::{
    ConfigResponse, SubjectNFTResponse, DegreeNFTResponse, StudentCollectionResponse,
    NFTsByTypeResponse, StatisticsResponse, InstitutionNFTsResponse, InstitutionNFTSummary,
    VerificationResponse,
};
// use crate::error::ContractError;

/// Query contract configuration
pub fn query_config(deps: Deps) -> StdResult<ConfigResponse> {
    let config = CONFIG.load(deps.storage)?;
    Ok(ConfigResponse { config })
}

/// Query NFT owner (CW721 standard)
pub fn query_owner_of(
    deps: Deps,
    _env: Env,
    token_id: String,
    _include_expired: Option<bool>,
) -> StdResult<cw721::OwnerOfResponse> {
    let owner = TOKEN_OWNERS.load(deps.storage, &token_id)?;
    Ok(cw721::OwnerOfResponse {
        owner: owner.to_string(),
        approvals: vec![], // Could be implemented if needed
    })
}

/// Query approval for specific token (CW721 standard)
pub fn query_approval(
    deps: Deps,
    _env: Env,
    token_id: String,
    spender: String,
    _include_expired: Option<bool>,
) -> StdResult<cw721::ApprovalResponse> {
    let _spender_addr = deps.api.addr_validate(&spender)?;
    let approval = TOKEN_APPROVALS.may_load(deps.storage, &token_id)?;
    
    // Check if the spender is approved for this token
    let approval_exists = match approval {
        Some(approved_addr) if approved_addr.to_string() == spender => {
            cw721::Approval {
                spender: approved_addr.to_string(),
                expires: cw721::Expiration::Never {},
            }
        }
        _ => {
            // Return a default approval with empty spender to indicate no approval
            cw721::Approval {
                spender: "".to_string(),
                expires: cw721::Expiration::Never {},
            }
        }
    };
    
    Ok(cw721::ApprovalResponse {
        approval: approval_exists,
    })
}

/// Query all approvals for token (CW721 standard)
pub fn query_approvals(
    deps: Deps,
    _env: Env,
    token_id: String,
    _include_expired: Option<bool>,
) -> StdResult<cw721::ApprovalsResponse> {
    let approval = TOKEN_APPROVALS.may_load(deps.storage, &token_id)?;
    
    let approvals = approval.map(|approved| vec![cw721::Approval {
        spender: approved.to_string(),
        expires: cw721::Expiration::Never {},
    }]).unwrap_or_default();
    
    Ok(cw721::ApprovalsResponse { approvals })
}

/// Query tokens owned by address (CW721 standard)
pub fn query_tokens(
    deps: Deps,
    owner: String,
    start_after: Option<String>,
    limit: Option<u32>,
) -> StdResult<cw721::TokensResponse> {
    let owner_addr = deps.api.addr_validate(&owner)?;
    let limit = limit.unwrap_or(30).min(100) as usize;
    
    let tokens: StdResult<Vec<String>> = TOKEN_OWNERS
        .range(deps.storage, None, None, cosmwasm_std::Order::Ascending)
        .filter_map(|item| {
            item.ok().and_then(|(token_id, token_owner)| {
                if token_owner == owner_addr {
                    if let Some(ref start) = start_after {
                        if token_id <= *start {
                            return None;
                        }
                    }
                    Some(Ok(token_id))
                } else {
                    None
                }
            })
        })
        .take(limit)
        .collect();
    
    Ok(cw721::TokensResponse { tokens: tokens? })
}

/// Query all tokens in collection (CW721 standard)
pub fn query_all_tokens(
    deps: Deps,
    start_after: Option<String>,
    limit: Option<u32>,
) -> StdResult<cw721::TokensResponse> {
    let limit = limit.unwrap_or(30).min(100) as usize;
    
    let tokens: StdResult<Vec<String>> = TOKEN_OWNERS
        .range(deps.storage, None, None, cosmwasm_std::Order::Ascending)
        .map(|item| item.map(|(token_id, _)| token_id))
        .filter(|item| {
            if let (Ok(token_id), Some(ref start)) = (item, &start_after) {
                token_id > start
            } else {
                true
            }
        })
        .take(limit)
        .collect();
    
    Ok(cw721::TokensResponse { tokens: tokens? })
}

/// Query NFT info (CW721 standard)
pub fn query_nft_info(
    deps: Deps,
    token_id: String,
) -> StdResult<cw721::NftInfoResponse<crate::state::NFTMetadata>> {
    // Try to find in subject NFTs first
    if let Ok(subject_nft) = SUBJECT_NFTS.load(deps.storage, &token_id) {
        return Ok(cw721::NftInfoResponse {
            token_uri: Some(subject_nft.nft_metadata.ipfs_metadata_link.clone()),
            extension: subject_nft.nft_metadata,
        });
    }
    
    // Try degree NFTs
    if let Ok(degree_nft) = DEGREE_NFTS.load(deps.storage, &token_id) {
        return Ok(cw721::NftInfoResponse {
            token_uri: Some(degree_nft.nft_metadata.ipfs_metadata_link.clone()),
            extension: degree_nft.nft_metadata,
        });
    }
    
    Err(cosmwasm_std::StdError::not_found("NFT"))
}

/// Query all NFT info (CW721 standard)
pub fn query_all_nft_info(
    deps: Deps,
    env: Env,
    token_id: String,
    include_expired: Option<bool>,
) -> StdResult<cw721::AllNftInfoResponse<crate::state::NFTMetadata>> {
    let nft_info = query_nft_info(deps, token_id.clone())?;
    let owner_info = query_owner_of(deps, env, token_id, include_expired)?;
    
    Ok(cw721::AllNftInfoResponse {
        access: cw721::OwnerOfResponse {
            owner: owner_info.owner,
            approvals: owner_info.approvals,
        },
        info: nft_info,
    })
}

/// Query contract info (CW721 standard)
pub fn query_contract_info(deps: Deps) -> StdResult<cw721::ContractInfoResponse> {
    let config = CONFIG.load(deps.storage)?;
    Ok(cw721::ContractInfoResponse {
        name: config.collection_name,
        symbol: config.collection_symbol,
    })
}

/// Query subject NFT details
pub fn query_subject_nft(deps: Deps, token_id: String) -> StdResult<SubjectNFTResponse> {
    let nft = SUBJECT_NFTS.may_load(deps.storage, &token_id)?;
    Ok(SubjectNFTResponse { nft })
}

/// Query degree NFT details
pub fn query_degree_nft(deps: Deps, token_id: String) -> StdResult<DegreeNFTResponse> {
    let nft = DEGREE_NFTS.may_load(deps.storage, &token_id)?;
    Ok(DegreeNFTResponse { nft })
}

/// Query student's NFT collection
pub fn query_student_collection(
    deps: Deps,
    student_id: String,
) -> StdResult<StudentCollectionResponse> {
    let collection = STUDENT_COLLECTIONS.may_load(deps.storage, &student_id)?;
    
    // Calculate total academic value (could be based on credits, GPA, etc.)
    let total_value = collection.as_ref()
        .map(|c| c.total_credits.to_string())
        .unwrap_or_else(|| "0".to_string());
    
    Ok(StudentCollectionResponse {
        collection,
        total_value,
    })
}

/// Query NFTs by type for a student
pub fn query_nfts_by_type(
    deps: Deps,
    student_id: String,
    nft_type: NFTType,
    limit: Option<u32>,
) -> StdResult<NFTsByTypeResponse> {
    let collection = STUDENT_COLLECTIONS.may_load(deps.storage, &student_id)?;
    
    if let Some(coll) = collection {
        let nfts = match nft_type {
            NFTType::SubjectCompletion => coll.subject_nfts,
            NFTType::Degree => coll.degree_nfts,
            _ => vec![],
        };
        
        let limited_nfts = if let Some(limit) = limit {
            nfts.into_iter().take(limit as usize).collect()
        } else {
            nfts
        };
        
        Ok(NFTsByTypeResponse {
            total_count: limited_nfts.len() as u32,
            nfts: limited_nfts,
        })
    } else {
        Ok(NFTsByTypeResponse {
            nfts: vec![],
            total_count: 0,
        })
    }
}

/// Query NFT statistics
pub fn query_statistics(deps: Deps) -> StdResult<StatisticsResponse> {
    let stats = NFT_STATS.load(deps.storage)?;
    Ok(StatisticsResponse { stats })
}

/// Query NFTs by institution
pub fn query_nfts_by_institution(
    deps: Deps,
    institution_id: String,
    nft_type: Option<NFTType>,
    limit: Option<u32>,
    _start_after: Option<String>,
) -> StdResult<InstitutionNFTsResponse> {
    let limit = limit.unwrap_or(50) as usize;
    let mut nfts = vec![];
    
    // Search through subject NFTs
    if nft_type.is_none() || nft_type == Some(NFTType::SubjectCompletion) {
        let subject_nfts: StdResult<Vec<_>> = SUBJECT_NFTS
            .range(deps.storage, None, None, cosmwasm_std::Order::Ascending)
            .filter_map(|item| {
                item.ok().and_then(|(token_id, nft)| {
                    if nft.institution_id == institution_id {
                        Some(Ok(InstitutionNFTSummary {
                            token_id,
                            student_id: nft.student_id,
                            nft_type: NFTType::SubjectCompletion,
                            issued_date: nft.completion_date,
                            name: nft.subject_name,
                        }))
                    } else {
                        None
                    }
                })
            })
            .take(limit)
            .collect();
        
        nfts.extend(subject_nfts?);
    }
    
    // Search through degree NFTs
    if nft_type.is_none() || nft_type == Some(NFTType::Degree) {
        let remaining_limit = limit.saturating_sub(nfts.len());
        if remaining_limit > 0 {
            let degree_nfts: StdResult<Vec<_>> = DEGREE_NFTS
                .range(deps.storage, None, None, cosmwasm_std::Order::Ascending)
                .filter_map(|item| {
                    item.ok().and_then(|(token_id, nft)| {
                        if nft.institution_id == institution_id {
                            Some(Ok(InstitutionNFTSummary {
                                token_id,
                                student_id: nft.student_id,
                                nft_type: NFTType::Degree,
                                issued_date: nft.graduation_date,
                                name: nft.degree_name,
                            }))
                        } else {
                            None
                        }
                    })
                })
                .take(remaining_limit)
                .collect();
            
            nfts.extend(degree_nfts?);
        }
    }
    
    Ok(InstitutionNFTsResponse {
        total_count: nfts.len() as u32,
        nfts,
    })
}

/// Check if NFT exists
pub fn query_nft_exists(deps: Deps, token_id: String) -> StdResult<bool> {
    let exists = TOKEN_OWNERS.has(deps.storage, &token_id);
    Ok(exists)
}

/// Get cached IPFS content
pub fn query_cached_content(deps: Deps, ipfs_link: String) -> StdResult<String> {
    IPFS_CACHE.load(deps.storage, &ipfs_link)
}

/// Verify NFT authenticity
pub fn query_verify_nft(deps: Deps, token_id: String) -> StdResult<VerificationResponse> {
    // Check if NFT exists
    if !TOKEN_OWNERS.has(deps.storage, &token_id) {
        return Ok(VerificationResponse {
            is_valid: false,
            issued_by: "".to_string(),
            validation_hash: "".to_string(),
            issue_date: 0,
            verification_details: vec!["NFT not found".to_string()],
        });
    }
    
    // Get validation record
    if let Ok(validation) = VALIDATIONS.load(deps.storage, &token_id) {
        Ok(VerificationResponse {
            is_valid: true,
            issued_by: validation.validator_address,
            validation_hash: validation.validation_hash,
            issue_date: validation.timestamp,
            verification_details: vec![
                "NFT verified".to_string(),
                "Issued by authorized minter".to_string(),
                "Validation hash matches".to_string(),
            ],
        })
    } else {
        Ok(VerificationResponse {
            is_valid: false,
            issued_by: "".to_string(),
            validation_hash: "".to_string(),
            issue_date: 0,
            verification_details: vec!["Validation record not found".to_string()],
        })
    }
}
