package types

import (
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSetGraduationRequirements{}

func NewMsgSetGraduationRequirements(
	creator string,
	curriculumIndex string,
	totalCreditsRequired string,
	minGpa string,
	requiredElectiveCredits string,
	minimumTimeYears string,
	maximumTimeYears string,
) *MsgSetGraduationRequirements {
	// Parse string values to appropriate types
	totalCredits, _ := strconv.ParseUint(totalCreditsRequired, 10, 64)
	minGpaFloat, _ := strconv.ParseFloat(minGpa, 32)
	requiredElective, _ := strconv.ParseUint(requiredElectiveCredits, 10, 64)
	minYears, _ := strconv.ParseFloat(minimumTimeYears, 32)
	maxYears, _ := strconv.ParseFloat(maximumTimeYears, 32)

	return &MsgSetGraduationRequirements{
		Creator:                 creator,
		CurriculumIndex:         curriculumIndex,
		TotalCreditsRequired:    totalCredits,
		MinGpa:                  float32(minGpaFloat),
		RequiredElectiveCredits: requiredElective,
		MinimumTimeYears:        float32(minYears),
		MaximumTimeYears:        float32(maxYears),
	}
}

func (msg *MsgSetGraduationRequirements) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Add validation for numeric fields
	if msg.TotalCreditsRequired == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "total credits required cannot be 0")
	}

	if msg.MinGpa <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "minimum GPA must be greater than 0")
	}

	if msg.MinimumTimeYears <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "minimum time years must be greater than 0")
	}

	if msg.MaximumTimeYears <= msg.MinimumTimeYears {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "maximum time years must be greater than minimum time years")
	}

	return nil
}
