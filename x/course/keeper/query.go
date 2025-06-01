package keeper

import (
	"academictoken/x/course/types"
)

var _ types.QueryServer = Keeper{}
