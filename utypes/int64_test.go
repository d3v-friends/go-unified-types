package utypes_test

import (
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/go-unified-types/utypes"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestInt64(test *testing.T) {
	test.Run("append_filter", func(t *testing.T) {
		var args = &utypes.Int64Args{
			Gt:       fnPointer.Make(int64(1)),
			Gte:      fnPointer.Make(int64(1)),
			Lt:       fnPointer.Make(int64(1)),
			Lte:      fnPointer.Make(int64(1)),
			Equal:    fnPointer.Make(int64(1)),
			NotEqual: fnPointer.Make(int64(1)),
		}

		var registry = bson.NewRegistry()
		for _, fn := range utypes.Registries {
			registry = fn(registry)
		}

		var filter = args.AppendFilter(bson.M{}, "test")

		var filterEncoder, filterBuffer = NewEncoder(t)
		var err = filterEncoder.Encode(filter)
		assert.NoError(t, err)

		var filterBson = make(bson.M)
		err = bson.Unmarshal(filterBuffer.Bytes(), &filterBson)
		assert.NoError(t, err)

		var answer = bson.M{
			"test": bson.M{
				"$gt":  int64(1),
				"$gte": int64(1),
				"$lt":  int64(1),
				"$lte": int64(1),
				"$eq":  int64(1),
				"$ne":  int64(1),
			},
		}

		var answerEncoder, answerBuffer = NewEncoder(t)
		err = answerEncoder.Encode(answer)
		assert.NoError(t, err)

		var answerBson = make(bson.M)
		err = bson.Unmarshal(answerBuffer.Bytes(), &answerBson)
		assert.NoError(t, err)

		assert.Equal(t, filterBson, answerBson)

	})
}
