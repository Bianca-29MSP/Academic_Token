package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"academictoken/x/tokendef/types"
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

// CreateTokenDefinition creates a new token definition for a subject
func (k msgServer) CreateTokenDefinition(goCtx context.Context, req *types.MsgCreateTokenDefinition) (*types.MsgCreateTokenDefinitionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.SubjectId == "" {
		return nil, fmt.Errorf("subject ID cannot be empty")
	}
	if req.TokenName == "" {
		return nil, fmt.Errorf("token name cannot be empty")
	}
	if req.TokenSymbol == "" {
		return nil, fmt.Errorf("token symbol cannot be empty")
	}
	if req.TokenType == "" {
		return nil, fmt.Errorf("token type cannot be empty")
	}

	// Validate token type
	validTokenTypes := []string{"NFT", "FUNGIBLE", "ACHIEVEMENT"}
	validTokenType := false
	for _, ttype := range validTokenTypes {
		if req.TokenType == ttype {
			validTokenType = true
			break
		}
	}
	if !validTokenType {
		return nil, fmt.Errorf("invalid token type: must be one of %v, got '%s'", validTokenTypes, req.TokenType)
	}

	// Check if subject exists and get its data
	var institutionId, courseId string
	if k.subjectKeeper != nil {
		subject, found := k.subjectKeeper.GetSubject(ctx, req.SubjectId)
		if !found {
			return nil, fmt.Errorf("subject with ID '%s' not found", req.SubjectId)
		}
		// Extract institution and course IDs from subject
		institutionId = subject.Institution
		courseId = subject.CourseId
	}

	// Check if token definition already exists for this subject
	allTokenDefs := k.GetAllTokenDefinition(ctx)
	for _, tokenDef := range allTokenDefs {
		if tokenDef.SubjectId == req.SubjectId {
			return nil, fmt.Errorf("token definition already exists for subject '%s'", req.SubjectId)
		}
	}

	// Generate next index
	index := k.GetNextTokenDefinitionIndex(ctx)

	// Create token metadata
	var attributes []*types.TokenAttribute
	for _, attrInput := range req.Attributes {
		attribute := &types.TokenAttribute{
			TraitType:   attrInput.TraitType,
			DisplayType: attrInput.DisplayType,
			IsDynamic:   attrInput.IsDynamic,
		}
		attributes = append(attributes, attribute)
	}

	metadata := &types.TokenMetadata{
		Description: req.Description,
		ImageUri:    req.ImageUri,
		Attributes:  attributes,
	}

	// Create token definition with institution and course IDs
	tokenDef := types.TokenDefinition{
		Index:          index,
		TokenDefId:     index, // Use index as token def ID
		SubjectId:      req.SubjectId,
		InstitutionId:  institutionId, // ✅ FIXED: Now populated from subject
		CourseId:       courseId,      // ✅ FIXED: Now populated from subject
		TokenName:      req.TokenName,
		TokenSymbol:    req.TokenSymbol,
		TokenType:      req.TokenType,
		IsTransferable: req.IsTransferable,
		IsBurnable:     req.IsBurnable,
		MaxSupply:      req.MaxSupply,
		Metadata:       metadata,
		Creator:        req.Creator,
		CreatedAt:      ctx.BlockTime().Format("2006-01-02T15:04:05Z"),
	}

	// Store token definition
	k.SetTokenDefinition(ctx, tokenDef)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"token_definition_created",
			sdk.NewAttribute("token_def_index", index),
			sdk.NewAttribute("subject_id", req.SubjectId),
			sdk.NewAttribute("token_name", req.TokenName),
			sdk.NewAttribute("token_symbol", req.TokenSymbol),
			sdk.NewAttribute("token_type", req.TokenType),
			sdk.NewAttribute("max_supply", fmt.Sprintf("%d", req.MaxSupply)),
			sdk.NewAttribute("is_transferable", fmt.Sprintf("%t", req.IsTransferable)),
			sdk.NewAttribute("is_burnable", fmt.Sprintf("%t", req.IsBurnable)),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Token definition created successfully",
		"token_def_index", index,
		"subject_id", req.SubjectId,
		"token_name", req.TokenName,
		"token_symbol", req.TokenSymbol,
		"creator", req.Creator,
	)

	return &types.MsgCreateTokenDefinitionResponse{
		Index: index,
	}, nil
}

