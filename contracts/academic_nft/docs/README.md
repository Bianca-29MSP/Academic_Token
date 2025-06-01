# Academic NFT Contract

This contract implements NFT functionality for the Academic Token system, handling the minting and management of academic achievement NFTs.

## Features

- **Subject Completion NFTs**: Tokenize individual subject/course completions
- **Degree NFTs**: Tokenize degree and diploma achievements  
- **CW721 Compatibility**: Fully compatible with Cosmos NFT standard
- **IPFS Integration**: Hybrid storage with metadata in IPFS
- **Verification System**: Built-in authenticity verification
- **Student Collections**: Track all NFTs owned by each student
- **Institution Queries**: Query NFTs by issuing institution

## NFT Types

### Subject Completion NFTs
- Represent completion of individual academic subjects/courses
- Include grade, credits, completion date, and full subject metadata
- Linked to syllabus content stored in IPFS

### Degree NFTs  
- Represent completion of full academic programs
- Include GPA, honors, graduation date, and comprehensive degree data
- Digitally signed by academic authorities
- Reference all completed subjects that contributed to the degree

## Contract Integration

This contract is designed to be called by:
- **AcademicNFT Module**: Primary minter for subject and degree NFTs
- **Progress Contract**: For updating student achievement records
- **Degree Contract**: For minting diplomas after validation
- **Student Module**: For querying student NFT collections

## Storage Architecture

Following the mandatory hybrid storage model:
- **On-chain**: NFT ownership, basic metadata, validation hashes
- **IPFS**: Complete academic metadata, syllabi, verification documents

## Key Operations

### Minting
- `MintSubjectNFT`: Create NFT for completed subject
- `MintDegreeNFT`: Create NFT for earned degree

### Management  
- `TransferNft`: Transfer NFT ownership
- `Approve/Revoke`: Manage NFT approvals
- `UpdateMetadata`: Update NFT metadata (admin only)
- `Burn`: Remove invalid NFTs (admin only)

### Queries
- `GetStudentCollection`: All NFTs owned by student
- `GetNFTsByInstitution`: NFTs issued by specific institution
- `VerifyNFT`: Verify NFT authenticity and validation
- Standard CW721 queries for ownership and approvals

## Verification System

Each NFT includes:
- Validation hash from issuing authority
- Timestamp of issuance
- Digital signatures from academic authorities
- Verification URL for external validation

## Example Usage

```rust
// Mint subject completion NFT
let mint_msg = ExecuteMsg::MintSubjectNFT {
    student_id: "student123".to_string(),
    subject_data: SubjectCompletionData {
        subject_id: "CS101".to_string(),
        subject_name: "Introduction to Computer Science".to_string(),
        final_grade: 85,
        credits: 3,
        // ... other fields
    },
    metadata: NFTMetadata {
        name: "CS101 Completion Certificate".to_string(),
        description: "Certificate for completing CS101".to_string(),
        image: "ipfs://QmImageHash".to_string(),
        // ... other metadata
    },
    validation_hash: "validation_hash_from_institution".to_string(),
};
```

## Security Features

- Only authorized minters can create NFTs
- All minting operations are validated and logged
- NFT metadata includes cryptographic proofs
- Integration with institution validation systems
- Support for multi-signature degree validation

## IPFS Integration

The contract supports rich metadata storage in IPFS including:
- Complete academic transcripts
- Syllabus content and learning outcomes  
- Verification documents and signatures
- Multilingual support for international institutions
- Media files (certificates, photos, videos)

This contract is a critical component of the Academic Token ecosystem, providing the NFT infrastructure needed to tokenize academic achievements in a secure, verifiable, and interoperable manner.
