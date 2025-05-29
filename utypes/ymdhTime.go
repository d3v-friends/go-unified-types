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

const ErrInvalidYMDHTime = "invalid_ymdh_time"

var NilYMDHTime = &YMDHTime{
	V: NilRFC3339TimeString,
}

func NewYMDHTime(vs ...time.Time) (res *YMDHTime) {
	var t time.Time
	if len(vs) == 1 {
		t = vs[0]
	} else {
		t = time.Now()
	}

	t = time.UnixMilli(t.UnixMilli() - t.UnixMilli()%time.Hour.Milliseconds())

	res = &YMDHTime{
		V: t.Format(time.RFC3339),
	}

	return
}

func NewYMDHTimeByString(str string) (res *YMDHTime, err error) {
	var t time.Time
	if t, err = time.Parse(time.RFC3339, str); err != nil {
		return
	}

	res = NewYMDHTime(t)
	return
}

func (x *YMDHTime) Time() (time.Time, error) {
	return time.Parse(time.RFC3339, x.V)
}

func (x *YMDHTime) TimeP() time.Time {
	return fnPanic.Value(x.Time())
}

func (x *YMDHTime) Scan(value any) (err error) {
	var body, ok = value.(time.Time)
	if !ok {
		return fnError.NewF(ErrInvalidYMDHTime)
	}

	*x = YMDHTime{
		V: body.Format(time.RFC3339),
	}

	return
}

func (x *YMDHTime) Value() (driver.Value, error) {
	return x.Time()
}

func (x *YMDHTime) YMDTime(stdHour ...int64) *YMDTime {
	return NewYMDTimeByTime(x.TimeP(), stdHour...)
}

func (x *YMDHTime) RFC3339Time() *RFC3339Time {
	return NewRFC3339Time(x.TimeP())
}

func (x *YMDHTime) UnixNanoTime() *UnixNanoTime {
	return NewUnixNanoTime(x.TimeP())
}

/* ------------------------------------------------------------------------------------------------------------ */

func MarshalYMDHTime(v *YMDHTime) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Quote(v.V)))
	})
}

func UnmarshalYMDHTime(v any) (res *YMDHTime, err error) {
	switch t := v.(type) {
	case string:
		return NewYMDHTimeByString(t)
	case []byte:
		return NewYMDHTimeByString(string(t))
	default:
		err = fnError.NewF(ErrInvalidYMDHTime)
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

type ymdhValueCodec struct {
}

func (x *ymdhValueCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(YMDHTime)
	if !isOk {
		err = fnError.NewF(ErrInvalidYMDHTime)
		return
	}

	var t time.Time
	if t, err = i.Time(); err != nil {
		return
	}

	return writer.WriteDateTime(t.UnixMilli())
}

func (x *ymdhValueCodec) DecodeValue(
	_ bsoncodec.DecodeContext,
	reader bsonrw.ValueReader,
	value reflect.Value,
) (err error) {
	if err = reader.ReadNull(); err == nil {
		value.Set(reflect.ValueOf(YMDHTime{V: ""}))
		return
	}

	var dateTime int64
	if dateTime, err = reader.ReadDateTime(); err != nil {
		return
	}

	value.Set(reflect.ValueOf(YMDHTime{
		V: time.UnixMilli(dateTime).Format(time.RFC3339),
	}))

	return
}

type ymdhPtrCodec struct {
}

func (x *ymdhPtrCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(*YMDHTime)
	if !isOk {
		err = fnError.NewF(ErrInvalidYMDHTime)
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

func (x *ymdhPtrCodec) DecodeValue(
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

	value.Set(reflect.ValueOf(&YMDHTime{
		V: time.UnixMilli(dateTime).Format(time.RFC3339),
	}))

	return
}

func YMDHTimeRegistry(register *bsoncodec.Registry) *bsoncodec.Registry {
	var value, valueTypeof = &ymdhValueCodec{}, reflect.TypeOf(YMDHTime{})
	register.RegisterTypeEncoder(valueTypeof, value)
	register.RegisterTypeDecoder(valueTypeof, value)

	var ptr, ptrTypeof = &ymdhPtrCodec{}, reflect.TypeOf(&YMDHTime{})
	register.RegisterTypeEncoder(ptrTypeof, ptr)
	register.RegisterTypeDecoder(ptrTypeof, ptr)
	return register
}

func (x *YMDHArgs) AppendFilter(filter bson.M, key string) bson.M {
	return mgQuery.AppendFilterCompareArgs[YMDHTime](filter, key, x)
}
