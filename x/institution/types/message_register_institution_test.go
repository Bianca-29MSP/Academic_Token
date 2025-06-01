package types

import (
	"testing"

	"academictoken/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgRegisterInstitution_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgRegisterInstitution
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgRegisterInstitution{
				Creator: "invalid_address",
				Name:    "Test University",
				Address: "123 Main St, City, Country",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "empty name",
			msg: MsgRegisterInstitution{
				Creator: sample.AccAddress(),
				Name:    "", // Empty name must fail
				Address: "123 Main St, City, Country",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "empty address",
			msg: MsgRegisterInstitution{
				Creator: sample.AccAddress(),
				Name:    "Test University",
				Address: "", // Empty addres must fail
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid message",
			msg: MsgRegisterInstitution{
				Creator: sample.AccAddress(),
				Name:    "Test University",
				Address: "123 Main St, City, Country",
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
