package academicnft

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "academictoken/api/academictoken/academicnft"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "MintSubjectToken",
					Use:            "mint-subject-token [token-id] [student] [grade] [semester] [professor-signature]",
					Short:          "Send a mint-subject-token tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "tokenId"}, {ProtoField: "student"}, {ProtoField: "grade"}, {ProtoField: "semester"}, {ProtoField: "professorSignature"}},
				},
				{
					RpcMethod:      "VerifyTokenInstance",
					Use:            "verify-token-instance [token-instance-id]",
					Short:          "Send a verify-token-instance tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "tokenInstanceId"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
