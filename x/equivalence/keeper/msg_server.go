package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"academictoken/x/equivalence/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// UpdateParams updates the module parameters
func (k msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the message signer is the authority
	if req.Authority != k.GetAuthority() {
		return nil, fmt.Errorf("invalid authority: expected %s, got %s", k.GetAuthority(), req.Authority)
	}

	// Validate the new params
	if err := req.Params.Validate(); err != nil {
		return nil, fmt.Errorf("invalid params: %w", err)
	}

	// Set the new params
	k.SetParams(ctx, req.Params)

	return &types.MsgUpdateParamsResponse{}, nil
}

// Note: All other message handler methods are implemented in msg_server_handlers.go
// to avoid duplication and maintain better organization.
