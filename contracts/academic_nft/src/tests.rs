#[cfg(test)]
mod tests {
    use crate::ContractError;
    use cosmwasm_std::testing::{mock_dependencies, mock_env, mock_info};
    use cosmwasm_std::{coins, from_json};
    use crate::contract::{instantiate, execute, query};
    use crate::msg::{InstantiateMsg, ExecuteMsg, QueryMsg, SubjectCompletionData, DegreeCompletionData};
    use crate::state::{NFTMetadata, NFTAttribute, NFTType};

    #[test]
    fn proper_initialization() {
        let mut deps = mock_dependencies();

        let msg = InstantiateMsg {
            admin: None,
            minter: "academic_nft_module".to_string(),
            ipfs_gateway: "https://ipfs.io".to_string(),
            collection_name: "Academic Token NFTs".to_string(),
            collection_symbol: "ATNFT".to_string(),
        };
        let info = mock_info("creator", &coins(1000, "earth"));

        // Instantiate contract
        let res = instantiate(deps.as_mut(), mock_env(), info, msg).unwrap();
        assert_eq!(0, res.messages.len());

        // Query config
        let res = query(deps.as_ref(), mock_env(), QueryMsg::GetConfig {}).unwrap();
        let config: crate::msg::ConfigResponse = from_json(&res).unwrap();
        assert_eq!("Academic Token NFTs", config.config.collection_name);
        assert_eq!("ATNFT", config.config.collection_symbol);
    }

    #[test]
    fn mint_subject_nft() {
        let mut deps = mock_dependencies();

        // Initialize contract
        let instantiate_msg = InstantiateMsg {
            admin: Some("admin".to_string()),
            minter: "minter".to_string(),
            ipfs_gateway: "https://ipfs.io".to_string(),
            collection_name: "Academic Token NFTs".to_string(),
            collection_symbol: "ATNFT".to_string(),
        };
        let info = mock_info("creator", &coins(1000, "earth"));
        instantiate(deps.as_mut(), mock_env(), info, instantiate_msg).unwrap();

        // Mint subject NFT
        let subject_data = SubjectCompletionData {
            subject_id: "CS101".to_string(),
            institution_id: "MIT".to_string(),
            course_id: "computer_science".to_string(),
            subject_name: "Introduction to Computer Science".to_string(),
            credits: 3,
            final_grade: 85,
            completion_date: "2024-05-31".to_string(),
            semester: "Spring 2024".to_string(),
            academic_year: "2023-2024".to_string(),
            instructor: Some("Dr. Smith".to_string()),
        };

        let metadata = NFTMetadata {
            name: "CS101 Completion".to_string(),
            description: "Certificate for completing CS101".to_string(),
            image: "ipfs://QmTest123".to_string(),
            external_url: None,
            animation_url: None,
            attributes: vec![
                NFTAttribute {
                    trait_type: "Grade".to_string(),
                    value: "85".to_string(),
                    display_type: Some("number".to_string()),
                }
            ],
            ipfs_metadata_link: "ipfs://QmMetadata123".to_string(),
            created_at: 1234567890,
            updated_at: 1234567890,
        };

        let mint_msg = ExecuteMsg::MintSubjectNFT {
            student_id: "student123".to_string(),
            subject_data,
            metadata,
            validation_hash: "hash123".to_string(),
        };

        let info = mock_info("minter", &[]);
        let res = execute(deps.as_mut(), mock_env(), info, mint_msg).unwrap();

        // Check response
        assert_eq!(4, res.attributes.len());
        assert_eq!("mint_subject_nft", res.attributes[0].value);
        assert_eq!("subject_1", res.attributes[1].value); // token_id
        assert_eq!("student123", res.attributes[2].value); // student_id
        assert_eq!("CS101", res.attributes[3].value); // subject_id

        // Query the minted NFT
        let query_msg = QueryMsg::GetSubjectNFT {
            token_id: "subject_1".to_string(),
        };
        let res = query(deps.as_ref(), mock_env(), query_msg).unwrap();
        let nft_response: crate::msg::SubjectNFTResponse = from_json(&res).unwrap();
        assert!(nft_response.nft.is_some());
        
        let nft = nft_response.nft.unwrap();
        assert_eq!("CS101", nft.subject_id);
        assert_eq!("student123", nft.student_id);
        assert_eq!(85, nft.final_grade);
    }

