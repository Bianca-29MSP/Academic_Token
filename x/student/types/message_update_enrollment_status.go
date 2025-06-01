package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateEnrollmentStatus{}

// Valid enrollment statuses
const (
	EnrollmentStatusPending   = "pending"
	EnrollmentStatusActive    = "active"
	EnrollmentStatusSuspended = "suspended"
	EnrollmentStatusCompleted = "completed"
	EnrollmentStatusCanceled  = "canceled"
	EnrollmentStatusExpired   = "expired"
)

func NewMsgUpdateEnrollmentStatus(creator string, enrollmentId string, status string) *MsgUpdateEnrollmentStatus {
	return &MsgUpdateEnrollmentStatus{
		Creator:      creator,
		EnrollmentId: enrollmentId,
		Status:       status,
	}
}

func (msg *MsgUpdateEnrollmentStatus) ValidateBasic() error {
	// Validate creator address
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate enrollment ID
	if strings.TrimSpace(msg.EnrollmentId) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "enrollment ID cannot be empty")
	}

	// Validate status
	if strings.TrimSpace(msg.Status) == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "status cannot be empty")
	}

	// Validate status is one of the allowed values
	if !isValidEnrollmentStatus(msg.Status) {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
			"invalid status: %s. Valid statuses are: %s",
			msg.Status,
			getValidEnrollmentStatuses())
	}

	return nil
}

// isValidEnrollmentStatus checks if the provided status is valid
func isValidEnrollmentStatus(status string) bool {
	validStatuses := []string{
		EnrollmentStatusPending,
		EnrollmentStatusActive,
		EnrollmentStatusSuspended,
		EnrollmentStatusCompleted,
		EnrollmentStatusCanceled,
		EnrollmentStatusExpired,
	}

	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// getValidEnrollmentStatuses returns a comma-separated string of valid statuses
func getValidEnrollmentStatuses() string {
	validStatuses := []string{
		EnrollmentStatusPending,
		EnrollmentStatusActive,
		EnrollmentStatusSuspended,
		EnrollmentStatusCompleted,
		EnrollmentStatusCanceled,
		EnrollmentStatusExpired,
	}
	return strings.Join(validStatuses, ", ")
}