// UpdateTokenDefinition updates an existing token definition
func (k msgServer) UpdateTokenDefinition(goCtx context.Context, req *types.MsgUpdateTokenDefinition) (*types.MsgUpdateTokenDefinitionResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.TokenDefId == "" {
		return nil, fmt.Errorf("token definition ID cannot be empty")
	}

	// Get existing token definition
	tokenDef, found := k.GetTokenDefinitionByIndex(ctx, req.TokenDefId)
	if !found {
		return nil, fmt.Errorf("token definition with ID '%s' not found", req.TokenDefId)
	}

	// Check permissions (only creator or authority can update)
	if req.Creator != tokenDef.Creator && req.Creator != k.GetAuthority() {
		return nil, fmt.Errorf("creator '%s' is not authorized to update token definition '%s'", req.Creator, req.TokenDefId)
	}

	// Update fields if provided
	updated := false

	if req.TokenName != "" && req.TokenName != tokenDef.TokenName {
		tokenDef.TokenName = req.TokenName
		updated = true
	}

	if req.TokenSymbol != "" && req.TokenSymbol != tokenDef.TokenSymbol {
		tokenDef.TokenSymbol = req.TokenSymbol
		updated = true
	}

	if req.Description != "" {
		if tokenDef.Metadata == nil {
			tokenDef.Metadata = &types.TokenMetadata{}
		}
		if req.Description != tokenDef.Metadata.Description {
			tokenDef.Metadata.Description = req.Description
			updated = true
		}
	}

	if req.IsTransferable != tokenDef.IsTransferable {
		tokenDef.IsTransferable = req.IsTransferable
		updated = true
	}

	if req.IsBurnable != tokenDef.IsBurnable {
		tokenDef.IsBurnable = req.IsBurnable
		updated = true
	}

	if req.MaxSupply != tokenDef.MaxSupply {
		// Note: CurrentSupply validation not available in this protobuf version
		tokenDef.MaxSupply = req.MaxSupply
		updated = true
	}

	if req.ImageUri != "" {
		if tokenDef.Metadata == nil {
			tokenDef.Metadata = &types.TokenMetadata{}
		}
		if req.ImageUri != tokenDef.Metadata.ImageUri {
			tokenDef.Metadata.ImageUri = req.ImageUri
			updated = true
		}
	}

	if len(req.Attributes) > 0 {
		if tokenDef.Metadata == nil {
			tokenDef.Metadata = &types.TokenMetadata{}
		}
		// Convert attribute inputs to token attributes
		var attributes []*types.TokenAttribute
		for _, attrInput := range req.Attributes {
			attribute := &types.TokenAttribute{
				TraitType:   attrInput.TraitType,
				DisplayType: attrInput.DisplayType,
				IsDynamic:   attrInput.IsDynamic,
			}
			attributes = append(attributes, attribute)
		}
		tokenDef.Metadata.Attributes = attributes
		updated = true
	}

	if !updated {
		return nil, fmt.Errorf("no valid updates provided")
	}

	// Store updated token definition
	k.SetTokenDefinition(ctx, tokenDef)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"token_definition_updated",
			sdk.NewAttribute("token_def_id", req.TokenDefId),
			sdk.NewAttribute("token_name", tokenDef.TokenName),
			sdk.NewAttribute("token_symbol", tokenDef.TokenSymbol),
			sdk.NewAttribute("max_supply", fmt.Sprintf("%d", tokenDef.MaxSupply)),
			sdk.NewAttribute("is_transferable", fmt.Sprintf("%t", tokenDef.IsTransferable)),
			sdk.NewAttribute("is_burnable", fmt.Sprintf("%t", tokenDef.IsBurnable)),
			sdk.NewAttribute("updater", req.Creator),
		),
	)

	k.Logger().Info("Token definition updated successfully",
		"token_def_id", req.TokenDefId,
		"token_name", tokenDef.TokenName,
		"updater", req.Creator,
	)

	return &types.MsgUpdateTokenDefinitionResponse{}, nil
}
