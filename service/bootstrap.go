package service

import (
	"flag"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/persistence"
	"github.com/go-zoo/bone"
	"github.com/spf13/viper"
	"sync"
)

var (
	oneServer,
	oneDataInit,
	oneConfig,
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
	oneConfig.Do(func() {
		var yamlPath string
		flag.StringVar(&yamlPath, "config", "./config.yaml", "Location for configuration file")
		flag.Parse()
		viper.SetConfigFile(yamlPath)
		err := viper.ReadInConfig()
		if nil != err {
			panic(err)
		}
	})
	oneDataInit.Do(func() {
		schemaRepo := persistence.GetSchemaRepository()
		for _, path := range viper.GetStringSlice("scim.stock.schema") {
			if schema, _, err := helper.LoadSchema(path); err != nil {
				panic(err)
			} else if err = schemaRepo.Create(schema, nil); err != nil {
				panic(err)
			}
		}

		resourceTypeRepo := persistence.GetResourceTypeRepository()
		for _, path := range viper.GetStringSlice("scim.stock.resource_type") {
			if resource, _, err := helper.LoadResource(path); err != nil {
				panic(err)
			} else if err = resourceTypeRepo.Create(resource, nil); err != nil {
				panic(err)
			}
		}

		serviceProviderConfigRepo := persistence.GetServiceProviderConfigRepository()
		for _, path := range viper.GetStringSlice("scim.stock.sp_config") {
			if resource, _, err := helper.LoadResource(path); err != nil {
				panic(err)
			} else if err = serviceProviderConfigRepo.Create(resource, nil); err != nil {
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
