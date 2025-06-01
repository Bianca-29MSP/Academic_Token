package keeper

import (
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"academictoken/x/subject/ipfs"
	"academictoken/x/subject/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger
		paramstore   paramtypes.Subspace
		ipfsClient   *ipfs.IPFSClient
		wasmQuerier  types.WasmQuerier

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority         string
		institutionKeeper types.InstitutionKeeper
		courseKeeper      types.CourseKeeper
	}
)

// AddPrerequisiteGroup adds a prerequisite group to the store
func (k Keeper) AddPrerequisiteGroup(ctx sdk.Context, group types.PrerequisiteGroup) error {
	store := k.storeService.OpenKVStore(ctx)

	// Create the key for this prerequisite group
	key := types.PrerequisiteGroupKey(group.Id)

	// Encode the prerequisite group
	value := k.cdc.MustMarshal(&group)

	// Store the prerequisite group
	if err := store.Set(key, value); err != nil {
		return err
	}

	// Add reference to subject's prerequisite groups
	subjectPrereqKey := types.SubjectPrerequisiteGroupKey(group.SubjectId, group.Id)
	if err := store.Set(subjectPrereqKey, []byte{1}); err != nil { // Just a marker
		return err
	}

	return nil
}

// GetNextPrerequisiteGroupIndex generates a new prerequisite group ID
func (k Keeper) GetNextPrerequisiteGroupIndex(ctx sdk.Context) string {
	store := k.storeService.OpenKVStore(ctx)

	// Get the current count
	countKey := []byte(types.PrerequisiteGroupCountKey)
	countBytes, err := store.Get(countKey)
	var count uint64 = 1

	if err == nil && countBytes != nil {
		count = sdk.BigEndianToUint64(countBytes) + 1
	}

	// Store the new count
	countBytes = sdk.Uint64ToBigEndian(count)
	if err := store.Set(countKey, countBytes); err != nil {
		k.logger.Error("failed to update prerequisite group counter", "error", err)
	}

	return fmt.Sprintf("prereq-group-%d", count)
}

func (k Keeper) SetSubjectWithoutIPFS(ctx sdk.Context, subject types.SubjectContent) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.SubjectKey(subject.Index)
	value := k.cdc.MustMarshal(&subject)

	if err := store.Set(key, value); err != nil {
		return err
	}

	return nil
}

// GetPrerequisiteGroup gets a prerequisite group by ID
func (k Keeper) GetPrerequisiteGroup(ctx sdk.Context, id string) (types.PrerequisiteGroup, bool) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.PrerequisiteGroupKey(id)

	bz, err := store.Get(key)
	if err != nil || bz == nil {
		return types.PrerequisiteGroup{}, false
	}

	var group types.PrerequisiteGroup
	k.cdc.MustUnmarshal(bz, &group)

	return group, true
}

// GetPrerequisiteGroupsBySubject gets all prerequisite groups for a subject
func (k Keeper) GetPrerequisiteGroupsBySubject(ctx sdk.Context, subjectId string) []types.PrerequisiteGroup {
	store := k.storeService.OpenKVStore(ctx)
	prefixKey := types.GetSubjectPrerequisiteGroupPrefix(subjectId)

	// Get an iterator over all keys with the prefix
	iterator, err := store.Iterator(prefixKey, nil)
	if err != nil {
		k.logger.Error("failed to create iterator", "error", err)
		return nil
	}
	defer iterator.Close()

	var groups []types.PrerequisiteGroup

	// Iterate through all references
	for ; iterator.Valid(); iterator.Next() {
		// Extract group ID from the key
		key := iterator.Key()
		groupId := string(key[len(prefixKey):])

		// Get the actual group
		group, found := k.GetPrerequisiteGroup(ctx, groupId)
		if found {
			groups = append(groups, group)
		}
	}

	return groups
}

