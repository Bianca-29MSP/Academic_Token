package keeper

import (
	"context"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"academictoken/x/academicnft/types"
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

// MintSubjectToken mints a new subject completion token
func (k msgServer) MintSubjectToken(goCtx context.Context, req *types.MsgMintSubjectToken) (*types.MsgMintSubjectTokenResponse, error) {
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
	if req.Student == "" {
		return nil, fmt.Errorf("student cannot be empty")
	}
	if req.CompletionDate == "" {
		return nil, fmt.Errorf("completion date cannot be empty")
	}
	if req.Grade == "" {
		return nil, fmt.Errorf("grade cannot be empty")
	}
	if req.IssuerInstitution == "" {
		return nil, fmt.Errorf("issuer institution cannot be empty")
	}
	if req.Semester == "" {
		return nil, fmt.Errorf("semester cannot be empty")
	}

	// Validate grade format and range
	grade, err := strconv.ParseFloat(req.Grade, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid grade format: %w", err)
	}
	if grade < 0 || grade > 100 {
		return nil, fmt.Errorf("grade must be between 0 and 100, got %f", grade)
	}

	// Check if token definition exists
	if k.tokenDefKeeper != nil {
		_, found := k.tokenDefKeeper.GetTokenDefinitionByIndex(ctx, req.TokenDefId)
		if !found {
			return nil, fmt.Errorf("token definition with ID '%s' not found", req.TokenDefId)
		}
	}

	// Check if student exists - convert string to AccAddress
	if k.studentKeeper != nil {
		studentAddr, err := sdk.AccAddressFromBech32(req.Student)
		if err != nil {
			return nil, fmt.Errorf("invalid student address format: %w", err)
		}
		_, found := k.studentKeeper.GetStudentByAddress(ctx, studentAddr)
		if !found {
			return nil, fmt.Errorf("student with address '%s' not found", req.Student)
		}
	}

	// Check if institution exists
	if k.institutionKeeper != nil {
		_, found := k.institutionKeeper.GetInstitution(ctx, req.IssuerInstitution)
		if !found {
			return nil, fmt.Errorf("institution with ID '%s' not found", req.IssuerInstitution)
		}
	}

	// Check if student already has a token for this subject (prevent duplicates)
	existingTokens := k.GetAllSubjectTokenInstances(ctx)
	for _, token := range existingTokens {
		if token.TokenDefId == req.TokenDefId && token.Student == req.Student {
			return nil, fmt.Errorf("student '%s' already has a token for token definition '%s'", req.Student, req.TokenDefId)
		}
	}

	// Generate a unique token instance ID
	tokenInstanceId := k.GenerateTokenInstanceID(ctx)

	// Create subject token instance
	tokenInstance := types.SubjectTokenInstance{
		Index:              tokenInstanceId,
		TokenDefId:         req.TokenDefId,
		Student:            req.Student,
		CompletionDate:     req.CompletionDate,
		Grade:              req.Grade,
		IssuerInstitution:  req.IssuerInstitution,
		Semester:           req.Semester,
		ProfessorSignature: req.ProfessorSignature,
	}

	// Store token instance
	k.SetSubjectTokenInstance(ctx, tokenInstance)

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"subject_token_minted",
			sdk.NewAttribute("token_instance_id", tokenInstanceId),
			sdk.NewAttribute("token_def_id", req.TokenDefId),
			sdk.NewAttribute("student", req.Student),
			sdk.NewAttribute("grade", req.Grade),
			sdk.NewAttribute("completion_date", req.CompletionDate),
			sdk.NewAttribute("semester", req.Semester),
			sdk.NewAttribute("issuer_institution", req.IssuerInstitution),
			sdk.NewAttribute("minted_date", time.Now().UTC().Format(time.RFC3339)),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Subject token minted successfully",
		"token_instance_id", tokenInstanceId,
		"token_def_id", req.TokenDefId,
		"student", req.Student,
		"grade", req.Grade,
		"issuer_institution", req.IssuerInstitution,
		"creator", req.Creator,
	)

	return &types.MsgMintSubjectTokenResponse{
		TokenInstanceId: tokenInstanceId,
	}, nil
}

// VerifyTokenInstance verifies the authenticity of a token instance
func (k msgServer) VerifyTokenInstance(goCtx context.Context, req *types.MsgVerifyTokenInstance) (*types.MsgVerifyTokenInstanceResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.TokenInstanceId == "" {
		return nil, fmt.Errorf("token instance ID cannot be empty")
	}

	// Get token instance
	tokenInstance, found := k.GetSubjectTokenInstance(ctx, req.TokenInstanceId)
	if !found {
		return &types.MsgVerifyTokenInstanceResponse{
			IsValid: false,
		}, nil
	}

	// Check if token definition exists
	isValid := true
	if k.tokenDefKeeper != nil {
		_, found := k.tokenDefKeeper.GetTokenDefinitionByIndex(ctx, tokenInstance.TokenDefId)
		if !found {
			isValid = false
		}
	}

	// Check if student exists - convert string to AccAddress
	if isValid && k.studentKeeper != nil {
		studentAddr, err := sdk.AccAddressFromBech32(tokenInstance.Student)
		if err != nil {
			isValid = false
		} else {
			_, found := k.studentKeeper.GetStudentByAddress(ctx, studentAddr)
			if !found {
				isValid = false
			}
		}
	}

	// Check if institution exists
	if isValid && k.institutionKeeper != nil {
		_, found := k.institutionKeeper.GetInstitution(ctx, tokenInstance.IssuerInstitution)
		if !found {
			isValid = false
		}
	}

	// Validate grade format
	if isValid {
		_, err := strconv.ParseFloat(tokenInstance.Grade, 64)
		if err != nil {
			isValid = false
		}
	}

	// Validate completion date format
	if isValid {
		_, err := time.Parse(time.RFC3339, tokenInstance.CompletionDate)
		if err != nil {
			// Try alternative date formats
			_, err = time.Parse("2006-01-02", tokenInstance.CompletionDate)
			if err != nil {
				isValid = false
			}
		}
	}

	// Update verification status if valid by setting token again
	// Note: Since there's no IsVerified field in the current struct,
	// we just perform validation without modifying the token

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"token_instance_verified",
			sdk.NewAttribute("token_instance_id", req.TokenInstanceId),
			sdk.NewAttribute("is_valid", fmt.Sprintf("%t", isValid)),
			sdk.NewAttribute("verifier", req.Creator),
		),
	)

	k.Logger().Info("Token instance verification completed",
		"token_instance_id", req.TokenInstanceId,
		"is_valid", isValid,
		"verifier", req.Creator,
	)

	return &types.MsgVerifyTokenInstanceResponse{
		IsValid: isValid,
	}, nil
}
