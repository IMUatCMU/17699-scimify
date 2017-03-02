package service

import (
	"github.com/go-zoo/bone"
	"net/http"
	"sync"
)

var (
	oneServer sync.Once
	mux       *bone.Mux
)

func Mux() *bone.Mux {
	oneServer.Do(func() {
		mux = bone.New()
		mux.Prefix("/v2")
		mux.GetFunc("/Schemas", http.HandlerFunc(getAllSchemas))
		mux.GetFunc("/Schemas/:schemaId", http.HandlerFunc(getSchemaById))
	})
	return mux
}