// NewKeeper creates a new Subject Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	ps paramtypes.Subspace,
	wasmQuerier types.WasmQuerier,
	institutionKeeper types.InstitutionKeeper,
	courseKeeper types.CourseKeeper,
	authority string,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	// Set default parameters if not set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	// Create the IPFS client with default parameters
	// These can be updated later based on module parameters
	ipfsClient := ipfs.NewIPFSClient(
		"http://localhost:5001",
		"/tmp/ipfs-cache",
		true,
	)

	return Keeper{
		cdc:               cdc,
		storeService:      storeService,
		authority:         authority,
		logger:            logger,
		paramstore:        ps,
		ipfsClient:        ipfsClient,
		wasmQuerier:       wasmQuerier,
		institutionKeeper: institutionKeeper,
		courseKeeper:      courseKeeper,
	}
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// SetIPFSClient sets the IPFS client (useful for testing or config updates)
func (k *Keeper) SetIPFSClient(client *ipfs.IPFSClient) {
	k.ipfsClient = client
}

// UpdateIPFSClientFromParams updates the IPFS client based on current parameters
func (k *Keeper) UpdateIPFSClientFromParams(ctx sdk.Context) {
	params := k.GetParams(ctx)

	// Only update if IPFS is enabled
	if params.IpfsEnabled {
		k.ipfsClient = ipfs.NewIPFSClient(
			params.IpfsGateway,
			"/tmp/ipfs-cache", // This could be parameterized in the future
			true,              // Use HTTP when enabled
		)
	}
}

// SetSubject stores a subject following the hybrid storage model (MANDATORY)
func (k Keeper) SetSubject(ctx sdk.Context, subject types.SubjectContent, fullContent []byte) error {
	// Store in IPFS (MANDATORY)
	contentHash, ipfsLink, err := k.ipfsClient.Add(fullContent)
	if err != nil {
		return fmt.Errorf("error storing in IPFS: %v", err)
	}

	// Update subject with references
	subject.ContentHash = contentHash
	subject.IpfsLink = ipfsLink

	// Store on-chain (only essential data)
	store := k.storeService.OpenKVStore(ctx)
	key := types.SubjectKey(subject.Index)
	value := k.cdc.MustMarshal(&subject)
	if err := store.Set(key, value); err != nil {
		return err
	}

	return nil
}

// GetSubject retrieves a subject's basic data from the blockchain
func (k Keeper) GetSubject(ctx sdk.Context, index string) (types.SubjectContent, bool) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.SubjectKey(index)

	bz, err := store.Get(key)
	if err != nil || bz == nil {
		return types.SubjectContent{}, false
	}

	var subject types.SubjectContent
	k.cdc.MustUnmarshal(bz, &subject)
	return subject, true
}

// GetSubjectFull retrieves a subject with content from IPFS
func (k Keeper) GetSubjectFull(ctx sdk.Context, index string) (types.SubjectContent, map[string]interface{}, error) {
	// Get on-chain data
	subject, found := k.GetSubject(ctx, index)
	if !found {
		return types.SubjectContent{}, nil, types.ErrSubjectNotFound
	}

	// If no IPFS link, return only on-chain data
	if subject.IpfsLink == "" {
		return subject, nil, nil
	}

	// Get data from IPFS (MANDATORY)
	extendedContentBytes, err := k.ipfsClient.Get(subject.IpfsLink)
	if err != nil {
		return subject, nil, fmt.Errorf("error retrieving from IPFS: %v", err)
	}

	// Deserialize
	var extendedContent map[string]interface{}
	if err := json.Unmarshal(extendedContentBytes, &extendedContent); err != nil {
		return subject, nil, fmt.Errorf("error deserializing content: %v", err)
	}

	return subject, extendedContent, nil
}

// GetSubjectWithPrerequisites retrieves a subject with its prerequisite groups
func (k Keeper) GetSubjectWithPrerequisites(ctx sdk.Context, subjectId string) (types.SubjectWithPrerequisites, error) {
	// Get basic subject data
	subject, found := k.GetSubject(ctx, subjectId)
	if !found {
		return types.SubjectWithPrerequisites{}, types.ErrSubjectNotFound
	}

	// Get prerequisite groups
	prerequisiteGroups := k.GetPrerequisiteGroupsBySubject(ctx, subjectId)

	// Convert to pointers for the response struct
	var prerequisiteGroupPtrs []*types.PrerequisiteGroup
	for i := range prerequisiteGroups {
		prerequisiteGroupPtrs = append(prerequisiteGroupPtrs, &prerequisiteGroups[i])
	}

	// Create a copy of the subject to use as a pointer
	subjectCopy := subject

	return types.SubjectWithPrerequisites{
		Subject:            &subjectCopy,
		PrerequisiteGroups: prerequisiteGroupPtrs,
	}, nil
}

