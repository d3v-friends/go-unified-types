package utypes_test

import (
	"bytes"
	"github.com/d3v-friends/go-unified-types/utypes"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"testing"
)

func NewEncoder(t *testing.T) (encoder *bson.Encoder, buffer *bytes.Buffer) {
	buffer = new(bytes.Buffer)
	var writer, err = bsonrw.NewBSONValueWriter(buffer)
	assert.NoError(t, err)

	encoder, err = bson.NewEncoder(writer)
	assert.NoError(t, err)

	var registry = bson.NewRegistry()
	for _, fn := range utypes.Registries {
		registry = fn(registry)
	}

	err = encoder.SetRegistry(registry)
	assert.NoError(t, err)

	return

}
