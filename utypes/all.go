package utypes

import (
	"github.com/d3v-friends/mango/mgConn"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
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

func NewRegistry() (registry *bsoncodec.Registry) {
	registry = bson.NewRegistry()
	for _, fn := range Registries {
		registry = fn(registry)
	}
	return
}
