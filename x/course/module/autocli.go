package course

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "academictoken/api/academictoken/course"
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
					RpcMethod:      "CreateCourse",
					Use:            "create-course [institution] [name] [code] [description] [total-credits] [degree-level]",
					Short:          "Send a create-course tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "institution"}, {ProtoField: "name"}, {ProtoField: "code"}, {ProtoField: "description"}, {ProtoField: "totalCredits"}, {ProtoField: "degreeLevel"}},
				},
				{
					RpcMethod:      "UpdateCourse",
					Use:            "update-course [index] [name] [description] [total-credits]",
					Short:          "Send a update-course tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{{ProtoField: "index"}, {ProtoField: "name"}, {ProtoField: "description"}, {ProtoField: "totalCredits"}},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
