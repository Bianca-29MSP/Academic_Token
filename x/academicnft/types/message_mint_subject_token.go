package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgMintSubjectToken{}

func NewMsgMintSubjectToken(
	creator string,
	tokenDefId string,
	student string,
	completionDate string,
	grade string,
	issuerInstitution string,
	semester string,
	professorSignature string,
) *MsgMintSubjectToken {
	return &MsgMintSubjectToken{
		Creator:            creator,
		TokenDefId:         tokenDefId,
		Student:            student,
		CompletionDate:     completionDate,
		Grade:              grade,
		IssuerInstitution:  issuerInstitution,
		Semester:           semester,
		ProfessorSignature: professorSignature,
	}
}

// NOTE: For passive mode with contract authorization, use NewExtendedMsgMintSubjectToken
// from the message_extensions.go file instead

func (msg *MsgMintSubjectToken) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate required fields
	if msg.TokenDefId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "tokenDefId cannot be empty")
	}

	if msg.Student == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "student cannot be empty")
	}

	if msg.CompletionDate == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "completionDate cannot be empty")
	}

	if msg.Grade == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "grade cannot be empty")
	}

	if msg.IssuerInstitution == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "issuerInstitution cannot be empty")
	}

	// NOTE: PASSIVE MODULE VALIDATION is handled by ExtendedMsgMintSubjectToken
	// This base message only validates core fields

	return nil
}

// IsContractAuthorized returns true if this message includes contract authorization
// Note: Base message doesn't have passive mode fields - use ExtendedMsgMintSubjectToken
func (msg *MsgMintSubjectToken) IsContractAuthorized() bool {
	return false // Base message is never contract authorized
}

// Route returns the message route for routing
func (msg *MsgMintSubjectToken) Route() string {
	return RouterKey
}

// Type returns the message type
func (msg *MsgMintSubjectToken) Type() string {
	return TypeMsgMintSubjectToken
}

// GetSigners returns the expected signers for this message
func (msg *MsgMintSubjectToken) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

// GetSignBytes returns the bytes for the message signer to sign on
func (msg *MsgMintSubjectToken) GetSignBytes() []byte {
	// Use legacy amino codec for backward compatibility
	// In production, this should use the proper module codec
	bz := sdk.MustSortJSON([]byte("{}")) // Placeholder
	return bz
}
