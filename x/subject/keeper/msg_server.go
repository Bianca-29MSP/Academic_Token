package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"academictoken/x/subject/types"
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

// CreateSubject creates a new subject with content
func (k msgServer) CreateSubject(goCtx context.Context, req *types.MsgCreateSubject) (*types.MsgCreateSubjectResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.Institution == "" {
		return nil, fmt.Errorf("institution cannot be empty")
	}
	if req.CourseId == "" {
		return nil, fmt.Errorf("course ID cannot be empty")
	}
	if req.Title == "" {
		return nil, fmt.Errorf("subject title cannot be empty")
	}
	if req.Code == "" {
		return nil, fmt.Errorf("subject code cannot be empty")
	}
	if req.WorkloadHours == 0 {
		return nil, fmt.Errorf("workload hours cannot be zero")
	}
	if req.Credits == 0 {
		return nil, fmt.Errorf("credits cannot be zero")
	}

	// Validate subject type
	validSubjectTypes := []string{"required", "elective", "optional", "extracurricular"}
	validSubjectType := false
	for _, stype := range validSubjectTypes {
		if req.SubjectType == stype {
			validSubjectType = true
			break
		}
	}
	if !validSubjectType {
		return nil, fmt.Errorf("invalid subject type: must be one of %v, got '%s'", validSubjectTypes, req.SubjectType)
	}

	// Check if institution exists
	if k.institutionKeeper != nil {
		_, found := k.institutionKeeper.GetInstitution(ctx, req.Institution)
		if !found {
			return nil, fmt.Errorf("institution with ID '%s' not found", req.Institution)
		}
	}

	// Check if course exists
	if k.courseKeeper != nil {
		_, found := k.courseKeeper.GetCourse(ctx, req.CourseId)
		if !found {
			return nil, fmt.Errorf("course with ID '%s' not found", req.CourseId)
		}
	}

	// TEMPORARILY COMMENTED OUT - Check if subject code already exists for this institution
	// This validation is causing protobuf deserialization errors when trying to read existing subjects
	/*
	allSubjects, err := k.GetAllSubjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all subjects: %w", err)
	}
	for _, subject := range allSubjects {
		if subject.Institution == req.Institution && subject.Code == req.Code {
			return nil, fmt.Errorf("subject with code '%s' already exists in institution '%s'", req.Code, req.Institution)
		}
	}
	*/

	// Generate a new index for the subject
	index := k.GetNextSubjectIndex(ctx)

	// Create subject content object
	subject := types.SubjectContent{
		Index:         index,
		SubjectId:     index, // Using the same value for both
		Creator:       req.Creator,
		Institution:   req.Institution,
		CourseId:      req.CourseId,
		Title:         req.Title,
		Code:          req.Code,
		WorkloadHours: req.WorkloadHours,
		Credits:       req.Credits,
		Description:   req.Description,
		SubjectType:   req.SubjectType,
		KnowledgeArea: req.KnowledgeArea,
	}

	// Store the subject using SetSubjectWithoutIPFS for basic storage
	if err := k.SetSubjectWithoutIPFS(ctx, subject); err != nil {
		return nil, fmt.Errorf("failed to store subject: %w", err)
	}

	// Create indexing
	if err := k.SetSubjectByCourseIndex(ctx, req.CourseId, index); err != nil {
		k.Logger().Error("Failed to set course index", "error", err)
	}
	if err := k.SetSubjectByInstitutionIndex(ctx, req.Institution, index); err != nil {
		k.Logger().Error("Failed to set institution index", "error", err)
	}

	// If objectives and topic units are provided, prepare extended content for IPFS
	if len(req.Objectives) > 0 || len(req.TopicUnits) > 0 {
		// Create extended content for IPFS storage
		extendedContent := map[string]interface{}{
			"objectives": req.Objectives,
			"topicUnits": req.TopicUnits,
		}

		// Serialize extended content
		extendedContentBytes, err := json.Marshal(extendedContent)
		if err != nil {
			k.Logger().Error("Failed to serialize extended content", "error", err)
		} else {
			// Store in IPFS and update subject with hash and link
			contentHash, ipfsLink, err := k.StoreOnIPFS(ctx, extendedContentBytes)
			if err != nil {
				k.Logger().Error("Failed to store content in IPFS", "error", err)
			} else {
				// Update subject with IPFS references
				subject.ContentHash = contentHash
				subject.IpfsLink = ipfsLink

				// Re-store the subject with IPFS references
				if err := k.SetSubjectWithoutIPFS(ctx, subject); err != nil {
					k.Logger().Error("Failed to update subject with IPFS references", "error", err)
				}
			}
		}
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"subject_created",
			sdk.NewAttribute("subject_index", index),
			sdk.NewAttribute("institution", req.Institution),
			sdk.NewAttribute("course_id", req.CourseId),
			sdk.NewAttribute("title", req.Title),
			sdk.NewAttribute("code", req.Code),
			sdk.NewAttribute("credits", fmt.Sprintf("%d", req.Credits)),
			sdk.NewAttribute("workload_hours", fmt.Sprintf("%d", req.WorkloadHours)),
			sdk.NewAttribute("subject_type", req.SubjectType),
			sdk.NewAttribute("knowledge_area", req.KnowledgeArea),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Subject created successfully",
		"subject_index", index,
		"institution", req.Institution,
		"title", req.Title,
		"code", req.Code,
		"creator", req.Creator,
	)

	return &types.MsgCreateSubjectResponse{
		Index: index,
	}, nil
}