    #[test]
    fn mint_degree_nft() {
        let mut deps = mock_dependencies();

        // Initialize contract
        let instantiate_msg = InstantiateMsg {
            admin: Some("admin".to_string()),
            minter: "minter".to_string(),
            ipfs_gateway: "https://ipfs.io".to_string(),
            collection_name: "Academic Token NFTs".to_string(),
            collection_symbol: "ATNFT".to_string(),
        };
        let info = mock_info("creator", &coins(1000, "earth"));
        instantiate(deps.as_mut(), mock_env(), info, instantiate_msg).unwrap();

        // Mint degree NFT
        let degree_data = DegreeCompletionData {
            institution_id: "MIT".to_string(),
            course_id: "computer_science".to_string(),
            degree_name: "Bachelor of Science in Computer Science".to_string(),
            degree_type: "Bachelor".to_string(),
            major: "Computer Science".to_string(),
            minor: None,
            graduation_date: "2024-05-31".to_string(),
            final_gpa: "3.8".to_string(),
            total_credits: 120,
            honors: Some("Magna Cum Laude".to_string()),
        };

        let metadata = NFTMetadata {
            name: "BS Computer Science - MIT".to_string(),
            description: "Bachelor's degree in Computer Science from MIT".to_string(),
            image: "ipfs://QmDegree123".to_string(),
            external_url: None,
            animation_url: None,
            attributes: vec![
                NFTAttribute {
                    trait_type: "GPA".to_string(),
                    value: "3.8".to_string(),
                    display_type: Some("number".to_string()),
                }
            ],
            ipfs_metadata_link: "ipfs://QmDegreeMetadata123".to_string(),
            created_at: 1234567890,
            updated_at: 1234567890,
        };

        let mint_msg = ExecuteMsg::MintDegreeNFT {
            student_id: "student123".to_string(),
            degree_data,
            metadata,
            validation_hash: "degreehash123".to_string(),
            signatures: vec!["signature1".to_string(), "signature2".to_string()],
        };

        let info = mock_info("minter", &[]);
        let res = execute(deps.as_mut(), mock_env(), info, mint_msg).unwrap();

        // Check response
        assert_eq!(4, res.attributes.len());
        assert_eq!("mint_degree_nft", res.attributes[0].value);
        assert_eq!("degree_1", res.attributes[1].value); // token_id
        assert_eq!("student123", res.attributes[2].value); // student_id
        assert_eq!("MIT", res.attributes[3].value); // institution_id

        // Query the minted NFT
        let query_msg = QueryMsg::GetDegreeNFT {
            token_id: "degree_1".to_string(),
        };
        let res = query(deps.as_ref(), mock_env(), query_msg).unwrap();
        let nft_response: crate::msg::DegreeNFTResponse = from_json(&res).unwrap();
        assert!(nft_response.nft.is_some());
        
        let nft = nft_response.nft.unwrap();
        assert_eq!("Computer Science", nft.major);
        assert_eq!("3.8", nft.final_gpa);
        assert_eq!(Some("Magna Cum Laude".to_string()), nft.honors);
    }

    #[test]
    fn unauthorized_minting() {
        let mut deps = mock_dependencies();

        // Initialize contract
        let instantiate_msg = InstantiateMsg {
            admin: Some("admin".to_string()),
            minter: "minter".to_string(),
            ipfs_gateway: "https://ipfs.io".to_string(),
            collection_name: "Academic Token NFTs".to_string(),
            collection_symbol: "ATNFT".to_string(),
        };
        let info = mock_info("creator", &coins(1000, "earth"));
        instantiate(deps.as_mut(), mock_env(), info, instantiate_msg).unwrap();

        // Try to mint with unauthorized user
        let subject_data = SubjectCompletionData {
            subject_id: "CS101".to_string(),
            institution_id: "MIT".to_string(),
            course_id: "computer_science".to_string(),
            subject_name: "Introduction to Computer Science".to_string(),
            credits: 3,
            final_grade: 85,
            completion_date: "2024-05-31".to_string(),
            semester: "Spring 2024".to_string(),
            academic_year: "2023-2024".to_string(),
            instructor: None,
        };

        let metadata = NFTMetadata {
            name: "CS101 Completion".to_string(),
            description: "Certificate for completing CS101".to_string(),
            image: "ipfs://QmTest123".to_string(),
            external_url: None,
            animation_url: None,
            attributes: vec![],
            ipfs_metadata_link: "ipfs://QmMetadata123".to_string(),
            created_at: 1234567890,
            updated_at: 1234567890,
        };

        let mint_msg = ExecuteMsg::MintSubjectNFT {
            student_id: "student123".to_string(),
            subject_data,
            metadata,
            validation_hash: "hash123".to_string(),
        };

        let info = mock_info("unauthorized", &[]);
        let err = execute(deps.as_mut(), mock_env(), info, mint_msg).unwrap_err();
        
        match err {
            ContractError::Unauthorized { reason } => {
                assert!(reason.contains("Only minter or admin can mint NFTs"));
            }
            _ => panic!("Expected Unauthorized error"),
        }
    }

