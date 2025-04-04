package utypes

import "github.com/d3v-friends/go-tools/fnPointer"

func (x *PageArgs) GetSkip() *int64 {
	return fnPointer.Make(x.Size * x.Page)
}

func (x *PageArgs) GetLimit() *int64 {
	return fnPointer.Make(x.Size)
}
