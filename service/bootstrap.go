package service

import (
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-zoo/bone"
	"github.com/spf13/viper"
	"sync"
)

var (
	oneServer,
	oneDataInit,
	oneConfigDefault sync.Once

	schemaSrv       *schemaService
	resourceTypeSrv *resourceTypeService
	spConfigSrv     *spConfigService
	userSrv         *userService

	mux *bone.Mux
)

func Bootstrap() *bone.Mux {
	oneConfigDefault.Do(func() {
		viper.SetDefault("mongo.address", "localhost:32768")
		viper.SetDefault("mongo.database", "test_db")
		viper.SetDefault("mongo.userCollectionName", "users")
		viper.SetDefault("mongo.groupCollectionName", "groups")
		viper.SetDefault("scim.itemsPerPage", 10)
		viper.SetDefault("server.rootPath", "http://localhost:8080/v2/")
	})
	oneDataInit.Do(func() {
		schemaRepo := persistence.GetSchemaRepository()
		for _, path := range []string{
			"./stock_data/schema/core.json",
			"./stock_data/schema/user_schema.json",
			"./stock_data/schema/group_schema.json",
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
		schemaSrv = &schemaService{}
		resourceTypeSrv = &resourceTypeService{}
		userSrv = &userService{}

		mux = bone.New()
		mux.Prefix("/v2")

		mux.GetFunc("/Schemas", endpoint(schemaSrv.getAllSchemas))
		mux.GetFunc("/Schemas/:schemaId", endpoint(schemaSrv.getSchemaById))

		mux.GetFunc("/ResourceTypes", endpoint(resourceTypeSrv.getAllResourceTypes))

		mux.GetFunc("/ServiceProviderConfig", endpoint(spConfigSrv.getServiceProviderConfig))

		mux.GetFunc("/Users/:userId", endpoint(userSrv.getUserById))
		mux.PostFunc("/Users", endpoint(userSrv.createUser))
		mux.PutFunc("/Users/:userId", endpoint(userSrv.replaceUserById))
		mux.PatchFunc("/Users/:userId", endpoint(userSrv.updateUserById))
		mux.DeleteFunc("/Users/:userId", endpoint(userSrv.deleteUserById))
		mux.GetFunc("/Users", endpoint(userSrv.queryUser))
		mux.PostFunc("/Users/.search", endpoint(userSrv.queryUser))
	})
	return mux
}
