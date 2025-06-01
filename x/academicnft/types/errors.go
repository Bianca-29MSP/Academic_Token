package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/academicnft module sentinel errors
var (
	ErrInvalidSigner          = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrTokenInstanceNotFound  = sdkerrors.Register(ModuleName, 1101, "token instance not found")
	ErrInvalidTokenDef        = sdkerrors.Register(ModuleName, 1102, "invalid token definition")
	ErrInvalidInstitution     = sdkerrors.Register(ModuleName, 1103, "invalid institution")
	ErrInvalidStudent         = sdkerrors.Register(ModuleName, 1104, "invalid student")
	ErrDuplicateTokenInstance = sdkerrors.Register(ModuleName, 1105, "token instance already exists")
	ErrInvalidGrade           = sdkerrors.Register(ModuleName, 1106, "invalid grade format")
	ErrInvalidDate            = sdkerrors.Register(ModuleName, 1107, "invalid date format")
	ErrMintingLimitExceeded   = sdkerrors.Register(ModuleName, 1108, "minting limit exceeded for this token definition")
	ErrUnauthorizedMinter     = sdkerrors.Register(ModuleName, 1109, "unauthorized to mint tokens for this institution")

	// PASSIVE MODULE ERRORS - Contract authorization required
	ErrUnauthorizedMinting            = sdkerrors.Register(ModuleName, 1110, "unauthorized minting: only authorized contracts can mint NFTs")
	ErrMissingContractAuthorization   = sdkerrors.Register(ModuleName, 1111, "missing contract authorization hash")
	ErrInvalidContractAuthorization   = sdkerrors.Register(ModuleName, 1112, "invalid contract authorization - cryptographic verification failed")
	ErrTokenDefMismatch              = sdkerrors.Register(ModuleName, 1113, "token definition subject mismatch with provided subject ID")
	ErrContractCallRequired          = sdkerrors.Register(ModuleName, 1114, "direct minting not allowed - must be called via authorized contract")
	ErrInvalidContractCaller         = sdkerrors.Register(ModuleName, 1115, "caller is not an authorized contract address")
	ErrPassiveModeViolation          = sdkerrors.Register(ModuleName, 1116, "passive mode violation - operation requires contract authorization")
)
