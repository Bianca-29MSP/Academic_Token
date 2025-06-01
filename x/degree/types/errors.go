package types

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/degree module sentinel errors
var (
	ErrInvalidDegreeRequest     = sdkerrors.Register(ModuleName, 1101, "invalid degree request")
	ErrStudentNotFound          = sdkerrors.Register(ModuleName, 1102, "student not found")
	ErrCurriculumNotCompleted   = sdkerrors.Register(ModuleName, 1103, "curriculum requirements not completed")
	ErrContractValidationFailed = sdkerrors.Register(ModuleName, 1104, "contract validation failed")
	ErrDegreeAlreadyIssued      = sdkerrors.Register(ModuleName, 1105, "degree already issued for this student and curriculum")
	ErrDegreeNotFound           = sdkerrors.Register(ModuleName, 1106, "degree not found")
	ErrDegreeRequestNotFound    = sdkerrors.Register(ModuleName, 1107, "degree request not found")
	ErrInvalidContractAddress   = sdkerrors.Register(ModuleName, 1108, "invalid contract address")
	ErrContractExecutionFailed  = sdkerrors.Register(ModuleName, 1109, "contract execution failed")
	ErrInsufficientCredits      = sdkerrors.Register(ModuleName, 1110, "insufficient credits for graduation")
	ErrInvalidGrade             = sdkerrors.Register(ModuleName, 1111, "invalid grade or GPA below minimum")
	ErrDegreeValidationPending  = sdkerrors.Register(ModuleName, 1112, "degree validation is still pending")
	ErrUnauthorizedDegreeIssuer = sdkerrors.Register(ModuleName, 1113, "unauthorized to issue degrees")
	ErrInvalidDegreeStatus      = sdkerrors.Register(ModuleName, 1114, "invalid degree status")
	ErrDuplicateDegreeRequest   = sdkerrors.Register(ModuleName, 1115, "duplicate degree request")
	ErrInvalidSigner            = sdkerrors.Register(ModuleName, 1116, "invalid signer for the operation")
	ErrInvalidAddress           = sdkerrors.Register(ModuleName, 1117, "invalid address format")
)
