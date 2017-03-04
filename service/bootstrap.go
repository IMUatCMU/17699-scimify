package service

import (
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-zoo/bone"
	"net/http"
	"sync"
)

var (
	oneServer, oneDataInit sync.Once
	mux                    *bone.Mux
)

func Bootstrap() *bone.Mux {
	oneDataInit.Do(func() {
		repo := persistence.GetSchemaRepository()
		for _, path := range []string{
			"./schemas/user_schema_all.json",
		} {
			if schema, _, err := helper.LoadSchema(path); err != nil {
				panic(err)
			} else if err = repo.Create(schema, nil); err != nil {
				panic(err)
			}
		}
	})
	oneServer.Do(func() {
		mux = bone.New()
		mux.Prefix("/v2")
		mux.GetFunc("/Schemas", http.HandlerFunc(getAllSchemas))
		mux.GetFunc("/Schemas/:schemaId", http.HandlerFunc(getSchemaById))
	})
	return mux
}
