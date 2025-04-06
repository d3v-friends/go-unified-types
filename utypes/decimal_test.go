package utypes_test

import (
	"github.com/d3v-friends/go-unified-types/utypes"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

func TestDecimalType(test *testing.T) {
	test.Run("append_filter", func(t *testing.T) {
		var args = &utypes.DecimalArgs{
			Gt:       utypes.NewDecimal(decimal.NewFromInt(1)),
			Gte:      utypes.NewDecimal(decimal.NewFromInt(1)),
			Lt:       utypes.NewDecimal(decimal.NewFromInt(1)),
			Lte:      utypes.NewDecimal(decimal.NewFromInt(1)),
			Equal:    utypes.NewDecimal(decimal.NewFromInt(1)),
			NotEqual: utypes.NewDecimal(decimal.NewFromInt(1)),
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
				"$gt":  utypes.NewDecimal(decimal.NewFromInt(1)),
				"$gte": utypes.NewDecimal(decimal.NewFromInt(1)),
				"$lt":  utypes.NewDecimal(decimal.NewFromInt(1)),
				"$lte": utypes.NewDecimal(decimal.NewFromInt(1)),
				"$eq":  utypes.NewDecimal(decimal.NewFromInt(1)),
				"$ne":  utypes.NewDecimal(decimal.NewFromInt(1)),
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