    #[test]
    fn query_student_collection() {
        let mut deps = mock_dependencies();

        // Initialize contract
        let instantiate_msg = InstantiateMsg {
            admin: Some("admin".to_string()),
            minter: "minter".to_string(),
            ipfs_gateway: "https://ipfs.io".to_string(),
            collection_name: "Academic Token NFTs".to_string(),
            collection_symbol: "ATNFT".to_string(),
        };
        let info = mock_info("creator", &coins(1000, "earth"));
        instantiate(deps.as_mut(), mock_env(), info, instantiate_msg).unwrap();

        // Mint a subject NFT first
        let subject_data = SubjectCompletionData {
            subject_id: "CS101".to_string(),
            institution_id: "MIT".to_string(),
            course_id: "computer_science".to_string(),
            subject_name: "Introduction to Computer Science".to_string(),
            credits: 3,
            final_grade: 85,
            completion_date: "2024-05-31".to_string(),
            semester: "Spring 2024".to_string(),
            academic_year: "2023-2024".to_string(),
            instructor: None,
        };

        let metadata = NFTMetadata {
            name: "CS101 Completion".to_string(),
            description: "Certificate for completing CS101".to_string(),
            image: "ipfs://QmTest123".to_string(),
            external_url: None,
            animation_url: None,
            attributes: vec![],
            ipfs_metadata_link: "ipfs://QmMetadata123".to_string(),
            created_at: 1234567890,
            updated_at: 1234567890,
        };

        let mint_msg = ExecuteMsg::MintSubjectNFT {
            student_id: "student123".to_string(),
            subject_data,
            metadata,
            validation_hash: "hash123".to_string(),
        };

        let info = mock_info("minter", &[]);
        execute(deps.as_mut(), mock_env(), info, mint_msg).unwrap();

        // Query student collection
        let query_msg = QueryMsg::GetStudentCollection {
            student_id: "student123".to_string(),
        };
        let res = query(deps.as_ref(), mock_env(), query_msg).unwrap();
        let collection_response: crate::msg::StudentCollectionResponse = from_json(&res).unwrap();
        
        assert!(collection_response.collection.is_some());
        let collection = collection_response.collection.unwrap();
        assert_eq!("student123", collection.student_id);
        assert_eq!(1, collection.subject_nfts.len());
        assert_eq!("subject_1", collection.subject_nfts[0]);
        assert_eq!(0, collection.degree_nfts.len());
    }

    #[test]
    fn query_nfts_by_type() {
        let mut deps = mock_dependencies();

        // Initialize and mint NFT (similar setup as above)
        let instantiate_msg = InstantiateMsg {
            admin: Some("admin".to_string()),
            minter: "minter".to_string(),
            ipfs_gateway: "https://ipfs.io".to_string(),
            collection_name: "Academic Token NFTs".to_string(),
            collection_symbol: "ATNFT".to_string(),
        };
        let info = mock_info("creator", &coins(1000, "earth"));
        instantiate(deps.as_mut(), mock_env(), info, instantiate_msg).unwrap();

        // Query NFTs by type for student with no NFTs
        let query_msg = QueryMsg::GetNFTsByType {
            student_id: "student123".to_string(),
            nft_type: NFTType::SubjectCompletion,
            limit: None,
        };
        let res = query(deps.as_ref(), mock_env(), query_msg).unwrap();
        let nfts_response: crate::msg::NFTsByTypeResponse = from_json(&res).unwrap();
        
        assert_eq!(0, nfts_response.total_count);
        assert_eq!(0, nfts_response.nfts.len());
    }

    #[test]
    fn test_ipfs_utils() {
        use crate::ipfs::ipfs_utils::*;

        // Test content hash generation
        let content = "test content";
        let hash = generate_content_hash(content);
        assert!(!hash.is_empty());
        assert_eq!(64, hash.len()); // SHA256 produces 64 character hex string

        // Test verification URL creation
        let url = create_verification_url("contract123", "token456", "cosmos");
        assert_eq!("https://academic-token.org/verify/cosmos/contract123/token456", url);

        // Test IPFS link validation
        assert!(validate_ipfs_link("ipfs://QmTest123"));
        assert!(validate_ipfs_link("https://ipfs.io/ipfs/QmTest123"));
        assert!(!validate_ipfs_link("http://example.com/file"));

        // Test IPFS to gateway URL conversion
        let gateway_url = ipfs_to_gateway_url("ipfs://QmTest123", "https://ipfs.io");
        assert_eq!("https://ipfs.io/ipfs/QmTest123", gateway_url);
        
        let gateway_url2 = ipfs_to_gateway_url("QmTest123", "https://gateway.pinata.cloud/");
        assert_eq!("https://gateway.pinata.cloud/ipfs/QmTest123", gateway_url2);
    }
}
