# AcademicNFT Module - Passive Mode Implementation

## Overview

The AcademicNFT Module has been successfully refactored to operate in **PASSIVE MODE**, meaning it no longer implements business logic directly. Instead, it acts as a "thin wrapper" that only executes operations when authorized by CosmWasm contracts.

## Key Changes Implemented

### 1. **Passive Mode Architecture**
- The module now only responds to contract-authorized operations
- All business logic has been moved to CosmWasm contracts
- The module validates contract authorization before executing any operations

### 2. **Contract Authorization System**
- **Authorization Hash**: Every mint operation requires a cryptographic hash proving contract authorization
- **Authorized Contract List**: Only pre-approved contract addresses can authorize operations
- **Strict Validation**: Multiple layers of verification ensure security

### 3. **Enhanced Error Handling**
New error types for passive mode:
- `ErrUnauthorizedMinting`: Only authorized contracts can mint NFTs
- `ErrMissingContractAuthorization`: Contract authorization hash required
- `ErrInvalidContractAuthorization`: Cryptographic verification failed
- `ErrContractCallRequired`: Direct minting not allowed in passive mode
- `ErrPassiveModeViolation`: General passive mode violations

### 4. **Extended Data Types**
- `ExtendedSubjectTokenInstance`: Adds passive mode fields to token instances
- `PassiveModeConfig`: Configuration for passive mode operations
- Contract authorization tracking in responses

## How Passive Mode Works

### **Before (Active Mode)**
```
User → Module → Business Logic → Store Token
```

### **After (Passive Mode)**
```
User → Contract → Business Logic → Contract Authorization → Module → Store Token
```

## Contract Authorization Flow

1. **Contract Processes Request**: A CosmWasm contract validates the business logic
2. **Authorization Hash Generated**: Contract creates a cryptographic hash of the operation
3. **Module Call**: Contract calls the AcademicNFT module with authorization
4. **Verification**: Module verifies the contract is authorized and hash is valid
5. **Execution**: Module executes the operation (minting, etc.)
6. **Audit Trail**: Operation is logged for transparency

## Key Methods Added

### Authorization Methods
- `isAuthorizedContractCall()`: Verifies if caller is an authorized contract
- `verifyContractAuthorization()`: Cryptographic verification of authorization
- `ValidatePassiveModeOperation()`: Validates passive mode operations

### Contract Management
- `AddAuthorizedContract()`: Add contract to authorized list (governance)
- `RemoveAuthorizedContract()`: Remove contract authorization
- `GetAuthorizedContracts()`: Get list of authorized contracts

### Passive Mode Helpers
- `RecordPassiveOperation()`: Audit trail for passive operations
- `ExtendTokenInstanceForPassiveMode()`: Add passive mode fields

## Configuration Parameters

New parameters for passive mode:
- `PassiveModeEnabled`: Enable/disable passive mode
- `RequireContractAuth`: Require contract authorization for all operations
- `AllowDirectMinting`: Allow direct minting (for development)
- `AuthorizedContracts`: List of authorized contract addresses
- `VerificationLevel`: "strict", "moderate", or "permissive"

## Security Features

### **Multiple Authorization Layers**
1. **Caller Verification**: Only authorized contract addresses can call
2. **Hash Verification**: Cryptographic proof of contract authorization
3. **Parameter Validation**: All data validated before processing
4. **Audit Logging**: All operations logged for transparency

### **Default Security Stance**
- Passive mode enabled by default
- Contract authorization required by default
- Direct minting disabled by default
- Strict verification level by default

## Usage Examples

### **Contract-Authorized Minting**
```go
// Contract calls with authorization
msg := types.NewMsgMintSubjectTokenPassive(
    contractAddress,           // Contract address
    tokenDefId,               // Token definition ID
    studentAddress,           // Student address
    completionDate,           // Completion date
    grade,                    // Grade achieved
    institutionId,            // Issuing institution
    semester,                 // Semester
    professorSignature,       // Professor signature
    subjectId,                // Subject ID
    authorizationHash,        // Contract authorization hash
)
```

### **Authorization Hash Generation**
```go
// In the contract:
authData := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s:%s",
    tokenDefId, studentAddress, subjectId, completionDate,
    grade, institutionId, semester, contractAddress)
authHash := sha256.Sum256([]byte(authData))
```

## Integration with Other Modules

The AcademicNFT Module now integrates with:

1. **Student Module**: Receives NFT minting requests from Student contract calls
2. **Degree Module**: Processes degree NFT minting when graduation requirements are met
3. **Academic Progress Contract**: Validates academic progress before NFT issuance
4. **NFT Minting Contract**: Authorizes individual subject NFT minting

## Benefits of Passive Mode

### **Security**
- No direct access to business logic
- All operations require contract authorization
- Cryptographic verification of requests
- Comprehensive audit trail

### **Flexibility**
- Business logic can be updated via contract upgrades
- No module code changes needed for logic updates
- Easy to add new contract integrations

### **Transparency**
- All operations are contract-authorized
- Clear audit trail of all actions
- Public verification of authorization

### **Interoperability**
- Consistent with other passive modules
- Standardized contract interface
- Easy integration with external contracts

## Migration Notes

### **From Active to Passive Mode**
1. All existing business logic removed from module
2. Contract authorization system implemented
3. Enhanced error handling and validation
4. Audit logging and transparency features added

### **Backward Compatibility**
- Basic token instance structure unchanged
- Query interfaces remain compatible
- Events enhanced with passive mode information

## Testing Passive Mode

### **Development Mode**
- Set `AllowDirectMinting: true` for testing
- Use mock contract addresses for development
- Verification level can be set to "permissive"

### **Production Mode**
- `AllowDirectMinting: false` (strict passive mode)
- Only authorized contract addresses allowed
- "strict" verification level enforced

## Future Enhancements

1. **Multi-Signature Authorization**: Require multiple contract approvals
2. **Time-Limited Authorization**: Authorization hashes with expiration
3. **Role-Based Authorization**: Different authorization levels for different operations
4. **Cross-Chain Authorization**: Support for contracts on other chains

---

## Summary

The AcademicNFT Module has been successfully transformed from an active module with embedded business logic into a secure, passive module that only executes contract-authorized operations. This architecture provides enhanced security, flexibility, and transparency while maintaining backward compatibility and ease of use.

The passive mode implementation ensures that all business logic is handled by specialized CosmWasm contracts, making the system more modular, upgradeable, and secure.
