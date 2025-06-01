package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"academictoken/x/tokendef/types"
)

var _ types.QueryServer = Keeper{}

// Params returns the total set of parameters.
/*func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryParamsResponse{Params: k.GetParams(ctx)}, nil
}*/

// GetTokenDefinition queries a token definition by index (QueryServer interface method)
func (k Keeper) GetTokenDefinition(goCtx context.Context, req *types.QueryGetTokenDefinitionRequest) (*types.QueryGetTokenDefinitionResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Use the internal keeper method
	val, found := k.GetTokenDefinitionByIndex(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetTokenDefinitionResponse{TokenDefinition: &val}, nil
}

// GetTokenDefinitionFull queries a token definition with full content (including IPFS data)
func (k Keeper) GetTokenDefinitionFull(goCtx context.Context, req *types.QueryGetTokenDefinitionFullRequest) (*types.QueryGetTokenDefinitionFullResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	tokenDef, found := k.GetTokenDefinitionByIndex(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.NotFound, "token definition not found")
	}

	// TODO: Retrieve extended content from IPFS using tokenDef.IpfsLink
	// For now, return empty extended content
	extendedContent := ""

	return &types.QueryGetTokenDefinitionFullResponse{
		TokenDefinition: &tokenDef,
		ExtendedContent: extendedContent,
	}, nil
}

// ListTokenDefinitions queries all token definitions with pagination
func (k Keeper) ListTokenDefinitions(goCtx context.Context, req *types.QueryListTokenDefinitionsRequest) (*types.QueryListTokenDefinitionsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var tokenDefinitions []*types.TokenDefinition
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	tokenDefinitionStore := prefix.NewStore(store, types.KeyPrefix(types.TokenDefinitionKeyPrefix))

	pageRes, err := query.Paginate(tokenDefinitionStore, req.Pagination, func(key []byte, value []byte) error {
		var tokenDefinition types.TokenDefinition
		if err := k.cdc.Unmarshal(value, &tokenDefinition); err != nil {
			return err
		}

		tokenDefinitions = append(tokenDefinitions, &tokenDefinition)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryListTokenDefinitionsResponse{
		TokenDefinitions: tokenDefinitions,
		Pagination:       pageRes,
	}, nil
}

// ListTokenDefinitionsBySubject queries token definitions by subject
func (k Keeper) ListTokenDefinitionsBySubject(goCtx context.Context, req *types.QueryListTokenDefinitionsBySubjectRequest) (*types.QueryListTokenDefinitionsBySubjectResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get all token definitions for the subject
	allTokenDefs := k.GetTokenDefinitionsBySubject(ctx, req.SubjectId)

	// Convert to pointers for response
	var tokenDefinitions []*types.TokenDefinition
	for i := range allTokenDefs {
		tokenDefinitions = append(tokenDefinitions, &allTokenDefs[i])
	}

	// Apply pagination manually since we're using a custom query
	var pagedTokenDefs []*types.TokenDefinition
	var pageRes *query.PageResponse

	if req.Pagination != nil {
		offset := uint64(0)
		limit := uint64(100) // Default limit

		if req.Pagination.Offset > 0 {
			offset = req.Pagination.Offset
		}
		if req.Pagination.Limit > 0 {
			limit = req.Pagination.Limit
		}

		start := offset
		end := offset + limit
		if end > uint64(len(tokenDefinitions)) {
			end = uint64(len(tokenDefinitions))
		}

		if start < uint64(len(tokenDefinitions)) {
			pagedTokenDefs = tokenDefinitions[start:end]
		}

		// Create page response
		pageRes = &query.PageResponse{
			Total: uint64(len(tokenDefinitions)),
		}
		if end < uint64(len(tokenDefinitions)) {
			pageRes.NextKey = []byte("next") // Simplified next key
		}
	} else {
		pagedTokenDefs = tokenDefinitions
		pageRes = &query.PageResponse{
			Total: uint64(len(tokenDefinitions)),
		}
	}

	return &types.QueryListTokenDefinitionsBySubjectResponse{
		TokenDefinitions: pagedTokenDefs,
		Pagination:       pageRes,
	}, nil
}
