package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	// this line is used by starport scaffolding # 1
)

// ModuleCdc references the global x/student module codec. Note, the codec
// should ONLY be used in certain instances of tests and for JSON encoding as Amino
// is still used for that purpose.
//
// The actual codec used for serialization should be provided to x/student and
// defined at the application level.
var ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRegisterStudent{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateEnrollment{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateEnrollmentStatus{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRequestSubjectEnrollment{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateAcademicTree{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCompleteSubject{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgRequestEquivalence{},
	)
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
