package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"academictoken/x/institution/types"
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

// RegisterInstitution creates a new institution
func (k msgServer) RegisterInstitution(goCtx context.Context, req *types.MsgRegisterInstitution) (*types.MsgRegisterInstitutionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.Name == "" {
		return nil, fmt.Errorf("institution name cannot be empty")
	}
	if req.Address == "" {
		return nil, fmt.Errorf("institution address cannot be empty")
	}

	// Check if institution already exists by name
	if k.InstitutionExistsByName(ctx, req.Name) {
		return nil, fmt.Errorf("institution with name '%s' already exists", req.Name)
	}

	// Check if institution already exists by address
	if k.InstitutionExistsByAddress(ctx, req.Address) {
		return nil, fmt.Errorf("institution with address '%s' already exists", req.Address)
	}

	// Generate new institution index
	index := k.GetNextInstitutionIndex(ctx)

	// Create the institution
	institution := types.Institution{
		Index:        index,
		Creator:      req.Creator,
		Name:         req.Name,
		Address:      req.Address,
		IsAuthorized: "false", // New institutions start as unauthorized
	}

	// Store the institution
	if err := k.SetInstitution(ctx, institution); err != nil {
		return nil, fmt.Errorf("failed to store institution: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"institution_registered",
			sdk.NewAttribute("index", index),
			sdk.NewAttribute("name", req.Name),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Institution registered successfully",
		"index", index,
		"name", req.Name,
		"creator", req.Creator,
	)

	return &types.MsgRegisterInstitutionResponse{
		Index: index,
	}, nil
}

// UpdateInstitution updates an existing institution
func (k msgServer) UpdateInstitution(goCtx context.Context, req *types.MsgUpdateInstitution) (*types.MsgUpdateInstitutionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.Index == "" {
		return nil, fmt.Errorf("institution index cannot be empty")
	}

	// Check if institution exists
	institution, found := k.GetInstitution(ctx, req.Index)
	if !found {
		return nil, fmt.Errorf("institution with index '%s' not found", req.Index)
	}

	// Check permissions
	if !k.CanUpdateInstitution(ctx, req.Index, req.Creator) {
		return nil, fmt.Errorf("creator '%s' is not authorized to update institution '%s'", req.Creator, req.Index)
	}

	// Update fields if provided
	updated := false

	if req.Name != "" && req.Name != institution.Name {
		// Check if the new name conflicts with existing institutions
		if k.InstitutionExistsByName(ctx, req.Name) {
			return nil, fmt.Errorf("institution with name '%s' already exists", req.Name)
		}
		institution.Name = req.Name
		updated = true
	}

	if req.Address != "" && req.Address != institution.Address {
		// Check if the new address conflicts with existing institutions
		if k.InstitutionExistsByAddress(ctx, req.Address) {
			return nil, fmt.Errorf("institution with address '%s' already exists", req.Address)
		}
		institution.Address = req.Address
		updated = true
	}

	// Handle authorization status change
	if req.IsAuthorized != "" && req.IsAuthorized != institution.IsAuthorized {
		// Check if the requester can change authorization status
		if !k.CanAuthorizeInstitution(ctx, req.Index, req.Creator) {
			return nil, fmt.Errorf("creator '%s' is not authorized to change authorization status of institution '%s'", req.Creator, req.Index)
		}

		// Validate the new authorization status
		if req.IsAuthorized != "true" && req.IsAuthorized != "false" {
			return nil, fmt.Errorf("invalid authorization status: must be 'true' or 'false', got '%s'", req.IsAuthorized)
		}

		institution.IsAuthorized = req.IsAuthorized
		updated = true
	}

	if !updated {
		return nil, fmt.Errorf("no valid updates provided")
	}

	// Store the updated institution
	if err := k.SetInstitution(ctx, institution); err != nil {
		return nil, fmt.Errorf("failed to update institution: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"institution_updated",
			sdk.NewAttribute("index", req.Index),
			sdk.NewAttribute("updater", req.Creator),
			sdk.NewAttribute("name", institution.Name),
			sdk.NewAttribute("is_authorized", institution.IsAuthorized),
		),
	)

	k.Logger().Info("Institution updated successfully",
		"index", req.Index,
		"updater", req.Creator,
		"name", institution.Name,
		"is_authorized", institution.IsAuthorized,
	)

	return &types.MsgUpdateInstitutionResponse{}, nil
}
