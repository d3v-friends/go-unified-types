package utypes

import (
	"database/sql/driver"
	"github.com/99designs/gqlgen/graphql"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnPanic"
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/mgQuery"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"io"
	"reflect"
	"strconv"
	"time"
)

const ErrInvalidRFC3339Time = "invalid_rfc3339_time"

func NewRFC3339Time(vs ...time.Time) *RFC3339Time {
	if len(vs) == 1 {
		return &RFC3339Time{
			V: vs[0].Format(time.RFC3339),
		}
	}
	return &RFC3339Time{
		V: time.Now().Format(time.RFC3339),
	}
}

func NewRFC3339TimeByString(str string) (res *RFC3339Time, err error) {
	var t time.Time
	if t, err = time.Parse(time.RFC3339, str); err != nil {
		return
	}

	res = NewRFC3339Time(t)
	return
}

func (x *RFC3339Time) Time() (time.Time, error) {
	return time.Parse(time.RFC3339, x.V)
}

func (x *RFC3339Time) TimeP() time.Time {
	return fnPanic.Value(x.Time())
}

func (x *RFC3339Time) Scan(value any) (err error) {
	var body, ok = value.(time.Time)
	if !ok {
		return fnError.NewF(ErrInvalidRFC3339Time)
	}

	*x = RFC3339Time{
		V: body.Format(time.RFC3339),
	}

	return nil
}

func (x *RFC3339Time) Value() (driver.Value, error) {
	return x.Time()
}

func (x *RFC3339Time) YMDHTime() *YMDHTime {
	return NewYMDHTime(x.TimeP())
}

func (x *RFC3339Time) YMDTime(stdHour ...int64) *YMDTime {
	return NewYMDTimeByTime(x.TimeP(), stdHour...)
}

func (x *RFC3339Time) RFC3339Time() *RFC3339Time {
	return NewRFC3339Time(x.TimeP())
}

func (x *RFC3339Time) UnixNanoTime() *UnixNanoTime {
	return NewUnixNanoTime(x.TimeP())
}

/* ------------------------------------------------------------------------------------------------------------ */

func MarshalRFC3339Time(v *RFC3339Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Quote(v.V)))
	})
}

func UnmarshalRFC3339Time(v any) (res *RFC3339Time, err error) {
	switch t := v.(type) {
	case string:
		return NewRFC3339TimeByString(t)
	case []byte:
		return NewRFC3339TimeByString(string(t))
	default:
		err = fnError.NewF(ErrInvalidRFC3339Time)
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

type rfc3339TimeValueCodec struct {
}

func (x *rfc3339TimeValueCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(RFC3339Time)
	if !isOk {
		err = fnError.NewF(ErrInvalidRFC3339Time)
		return
	}

	var t time.Time
	if t, err = i.Time(); err != nil {
		return
	}

	return writer.WriteDateTime(t.UnixMilli())
}

func (x *rfc3339TimeValueCodec) DecodeValue(
	_ bsoncodec.DecodeContext,
	reader bsonrw.ValueReader,
	value reflect.Value,
) (err error) {
	if err = reader.ReadNull(); err == nil {
		value.Set(reflect.ValueOf(RFC3339Time{
			V: time.RFC3339,
		}))
		return
	}

	var dateTime int64
	if dateTime, err = reader.ReadDateTime(); err != nil {
		return
	}

	value.Set(reflect.ValueOf(RFC3339Time{
		V: time.UnixMilli(dateTime).Format(time.RFC3339),
	}))

	return
}

type rfc3339TimePtrCodec struct {
}

func (x *rfc3339TimePtrCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(*RFC3339Time)
	if !isOk {
		err = fnError.NewF(ErrInvalidRFC3339Time)
		return
	}

	if fnPointer.IsNil(i) {
		return writer.WriteNull()
	}

	var t time.Time
	if t, err = i.Time(); err != nil {
		return
	}

	return writer.WriteDateTime(t.UnixMilli())
}

func (x *rfc3339TimePtrCodec) DecodeValue(
	_ bsoncodec.DecodeContext,
	reader bsonrw.ValueReader,
	value reflect.Value,
) (err error) {
	if err = reader.ReadNull(); err == nil {
		value.Set(reflect.Zero(value.Type()))
		return
	}

	var dateTime int64
	if dateTime, err = reader.ReadDateTime(); err != nil {
		return
	}

	value.Set(reflect.ValueOf(&RFC3339Time{
		V: time.UnixMilli(dateTime).Format(time.RFC3339),
	}))

	return
}

func RFC3339TimeRegistry(registry *bsoncodec.Registry) *bsoncodec.Registry {
	var value, valueTypeof = &rfc3339TimeValueCodec{}, reflect.TypeOf(RFC3339Time{})
	registry.RegisterTypeEncoder(valueTypeof, value)
	registry.RegisterTypeDecoder(valueTypeof, value)

	var ptr, ptrTypeof = &rfc3339TimePtrCodec{}, reflect.TypeOf(&RFC3339Time{})
	registry.RegisterTypeEncoder(ptrTypeof, ptr)
	registry.RegisterTypeDecoder(ptrTypeof, ptr)
	return registry
}

/* ------------------------------------------------------------------------------------------------------------ */

func (x *RFC3339TimeArgs) AppendFilter(filter bson.M, key string) bson.M {
	return mgQuery.AppendFilterCompareArgs[RFC3339Time](filter, key, x)
}
