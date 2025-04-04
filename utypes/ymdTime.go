package utypes

import (
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

const ErrInvalidYMDTime = "invalid_ymd_time"

func NewYMDTime(stdHour ...int64) (res *YMDTime) {
	return NewYMDTimeByTime(time.Now(), stdHour...)
}

func NewYMDTimeByTime(now time.Time, stdHour ...int64) (res *YMDTime) {
	now = time.UnixMilli(now.UnixMilli() - now.UnixMilli()%(24*time.Hour.Milliseconds()))

	if len(stdHour) == 1 && stdHour[0] < 24 {
		now = now.Add(time.Duration(stdHour[0]) * time.Hour)
	}

	res = &YMDTime{
		V: now.Format(time.RFC3339),
	}

	return
}

func (x *YMDTime) Time() (time.Time, error) {
	return time.Parse(time.RFC3339, x.V)
}

func (x *YMDTime) TimeP() time.Time {
	return fnPanic.Value(x.Time())
}

func (x *YMDTime) YMDHTime() *YMDHTime {
	return NewYMDHTime(x.TimeP())
}

func (x *YMDTime) RFC3339Time() *RFC3339Time {
	return NewRFC3339Time(x.TimeP())
}

func (x *YMDTime) UnixNanoTime() *UnixNanoTime {
	return NewUnixNanoTime(x.TimeP())
}

func MarshalYMDTime(v *YMDTime) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Quote(v.V)))
	})
}

func UnmarshalYMDTime(v any) (res *YMDTime, err error) {
	switch t := v.(type) {
	case string:
		var to time.Time
		if to, err = time.Parse(time.RFC3339, t); err == nil {
			return
		}

		res = &YMDTime{
			V: to.Format(time.RFC3339),
		}

		return
	case []byte:
		return UnmarshalYMDTime(string(t))
	default:
		err = fnError.NewF(ErrInvalidYMDTime)
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

type ymdValueCodec struct {
}

func (x *ymdValueCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(YMDTime)
	if !isOk {
		err = fnError.NewF(ErrInvalidYMDTime)
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

func (x *ymdValueCodec) DecodeValue(
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

	value.Set(reflect.ValueOf(YMDTime{
		V: time.UnixMilli(dateTime).Format(time.RFC3339),
	}))

	return
}

type ymdPtrCodec struct {
}

func (x *ymdPtrCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(*YMDTime)
	if !isOk {
		err = fnError.NewF(ErrInvalidYMDTime)
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

func (x *ymdPtrCodec) DecodeValue(
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

	value.Set(reflect.ValueOf(&YMDTime{
		V: time.UnixMilli(dateTime).Format(time.RFC3339),
	}))

	return
}

func YMDTimeRegistry(register *bsoncodec.Registry) *bsoncodec.Registry {
	var value, valueTypeof = &ymdValueCodec{}, reflect.TypeOf(YMDTime{})
	register.RegisterTypeEncoder(valueTypeof, value)
	register.RegisterTypeDecoder(valueTypeof, value)

	var ptr, ptrTypeof = &ymdPtrCodec{}, reflect.TypeOf(&YMDTime{})
	register.RegisterTypeEncoder(ptrTypeof, ptr)
	register.RegisterTypeDecoder(ptrTypeof, ptr)

	return register
}

func (x *YMDArgs) AppendFilter(filter bson.M, key string) bson.M {
	return mgQuery.AppendFilterCompareArgs[YMDTime](filter, key, x)
}
