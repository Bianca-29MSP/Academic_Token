package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"academictoken/x/student/types"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Student transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdUpdateParams())
	return cmd
}

func CmdUpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [ipfs-gateway] [ipfs-enabled] [admin] [prerequisites-contract] [equivalence-contract] [academic-progress-contract] [degree-contract] [nft-minting-contract]",
		Short: "Update ALL student module parameters including ALL contract addresses",
		Args:  cobra.ExactArgs(8),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ipfsGateway := args[0]
			ipfsEnabled, err := strconv.ParseBool(args[1])
			if err != nil {
				return err
			}
			admin := args[2]
			prerequisitesContract := args[3]
			equivalenceContract := args[4]
			academicProgressContract := args[5]
			degreeContract := args[6]
			nftMintingContract := args[7]

			params := types.Params{
				IpfsGateway:                   ipfsGateway,
				IpfsEnabled:                   ipfsEnabled,
				Admin:                        admin,
				PrerequisitesContractAddr:    prerequisitesContract,
				EquivalenceContractAddr:      equivalenceContract,
				AcademicProgressContractAddr: academicProgressContract,
				DegreeContractAddr:           degreeContract,
				NftMintingContractAddr:       nftMintingContract,
			}

			msg := &types.MsgUpdateParams{
				Authority: clientCtx.GetFromAddress().String(),
				Params:    params,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
