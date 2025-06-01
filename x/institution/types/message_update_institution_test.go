package types

import (
	"testing"

	"academictoken/testutil/sample"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgUpdateInstitution_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUpdateInstitution
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgUpdateInstitution{
				Creator:      "invalid_address",
				Index:        "institution-1",
				Name:         "Updated University",
				Address:      "456 New St, City, Country",
				IsAuthorized: "true",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "empty index",
			msg: MsgUpdateInstitution{
				Creator:      sample.AccAddress(),
				Index:        "", // Index vazio deve falhar
				Name:         "Updated University",
				Address:      "456 New St, City, Country",
				IsAuthorized: "true",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "empty name",
			msg: MsgUpdateInstitution{
				Creator:      sample.AccAddress(),
				Index:        "institution-1",
				Name:         "", // Nome vazio deve falhar
				Address:      "456 New St, City, Country",
				IsAuthorized: "true",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "empty address",
			msg: MsgUpdateInstitution{
				Creator:      sample.AccAddress(),
				Index:        "institution-1",
				Name:         "Updated University",
				Address:      "", // Endereço vazio deve falhar
				IsAuthorized: "true",
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid message",
			msg: MsgUpdateInstitution{
				Creator:      sample.AccAddress(),
				Index:        "institution-1",
				Name:         "Updated University",
				Address:      "456 New St, City, Country",
				IsAuthorized: "true",
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
