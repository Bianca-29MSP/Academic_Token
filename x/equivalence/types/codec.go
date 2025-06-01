package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	// this line is used by starport scaffolding # 1
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUpdateParams{}, "equivalence/MsgUpdateParams", nil)
	cdc.RegisterConcrete(&MsgRequestEquivalence{}, "equivalence/MsgRequestEquivalence", nil)
	cdc.RegisterConcrete(&MsgExecuteEquivalenceAnalysis{}, "equivalence/MsgExecuteEquivalenceAnalysis", nil)
	cdc.RegisterConcrete(&MsgBatchRequestEquivalence{}, "equivalence/MsgBatchRequestEquivalence", nil)
	cdc.RegisterConcrete(&MsgUpdateContractAddress{}, "equivalence/MsgUpdateContractAddress", nil)
	cdc.RegisterConcrete(&MsgReanalyzeEquivalence{}, "equivalence/MsgReanalyzeEquivalence", nil)
	// Removed MsgApproveEquivalence and MsgRejectEquivalence as they are not in protobuf
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
		&MsgRequestEquivalence{},
		&MsgExecuteEquivalenceAnalysis{},
		&MsgBatchRequestEquivalence{},
		&MsgUpdateContractAddress{},
		&MsgReanalyzeEquivalence{},
		// Removed MsgApproveEquivalence and MsgRejectEquivalence as they are not in protobuf
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(Amino)
	Amino.Seal()
}
