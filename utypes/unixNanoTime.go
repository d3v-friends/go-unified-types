package utypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/mgQuery"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"reflect"
	"time"
)

const ErrInvalidUnixNanoTime = "invalid_unix_nano_time"

var ZeroUnixNanoTime = &UnixNanoTime{V: 0}

func NewUnixNanoTime(vs ...time.Time) *UnixNanoTime {
	if len(vs) == 1 {
		return &UnixNanoTime{V: vs[0].UnixNano()}
	}
	return &UnixNanoTime{V: time.Now().UnixNano()}
}

func NewUnixNanoTimeByString(v string) (res *UnixNanoTime, err error) {
	var i decimal.Decimal
	if i, err = decimal.NewFromString(v); err != nil {
		return
	}
	res = &UnixNanoTime{V: i.IntPart()}
	return
}

func (x *UnixNanoTime) Scan(value any) (err error) {
	var body, ok = value.(int64)
	if !ok {
		err = fnError.NewF(ErrInvalidUnixNanoTime)
		return
	}

	*x = UnixNanoTime{V: body}
	return
}

func (x *UnixNanoTime) Value() (res driver.Value, err error) {
	res = x.V
	return
}

func (x *UnixNanoTime) YMDHTime() *YMDHTime {
	return NewYMDHTime(time.UnixMilli(x.V))
}

func (x *UnixNanoTime) YMDTime(stdHour ...int64) *YMDTime {
	return NewYMDTimeByTime(time.UnixMilli(x.V), stdHour...)
}

func (x *UnixNanoTime) RFC3339Time() *RFC3339Time {
	return NewRFC3339Time(time.UnixMilli(x.V))
}

/* ------------------------------------------------------------------------------------------------------------ */

func MarshalUnixNanoTime(v *UnixNanoTime) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(fmt.Sprintf("%d", v.V)))
	})
}

func UnmarshalUnixNanoTime(v any) (res *UnixNanoTime, err error) {
	switch t := v.(type) {
	case string:
		return NewUnixNanoTimeByString(t)
	case []byte:
		return NewUnixNanoTimeByString(string(t))
	case json.Number:
		return NewUnixNanoTimeByString(t.String())
	default:
		err = fnError.NewF(ErrInvalidUnixNanoTime)
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

type unixNanoTimeValueCodec struct {
}

func (x *unixNanoTimeValueCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(UnixNanoTime)
	if !isOk {
		err = fnError.NewF(ErrInvalidUnixNanoTime)
		return
	}

	var dec primitive.Decimal128
	if dec, err = primitive.ParseDecimal128(fmt.Sprintf("%d", i.V)); err != nil {
		return
	}

	return writer.WriteDecimal128(dec)
}

func (x *unixNanoTimeValueCodec) DecodeValue(
	_ bsoncodec.DecodeContext,
	reader bsonrw.ValueReader,
	value reflect.Value,
) (err error) {
	if err = reader.ReadNull(); err == nil {
		value.Set(reflect.ValueOf(UnixNanoTime{V: 0}))
		return
	}

	var dec primitive.Decimal128
	if dec, err = reader.ReadDecimal128(); err != nil {
		return
	}

	var d decimal.Decimal
	if d, err = decimal.NewFromString(dec.String()); err != nil {
		return
	}

	value.Set(reflect.ValueOf(UnixNanoTime{V: d.IntPart()}))

	return
}

type unixNanoTimePtrCodec struct {
}

func (x *unixNanoTimePtrCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(*UnixNanoTime)
	if !isOk {
		err = fnError.NewF(ErrInvalidUnixNanoTime)
		return
	}

	if fnPointer.IsNil(i) {
		return writer.WriteNull()
	}

	var dec primitive.Decimal128
	if dec, err = primitive.ParseDecimal128(fmt.Sprintf("%d", i.V)); err != nil {
		return
	}

	return writer.WriteDecimal128(dec)
}

func (x *unixNanoTimePtrCodec) DecodeValue(
	_ bsoncodec.DecodeContext,
	reader bsonrw.ValueReader,
	value reflect.Value,
) (err error) {
	if err = reader.ReadNull(); err == nil {
		value.Set(reflect.Zero(value.Type()))
		return
	}

	var dec primitive.Decimal128
	if dec, err = reader.ReadDecimal128(); err != nil {
		return
	}

	var d decimal.Decimal
	if d, err = decimal.NewFromString(dec.String()); err != nil {
		return
	}

	value.Set(reflect.ValueOf(&UnixNanoTime{V: d.IntPart()}))

	return
}

func UnixNanoTimeRegistry(registry *bsoncodec.Registry) *bsoncodec.Registry {
	var value, valueTypeof = &unixNanoTimeValueCodec{}, reflect.TypeOf(UnixNanoTime{})
	registry.RegisterTypeEncoder(valueTypeof, value)
	registry.RegisterTypeDecoder(valueTypeof, value)

	var ptr, ptrTypeof = &unixNanoTimePtrCodec{}, reflect.TypeOf(&UnixNanoTime{})
	registry.RegisterTypeEncoder(ptrTypeof, ptr)
	registry.RegisterTypeDecoder(ptrTypeof, ptr)
	return registry
}

func (x *UnixNanoTimeArgs) AppendFilter(filter bson.M, key string) bson.M {
	return mgQuery.AppendFilterCompareArgs[UnixNanoTime](filter, key, x)
}
