package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"academictoken/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestMsgUpdateDegreeContract_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateDegreeContract
		err  error
	}{
		{
			name: "invalid authority address",
			msg: MsgUpdateDegreeContract{
				Authority:          "invalid_address",
				NewContractAddress: "wasm1contract",
				ContractVersion:    "1.0.0",
				MigrationReason:    "upgrade",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "empty new contract address",
			msg: MsgUpdateDegreeContract{
				Authority:          sample.AccAddress(),
				NewContractAddress: "",
				ContractVersion:    "1.0.0",
				MigrationReason:    "upgrade",
			},
			err: ErrInvalidContractAddress,
		},
		{
			name: "empty contract version",
			msg: MsgUpdateDegreeContract{
				Authority:          sample.AccAddress(),
				NewContractAddress: "wasm1contract",
				ContractVersion:    "",
				MigrationReason:    "upgrade",
			},
			err: ErrInvalidDegreeRequest,
		},
		{
			name: "valid message",
			msg: MsgUpdateDegreeContract{
				Authority:          sample.AccAddress(),
				NewContractAddress: "wasm1contract",
				ContractVersion:    "1.0.0",
				MigrationReason:    "upgrade",
			},
		},
		{
			name: "valid message without migration reason",
			msg: MsgUpdateDegreeContract{
				Authority:          sample.AccAddress(),
				NewContractAddress: "wasm1contract",
				ContractVersion:    "1.0.0",
				MigrationReason:    "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
