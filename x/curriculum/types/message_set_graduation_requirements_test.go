package types

import (
	"testing"

	"academictoken/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgSetGraduationRequirements_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSetGraduationRequirements
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSetGraduationRequirements{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSetGraduationRequirements{
				Creator: sample.AccAddress(),
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