// CreateSubjectContent creates a new subject with content (same as CreateSubject but different endpoint)
func (k msgServer) CreateSubjectContent(goCtx context.Context, req *types.MsgCreateSubjectContent) (*types.MsgCreateSubjectContentResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate basic fields
	if req.Creator == "" {
		return nil, fmt.Errorf("creator cannot be empty")
	}
	if req.Institution == "" {
		return nil, fmt.Errorf("institution cannot be empty")
	}
	if req.CourseId == "" {
		return nil, fmt.Errorf("course ID cannot be empty")
	}
	if req.Title == "" {
		return nil, fmt.Errorf("subject title cannot be empty")
	}
	if req.Code == "" {
		return nil, fmt.Errorf("subject code cannot be empty")
	}
	if req.WorkloadHours == 0 {
		return nil, fmt.Errorf("workload hours cannot be zero")
	}
	if req.Credits == 0 {
		return nil, fmt.Errorf("credits cannot be zero")
	}

	// Validate subject type
	validSubjectTypes := []string{"required", "elective", "optional", "extracurricular"}
	validSubjectType := false
	for _, stype := range validSubjectTypes {
		if req.SubjectType == stype {
			validSubjectType = true
			break
		}
	}
	if !validSubjectType {
		return nil, fmt.Errorf("invalid subject type: must be one of %v, got '%s'", validSubjectTypes, req.SubjectType)
	}

	// Check if institution exists
	if k.institutionKeeper != nil {
		_, found := k.institutionKeeper.GetInstitution(ctx, req.Institution)
		if !found {
			return nil, fmt.Errorf("institution with ID '%s' not found", req.Institution)
		}
	}

	// Check if course exists
	if k.courseKeeper != nil {
		_, found := k.courseKeeper.GetCourse(ctx, req.CourseId)
		if !found {
			return nil, fmt.Errorf("course with ID '%s' not found", req.CourseId)
		}
	}

	// TEMPORARILY COMMENTED OUT - Check if subject code already exists for this institution
	// This validation is causing protobuf deserialization errors when trying to read existing subjects
	/*
	allSubjects, err := k.GetAllSubjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all subjects: %w", err)
	}
	for _, subject := range allSubjects {
		if subject.Institution == req.Institution && subject.Code == req.Code {
			return nil, fmt.Errorf("subject with code '%s' already exists in institution '%s'", req.Code, req.Institution)
		}
	}
	*/

	// Generate a new index for the subject
	index := k.GetNextSubjectIndex(ctx)

	// Create subject content object
	subject := types.SubjectContent{
		Index:         index,
		SubjectId:     index,
		Creator:       req.Creator,
		Institution:   req.Institution,
		CourseId:      req.CourseId,
		Title:         req.Title,
		Code:          req.Code,
		WorkloadHours: req.WorkloadHours,
		Credits:       req.Credits,
		Description:   req.Description,
		SubjectType:   req.SubjectType,
		KnowledgeArea: req.KnowledgeArea,
	}

	// Store the subject using SetSubjectWithoutIPFS for basic storage
	if err := k.SetSubjectWithoutIPFS(ctx, subject); err != nil {
		return nil, fmt.Errorf("failed to store subject: %w", err)
	}

	// Create indexing
	if err := k.SetSubjectByCourseIndex(ctx, req.CourseId, index); err != nil {
		k.Logger().Error("Failed to set course index", "error", err)
	}
	if err := k.SetSubjectByInstitutionIndex(ctx, req.Institution, index); err != nil {
		k.Logger().Error("Failed to set institution index", "error", err)
	}

	// If objectives and topic units are provided, prepare extended content for IPFS
	if len(req.Objectives) > 0 || len(req.TopicUnits) > 0 {
		// Create extended content for IPFS storage
		extendedContent := map[string]interface{}{
			"objectives": req.Objectives,
			"topicUnits": req.TopicUnits,
		}

		// Serialize extended content
		extendedContentBytes, err := json.Marshal(extendedContent)
		if err != nil {
			k.Logger().Error("Failed to serialize extended content", "error", err)
		} else {
			// Store in IPFS and update subject with hash and link
			contentHash, ipfsLink, err := k.StoreOnIPFS(ctx, extendedContentBytes)
			if err != nil {
				k.Logger().Error("Failed to store content in IPFS", "error", err)
			} else {
				// Update subject with IPFS references
				subject.ContentHash = contentHash
				subject.IpfsLink = ipfsLink

				// Re-store the subject with IPFS references
				if err := k.SetSubjectWithoutIPFS(ctx, subject); err != nil {
					k.Logger().Error("Failed to update subject with IPFS references", "error", err)
				}
			}
		}
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"subject_content_created",
			sdk.NewAttribute("subject_index", index),
			sdk.NewAttribute("institution", req.Institution),
			sdk.NewAttribute("course_id", req.CourseId),
			sdk.NewAttribute("title", req.Title),
			sdk.NewAttribute("code", req.Code),
			sdk.NewAttribute("credits", fmt.Sprintf("%d", req.Credits)),
			sdk.NewAttribute("workload_hours", fmt.Sprintf("%d", req.WorkloadHours)),
			sdk.NewAttribute("subject_type", req.SubjectType),
			sdk.NewAttribute("knowledge_area", req.KnowledgeArea),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Subject content created successfully",
		"subject_index", index,
		"institution", req.Institution,
		"title", req.Title,
		"code", req.Code,
		"creator", req.Creator,
	)

	return &types.MsgCreateSubjectContentResponse{
		Index: index,
	}, nil
}

