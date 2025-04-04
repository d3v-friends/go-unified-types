package utypes

import "github.com/d3v-friends/mango/mgCodec"

func (x SortDirection) GetDirection() int64 {
	switch x {
	case SortDirection_SD_ASC:
		return 1
	case SortDirection_SD_DESC:
		return -1
	default:
		return 1
	}
}

func (x SortDirection) New(i int32) SortDirection {
	return SortDirection(i)
}

var SortDirectionRegistry = mgCodec.NewGrpcEnum[SortDirection](SortDirection_value)
