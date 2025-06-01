package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRequestDegree{}, "degree/RequestDegree", nil)
	cdc.RegisterConcrete(&MsgValidateDegreeRequirements{}, "degree/ValidateDegreeRequirements", nil)
	cdc.RegisterConcrete(&MsgIssueDegree{}, "degree/IssueDegree", nil)
	cdc.RegisterConcrete(&MsgUpdateDegreeContract{}, "degree/UpdateDegreeContract", nil)
	cdc.RegisterConcrete(&MsgCancelDegreeRequest{}, "degree/CancelDegreeRequest", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRequestDegree{},
		&MsgValidateDegreeRequirements{},
		&MsgIssueDegree{},
		&MsgUpdateDegreeContract{},
		&MsgCancelDegreeRequest{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(types.NewInterfaceRegistry())
)

func init() {
	RegisterCodec(Amino)
	Amino.Seal()
}
