package utypes

import (
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/mgOp"
	"go.mongodb.org/mongo-driver/bson"
)

func (x *StringArgs) AppendFilter(filter bson.M, key string) bson.M {
	if !fnPointer.IsNil(x.Exact) {
		filter[key] = *x.Exact
		return filter
	}

	if !fnPointer.IsNil(x.Like) {
		filter[key] = bson.M{
			mgOp.Regex: *x.Like,
		}
		return filter
	}

	if 0 < len(x.In) {
		filter[key] = bson.M{
			mgOp.In: x.In,
		}
	}

	return filter
}
