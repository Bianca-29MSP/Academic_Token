use cosmwasm_std::{DepsMut, Env, MessageInfo, Response, StdResult};
use crate::error::ContractError;
use crate::state::{
    SubjectNFT, DegreeNFT, NFTMetadata, StudentNFTCollection, MintValidation,
    CONFIG, SUBJECT_NFTS, DEGREE_NFTS, STUDENT_COLLECTIONS, TOKEN_OWNERS,
    TOKEN_APPROVALS, OPERATOR_APPROVALS, NFT_COUNTER, NFT_STATS, IPFS_CACHE, VALIDATIONS,
    NFTType,
};
use crate::msg::{SubjectCompletionData, DegreeCompletionData};
// use sha2::{Sha256, Digest};

/// Mint subject completion NFT
pub fn execute_mint_subject_nft(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    student_id: String,
    subject_data: SubjectCompletionData,
    metadata: NFTMetadata,
    validation_hash: String,
) -> Result<Response, ContractError> {
    let config = CONFIG.load(deps.storage)?;
    
    // Only minter (AcademicNFT module) can mint NFTs
    if info.sender != config.minter && info.sender != config.admin {
        return Err(ContractError::Unauthorized {
            reason: "Only minter or admin can mint NFTs".to_string(),
        });
    }

    // Generate unique token ID
    let mut counter = NFT_COUNTER.load(deps.storage)?;
    counter += 1;
    let token_id = format!("subject_{}", counter);
    NFT_COUNTER.save(deps.storage, &counter)?;

    // Validate that this subject hasn't been completed before
    let existing_subject_key = format!("{}_{}", student_id, subject_data.subject_id);
    if let Ok(Some(_)) = SUBJECT_NFTS.may_load(deps.storage, &existing_subject_key) {
        return Err(ContractError::SubjectAlreadyCompleted {
            subject_id: subject_data.subject_id,
        });
    }

    // Create subject NFT
    let subject_nft = SubjectNFT {
        token_id: token_id.clone(),
        student_id: student_id.clone(),
        subject_id: subject_data.subject_id.clone(),
        institution_id: subject_data.institution_id.clone(),
        course_id: subject_data.course_id.clone(),
        subject_name: subject_data.subject_name,
        credits: subject_data.credits,
        final_grade: subject_data.final_grade,
        completion_date: subject_data.completion_date,
        semester: subject_data.semester,
        academic_year: subject_data.academic_year,
        instructor: subject_data.instructor,
        nft_metadata: metadata,
    };

    // Save NFT data
    SUBJECT_NFTS.save(deps.storage, &token_id, &subject_nft)?;
    
    // Set token owner
    let student_addr = deps.api.addr_validate(&student_id)?;
    TOKEN_OWNERS.save(deps.storage, &token_id, &student_addr)?;

    // Update student collection
    update_student_collection(deps.storage, &student_id, &token_id, NFTType::SubjectCompletion)?;

    // Update statistics
    update_statistics(deps.storage, NFTType::SubjectCompletion)?;

    // Save validation record
    let validation = MintValidation {
        validator_address: info.sender.to_string(),
        validation_hash,
        timestamp: env.block.time.seconds(),
        signature: "".to_string(), // Could be implemented for additional security
    };
    VALIDATIONS.save(deps.storage, &token_id, &validation)?;

    Ok(Response::new()
        .add_attribute("method", "mint_subject_nft")
        .add_attribute("token_id", token_id)
        .add_attribute("student_id", student_id)
        .add_attribute("subject_id", subject_data.subject_id))
}