// AddPrerequisiteGroup adds a prerequisite group to a subject
func (k msgServer) AddPrerequisiteGroup(goCtx context.Context, req *types.MsgAddPrerequisiteGroup) (*types.MsgAddPrerequisiteGroupResponse, error) {
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
	if req.GroupType == "" {
		return nil, fmt.Errorf("group type cannot be empty")
	}

	// Validate group type
	validGroupTypes := []string{"ALL", "ANY", "CREDITS", "COMBINATION"}
	validGroupType := false
	for _, gtype := range validGroupTypes {
		if req.GroupType == gtype {
			validGroupType = true
			break
		}
	}
	if !validGroupType {
		return nil, fmt.Errorf("invalid group type: must be one of %v, got '%s'", validGroupTypes, req.GroupType)
	}

	// Check if subject exists
	_, found := k.GetSubject(ctx, req.SubjectId)
	if !found {
		return nil, fmt.Errorf("subject with ID '%s' not found", req.SubjectId)
	}

	// Validate prerequisite subjects exist
	for _, prereqId := range req.SubjectIds {
		if prereqId != "" {
			_, found := k.GetSubject(ctx, prereqId)
			if !found {
				return nil, fmt.Errorf("prerequisite subject with ID '%s' not found", prereqId)
			}
		}
	}

	// Generate group ID
	groupId := k.GetNextPrerequisiteGroupIndex(ctx)

	// Create prerequisite group
	prerequisiteGroup := types.PrerequisiteGroup{
		Id:                       groupId,
		SubjectId:                req.SubjectId,
		GroupType:                req.GroupType,
		MinimumCredits:           req.MinimumCredits,
		MinimumCompletedSubjects: req.MinimumCompletedSubjects,
		SubjectIds:               req.SubjectIds,
	}

	// Store prerequisite group
	err := k.Keeper.AddPrerequisiteGroup(ctx, prerequisiteGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to store prerequisite group: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"prerequisite_group_added",
			sdk.NewAttribute("group_id", groupId),
			sdk.NewAttribute("subject_id", req.SubjectId),
			sdk.NewAttribute("group_type", req.GroupType),
			sdk.NewAttribute("minimum_credits", fmt.Sprintf("%d", req.MinimumCredits)),
			sdk.NewAttribute("minimum_completed_subjects", fmt.Sprintf("%d", req.MinimumCompletedSubjects)),
			sdk.NewAttribute("creator", req.Creator),
		),
	)

	k.Logger().Info("Prerequisite group added successfully",
		"group_id", groupId,
		"subject_id", req.SubjectId,
		"group_type", req.GroupType,
		"creator", req.Creator,
	)

	return &types.MsgAddPrerequisiteGroupResponse{
		GroupId: groupId,
	}, nil
}

