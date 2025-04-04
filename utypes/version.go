package utypes

import (
	"database/sql/driver"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/d3v-friends/go-tools/fnError"
	"github.com/d3v-friends/go-tools/fnPointer"
	"github.com/d3v-friends/mango/mgOp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"io"
	"reflect"
	"regexp"
	"strconv"
)

const ErrInvalidVersion = "invalid_version"
const regexpVersion = "v[0-9]{1,}\\.[0-9]{1,}\\.[0-9]{1,}(-[0-9]{14})?"

var NilVersion = &Version{V: "v0.0.0"}

func NewVersion(
	major, minor, patch uint64,
) *Version {
	return &Version{
		V: fmt.Sprintf("v%d.%d.%d", major, minor, patch),
	}
}

func (x *Version) Scan(value any) (err error) {
	var body, ok = value.(string)
	if !ok {
		return fnError.NewF(ErrInvalidVersion)
	}

	*x = Version{
		V: body,
	}

	var matched, _ = regexp.MatchString(regexpVersion, body)
	if !matched {
		err = fnError.NewF(ErrInvalidVersion)
		return
	}

	return nil
}

func (x *Version) Value() (res driver.Value, err error) {
	var matched, _ = regexp.MatchString(regexpVersion, x.V)
	if !matched {
		err = fnError.NewF(ErrInvalidVersion)
		return
	}

	res = x.V
	return
}

func (x *Version) GormDBDataType(
	_ *gorm.DB,
	_ *schema.Field,
) string {
	return "char(50)"
}

/* ------------------------------------------------------------------------------------------------------------ */

func MarshalVersion(v *Version) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		_, _ = w.Write([]byte(strconv.Quote(v.V)))
	})
}

func UnmarshalVersion(v any) (res *Version, err error) {
	switch t := v.(type) {
	case string:
		var matched, _ = regexp.MatchString(regexpVersion, t)
		if !matched {
			err = fnError.NewF(ErrInvalidVersion)
			return
		}
		return &Version{V: t}, nil
	case []byte:
		return UnmarshalVersion(string(t))
	default:
		err = fnError.NewF(ErrInvalidVersion)
		return
	}
}

/* ------------------------------------------------------------------------------------------------------------ */

func (x *VersionArgs) AppendFilter(filter bson.M, key string) bson.M {
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

/* ------------------------------------------------------------------------------------------------------------ */

type versionValueCodec struct {
}

func (x *versionValueCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(Version)
	if !isOk {
		err = fnError.NewF(ErrInvalidVersion)
		return
	}

	return writer.WriteString(i.V)
}

func (x *versionValueCodec) DecodeValue(
	_ bsoncodec.DecodeContext,
	reader bsonrw.ValueReader,
	value reflect.Value,
) (err error) {
	if err = reader.ReadNull(); err == nil {
		value.Set(reflect.ValueOf(Version{V: ""}))
		return
	}

	var str string
	if str, err = reader.ReadString(); err != nil {
		return
	}

	value.Set(reflect.ValueOf(Version{
		V: str,
	}))

	return
}

type versionPtrCodec struct {
}

func (x *versionPtrCodec) EncodeValue(
	_ bsoncodec.EncodeContext,
	writer bsonrw.ValueWriter,
	value reflect.Value,
) (err error) {
	var i, isOk = value.Interface().(*Version)
	if !isOk {
		err = fnError.NewF(ErrInvalidVersion)
		return
	}

	if fnPointer.IsNil(i) {
		return writer.WriteNull()
	}

	return writer.WriteString(i.V)
}

func (x *versionPtrCodec) DecodeValue(
	_ bsoncodec.DecodeContext,
	reader bsonrw.ValueReader,
	value reflect.Value,
) (err error) {
	if err = reader.ReadNull(); err == nil {
		value.Set(reflect.Zero(value.Type()))
		return
	}

	var str string
	if str, err = reader.ReadString(); err != nil {
		return
	}

	value.Set(reflect.ValueOf(&Version{
		V: str,
	}))

	return
}

func VersionRegistry(register *bsoncodec.Registry) *bsoncodec.Registry {
	var value, valueTypeof = &versionValueCodec{}, reflect.TypeOf(Version{})
	register.RegisterTypeEncoder(valueTypeof, value)
	register.RegisterTypeDecoder(valueTypeof, value)

	var ptr, ptrTypeof = &versionPtrCodec{}, reflect.TypeOf(&Version{})
	register.RegisterTypeEncoder(ptrTypeof, ptr)
	register.RegisterTypeDecoder(ptrTypeof, ptr)
	return register
}