/// Mint degree NFT
pub fn execute_mint_degree_nft(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    student_id: String,
    degree_data: DegreeCompletionData,
    metadata: NFTMetadata,
    validation_hash: String,
    signatures: Vec<String>,
) -> Result<Response, ContractError> {
    let config = CONFIG.load(deps.storage)?;
    
    // Only minter (AcademicNFT module) can mint NFTs
    if info.sender != config.minter && info.sender != config.admin {
        return Err(ContractError::Unauthorized {
            reason: "Only minter or admin can mint NFTs".to_string(),
        });
    }

    // Generate unique token ID
    let mut counter = NFT_COUNTER.load(deps.storage)?;
    counter += 1;
    let token_id = format!("degree_{}", counter);
    NFT_COUNTER.save(deps.storage, &counter)?;

    // Create degree NFT
    let degree_nft = DegreeNFT {
        token_id: token_id.clone(),
        student_id: student_id.clone(),
        institution_id: degree_data.institution_id.clone(),
        course_id: degree_data.course_id.clone(),
        degree_name: degree_data.degree_name,
        degree_type: degree_data.degree_type,
        major: degree_data.major,
        minor: degree_data.minor,
        graduation_date: degree_data.graduation_date,
        final_gpa: degree_data.final_gpa,
        total_credits: degree_data.total_credits,
        honors: degree_data.honors,
        validation_hash: validation_hash.clone(),
        signatures,
        nft_metadata: metadata,
    };

    // Save NFT data
    DEGREE_NFTS.save(deps.storage, &token_id, &degree_nft)?;
    
    // Set token owner
    let student_addr = deps.api.addr_validate(&student_id)?;
    TOKEN_OWNERS.save(deps.storage, &token_id, &student_addr)?;

    // Update student collection
    update_student_collection(deps.storage, &student_id, &token_id, NFTType::Degree)?;

    // Update statistics
    update_statistics(deps.storage, NFTType::Degree)?;

    // Save validation record
    let validation = MintValidation {
        validator_address: info.sender.to_string(),
        validation_hash,
        timestamp: env.block.time.seconds(),
        signature: "".to_string(),
    };
    VALIDATIONS.save(deps.storage, &token_id, &validation)?;

    Ok(Response::new()
        .add_attribute("method", "mint_degree_nft")
        .add_attribute("token_id", token_id)
        .add_attribute("student_id", student_id)
        .add_attribute("institution_id", degree_data.institution_id))
}

/// Transfer NFT
pub fn execute_transfer_nft(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    recipient: String,
    token_id: String,
) -> Result<Response, ContractError> {
    let owner = TOKEN_OWNERS.load(deps.storage, &token_id)
        .map_err(|_| ContractError::NFTNotFound { token_id: token_id.clone() })?;

    // Check if sender is owner or approved
    if info.sender != owner {
        // Check for specific token approval
        if let Some(approved) = TOKEN_APPROVALS.may_load(deps.storage, &token_id)? {
            if approved != info.sender {
                // Check for operator approval
                if !OPERATOR_APPROVALS.may_load(deps.storage, (owner.as_str(), info.sender.as_str()))?.unwrap_or(false) {
                    return Err(ContractError::Unauthorized {
                        reason: "Not owner or approved".to_string(),
                    });
                }
            }
        }
    }

    let recipient_addr = deps.api.addr_validate(&recipient)?;
    TOKEN_OWNERS.save(deps.storage, &token_id, &recipient_addr)?;

    // Clear any existing approvals for this token
    TOKEN_APPROVALS.remove(deps.storage, &token_id);

    Ok(Response::new()
        .add_attribute("method", "transfer_nft")
        .add_attribute("token_id", token_id)
        .add_attribute("from", owner)
        .add_attribute("to", recipient))
}

/// Approve spender for specific token
pub fn execute_approve(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    spender: String,
    token_id: String,
    _expires: Option<cw721::Expiration>,
) -> Result<Response, ContractError> {
    let owner = TOKEN_OWNERS.load(deps.storage, &token_id)
        .map_err(|_| ContractError::NFTNotFound { token_id: token_id.clone() })?;

    if info.sender != owner {
        return Err(ContractError::Unauthorized {
            reason: "Only owner can approve".to_string(),
        });
    }

    let spender_addr = deps.api.addr_validate(&spender)?;
    TOKEN_APPROVALS.save(deps.storage, &token_id, &spender_addr)?;

    Ok(Response::new()
        .add_attribute("method", "approve")
        .add_attribute("token_id", token_id)
        .add_attribute("spender", spender))
}

