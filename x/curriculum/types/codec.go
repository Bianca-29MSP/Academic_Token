package types

import (
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	// this line is used by starport scaffolding # 1
)

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateCurriculumTree{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddSemesterToCurriculum{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgAddElectiveGroup{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSetGraduationRequirements{},
	)
	// this line is used by starport scaffolding # 3

	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
