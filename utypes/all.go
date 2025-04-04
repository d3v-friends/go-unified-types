package utypes

import (
	"github.com/d3v-friends/mango/mgConn"
)

var Registries = []mgConn.CodecRegistry{
	DecimalRegistry,
	ObjectIDRegistry,
	RFC3339TimeRegistry,
	UnixNanoTimeRegistry,
	VersionRegistry,
	YMDHTimeRegistry,
	YMDTimeRegistry,
	SortDirectionRegistry,
}