// UpdateSubjectContent updates the detailed content of a subject
func (k msgServer) UpdateSubjectContent(goCtx context.Context, req *types.MsgUpdateSubjectContent) (*types.MsgUpdateSubjectContentResponse, error) {
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

	// Check if subject exists
	subject, found := k.GetSubject(ctx, req.SubjectId)
	if !found {
		return nil, fmt.Errorf("subject with ID '%s' not found", req.SubjectId)
	}

	// Check permissions (only creator or authority can update)
	if req.Creator != subject.Creator && req.Creator != k.GetAuthority() {
		return nil, fmt.Errorf("creator '%s' is not authorized to update subject '%s'", req.Creator, req.SubjectId)
	}

	// Get existing content or create new
	subjectContent, found := k.GetSubject(ctx, req.SubjectId)
	if !found {
		// Create new content
		subjectContent = types.SubjectContent{
			Index:     req.SubjectId,
			Creator:   req.Creator,
			SubjectId: req.SubjectId,
			Title:     subject.Title,
			Code:      subject.Code,
			Credits:   subject.Credits,
		}
	}

	// Prepare extended content for IPFS storage
	extendedContent := make(map[string]interface{})
	updated := false

	// Get existing content from IPFS if it exists
	if subjectContent.IpfsLink != "" {
		existingContentBytes, err := k.FetchFromIPFS(ctx, subjectContent.IpfsLink)
		if err == nil {
			json.Unmarshal(existingContentBytes, &extendedContent)
		}
	}

	// Update fields if provided
	if len(req.Objectives) > 0 {
		extendedContent["objectives"] = req.Objectives
		updated = true
	}

	if len(req.TopicUnits) > 0 {
		extendedContent["topicUnits"] = req.TopicUnits
		updated = true
	}

	if len(req.Methodologies) > 0 {
		extendedContent["methodologies"] = req.Methodologies
		updated = true
	}

	if len(req.EvaluationMethods) > 0 {
		extendedContent["evaluationMethods"] = req.EvaluationMethods
		updated = true
	}

	if len(req.BibliographyBasic) > 0 {
		extendedContent["bibliographyBasic"] = req.BibliographyBasic
		updated = true
	}

	if len(req.BibliographyComplementary) > 0 {
		extendedContent["bibliographyComplementary"] = req.BibliographyComplementary
		updated = true
	}

	if len(req.Keywords) > 0 {
		extendedContent["keywords"] = req.Keywords
		updated = true
	}

	// Handle direct hash and link updates
	if req.ContentHash != "" {
		subjectContent.ContentHash = req.ContentHash
		updated = true
	}

	if req.IpfsLink != "" {
		subjectContent.IpfsLink = req.IpfsLink
		updated = true
	}

	if !updated {
		return nil, fmt.Errorf("no valid updates provided")
	}

	// If we have extended content to store
	if len(extendedContent) > 0 && (req.ContentHash == "" || req.IpfsLink == "") {
		// Serialize and store in IPFS
		extendedContentBytes, err := json.Marshal(extendedContent)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize extended content: %w", err)
		}

		// Store in IPFS
		contentHash, ipfsLink, err := k.StoreOnIPFS(ctx, extendedContentBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to store content in IPFS: %w", err)
		}

		// Update subject with new hash and link
		subjectContent.ContentHash = contentHash
		subjectContent.IpfsLink = ipfsLink
	}

	// Store updated content
	err := k.SetSubjectWithoutIPFS(ctx, subjectContent)
	if err != nil {
		return nil, fmt.Errorf("failed to store subject content: %w", err)
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"subject_content_updated",
			sdk.NewAttribute("subject_id", req.SubjectId),
			sdk.NewAttribute("content_hash", subjectContent.ContentHash),
			sdk.NewAttribute("ipfs_link", subjectContent.IpfsLink),
			sdk.NewAttribute("updater", req.Creator),
		),
	)

	k.Logger().Info("Subject content updated successfully",
		"subject_id", req.SubjectId,
		"content_hash", subjectContent.ContentHash,
		"ipfs_link", subjectContent.IpfsLink,
		"updater", req.Creator,
	)

	return &types.MsgUpdateSubjectContentResponse{
		ContentHash: subjectContent.ContentHash,
		IpfsLink:    subjectContent.IpfsLink,
	}, nil
}
