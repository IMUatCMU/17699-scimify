package service

import (
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-zoo/bone"
	"github.com/spf13/viper"
	"net/http"
	"sync"
)

var (
	oneServer,
	oneDataInit,
	oneConfigDefault sync.Once
	mux *bone.Mux
)

func Bootstrap() *bone.Mux {
	oneConfigDefault.Do(func() {
		viper.SetDefault("scim.itemsPerPage", 10)
	})
	oneDataInit.Do(func() {
		schemaRepo := persistence.GetSchemaRepository()
		for _, path := range []string{
			"./schemas/user_schema_all.json",
		} {
			if schema, _, err := helper.LoadSchema(path); err != nil {
				panic(err)
			} else if err = schemaRepo.Create(schema, nil); err != nil {
				panic(err)
			}
		}

		resourceTypeRepo := persistence.GetResourceTypeRepository()
		serviceProviderConfigRepo := persistence.GetServiceProviderConfigRepository()
		for _, each := range []struct {
			repo persistence.Repository
			path string
		}{
			{resourceTypeRepo, "./stock_data/resource_type/user_resource_type.json"},
			{resourceTypeRepo, "./stock_data/resource_type/group_resource_type.json"},
			{serviceProviderConfigRepo, "./stock_data/sp_config/sp_config.json"},
		} {
			if resource, _, err := helper.LoadResource(each.path); err != nil {
				panic(err)
			} else if err = each.repo.Create(resource, nil); err != nil {
				panic(err)
			}
		}
	})
	oneServer.Do(func() {
		mux = bone.New()
		mux.Prefix("/v2")

		mux.GetFunc("/Schemas", http.HandlerFunc(getAllSchemas))
		mux.GetFunc("/Schemas/:schemaId", http.HandlerFunc(getSchemaById))

		mux.GetFunc("/ResourceTypes", http.HandlerFunc(getAllResourceTypes))

		mux.GetFunc("/ServiceProviderConfig", http.HandlerFunc(getServiceProviderConfig))

		mux.GetFunc("/Users/:userId", http.HandlerFunc(getUserById))
		mux.PostFunc("/Users", http.HandlerFunc(createUser))
		mux.PutFunc("/Users/:userId", http.HandlerFunc(replaceUserById))
		mux.PatchFunc("/Users/:userId", http.HandlerFunc(updateUserById))
		mux.DeleteFunc("/Users/:userId", http.HandlerFunc(deleteUserById))
		mux.GetFunc("/Users", http.HandlerFunc(queryUser))
		mux.PostFunc("/Users/.search", http.HandlerFunc(queryUser))
	})
	return mux
}
