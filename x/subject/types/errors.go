package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/subject module sentinel errors
var (
	ErrInvalidSigner   = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample          = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrSubjectNotFound = sdkerrors.Register(ModuleName, 1102, "subject not found")
	ErrUnauthorized    = sdkerrors.Register(ModuleName, 1103, "unauthorized")
	ErrInvalidIPFS     = sdkerrors.Register(ModuleName, 1104, "invalid IPFS operation")
)
