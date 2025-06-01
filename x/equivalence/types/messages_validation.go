package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Implement sdk.Msg interface for messages that don't have separate files

// MsgExecuteEquivalenceAnalysis implementations
var _ sdk.Msg = &MsgExecuteEquivalenceAnalysis{}

func (msg *MsgExecuteEquivalenceAnalysis) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.EquivalenceId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "equivalence_id cannot be empty")
	}

	if msg.ContractAddress == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "contract_address cannot be empty")
	}

	return nil
}

func (msg *MsgExecuteEquivalenceAnalysis) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgExecuteEquivalenceAnalysis) Route() string { return ModuleName }
func (msg *MsgExecuteEquivalenceAnalysis) Type() string  { return "execute_equivalence_analysis" }

func (msg *MsgExecuteEquivalenceAnalysis) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgBatchRequestEquivalence implementations
var _ sdk.Msg = &MsgBatchRequestEquivalence{}

func (msg *MsgBatchRequestEquivalence) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if len(msg.Requests) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "requests cannot be empty")
	}

	if len(msg.Requests) > 100 { // Limit batch size
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "batch size cannot exceed 100 requests")
	}

	// Validate each request
	for i, req := range msg.Requests {
		if req.SourceSubjectId == "" {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "request[%d]: source_subject_id cannot be empty", i)
		}
		if req.TargetInstitution == "" {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "request[%d]: target_institution cannot be empty", i)
		}
		if req.TargetSubjectId == "" {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "request[%d]: target_subject_id cannot be empty", i)
		}
	}

	return nil
}

func (msg *MsgBatchRequestEquivalence) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBatchRequestEquivalence) Route() string { return ModuleName }
func (msg *MsgBatchRequestEquivalence) Type() string  { return "batch_request_equivalence" }

func (msg *MsgBatchRequestEquivalence) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgUpdateContractAddress implementations
var _ sdk.Msg = &MsgUpdateContractAddress{}

func (msg *MsgUpdateContractAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if msg.NewContractAddress == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "new_contract_address cannot be empty")
	}

	return nil
}

func (msg *MsgUpdateContractAddress) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateContractAddress) Route() string { return ModuleName }
func (msg *MsgUpdateContractAddress) Type() string  { return "update_contract_address" }

func (msg *MsgUpdateContractAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// MsgReanalyzeEquivalence implementations
var _ sdk.Msg = &MsgReanalyzeEquivalence{}

func (msg *MsgReanalyzeEquivalence) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.EquivalenceId == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "equivalence_id cannot be empty")
	}

	return nil
}

func (msg *MsgReanalyzeEquivalence) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgReanalyzeEquivalence) Route() string { return ModuleName }
func (msg *MsgReanalyzeEquivalence) Type() string  { return "reanalyze_equivalence" }

func (msg *MsgReanalyzeEquivalence) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