/// Revoke approval for specific token
pub fn execute_revoke(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    _spender: String,
    token_id: String,
) -> Result<Response, ContractError> {
    let owner = TOKEN_OWNERS.load(deps.storage, &token_id)
        .map_err(|_| ContractError::NFTNotFound { token_id: token_id.clone() })?;

    if info.sender != owner {
        return Err(ContractError::Unauthorized {
            reason: "Only owner can revoke".to_string(),
        });
    }

    TOKEN_APPROVALS.remove(deps.storage, &token_id);

    Ok(Response::new()
        .add_attribute("method", "revoke")
        .add_attribute("token_id", token_id))
}

/// Approve all tokens for operator
pub fn execute_approve_all(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    operator: String,
    _expires: Option<cw721::Expiration>,
) -> Result<Response, ContractError> {
    let operator_addr = deps.api.addr_validate(&operator)?;
    OPERATOR_APPROVALS.save(deps.storage, (info.sender.as_str(), operator_addr.as_str()), &true)?;

    Ok(Response::new()
        .add_attribute("method", "approve_all")
        .add_attribute("owner", info.sender)
        .add_attribute("operator", operator))
}

/// Revoke all approvals for operator
pub fn execute_revoke_all(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    operator: String,
) -> Result<Response, ContractError> {
    let operator_addr = deps.api.addr_validate(&operator)?;
    OPERATOR_APPROVALS.remove(deps.storage, (info.sender.as_str(), operator_addr.as_str()));

    Ok(Response::new()
        .add_attribute("method", "revoke_all")
        .add_attribute("owner", info.sender)
        .add_attribute("operator", operator))
}

/// Update NFT metadata (admin only)
pub fn execute_update_metadata(
    deps: DepsMut,
    info: MessageInfo,
    token_id: String,
    metadata: NFTMetadata,
) -> Result<Response, ContractError> {
    let config = CONFIG.load(deps.storage)?;
    
    if info.sender != config.admin {
        return Err(ContractError::Unauthorized {
            reason: "Only admin can update metadata".to_string(),
        });
    }

    // Update subject NFT metadata if it exists
    if let Ok(mut subject_nft) = SUBJECT_NFTS.load(deps.storage, &token_id) {
        subject_nft.nft_metadata = metadata;
        SUBJECT_NFTS.save(deps.storage, &token_id, &subject_nft)?;
    }
    // Update degree NFT metadata if it exists
    else if let Ok(mut degree_nft) = DEGREE_NFTS.load(deps.storage, &token_id) {
        degree_nft.nft_metadata = metadata;
        DEGREE_NFTS.save(deps.storage, &token_id, &degree_nft)?;
    } else {
        return Err(ContractError::NFTNotFound { token_id });
    }

    Ok(Response::new()
        .add_attribute("method", "update_metadata")
        .add_attribute("token_id", token_id))
}

/// Burn NFT (admin only)
pub fn execute_burn(
    deps: DepsMut,
    info: MessageInfo,
    token_id: String,
    reason: String,
) -> Result<Response, ContractError> {
    let config = CONFIG.load(deps.storage)?;
    
    if info.sender != config.admin {
        return Err(ContractError::Unauthorized {
            reason: "Only admin can burn NFTs".to_string(),
        });
    }

    // Remove from storage
    SUBJECT_NFTS.remove(deps.storage, &token_id);
    DEGREE_NFTS.remove(deps.storage, &token_id);
    TOKEN_OWNERS.remove(deps.storage, &token_id);
    TOKEN_APPROVALS.remove(deps.storage, &token_id);
    VALIDATIONS.remove(deps.storage, &token_id);

    Ok(Response::new()
        .add_attribute("method", "burn")
        .add_attribute("token_id", token_id)
        .add_attribute("reason", reason))
}

