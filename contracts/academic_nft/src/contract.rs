#[cfg(not(feature = "library"))]
use cosmwasm_std::entry_point;
use cosmwasm_std::{
    to_json_binary, Binary, Deps, DepsMut, Env, MessageInfo, Response, StdResult,
};
use cw2::set_contract_version;

use crate::error::ContractError;
use crate::msg::{ExecuteMsg, InstantiateMsg, QueryMsg};
use crate::state::{Config, NFTStatistics, CONFIG, NFT_STATS, NFT_COUNTER};
use crate::execute::{
    execute_mint_subject_nft, execute_mint_degree_nft, execute_transfer_nft,
    execute_approve, execute_revoke, execute_approve_all, execute_revoke_all,
    execute_update_metadata, execute_burn, execute_update_config,
    execute_cache_ipfs_content, execute_validate_nft_data,
};
use crate::query::{
    query_config, query_owner_of, query_approval, query_approvals,
    query_tokens, query_all_tokens, query_nft_info, query_all_nft_info,
    query_contract_info, query_subject_nft, query_degree_nft,
    query_student_collection, query_nfts_by_type, query_statistics,
    query_nfts_by_institution, query_nft_exists, query_cached_content,
    query_verify_nft,
};

// Contract name and version
const CONTRACT_NAME: &str = "crates.io:academic-nft";
const CONTRACT_VERSION: &str = env!("CARGO_PKG_VERSION");

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn instantiate(
    deps: DepsMut,
    _env: Env,
    info: MessageInfo,
    msg: InstantiateMsg,
) -> Result<Response, ContractError> {
    set_contract_version(deps.storage, CONTRACT_NAME, CONTRACT_VERSION)?;

    let admin = msg.admin
        .map(|a| deps.api.addr_validate(&a))
        .transpose()?
        .unwrap_or(info.sender);

    let minter = deps.api.addr_validate(&msg.minter)?;

    let config = Config {
        admin: admin.clone(),
        minter,
        ipfs_gateway: msg.ipfs_gateway,
        collection_name: msg.collection_name,
        collection_symbol: msg.collection_symbol,
    };

    CONFIG.save(deps.storage, &config)?;
    
    // Initialize NFT counter
    NFT_COUNTER.save(deps.storage, &0u64)?;
    
    // Initialize statistics
    let initial_stats = NFTStatistics {
        total_subject_nfts: 0,
        total_degree_nfts: 0,
        total_students: 0,
        total_institutions: 0,
        nfts_by_type: vec![],
    };
    NFT_STATS.save(deps.storage, &initial_stats)?;

    Ok(Response::new()
        .add_attribute("method", "instantiate")
        .add_attribute("admin", admin)
        .add_attribute("minter", config.minter)
        .add_attribute("collection_name", config.collection_name))
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn execute(
    deps: DepsMut,
    env: Env,
    info: MessageInfo,
    msg: ExecuteMsg,
) -> Result<Response, ContractError> {
    match msg {
        ExecuteMsg::MintSubjectNFT {
            student_id,
            subject_data,
            metadata,
            validation_hash,
        } => execute_mint_subject_nft(deps, env, info, student_id, subject_data, metadata, validation_hash),
        
        ExecuteMsg::MintDegreeNFT {
            student_id,
            degree_data,
            metadata,
            validation_hash,
            signatures,
        } => execute_mint_degree_nft(deps, env, info, student_id, degree_data, metadata, validation_hash, signatures),
        
        ExecuteMsg::TransferNft {
            recipient,
            token_id,
        } => execute_transfer_nft(deps, env, info, recipient, token_id),
        
        ExecuteMsg::Approve {
            spender,
            token_id,
            expires,
        } => execute_approve(deps, env, info, spender, token_id, expires),
        
        ExecuteMsg::Revoke {
            spender,
            token_id,
        } => execute_revoke(deps, env, info, spender, token_id),
        
        ExecuteMsg::ApproveAll {
            operator,
            expires,
        } => execute_approve_all(deps, env, info, operator, expires),
        
        ExecuteMsg::RevokeAll {
            operator,
        } => execute_revoke_all(deps, env, info, operator),
        
        ExecuteMsg::UpdateMetadata {
            token_id,
            metadata,
        } => execute_update_metadata(deps, info, token_id, metadata),
        
        ExecuteMsg::Burn {
            token_id,
            reason,
        } => execute_burn(deps, info, token_id, reason),
        
        ExecuteMsg::UpdateConfig {
            admin,
            minter,
            ipfs_gateway,
        } => execute_update_config(deps, info, admin, minter, ipfs_gateway),
        
        ExecuteMsg::CacheIpfsContent {
            ipfs_link,
            content,
        } => execute_cache_ipfs_content(deps, ipfs_link, content),
        
        ExecuteMsg::ValidateNFTData {
            nft_type,
            data,
        } => execute_validate_nft_data(deps, nft_type, data),
    }
}

#[cfg_attr(not(feature = "library"), entry_point)]
pub fn query(deps: Deps, env: Env, msg: QueryMsg) -> StdResult<Binary> {
    match msg {
        QueryMsg::GetConfig {} => to_json_binary(&query_config(deps)?),
        
        QueryMsg::OwnerOf {
            token_id,
            include_expired,
        } => to_json_binary(&query_owner_of(deps, env, token_id, include_expired)?),
        
        QueryMsg::Approval {
            token_id,
            spender,
            include_expired,
        } => to_json_binary(&query_approval(deps, env, token_id, spender, include_expired)?),
        
        QueryMsg::Approvals {
            token_id,
            include_expired,
        } => to_json_binary(&query_approvals(deps, env, token_id, include_expired)?),
        
        QueryMsg::Tokens {
            owner,
            start_after,
            limit,
        } => to_json_binary(&query_tokens(deps, owner, start_after, limit)?),
        
        QueryMsg::AllTokens {
            start_after,
            limit,
        } => to_json_binary(&query_all_tokens(deps, start_after, limit)?),
        
        QueryMsg::NftInfo {
            token_id,
        } => to_json_binary(&query_nft_info(deps, token_id)?),
        
        QueryMsg::AllNftInfo {
            token_id,
            include_expired,
        } => to_json_binary(&query_all_nft_info(deps, env, token_id, include_expired)?),
        
        QueryMsg::ContractInfo {} => to_json_binary(&query_contract_info(deps)?),
        
        QueryMsg::GetSubjectNFT {
            token_id,
        } => to_json_binary(&query_subject_nft(deps, token_id)?),
        
        QueryMsg::GetDegreeNFT {
            token_id,
        } => to_json_binary(&query_degree_nft(deps, token_id)?),
        
        QueryMsg::GetStudentCollection {
            student_id,
        } => to_json_binary(&query_student_collection(deps, student_id)?),
        
        QueryMsg::GetNFTsByType {
            student_id,
            nft_type,
            limit,
        } => to_json_binary(&query_nfts_by_type(deps, student_id, nft_type, limit)?),
        
        QueryMsg::GetStatistics {} => to_json_binary(&query_statistics(deps)?),
        
        QueryMsg::GetNFTsByInstitution {
            institution_id,
            nft_type,
            limit,
            start_after,
        } => to_json_binary(&query_nfts_by_institution(deps, institution_id, nft_type, limit, start_after)?),
        
        QueryMsg::NFTExists {
            token_id,
        } => to_json_binary(&query_nft_exists(deps, token_id)?),
        
        QueryMsg::GetCachedContent {
            ipfs_link,
        } => to_json_binary(&query_cached_content(deps, ipfs_link)?),
        
        QueryMsg::VerifyNFT {
            token_id,
        } => to_json_binary(&query_verify_nft(deps, token_id)?),
    }
}
