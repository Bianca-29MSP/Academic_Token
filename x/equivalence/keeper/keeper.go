package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/core/store"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"academictoken/x/equivalence/types"
)

type (
	Keeper struct {
		cdc          codec.BinaryCodec
		storeService store.KVStoreService
		logger       log.Logger

		// the address capable of executing a MsgUpdateParams message. Typically, this
		// should be the x/gov module account.
		authority string

		subjectKeeper     types.SubjectKeeper
		institutionKeeper types.InstitutionKeeper
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeService store.KVStoreService,
	logger log.Logger,
	authority string,
) Keeper {
	if _, err := sdk.AccAddressFromBech32(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address: %s", authority))
	}

	return Keeper{
		cdc:          cdc,
		storeService: storeService,
		authority:    authority,
		logger:       logger,
	}
}

// SetSubjectKeeper sets the subject keeper
func (k *Keeper) SetSubjectKeeper(subjectKeeper types.SubjectKeeper) {
	k.subjectKeeper = subjectKeeper
}

// SetInstitutionKeeper sets the institution keeper
func (k *Keeper) SetInstitutionKeeper(institutionKeeper types.InstitutionKeeper) {
	k.institutionKeeper = institutionKeeper
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger() log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetStore returns the module's KVStore
func (k Keeper) GetStore(ctx context.Context) storetypes.KVStore {
	storeAdapter := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return storeAdapter
}

// ============================================================================
// SubjectEquivalence CRUD operations
// ============================================================================

// SetSubjectEquivalence sets a subject equivalence in the store
func (k Keeper) SetSubjectEquivalence(ctx context.Context, equivalence types.SubjectEquivalence) {
	store := k.GetStore(ctx)
	b := k.cdc.MustMarshal(&equivalence)
	store.Set(types.KeyPrefix(types.SubjectEquivalenceKeyPrefix+equivalence.Index), b)

	// Set secondary indexes
	k.setEquivalenceIndexes(ctx, equivalence)
}

// GetSubjectEquivalence returns a subject equivalence from its index
func (k Keeper) GetSubjectEquivalence(ctx context.Context, index string) (val types.SubjectEquivalence, found bool) {
	store := k.GetStore(ctx)

	b := store.Get(types.KeyPrefix(types.SubjectEquivalenceKeyPrefix + index))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveSubjectEquivalence removes a subject equivalence from the store
func (k Keeper) RemoveSubjectEquivalence(ctx context.Context, index string) {
	// Get the equivalence first to remove indexes
	equivalence, found := k.GetSubjectEquivalence(ctx, index)
	if !found {
		return
	}

	store := k.GetStore(ctx)
	store.Delete(types.KeyPrefix(types.SubjectEquivalenceKeyPrefix + index))

	// Remove secondary indexes
	k.removeEquivalenceIndexes(ctx, equivalence)
}

// GetAllSubjectEquivalences returns all subject equivalences
func (k Keeper) GetAllSubjectEquivalences(ctx context.Context) (list []types.SubjectEquivalence) {
	store := k.GetStore(ctx)
	iterator := storetypes.KVStorePrefixIterator(store, types.KeyPrefix(types.SubjectEquivalenceKeyPrefix))

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.SubjectEquivalence
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// ============================================================================
// SubjectEquivalence Query operations with pagination
// ============================================================================

// GetSubjectEquivalencesPaginated returns paginated subject equivalences
func (k Keeper) GetSubjectEquivalencesPaginated(ctx context.Context, pageReq *query.PageRequest, statusFilter string) ([]types.SubjectEquivalence, *query.PageResponse, error) {
	store := k.GetStore(ctx)
	var equivalences []types.SubjectEquivalence

	pageRes, err := query.Paginate(store, pageReq, func(key []byte, value []byte) error {
		var equivalence types.SubjectEquivalence
		if err := k.cdc.Unmarshal(value, &equivalence); err != nil {
			return err
		}
		equivalences = append(equivalences, equivalence)
		return nil
	})

	return equivalences, pageRes, err
}

// GetEquivalencesBySourceSubjectInternal returns equivalences for a source subject
func (k Keeper) GetEquivalencesBySourceSubjectInternal(ctx context.Context, sourceSubjectId string, pageReq *query.PageRequest, statusFilter string) ([]types.SubjectEquivalence, *query.PageResponse, error) {
	store := k.GetStore(ctx)
	var equivalences []types.SubjectEquivalence

	pageRes, err := query.Paginate(store, pageReq, func(key []byte, value []byte) error {
		// The value in secondary index is the equivalence index, not the full object
		equivalenceIndex := string(value)
		equivalence, found := k.GetSubjectEquivalence(ctx, equivalenceIndex)
		if !found {
			return nil // Skip if not found
		}

		// Apply status filter if specified
		if statusFilter != "" && equivalence.EquivalenceStatus != statusFilter {
			return nil
		}

		equivalences = append(equivalences, equivalence)
		return nil
	})

	return equivalences, pageRes, err
}

// GetEquivalencesByTargetSubjectInternal returns equivalences for a target subject
func (k Keeper) GetEquivalencesByTargetSubjectInternal(ctx context.Context, targetSubjectId string, pageReq *query.PageRequest, statusFilter string) ([]types.SubjectEquivalence, *query.PageResponse, error) {
	store := k.GetStore(ctx)
	var equivalences []types.SubjectEquivalence

	pageRes, err := query.Paginate(store, pageReq, func(key []byte, value []byte) error {
		equivalenceIndex := string(value)
		equivalence, found := k.GetSubjectEquivalence(ctx, equivalenceIndex)
		if !found {
			return nil
		}

		if statusFilter != "" && equivalence.EquivalenceStatus != statusFilter {
			return nil
		}

		equivalences = append(equivalences, equivalence)
		return nil
	})

	return equivalences, pageRes, err
}

// GetEquivalencesByInstitutionInternal returns equivalences for an institution
func (k Keeper) GetEquivalencesByInstitutionInternal(ctx context.Context, institutionId string, pageReq *query.PageRequest, statusFilter string) ([]types.SubjectEquivalence, *query.PageResponse, error) {
	store := k.GetStore(ctx)
	var equivalences []types.SubjectEquivalence

	pageRes, err := query.Paginate(store, pageReq, func(key []byte, value []byte) error {
		equivalenceIndex := string(value)
		equivalence, found := k.GetSubjectEquivalence(ctx, equivalenceIndex)
		if !found {
			return nil
		}

		if statusFilter != "" && equivalence.EquivalenceStatus != statusFilter {
			return nil
		}

		equivalences = append(equivalences, equivalence)
		return nil
	})

	return equivalences, pageRes, err
}

// GetEquivalencesByStatusInternal returns equivalences by status
func (k Keeper) GetEquivalencesByStatusInternal(ctx context.Context, status string, pageReq *query.PageRequest) ([]types.SubjectEquivalence, *query.PageResponse, error) {
	store := k.GetStore(ctx)
	var equivalences []types.SubjectEquivalence

	pageRes, err := query.Paginate(store, pageReq, func(key []byte, value []byte) error {
		equivalenceIndex := string(value)
		equivalence, found := k.GetSubjectEquivalence(ctx, equivalenceIndex)
		if !found {
			return nil
		}

		equivalences = append(equivalences, equivalence)
		return nil
	})

	return equivalences, pageRes, err
}

// GetEquivalencesByContractInternal returns equivalences analyzed by a specific contract
func (k Keeper) GetEquivalencesByContractInternal(ctx context.Context, contractAddress string, pageReq *query.PageRequest) ([]types.SubjectEquivalence, *query.PageResponse, error) {
	store := k.GetStore(ctx)
	var equivalences []types.SubjectEquivalence

	pageRes, err := query.Paginate(store, pageReq, func(key []byte, value []byte) error {
		equivalenceIndex := string(value)
		equivalence, found := k.GetSubjectEquivalence(ctx, equivalenceIndex)
		if !found {
			return nil
		}

		equivalences = append(equivalences, equivalence)
		return nil
	})

	return equivalences, pageRes, err
}

// GetEquivalencesByContractVersionInternal returns equivalences by contract version
func (k Keeper) GetEquivalencesByContractVersionInternal(ctx context.Context, contractVersion string, pageReq *query.PageRequest) ([]types.SubjectEquivalence, *query.PageResponse, error) {
	store := k.GetStore(ctx)
	var equivalences []types.SubjectEquivalence

	pageRes, err := query.Paginate(store, pageReq, func(key []byte, value []byte) error {
		equivalenceIndex := string(value)
		equivalence, found := k.GetSubjectEquivalence(ctx, equivalenceIndex)
		if !found {
			return nil
		}

		equivalences = append(equivalences, equivalence)
		return nil
	})

	return equivalences, pageRes, err
}

// ============================================================================
// Business Logic Functions
// ============================================================================

// CheckEquivalenceStatusInternal checks if equivalence exists between two subjects
func (k Keeper) CheckEquivalenceStatusInternal(ctx context.Context, sourceSubjectId, targetSubjectId string) (bool, string, string, *types.SubjectEquivalence) {
	index := types.GenerateEquivalenceIndex(sourceSubjectId, targetSubjectId)
	equivalence, found := k.GetSubjectEquivalence(ctx, index)

	if !found {
		return false, "", "", nil
	}

	return true, equivalence.EquivalenceStatus, equivalence.EquivalencePercent, &equivalence
}

// CreateEquivalenceRequest creates a new equivalence request
func (k Keeper) CreateEquivalenceRequest(ctx context.Context, sourceSubjectId, targetInstitution, targetSubjectId string, forceRecalculation bool) (string, error) {
	index := types.GenerateEquivalenceIndex(sourceSubjectId, targetSubjectId)

	// Check if equivalence already exists
	existingEquivalence, found := k.GetSubjectEquivalence(ctx, index)
	if found && !forceRecalculation {
		return index, fmt.Errorf("equivalence already exists between subjects")
	}

	// Create or update equivalence
	now := strconv.FormatInt(time.Now().Unix(), 10)

	equivalence := types.SubjectEquivalence{
		Index:               index,
		SourceSubjectId:     sourceSubjectId,
		TargetInstitution:   targetInstitution,
		TargetSubjectId:     targetSubjectId,
		EquivalenceStatus:   types.EquivalenceStatusPending,
		AnalysisCount:       1,
		EquivalencePercent:  "",
		AnalysisMetadata:    "",
		ContractAddress:     "",
		LastUpdateTimestamp: now,
		RequestTimestamp:    now,
		AnalysisHash:        "",
		ContractVersion:     "",
	}

	// If updating existing equivalence
	if found {
		equivalence.AnalysisCount = existingEquivalence.AnalysisCount + 1
		equivalence.RequestTimestamp = existingEquivalence.RequestTimestamp
	}

	k.SetSubjectEquivalence(ctx, equivalence)
	return index, nil
}

// UpdateEquivalenceAnalysis updates equivalence after contract analysis
func (k Keeper) UpdateEquivalenceAnalysis(ctx context.Context, equivalenceId, contractAddress, equivalencePercent, analysisMetadata, contractVersion string) error {
	equivalence, found := k.GetSubjectEquivalence(ctx, equivalenceId)
	if !found {
		return fmt.Errorf("equivalence not found")
	}

	// Calculate analysis hash for integrity
	analysisHash := k.calculateAnalysisHash(analysisMetadata, equivalencePercent, contractAddress)

	// Determine status based on equivalence percentage
	status := k.determineStatusFromPercent(equivalencePercent)

	// Update equivalence
	now := strconv.FormatInt(time.Now().Unix(), 10)
	equivalence.EquivalenceStatus = status
	equivalence.EquivalencePercent = equivalencePercent
	equivalence.AnalysisMetadata = analysisMetadata
	equivalence.ContractAddress = contractAddress
	equivalence.ContractVersion = contractVersion
	equivalence.AnalysisHash = analysisHash
	equivalence.LastUpdateTimestamp = now

	k.SetSubjectEquivalence(ctx, equivalence)
	return nil
}

// GetEquivalenceStatsInternal returns statistics about equivalences
func (k Keeper) GetEquivalenceStatsInternal(ctx context.Context) types.QueryGetEquivalenceStatsResponse {
	allEquivalences := k.GetAllSubjectEquivalences(ctx)

	stats := types.QueryGetEquivalenceStatsResponse{
		TotalEquivalences:             uint64(len(allEquivalences)),
		PendingAnalysis:               0,
		ApprovedEquivalences:          0,
		RejectedEquivalences:          0,
		ErrorEquivalences:             0,
		AverageEquivalencePercent:     "0.00",
		TotalInstitutionsInvolved:     0,
		TotalSubjectsWithEquivalences: 0,
		TotalContractAnalyses:         0,
		ActiveContractVersions:        []string{},
	}

	if len(allEquivalences) == 0 {
		return stats
	}

	// Count by status and collect data
	institutionSet := make(map[string]bool)
	subjectSet := make(map[string]bool)
	contractVersionSet := make(map[string]bool)
	var totalPercent float64
	var countWithPercent int
	var totalAnalyses uint64

	for _, eq := range allEquivalences {
		// Count by status
		switch eq.EquivalenceStatus {
		case types.EquivalenceStatusPending:
			stats.PendingAnalysis++
		case types.EquivalenceStatusApproved:
			stats.ApprovedEquivalences++
		case types.EquivalenceStatusRejected:
			stats.RejectedEquivalences++
		case types.EquivalenceStatusError:
			stats.ErrorEquivalences++
		}

		// Collect unique institutions
		institutionSet[eq.TargetInstitution] = true

		// Collect unique subjects
		subjectSet[eq.SourceSubjectId] = true
		subjectSet[eq.TargetSubjectId] = true

		// Collect contract versions
		if eq.ContractVersion != "" {
			contractVersionSet[eq.ContractVersion] = true
		}

		// Calculate average percentage
		if eq.EquivalencePercent != "" {
			if percent, err := strconv.ParseFloat(eq.EquivalencePercent, 64); err == nil {
				totalPercent += percent
				countWithPercent++
			}
		}

		// Count total analyses
		totalAnalyses += eq.AnalysisCount
	}

	// Set calculated values
	stats.TotalInstitutionsInvolved = uint64(len(institutionSet))
	stats.TotalSubjectsWithEquivalences = uint64(len(subjectSet))
	stats.TotalContractAnalyses = totalAnalyses

	// Calculate average percentage
	if countWithPercent > 0 {
		avgPercent := totalPercent / float64(countWithPercent)
		stats.AverageEquivalencePercent = fmt.Sprintf("%.2f", avgPercent)
	}

	// Convert contract versions to slice
	for version := range contractVersionSet {
		stats.ActiveContractVersions = append(stats.ActiveContractVersions, version)
	}

	return stats
}

// VerifyAnalysisIntegrityInternal verifies the integrity of an equivalence analysis
func (k Keeper) VerifyAnalysisIntegrityInternal(ctx context.Context, equivalenceId string) (bool, string, string) {
	equivalence, found := k.GetSubjectEquivalence(ctx, equivalenceId)
	if !found {
		return false, "", ""
	}

	storedHash := equivalence.AnalysisHash
	calculatedHash := k.calculateAnalysisHash(equivalence.AnalysisMetadata, equivalence.EquivalencePercent, equivalence.ContractAddress)

	return storedHash == calculatedHash, storedHash, calculatedHash
}

// SetSubjectEquivalenceCount sets the subject equivalence count
func (k Keeper) SetSubjectEquivalenceCount(ctx sdk.Context, count uint64) {
	store := k.GetStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(types.KeyPrefix(types.SubjectEquivalenceCountKey), bz)
}

// GetSubjectEquivalenceCount gets the subject equivalence count
func (k Keeper) GetSubjectEquivalenceCount(ctx sdk.Context) uint64 {
	store := k.GetStore(ctx)
	bz := store.Get(types.KeyPrefix(types.SubjectEquivalenceCountKey))
	if bz == nil {
		return 0
	}
	return binary.BigEndian.Uint64(bz)
}

// ============================================================================
// Helper Functions
// ============================================================================

// setEquivalenceIndexes sets secondary indexes for the equivalence
func (k Keeper) setEquivalenceIndexes(ctx context.Context, equivalence types.SubjectEquivalence) {
	store := k.GetStore(ctx)

	// Index by source subject
	sourceKey := append(types.SubjectEquivalenceBySourceKey(equivalence.SourceSubjectId), []byte(equivalence.Index)...)
	store.Set(sourceKey, []byte(equivalence.Index))

	// Index by target subject
	targetKey := append(types.SubjectEquivalenceByTargetKey(equivalence.TargetSubjectId), []byte(equivalence.Index)...)
	store.Set(targetKey, []byte(equivalence.Index))

	// Index by institution
	institutionKey := append(types.SubjectEquivalenceByInstitutionKey(equivalence.TargetInstitution), []byte(equivalence.Index)...)
	store.Set(institutionKey, []byte(equivalence.Index))

	// Index by status
	statusKey := append(types.SubjectEquivalenceByStatusKey(equivalence.EquivalenceStatus), []byte(equivalence.Index)...)
	store.Set(statusKey, []byte(equivalence.Index))

	// Index by contract address (if set)
	if equivalence.ContractAddress != "" {
		contractKey := append(types.SubjectEquivalenceByContractKey(equivalence.ContractAddress), []byte(equivalence.Index)...)
		store.Set(contractKey, []byte(equivalence.Index))
	}

	// Index by contract version (if set)
	if equivalence.ContractVersion != "" {
		versionKey := append(types.SubjectEquivalenceByContractVersionKey(equivalence.ContractVersion), []byte(equivalence.Index)...)
		store.Set(versionKey, []byte(equivalence.Index))
	}
}

// removeEquivalenceIndexes removes secondary indexes for the equivalence
func (k Keeper) removeEquivalenceIndexes(ctx context.Context, equivalence types.SubjectEquivalence) {
	store := k.GetStore(ctx)

	// Remove all indexes
	sourceKey := append(types.SubjectEquivalenceBySourceKey(equivalence.SourceSubjectId), []byte(equivalence.Index)...)
	store.Delete(sourceKey)

	targetKey := append(types.SubjectEquivalenceByTargetKey(equivalence.TargetSubjectId), []byte(equivalence.Index)...)
	store.Delete(targetKey)

	institutionKey := append(types.SubjectEquivalenceByInstitutionKey(equivalence.TargetInstitution), []byte(equivalence.Index)...)
	store.Delete(institutionKey)

	statusKey := append(types.SubjectEquivalenceByStatusKey(equivalence.EquivalenceStatus), []byte(equivalence.Index)...)
	store.Delete(statusKey)

	if equivalence.ContractAddress != "" {
		contractKey := append(types.SubjectEquivalenceByContractKey(equivalence.ContractAddress), []byte(equivalence.Index)...)
		store.Delete(contractKey)
	}

	if equivalence.ContractVersion != "" {
		versionKey := append(types.SubjectEquivalenceByContractVersionKey(equivalence.ContractVersion), []byte(equivalence.Index)...)
		store.Delete(versionKey)
	}
}

// calculateAnalysisHash calculates a hash for analysis integrity verification
func (k Keeper) calculateAnalysisHash(metadata, percent, contractAddress string) string {
	data := metadata + percent + contractAddress
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// determineStatusFromPercent determines equivalence status based on percentage
func (k Keeper) determineStatusFromPercent(percentStr string) string {
	if percentStr == "" {
		return types.EquivalenceStatusPending
	}

	percent, err := strconv.ParseFloat(percentStr, 64)
	if err != nil {
		return types.EquivalenceStatusError
	}

	threshold, _ := strconv.ParseFloat(types.MinApprovalThreshold, 64)
	if percent >= threshold {
		return types.EquivalenceStatusApproved
	}

	return types.EquivalenceStatusRejected
}

// ============================================================================
// HARDCODED CONFIGURATION METHODS (FIXED)
// ============================================================================

// GetHardcodedConfiguration returns all hardcoded configuration values
func (k Keeper) GetHardcodedConfiguration(ctx context.Context) map[string]interface{} {
	config := map[string]interface{}{
		"equivalence_contract_address": types.GetHardcodedEquivalenceContractAddress(),
		"ipfs_gateway":                types.GetHardcodedIPFSGateway(),
		"ipfs_enabled":                types.IsHardcodedIPFSEnabled(),
		"min_approval_threshold":      types.GetHardcodedMinApprovalThreshold(),
		"max_analysis_retries":        types.GetHardcodedMaxAnalysisRetries(),
		"analysis_timeout_seconds":    types.GetHardcodedAnalysisTimeoutSeconds(),
		"require_contract_auth":       types.IsHardcodedContractAuthRequired(),
		"admin":                       types.GetHardcodedAdmin(),
	}
	
	k.Logger().Debug("Retrieved hardcoded configuration", "config", config)
	return config
}

// GetContractConfigurationInternal returns contract-specific hardcoded configuration
func (k Keeper) GetContractConfigurationInternal(ctx context.Context) types.ContractInfo {
	info := types.ContractInfo{
		Address:         types.GetHardcodedEquivalenceContractAddress(),
		Version:         "v1.0.0",
		Status:          "active",
		LastUpdated:     time.Now().Format(time.RFC3339),
		AnalysisCount:   0,
		SuccessRate:     "100.0",
		AverageGasUsed:  150000,
	}
	
	k.Logger().Debug("Retrieved hardcoded contract configuration", "info", info)
	return info
}

// GetAnalysisConfigurationInternal returns analysis-specific hardcoded configuration
func (k Keeper) GetAnalysisConfigurationInternal(ctx context.Context) types.AnalysisConfig {
	config := types.GetDefaultAnalysisConfig()
	k.Logger().Debug("Retrieved hardcoded analysis configuration", "config", config)
	return config
}

// ValidateHardcodedConfigurationInternal validates all hardcoded configuration values
func (k Keeper) ValidateHardcodedConfigurationInternal(ctx context.Context) error {
	k.Logger().Info("Validating hardcoded configuration")
	
	// Validate contract address
	contractAddr := types.GetHardcodedEquivalenceContractAddress()
	if !types.IsValidContractAddress(contractAddr) {
		k.Logger().Error("Invalid hardcoded contract address", "address", contractAddr)
		return fmt.Errorf("invalid hardcoded contract address")
	}
	
	// Validate IPFS gateway
	ipfsGateway := types.GetHardcodedIPFSGateway()
	if ipfsGateway == "" {
		k.Logger().Error("Invalid hardcoded IPFS gateway", "gateway", ipfsGateway)
		return fmt.Errorf("invalid hardcoded IPFS gateway: cannot be empty")
	}
	
	// Validate min approval threshold
	threshold := types.GetHardcodedMinApprovalThreshold()
	if !types.IsValidEquivalencePercent(threshold) {
		k.Logger().Error("Invalid hardcoded min approval threshold", "threshold", threshold)
		return fmt.Errorf("invalid hardcoded min approval threshold")
	}
	
	// Validate max retries
	maxRetries := types.GetHardcodedMaxAnalysisRetries()
	if maxRetries == 0 {
		k.Logger().Error("Invalid hardcoded max analysis retries", "retries", maxRetries)
		return fmt.Errorf("invalid hardcoded max analysis retries: cannot be zero")
	}
	
	// Validate timeout
	timeout := types.GetHardcodedAnalysisTimeoutSeconds()
	if timeout == 0 {
		k.Logger().Error("Invalid hardcoded analysis timeout", "timeout", timeout)
		return fmt.Errorf("invalid hardcoded analysis timeout: cannot be zero")
	}
	
	k.Logger().Info("All hardcoded configuration values validated successfully")
	return nil
}

// GetSystemStatus returns the overall system status based on hardcoded configuration
func (k Keeper) GetSystemStatus(ctx context.Context) map[string]string {
	status := map[string]string{
		"module_status":          "active",
		"configuration_source":   "hardcoded",
		"contract_available":     fmt.Sprintf("%t", k.IsContractAvailable(ctx)),
		"ipfs_enabled":          fmt.Sprintf("%t", types.IsHardcodedIPFSEnabled()),
		"contract_auth_required": fmt.Sprintf("%t", types.IsHardcodedContractAuthRequired()),
		"last_validation":       time.Now().Format(time.RFC3339),
	}
	
	// Add contract information
	contractAddr := types.GetHardcodedEquivalenceContractAddress()
	if contractAddr != "" {
		status["contract_address"] = contractAddr
		status["contract_version"] = "v1.0.0"
	}
	
	// Add IPFS information
	if types.IsHardcodedIPFSEnabled() {
		status["ipfs_gateway"] = types.GetHardcodedIPFSGateway()
	}
	
	k.Logger().Debug("Retrieved system status", "status", status)
	return status
}

// LogHardcodedConfiguration logs all hardcoded configuration for debugging
func (k Keeper) LogHardcodedConfiguration(ctx context.Context) {
	k.Logger().Info("=== HARDCODED CONFIGURATION ===")
	k.Logger().Info("Contract Address", "value", types.GetHardcodedEquivalenceContractAddress())
	k.Logger().Info("IPFS Gateway", "value", types.GetHardcodedIPFSGateway())
	k.Logger().Info("IPFS Enabled", "value", types.IsHardcodedIPFSEnabled())
	k.Logger().Info("Min Approval Threshold", "value", types.GetHardcodedMinApprovalThreshold())
	k.Logger().Info("Max Analysis Retries", "value", types.GetHardcodedMaxAnalysisRetries())
	k.Logger().Info("Analysis Timeout Seconds", "value", types.GetHardcodedAnalysisTimeoutSeconds())
	k.Logger().Info("Contract Auth Required", "value", types.IsHardcodedContractAuthRequired())
	k.Logger().Info("Admin Address", "value", types.GetHardcodedAdmin())
	k.Logger().Info("=== END HARDCODED CONFIGURATION ===")
}
