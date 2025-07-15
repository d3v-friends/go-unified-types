package utypes

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/99designs/gqlgen/graphql"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnPanic"
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/mgQuery"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"reflect"
)

const ErrInvalidDecimal = "invalid_decimal"

var ZeroDecimal = &Decimal{
	V: "0",
}

func NewDecimal(vs ...decimal.Decimal) *Decimal {
	if len(vs) == 1 {
		return &Decimal{
			V: vs[0].String(),
		}
	}

	return &Decimal{
		V: "0",
	}
}

func NewDecimalByString(str string) (res *Decimal, err error) {
	var dec decimal.Decimal
	if dec, err = decimal.NewFromString(str); err != nil {
		return
	}

	res = &Decimal{
		V: dec.String(),
	}

	return
}

func (x *Decimal) Decimal() (decimal.Decimal, error) {
	return decimal.NewFromString(x.V)
}

func (x *Decimal) DecimalP() decimal.Decimal {
	return fnPanic.Value(x.Decimal())
}

func (x *Decimal) Scan(value any) (err error) {
	var body, ok = value.(string)
	if !ok {
		return fnError.NewF(ErrInvalidDecimal)
	}

	var dec decimal.Decimal
	if dec, err = decimal.NewFromString(body); err != nil {
		return
	}

	*x = Decimal{
		V: dec.String(),
	}

	return nil
}

func (x *Decimal) Value() (res driver.Value, err error) {
	var dec decimal.Decimal
	if dec, err = decimal.NewFromString(x.V); err != nil {
		return
	}
	res = dec.String()
	return
}

func (x *Decimal) Neg() (dec *Decimal, err error) {
	var d decimal.Decimal
	if d, err = x.Decimal(); err != nil {
		return
	}

	dec = NewDecimal(d.Neg())
	return
}

/* ------------------------------------------------------------------------------------------------------------ */

func MarshalDecimal(v *Decimal) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(v.V))
	})
}

func UnmarshalDecimal(v any) (res *Decimal, err error) {
	switch t := v.(type) {
	case string:
		return NewDecimalByString(t)
	case []byte:
		return NewDecimalByString(string(t))
	case json.Number:
		return NewDecimalByString(t.String())
	default:
		err = fnError.NewF(ErrInvalidDecimal)
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

type decimalValueCodec struct {
}

func (x *decimalValueCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(Decimal)
	if !isOk {
		err = fnError.NewF(ErrInvalidDecimal)
		return
	}

	var dec primitive.Decimal128
	if dec, err = primitive.ParseDecimal128(i.V); err != nil {
		return
	}

	return writer.WriteDecimal128(dec)
}

func (x *decimalValueCodec) DecodeValue(
	_ bsoncodec.DecodeContext,
	reader bsonrw.ValueReader,
	value reflect.Value,
) (err error) {
	if err = reader.ReadNull(); err == nil {
		value.Set(reflect.ValueOf(Decimal{V: "0"}))
		return
	}

	var dec primitive.Decimal128
	if dec, err = reader.ReadDecimal128(); err != nil {
		return
	}

	value.Set(reflect.ValueOf(Decimal{
		V: dec.String(),
	}))

	return
}

type decimalPointerCodec struct{}

func (x *decimalPointerCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(*Decimal)
	if !isOk {
		err = fnError.NewF(ErrInvalidDecimal)
		return
	}

	if fnPointer.IsNil(i) {
		return writer.WriteNull()
	}

	var dec primitive.Decimal128
	if dec, err = primitive.ParseDecimal128(i.V); err != nil {
		return
	}

	return writer.WriteDecimal128(dec)
}

func (x *decimalPointerCodec) DecodeValue(
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

	value.Set(reflect.ValueOf(&Decimal{
		V: dec.String(),
	}))

	return
}

/* ------------------------------------------------------------------------------------------------------------ */

func DecimalRegistry(registry *bsoncodec.Registry) *bsoncodec.Registry {
	var value, valueTypeof = &decimalValueCodec{}, reflect.TypeOf(Decimal{})
	registry.RegisterTypeEncoder(valueTypeof, value)
	registry.RegisterTypeDecoder(valueTypeof, value)

	var ptr, ptrTypeof = &decimalPointerCodec{}, reflect.TypeOf(&Decimal{})
	registry.RegisterTypeEncoder(ptrTypeof, ptr)
	registry.RegisterTypeDecoder(ptrTypeof, ptr)
	return registry
}

func (x *DecimalArgs) AppendFilter(filter bson.M, key string) bson.M {
	return mgQuery.AppendFilterCompareArgs[Decimal](filter, key, x)
}