// GetNextSubjectIndex generates a new unique index for subjects
func (k Keeper) GetNextSubjectIndex(ctx sdk.Context) string {
	store := k.storeService.OpenKVStore(ctx)
	key := []byte(types.SubjectCountKey)

	// Get current counter
	bz, err := store.Get(key)
	var count uint64 = 1
	if err == nil && bz != nil {
		count = sdk.BigEndianToUint64(bz) + 1
	}

	// Update counter
	if err := store.Set(key, sdk.Uint64ToBigEndian(count)); err != nil {
		// Handle error - in production code, this should be propagated
		k.Logger().Error("failed to set subject counter", "error", err)
	}

	// Generate index based on counter
	return fmt.Sprintf("subject-%d", count)
}

// StoreOnIPFS allows other modules to store content in IPFS
func (k Keeper) StoreOnIPFS(ctx sdk.Context, content []byte) (string, string, error) {
	// Update IPFS client from parameters if needed
	k.UpdateIPFSClientFromParams(ctx)

	// Store in IPFS
	return k.ipfsClient.Add(content)
}

// FetchFromIPFS allows other modules to fetch content from IPFS
func (k Keeper) FetchFromIPFS(ctx sdk.Context, ipfsLink string) ([]byte, error) {
	// Update IPFS client from parameters if needed
	k.UpdateIPFSClientFromParams(ctx)

	// Fetch from IPFS
	return k.ipfsClient.Get(ipfsLink)
}

// CheckPrerequisitesViaContract verifies prerequisites using CosmWasm contract
func (k Keeper) CheckPrerequisitesViaContract(ctx sdk.Context, studentID string, subjectID string) (bool, []string, error) {
	// Get contract address from parameters
	params := k.GetParams(ctx)
	contractAddr := params.PrerequisiteValidatorContract

	if contractAddr == "" {
		return false, nil, fmt.Errorf("prerequisite validator contract not configured")
	}

	// Create query message for the contract
	msg := types.CheckPrerequisitesMsg{
		StudentId: studentID,
		SubjectId: subjectID,
	}

	msgBz, err := json.Marshal(msg)
	if err != nil {
		return false, nil, fmt.Errorf("error serializing message: %v", err)
	}

	// Execute the CosmWasm contract (MANDATORY)
	queryMsg := wasmtypes.QuerySmartContractStateRequest{
		Address:   contractAddr,
		QueryData: msgBz,
	}

	res, err := k.wasmQuerier.SmartContractState(sdk.WrapSDKContext(ctx), &queryMsg)
	if err != nil {
		return false, nil, fmt.Errorf("error executing contract: %v", err)
	}

	// Parse response
	var checkResult types.PrerequisiteCheckResponse
	err = json.Unmarshal(res.Data, &checkResult)
	if err != nil {
		return false, nil, fmt.Errorf("error interpreting response: %v", err)
	}

	return checkResult.IsEligible, checkResult.MissingPrerequisites, nil
}

// CheckEquivalenceViaContract verifies equivalence between subjects using CosmWasm contract
func (k Keeper) CheckEquivalenceViaContract(ctx sdk.Context, sourceSubjectID string, targetSubjectID string, forceRecalculate bool) (uint64, string, error) {
	// Get contract address from parameters
	params := k.GetParams(ctx)
	contractAddr := params.EquivalenceValidatorContract

	if contractAddr == "" {
		return 0, "", fmt.Errorf("equivalence validator contract not configured")
	}

	// Create query message for the contract
	msg := types.EquivalenceCheckMsg{
		SourceSubjectId:  sourceSubjectID,
		TargetSubjectId:  targetSubjectID,
		ForceRecalculate: forceRecalculate,
	}

	msgBz, err := json.Marshal(msg)
	if err != nil {
		return 0, "", fmt.Errorf("error serializing message: %v", err)
	}

	// Execute the CosmWasm contract (MANDATORY)
	queryMsg := wasmtypes.QuerySmartContractStateRequest{
		Address:   contractAddr,
		QueryData: msgBz,
	}

	res, err := k.wasmQuerier.SmartContractState(sdk.WrapSDKContext(ctx), &queryMsg)
	if err != nil {
		return 0, "", fmt.Errorf("error executing contract: %v", err)
	}

	// Parse response
	var equivalenceResponse types.EquivalenceResponse
	err = json.Unmarshal(res.Data, &equivalenceResponse)
	if err != nil {
		return 0, "", fmt.Errorf("error interpreting response: %v", err)
	}

	return equivalenceResponse.EquivalencePercent, equivalenceResponse.Status, nil
}

