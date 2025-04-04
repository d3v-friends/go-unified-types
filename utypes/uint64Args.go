package utypes

import (
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/mgOp"
	"go.mongodb.org/mongo-driver/bson"
)

func (x *Uint64Args) AppendFilter(filter bson.M, key string) bson.M {
	if fnPointer.IsNil(x) {
		return filter
	}

	var compare = bson.M{}
	if gt := x.Gt; !fnPointer.IsNil(gt) {
		compare[mgOp.Gt] = *gt
	}

	if gte := x.Gte; !fnPointer.IsNil(gte) {
		compare[mgOp.Gte] = *gte
	}

	if lt := x.Lt; !fnPointer.IsNil(lt) {
		compare[mgOp.Lt] = *lt
	}

	if lte := x.Lte; !fnPointer.IsNil(lte) {
		compare[mgOp.Lte] = *lte
	}

	if equal := x.Equal; !fnPointer.IsNil(equal) {
		compare[mgOp.Eq] = *equal
	}

	if notEqual := x.NotEqual; !fnPointer.IsNil(notEqual) {
		compare[mgOp.Ne] = *notEqual
	}

	if len(compare) == 0 {
		return filter
	}

	filter[key] = compare

	return filter
}
