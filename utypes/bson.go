package utypes

import (
	"bytes"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnPointer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
)

type BsonType interface {
	bson.M | bson.E | bson.D | bson.A
}

func NewBson[T BsonType](
	args T,
	registries ...*bsoncodec.Registry,
) (res *BSON, err error) {
	var buf = new(bytes.Buffer)
	var vw bsonrw.ValueWriter
	if vw, err = bsonrw.NewBSONValueWriter(buf); err != nil {
		return
	}

	var enc *bson.Encoder
	if enc, err = bson.NewEncoder(vw); err != nil {
		return
	}

	if err = enc.Encode(args); err != nil {
		return
	}

	res = &BSON{}
	if len(registries) == 1 {
		if err = enc.SetRegistry(registries[0]); err != nil {
			return
		}
	}

	if err = enc.Encode(args); err != nil {
		return
	}

	res.V = buf.Bytes()

	return
}

func (x *BSON) Decode(
	value any,
	registries ...*bsoncodec.Registry,
) (err error) {
	if fnPointer.IsNil(x) {
		err = fnError.New("bson_is_nil")
		return
	}

	var dec *bson.Decoder
	if dec, err = bson.NewDecoder(bsonrw.NewBSONDocumentReader(x.V)); err != nil {
		return
	}

	if len(registries) == 1 {
		if err = dec.SetRegistry(registries[0]); err != nil {
			return
		}
	}

	if err = dec.Decode(value); err != nil {
		return
	}

	return
}

func (x *BSON) Raw(registries ...*bsoncodec.Registry) (raw bson.Raw, err error) {
	raw = bson.Raw{}
	if err = x.Decode(&raw, registries...); err != nil {
		return
	}
	return
}
