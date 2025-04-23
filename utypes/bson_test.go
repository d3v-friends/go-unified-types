package utypes_test

import (
	"github.com/d3v-friends/go-unified-types/utypes"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestBson(test *testing.T) {
	test.Run("encode", func(t *testing.T) {
		var id = primitive.NewObjectID()
		var key = "id"
		var b = bson.M{
			key: id,
		}

		var value, err = utypes.NewBson(b)
		assert.NoError(t, err)

		var raw bson.Raw
		raw, err = value.Raw()
		assert.NoError(t, err)

		var rawValue = raw.Lookup(key)
		assert.Equal(t, id.Hex(), rawValue.ObjectID().Hex())
	})

	test.Run("encode_registry", func(t *testing.T) {
		var id = utypes.NewObjectID()
		var key = "id"
		var b = bson.M{
			key: id,
		}

		var registry = utypes.NewRegistry()
		var value, err = utypes.NewBson(b, registry)
		assert.NoError(t, err)

		var raw bson.Raw
		raw, err = value.Raw(registry)
		assert.NoError(t, err)

		var parsedId = &utypes.ObjectID{}
		err = raw.Lookup(key).Unmarshal(parsedId)
		assert.NoError(t, err)

		assert.Equal(t, id.V, parsedId.V)
	})
}
