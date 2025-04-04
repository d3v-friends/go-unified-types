package utypes

import (
	"database/sql/driver"
	"github.com/99designs/gqlgen/graphql"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnPanic"
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/mgOp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"reflect"
	"strconv"
)

const ErrInvalidObjectID = "invalid_object_id"

var NilObjectID = &ObjectID{
	V: primitive.NilObjectID.Hex(),
}

func NewObjectID(vs ...primitive.ObjectID) *ObjectID {
	if len(vs) == 1 {
		return &ObjectID{
			V: vs[0].Hex(),
		}
	}
	return &ObjectID{
		V: primitive.NewObjectID().Hex(),
	}
}

func NewObjectIDByString(str string) (res *ObjectID, err error) {
	var id primitive.ObjectID
	if id, err = primitive.ObjectIDFromHex(str); err != nil {
		return
	}
	res = NewObjectID(id)
	return
}

func (x *ObjectID) ObjectID() (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(x.V)
}

func (x *ObjectID) ObjectIDP() primitive.ObjectID {
	return fnPanic.Value(x.ObjectID())
}

func (x *ObjectID) Scan(value any) (err error) {
	var body, ok = value.(string)
	if !ok {
		return fnError.NewF(ErrInvalidObjectID)
	}

	var id primitive.ObjectID
	if id, err = primitive.ObjectIDFromHex(body); err != nil {
		return
	}

	*x = ObjectID{
		V: id.Hex(),
	}

	return nil
}

func (x *ObjectID) Value() (res driver.Value, err error) {
	var id primitive.ObjectID
	if id, err = x.ObjectID(); err != nil {
		return
	}
	res = id.Hex()
	return
}

func (x *ObjectID) GormDBDataType(
	_ *gorm.DB,
	_ *schema.Field,
) string {
	return "char(24)"
}

/* ------------------------------------------------------------------------------------------------------------ */

func MarshalObjectID(v *ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Quote(v.V)))
	})
}

func UnmarshalObjectID(v any) (res *ObjectID, err error) {
	switch t := v.(type) {
	case string:
		return NewObjectIDByString(t)
	case []byte:
		return NewObjectIDByString(string(t))
	default:
		err = fnError.NewF(ErrInvalidObjectID)
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

type objectIdValueCodec struct {
}

func (x *objectIdValueCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(ObjectID)
	if !isOk {
		err = fnError.NewF(ErrInvalidObjectID)
		return
	}

	var dec primitive.ObjectID
	if dec, err = i.ObjectID(); err != nil {
		return
	}

	return writer.WriteObjectID(dec)
}

func (x *objectIdValueCodec) DecodeValue(
	_ bsoncodec.DecodeContext,
	reader bsonrw.ValueReader,
	value reflect.Value,
) (err error) {
	if err = reader.ReadNull(); err == nil {
		value.Set(reflect.ValueOf(ObjectID{V: primitive.NilObjectID.Hex()}))
		return
	}

	var dec primitive.ObjectID
	if dec, err = reader.ReadObjectID(); err != nil {
		return
	}

	value.Set(reflect.ValueOf(ObjectID{V: dec.Hex()}))

	return
}

type objectIdPointerCodec struct {
}

func (x *objectIdPointerCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(*ObjectID)
	if !isOk {
		err = fnError.NewF(ErrInvalidObjectID)
		return
	}

	if fnPointer.IsNil(i) {
		return writer.WriteNull()
	}

	var dec primitive.ObjectID
	if dec, err = i.ObjectID(); err != nil {
		return
	}

	return writer.WriteObjectID(dec)
}

func (x *objectIdPointerCodec) DecodeValue(
	_ bsoncodec.DecodeContext,
	reader bsonrw.ValueReader,
	value reflect.Value,
) (err error) {
	if err = reader.ReadNull(); err == nil {
		value.Set(reflect.Zero(value.Type()))
		return
	}

	var dec primitive.ObjectID
	if dec, err = reader.ReadObjectID(); err != nil {
		return
	}

	value.Set(reflect.ValueOf(&ObjectID{V: dec.Hex()}))

	return
}

func ObjectIDRegistry(registry *bsoncodec.Registry) *bsoncodec.Registry {
	var value = &objectIdValueCodec{}
	registry.RegisterTypeEncoder(reflect.TypeOf(ObjectID{}), value)
	registry.RegisterTypeDecoder(reflect.TypeOf(ObjectID{}), value)

	var ptr = &objectIdPointerCodec{}
	registry.RegisterTypeEncoder(reflect.TypeOf(&ObjectID{}), ptr)
	registry.RegisterTypeDecoder(reflect.TypeOf(&ObjectID{}), ptr)
	return registry
}

/* ------------------------------------------------------------------------------------------------------------ */

func (x *ObjectIDArgs) AppendFilter(filter bson.M, key string) bson.M {
	if fnPointer.IsNil(x.Equal) {
		filter[key] = x.Equal
		return filter
	}

	if len(x.In) != 0 {
		filter[key] = bson.M{
			mgOp.In: x.In,
		}
		return filter
	}

	if len(x.HasAll) != 0 {
		filter[key] = x.HasAll
		return filter
	}

	return filter
}