// SetSubjectByCourseIndex creates an index for subjects by course
func (k Keeper) SetSubjectByCourseIndex(ctx sdk.Context, courseId string, subjectIndex string) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.SubjectByCourseKey(courseId, subjectIndex)
	return store.Set(key, []byte{1}) // Just a marker
}

// SetSubjectByInstitutionIndex creates an index for subjects by institution
func (k Keeper) SetSubjectByInstitutionIndex(ctx sdk.Context, institutionId string, subjectIndex string) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.SubjectByInstitutionKey(institutionId, subjectIndex)
	return store.Set(key, []byte{1}) // Just a marker
}

// GetSubjectsByCourse retrieves all subjects for a specific course with pagination
func (k Keeper) GetSubjectsByCourse(ctx sdk.Context, courseId string) ([]types.SubjectContent, error) {
	store := k.storeService.OpenKVStore(ctx)
	prefix := types.GetSubjectByCoursePrefix(courseId)

	iterator, err := store.Iterator(prefix, nil)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var subjects []types.SubjectContent
	for ; iterator.Valid(); iterator.Next() {
		// Extract subject index from key
		key := iterator.Key()
		subjectIndex := string(key[len(prefix):])

		// Get the actual subject
		subject, found := k.GetSubject(ctx, subjectIndex)
		if found {
			subjects = append(subjects, subject)
		}
	}

	return subjects, nil
}

// GetSubjectsByInstitution retrieves all subjects for a specific institution
func (k Keeper) GetSubjectsByInstitution(ctx sdk.Context, institutionId string) ([]types.SubjectContent, error) {
	store := k.storeService.OpenKVStore(ctx)
	prefix := types.GetSubjectByInstitutionPrefix(institutionId)

	iterator, err := store.Iterator(prefix, nil)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var subjects []types.SubjectContent
	for ; iterator.Valid(); iterator.Next() {
		// Extract subject index from key
		key := iterator.Key()
		subjectIndex := string(key[len(prefix):])

		// Get the actual subject
		subject, found := k.GetSubject(ctx, subjectIndex)
		if found {
			subjects = append(subjects, subject)
		}
	}

	return subjects, nil
}

// GetAllSubjects retrieves all subjects with pagination support
func (k Keeper) GetAllSubjects(ctx sdk.Context) ([]types.SubjectContent, error) {
	store := k.storeService.OpenKVStore(ctx)
	prefix := types.GetSubjectPrefix()

	iterator, err := store.Iterator(prefix, nil)
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var subjects []types.SubjectContent
	for ; iterator.Valid(); iterator.Next() {
		var subject types.SubjectContent
		k.cdc.MustUnmarshal(iterator.Value(), &subject)
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

// RemoveSubjectByCourseIndex removes subject from course index
func (k Keeper) RemoveSubjectByCourseIndex(ctx sdk.Context, courseId string, subjectIndex string) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.SubjectByCourseKey(courseId, subjectIndex)
	return store.Delete(key)
}

// RemoveSubjectByInstitutionIndex removes subject from institution index
func (k Keeper) RemoveSubjectByInstitutionIndex(ctx sdk.Context, institutionId string, subjectIndex string) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.SubjectByInstitutionKey(institutionId, subjectIndex)
	return store.Delete(key)
}

// SubjectExists checks if a subject exists by its index
func (k Keeper) SubjectExists(ctx sdk.Context, index string) bool {
	_, found := k.GetSubject(ctx, index)
	return found
}

// GetSubjectsByCourseKeeper returns subjects for a course (implementing CourseKeeper interface)
func (k Keeper) GetSubjectsByCourseKeeper(ctx sdk.Context, courseID string) []types.SubjectContent {
	subjects, err := k.GetSubjectsByCourse(ctx, courseID)
	if err != nil {
		k.Logger().Error("failed to get subjects by course", "courseID", courseID, "error", err)
		return []types.SubjectContent{}
	}
	return subjects
}
