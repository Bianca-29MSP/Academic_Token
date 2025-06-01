package keeper

import (
	"context"
	"strconv"
	"time"

	"academictoken/x/equivalence/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Ensure the Keeper implements the QueryServer interface
var _ types.QueryServer = Keeper{}

// Params returns the module parameters
func (k Keeper) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// ListEquivalences queries all subject equivalences with pagination
func (k Keeper) ListEquivalences(goCtx context.Context, req *types.QueryListEquivalencesRequest) (*types.QueryListEquivalencesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	equivalences, pageRes, err := k.GetSubjectEquivalencesPaginated(ctx, req.Pagination, req.StatusFilter)
	if err != nil {
		return nil, err
	}

	return &types.QueryListEquivalencesResponse{
		Equivalences: equivalences,
		Pagination:   pageRes,
	}, nil
}

// GetEquivalence queries a specific subject equivalence by index
func (k Keeper) GetEquivalence(goCtx context.Context, req *types.QueryGetEquivalenceRequest) (*types.QueryGetEquivalenceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	equivalence, found := k.GetSubjectEquivalence(ctx, req.Index)
	if !found {
		return nil, types.ErrEquivalenceNotFound
	}

	return &types.QueryGetEquivalenceResponse{
		Equivalence: equivalence,
	}, nil
}

// GetEquivalencesBySourceSubject queries equivalences by source subject ID
func (k Keeper) GetEquivalencesBySourceSubject(goCtx context.Context, req *types.QueryGetEquivalencesBySourceSubjectRequest) (*types.QueryGetEquivalencesBySourceSubjectResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	equivalences, pageRes, err := k.GetEquivalencesBySourceSubjectInternal(ctx, req.SourceSubjectId, req.Pagination, req.StatusFilter)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetEquivalencesBySourceSubjectResponse{
		Equivalences: equivalences,
		Pagination:   pageRes,
	}, nil
}

// GetEquivalencesByTargetSubject queries equivalences by target subject ID
func (k Keeper) GetEquivalencesByTargetSubject(goCtx context.Context, req *types.QueryGetEquivalencesByTargetSubjectRequest) (*types.QueryGetEquivalencesByTargetSubjectResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	equivalences, pageRes, err := k.GetEquivalencesByTargetSubjectInternal(ctx, req.TargetSubjectId, req.Pagination, req.StatusFilter)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetEquivalencesByTargetSubjectResponse{
		Equivalences: equivalences,
		Pagination:   pageRes,
	}, nil
}

// GetEquivalencesByInstitution queries equivalences by target institution
func (k Keeper) GetEquivalencesByInstitution(goCtx context.Context, req *types.QueryGetEquivalencesByInstitutionRequest) (*types.QueryGetEquivalencesByInstitutionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	equivalences, pageRes, err := k.GetEquivalencesByInstitutionInternal(ctx, req.InstitutionId, req.Pagination, req.StatusFilter)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetEquivalencesByInstitutionResponse{
		Equivalences: equivalences,
		Pagination:   pageRes,
	}, nil
}

// CheckEquivalenceStatus checks if two subjects have an established equivalence
func (k Keeper) CheckEquivalenceStatus(goCtx context.Context, req *types.QueryCheckEquivalenceStatusRequest) (*types.QueryCheckEquivalenceStatusResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	exists, status, percent, equivalence := k.CheckEquivalenceStatusInternal(ctx, req.SourceSubjectId, req.TargetSubjectId)

	response := &types.QueryCheckEquivalenceStatusResponse{
		HasEquivalence:     exists,
		Status:             status,
		EquivalencePercent: percent,
	}

	if equivalence != nil {
		response.Equivalence = equivalence
		response.ContractVersion = equivalence.ContractVersion
		response.AnalysisTimestamp = equivalence.LastUpdateTimestamp
	}

	return response, nil
}

// GetPendingAnalysis queries equivalences awaiting contract analysis
func (k Keeper) GetPendingAnalysis(goCtx context.Context, req *types.QueryGetPendingAnalysisRequest) (*types.QueryGetPendingAnalysisResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	equivalences, pageRes, err := k.GetEquivalencesByStatusInternal(ctx, types.EquivalenceStatusPending, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetPendingAnalysisResponse{
		Equivalences: equivalences,
		Pagination:   pageRes,
	}, nil
}

// GetApprovedEquivalences queries equivalences with approved status
func (k Keeper) GetApprovedEquivalences(goCtx context.Context, req *types.QueryGetApprovedEquivalencesRequest) (*types.QueryGetApprovedEquivalencesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	equivalences, pageRes, err := k.GetEquivalencesByStatusInternal(ctx, types.EquivalenceStatusApproved, req.Pagination)
	if err != nil {
		return nil, err
	}

	// Apply minimum equivalence percent filter if specified
	if req.MinEquivalencePercent != "" {
		minPercent, parseErr := strconv.ParseFloat(req.MinEquivalencePercent, 64)
		if parseErr == nil {
			filteredEquivalences := make([]types.SubjectEquivalence, 0)
			for _, eq := range equivalences {
				if eq.EquivalencePercent != "" {
					if percent, err := strconv.ParseFloat(eq.EquivalencePercent, 64); err == nil && percent >= minPercent {
						filteredEquivalences = append(filteredEquivalences, eq)
					}
				}
			}
			equivalences = filteredEquivalences
		}
	}

	return &types.QueryGetApprovedEquivalencesResponse{
		Equivalences: equivalences,
		Pagination:   pageRes,
	}, nil
}

// GetRejectedEquivalences queries equivalences rejected by contract analysis
func (k Keeper) GetRejectedEquivalences(goCtx context.Context, req *types.QueryGetRejectedEquivalencesRequest) (*types.QueryGetRejectedEquivalencesResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	equivalences, pageRes, err := k.GetEquivalencesByStatusInternal(ctx, types.EquivalenceStatusRejected, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetRejectedEquivalencesResponse{
		Equivalences: equivalences,
		Pagination:   pageRes,
	}, nil
}

// GetEquivalencesByContract queries equivalences analyzed by a specific contract
func (k Keeper) GetEquivalencesByContract(goCtx context.Context, req *types.QueryGetEquivalencesByContractRequest) (*types.QueryGetEquivalencesByContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	equivalences, pageRes, err := k.GetEquivalencesByContractInternal(ctx, req.ContractAddress, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetEquivalencesByContractResponse{
		Equivalences: equivalences,
		Pagination:   pageRes,
	}, nil
}

// GetEquivalencesByContractVersion queries equivalences by contract version
func (k Keeper) GetEquivalencesByContractVersion(goCtx context.Context, req *types.QueryGetEquivalencesByContractVersionRequest) (*types.QueryGetEquivalencesByContractVersionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	equivalences, pageRes, err := k.GetEquivalencesByContractVersionInternal(ctx, req.ContractVersion, req.Pagination)
	if err != nil {
		return nil, err
	}

	return &types.QueryGetEquivalencesByContractVersionResponse{
		Equivalences: equivalences,
		Pagination:   pageRes,
	}, nil
}

// GetEquivalenceHistory queries the analysis history of equivalence requests for a subject
func (k Keeper) GetEquivalenceHistory(goCtx context.Context, req *types.QueryGetEquivalenceHistoryRequest) (*types.QueryGetEquivalenceHistoryResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Get equivalences for the subject as both source and target
	sourceEquivalences, _, err1 := k.GetEquivalencesBySourceSubjectInternal(ctx, req.SubjectId, req.Pagination, "")
	targetEquivalences, _, err2 := k.GetEquivalencesByTargetSubjectInternal(ctx, req.SubjectId, nil, "")

	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	// Combine and deduplicate results
	equivalenceMap := make(map[string]types.SubjectEquivalence)
	for _, eq := range sourceEquivalences {
		equivalenceMap[eq.Index] = eq
	}
	for _, eq := range targetEquivalences {
		equivalenceMap[eq.Index] = eq
	}

	var allEquivalences []types.SubjectEquivalence
	for _, eq := range equivalenceMap {
		allEquivalences = append(allEquivalences, eq)
	}

	return &types.QueryGetEquivalenceHistoryResponse{
		Equivalences: allEquivalences,
		Pagination:   nil, // We'll use the pagination from the first query
	}, nil
}

// GetEquivalenceStats queries statistics about automated equivalence analysis
func (k Keeper) GetEquivalenceStats(goCtx context.Context, req *types.QueryGetEquivalenceStatsRequest) (*types.QueryGetEquivalenceStatsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	stats := k.GetEquivalenceStatsInternal(ctx)
	return &stats, nil
}

// GetAnalysisMetadata queries detailed analysis metadata for an equivalence
func (k Keeper) GetAnalysisMetadata(goCtx context.Context, req *types.QueryGetAnalysisMetadataRequest) (*types.QueryGetAnalysisMetadataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	equivalence, found := k.GetSubjectEquivalence(ctx, req.EquivalenceId)
	if !found {
		return nil, types.ErrEquivalenceNotFound
	}

	return &types.QueryGetAnalysisMetadataResponse{
		AnalysisMetadata:  equivalence.AnalysisMetadata,
		ContractAddress:   equivalence.ContractAddress,
		ContractVersion:   equivalence.ContractVersion,
		AnalysisHash:      equivalence.AnalysisHash,
		AnalysisTimestamp: equivalence.LastUpdateTimestamp,
		AnalysisCount:     equivalence.AnalysisCount,
	}, nil
}

// VerifyAnalysisIntegrity verifies the integrity of an equivalence analysis
func (k Keeper) VerifyAnalysisIntegrity(goCtx context.Context, req *types.QueryVerifyAnalysisIntegrityRequest) (*types.QueryVerifyAnalysisIntegrityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	isValid, storedHash, calculatedHash := k.VerifyAnalysisIntegrityInternal(ctx, req.EquivalenceId)

	return &types.QueryVerifyAnalysisIntegrityResponse{
		IntegrityValid:        isValid,
		StoredHash:            storedHash,
		CalculatedHash:        calculatedHash,
		VerificationTimestamp: strconv.FormatInt(time.Now().Unix(), 10),
	}, nil
}
