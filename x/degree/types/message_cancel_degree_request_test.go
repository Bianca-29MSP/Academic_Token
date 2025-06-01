package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"academictoken/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestMsgCancelDegreeRequest_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgCancelDegreeRequest
		err  error
	}{
		{
			name: "invalid creator address",
			msg: MsgCancelDegreeRequest{
				Creator:            "invalid_address",
				DegreeRequestId:    "degree_req_123",
				CancellationReason: "Student changed mind",
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "empty degree request ID",
			msg: MsgCancelDegreeRequest{
				Creator:            sample.AccAddress(),
				DegreeRequestId:    "",
				CancellationReason: "Student changed mind",
			},
			err: ErrDegreeRequestNotFound,
		},
		{
			name: "empty cancellation reason",
			msg: MsgCancelDegreeRequest{
				Creator:            sample.AccAddress(),
				DegreeRequestId:    "degree_req_123",
				CancellationReason: "",
			},
			err: ErrInvalidDegreeRequest,
		},
		{
			name: "valid message",
			msg: MsgCancelDegreeRequest{
				Creator:            sample.AccAddress(),
				DegreeRequestId:    "degree_req_123",
				CancellationReason: "Student changed mind",
			},
		},
		{
			name: "valid message with long reason",
			msg: MsgCancelDegreeRequest{
				Creator:            sample.AccAddress(),
				DegreeRequestId:    "degree_req_456",
				CancellationReason: "Student decided to transfer to another institution and will complete degree there",
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
