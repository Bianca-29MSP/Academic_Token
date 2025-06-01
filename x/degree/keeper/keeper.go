package keeper

import (
	"encoding/binary"
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"academictoken/x/degree/types"
)

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace

		// Keeper dependencies
		studentKeeper     types.StudentKeeper
		curriculumKeeper  types.CurriculumKeeper
		academicNFTKeeper types.AcademicNFTKeeper
		wasmKeeper        types.WasmKeeper
		institutionKeeper types.InstitutionKeeper

		// authority is the address capable of executing a MsgUpdateParams message
		authority string
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	studentKeeper types.StudentKeeper,
	curriculumKeeper types.CurriculumKeeper,
	academicNFTKeeper types.AcademicNFTKeeper,
	wasmKeeper types.WasmKeeper,
	institutionKeeper types.InstitutionKeeper,
	authority string,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,

		studentKeeper:     studentKeeper,
		curriculumKeeper:  curriculumKeeper,
		academicNFTKeeper: academicNFTKeeper,
		wasmKeeper:        wasmKeeper,
		institutionKeeper: institutionKeeper,
		authority:         authority,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAuthority returns the module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Degree methods
func (k Keeper) SetDegree(ctx sdk.Context, degree types.Degree) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&degree)
	store.Set(append(types.DegreePrefix, []byte(degree.Index)...), b)
}

func (k Keeper) GetDegree(ctx sdk.Context, index string) (val types.Degree, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(append(types.DegreePrefix, []byte(index)...))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveDegree(ctx sdk.Context, index string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(append(types.DegreePrefix, []byte(index)...))
}

func (k Keeper) GetAllDegree(ctx sdk.Context) (list []types.Degree) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.DegreePrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Degree
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetDegreeByStudentAndCurriculum(ctx sdk.Context, studentId, curriculumId string) (val types.Degree, found bool) {
	degrees := k.GetAllDegree(ctx)
	for _, degree := range degrees {
		if degree.Student == studentId && degree.CourseId == curriculumId {
			return degree, true
		}
	}
	return val, false
}

// DegreeRequest methods
func (k Keeper) SetDegreeRequest(ctx sdk.Context, degreeRequest types.DegreeRequest) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&degreeRequest)
	store.Set(append(types.DegreeRequestPrefix, []byte(degreeRequest.Id)...), b)
}

func (k Keeper) GetDegreeRequest(ctx sdk.Context, id string) (val types.DegreeRequest, found bool) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(append(types.DegreeRequestPrefix, []byte(id)...))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveDegreeRequest(ctx sdk.Context, id string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(append(types.DegreeRequestPrefix, []byte(id)...))
}

func (k Keeper) GetAllDegreeRequest(ctx sdk.Context) (list []types.DegreeRequest) {
	store := ctx.KVStore(k.storeKey)
	iterator := storetypes.KVStorePrefixIterator(store, types.DegreeRequestPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.DegreeRequest
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// Helper methods for indexing
func (k Keeper) GetDegreesByStudent(ctx sdk.Context, studentId string) []types.Degree {
	var degrees []types.Degree
	allDegrees := k.GetAllDegree(ctx)

	for _, degree := range allDegrees {
		if degree.Student == studentId {
			degrees = append(degrees, degree)
		}
	}

	return degrees
}

func (k Keeper) GetDegreesByInstitution(ctx sdk.Context, institutionId string) []types.Degree {
	var degrees []types.Degree
	allDegrees := k.GetAllDegree(ctx)

	for _, degree := range allDegrees {
		if degree.Institution == institutionId {
			degrees = append(degrees, degree)
		}
	}

	return degrees
}

func (k Keeper) GetDegreeRequestsByStatus(ctx sdk.Context, status string) []types.DegreeRequest {
	var requests []types.DegreeRequest
	allRequests := k.GetAllDegreeRequest(ctx)

	for _, request := range allRequests {
		if request.Status == status {
			requests = append(requests, request)
		}
	}

	return requests
}

// AppendDegreeRequest stores a new degree request and returns its ID
func (k Keeper) AppendDegreeRequest(ctx sdk.Context, degreeRequest types.DegreeRequest) (uint64, error) {
	// Get current count
	store := ctx.KVStore(k.storeKey)
	countBz := store.Get(types.DegreeRequestCountPrefix)
	
	var count uint64 = 0
	if countBz != nil {
		count = binary.BigEndian.Uint64(countBz)
	}
	
	// Set ID and store
	degreeRequest.Id = fmt.Sprintf("%d", count)
	k.SetDegreeRequest(ctx, degreeRequest)
	
	// Update count
	count++
	countBz = make([]byte, 8)
	binary.BigEndian.PutUint64(countBz, count)
	store.Set(types.DegreeRequestCountPrefix, countBz)
	
	return count - 1, nil
}

// AppendDegree stores a new degree and returns its ID
func (k Keeper) AppendDegree(ctx sdk.Context, degree types.Degree) (uint64, error) {
	store := ctx.KVStore(k.storeKey)
	countBz := store.Get(types.DegreeCountPrefix)
	
	var count uint64 = 0
	if countBz != nil {
		count = binary.BigEndian.Uint64(countBz)
	}
	
	// Set index and store
	degree.Index = fmt.Sprintf("%d", count)
	k.SetDegree(ctx, degree)
	
	// Update count
	count++
	countBz = make([]byte, 8)
	binary.BigEndian.PutUint64(countBz, count)
	store.Set(types.DegreeCountPrefix, countBz)
	
	return count - 1, nil
}
