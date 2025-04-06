package utypes_test

import (
	"github.com/d3v-friends/go-unified-types/utypes"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestObjectID(test *testing.T) {
	test.Run("append_filter_equal", func(t *testing.T) {
		var id1 = utypes.NewObjectID()

		var args = &utypes.ObjectIDArgs{
			Equal: id1,
		}

		var filter = args.AppendFilter(bson.M{}, "test")
		var filterEncoder, filterBuffer = NewEncoder(t)
		var err = filterEncoder.Encode(filter)
		assert.NoError(t, err)

		var filterBson = make(bson.M)
		err = bson.Unmarshal(filterBuffer.Bytes(), &filterBson)
		assert.NoError(t, err)

		var answer = bson.M{
			"test": id1,
		}

		var answerEncoder, answerBuffer = NewEncoder(t)
		err = answerEncoder.Encode(answer)
		assert.NoError(t, err)

		var answerBson = make(bson.M)
		err = bson.Unmarshal(answerBuffer.Bytes(), &answerBson)
		assert.NoError(t, err)

		assert.Equal(t, filterBson, answerBson)
	})

	test.Run("append_filter_in", func(t *testing.T) {
		var id1 = utypes.NewObjectID()
		var id2 = utypes.NewObjectID()
		var id3 = utypes.NewObjectID()

		var args = &utypes.ObjectIDArgs{
			In: []*utypes.ObjectID{
				id1, id2, id3,
			},
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
				"$in": []primitive.ObjectID{
					id1.ObjectIDP(),
					id2.ObjectIDP(),
					id3.ObjectIDP(),
				},
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

	test.Run("append_filter_has_all", func(t *testing.T) {
		var id1 = utypes.NewObjectID()
		var id2 = utypes.NewObjectID()
		var id3 = utypes.NewObjectID()

		var args = &utypes.ObjectIDArgs{
			HasAll: []*utypes.ObjectID{
				id1, id2, id3,
			},
		}

		var filter = args.AppendFilter(bson.M{}, "test")
		var filterEncoder, filterBuffer = NewEncoder(t)
		var err = filterEncoder.Encode(filter)
		assert.NoError(t, err)

		var filterBson = make(bson.M)
		err = bson.Unmarshal(filterBuffer.Bytes(), &filterBson)
		assert.NoError(t, err)

		var answer = bson.M{
			"test": []primitive.ObjectID{
				id1.ObjectIDP(),
				id2.ObjectIDP(),
				id3.ObjectIDP(),
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

	// todo 우선순위 테스트 하기
}