/// Update contract configuration
pub fn execute_update_config(
    deps: DepsMut,
    info: MessageInfo,
    admin: Option<String>,
    minter: Option<String>,
    ipfs_gateway: Option<String>,
) -> Result<Response, ContractError> {
    let mut config = CONFIG.load(deps.storage)?;
    
    if info.sender != config.admin {
        return Err(ContractError::Unauthorized {
            reason: "Only admin can update config".to_string(),
        });
    }

    if let Some(new_admin) = admin {
        config.admin = deps.api.addr_validate(&new_admin)?;
    }
    if let Some(new_minter) = minter {
        config.minter = deps.api.addr_validate(&new_minter)?;
    }
    if let Some(new_gateway) = ipfs_gateway {
        config.ipfs_gateway = new_gateway;
    }

    CONFIG.save(deps.storage, &config)?;

    Ok(Response::new()
        .add_attribute("method", "update_config"))
}

/// Cache IPFS content
pub fn execute_cache_ipfs_content(
    deps: DepsMut,
    ipfs_link: String,
    content: String,
) -> Result<Response, ContractError> {
    IPFS_CACHE.save(deps.storage, &ipfs_link, &content)?;

    Ok(Response::new()
        .add_attribute("method", "cache_ipfs_content")
        .add_attribute("ipfs_link", ipfs_link))
}

/// Validate NFT data before minting
pub fn execute_validate_nft_data(
    _deps: DepsMut,
    nft_type: NFTType,
    _data: String,
) -> Result<Response, ContractError> {
    // Perform validation logic here
    // This could include checking data format, required fields, etc.
    
    match nft_type {
        NFTType::SubjectCompletion => {
            // Validate subject completion data
            // Check if all required fields are present
        }
        NFTType::Degree => {
            // Validate degree data
            // Check if GPA is valid, credits are sufficient, etc.
        }
        _ => {
            return Err(ContractError::InvalidNFTType {
                nft_type: format!("{:?}", nft_type),
            });
        }
    }

    Ok(Response::new()
        .add_attribute("method", "validate_nft_data")
        .add_attribute("nft_type", format!("{:?}", nft_type))
        .add_attribute("is_valid", "true"))
}

// Helper functions

/// Update student collection when new NFT is minted
fn update_student_collection(
    storage: &mut dyn cosmwasm_std::Storage,
    student_id: &str,
    token_id: &str,
    nft_type: NFTType,
) -> StdResult<()> {
    let mut collection = STUDENT_COLLECTIONS
        .may_load(storage, student_id)?
        .unwrap_or_else(|| StudentNFTCollection {
            student_id: student_id.to_string(),
            subject_nfts: vec![],
            degree_nfts: vec![],
            total_credits: 0,
            current_gpa: "0.0".to_string(),
            completion_percentage: 0,
        });

    match nft_type {
        NFTType::SubjectCompletion => {
            collection.subject_nfts.push(token_id.to_string());
        }
        NFTType::Degree => {
            collection.degree_nfts.push(token_id.to_string());
        }
        _ => {}
    }

    STUDENT_COLLECTIONS.save(storage, student_id, &collection)?;
    Ok(())
}

/// Update global statistics
fn update_statistics(
    storage: &mut dyn cosmwasm_std::Storage,
    nft_type: NFTType,
) -> StdResult<()> {
    let mut stats = NFT_STATS.load(storage)?;

    match nft_type {
        NFTType::SubjectCompletion => {
            stats.total_subject_nfts += 1;
        }
        NFTType::Degree => {
            stats.total_degree_nfts += 1;
        }
        _ => {}
    }

    NFT_STATS.save(storage, &stats)?;
    Ok(())
}
