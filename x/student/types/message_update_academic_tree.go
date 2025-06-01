package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateAcademicTree{}

func NewMsgUpdateAcademicTree(
	creator string,
	studentId string,
	completedTokens []string,
	inProgressTokens []string,
	availableTokens []string,
) *MsgUpdateAcademicTree {
	return &MsgUpdateAcademicTree{
		Creator:          creator,
		StudentId:        studentId,
		CompletedTokens:  completedTokens,
		InProgressTokens: inProgressTokens,
		AvailableTokens:  availableTokens,
	}
}

func (msg *MsgUpdateAcademicTree) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate student ID
	if strings.TrimSpace(msg.StudentId) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "student ID cannot be empty")
	}

	// Validate that tokens are not duplicated across different states
	if err := msg.validateTokenUniqueness(); err != nil {
		return err
	}

	// Validate individual token IDs
	if err := msg.validateTokenIds(); err != nil {
		return err
	}

	return nil
}

// validateTokenUniqueness ensures no token appears in multiple states
func (msg *MsgUpdateAcademicTree) validateTokenUniqueness() error {
	tokenMap := make(map[string]string) // token -> state

	// Check completed tokens
	for _, token := range msg.CompletedTokens {
		if state, exists := tokenMap[token]; exists {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
				"token %s appears in both %s and completed states", token, state)
		}
		tokenMap[token] = "completed"
	}

	// Check in-progress tokens
	for _, token := range msg.InProgressTokens {
		if state, exists := tokenMap[token]; exists {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
				"token %s appears in both %s and in-progress states", token, state)
		}
		tokenMap[token] = "in-progress"
	}

	// Check available tokens
	for _, token := range msg.AvailableTokens {
		if state, exists := tokenMap[token]; exists {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
				"token %s appears in both %s and available states", token, state)
		}
		tokenMap[token] = "available"
	}

	return nil
}

// validateTokenIds ensures all token IDs are valid (non-empty strings)
func (msg *MsgUpdateAcademicTree) validateTokenIds() error {
	allTokens := append(append(msg.CompletedTokens, msg.InProgressTokens...), msg.AvailableTokens...)

	for _, token := range allTokens {
		if strings.TrimSpace(token) == "" {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "token ID cannot be empty")
		}
	}

	return nil
}

// GetSigners returns the expected signers for the message
func (msg *MsgUpdateAcademicTree) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}
