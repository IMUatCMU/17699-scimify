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
		viper.SetDefault("scim.resourceType.user", "User")
		viper.SetDefault("scim.resourceType.group", "Group")
		viper.SetDefault("scim.resourceTypeUri.user", "/Users")
		viper.SetDefault("scim.resourceTypeUri.group", "/Groups")
		viper.SetDefault("scim.internalSchemaId.root", "")
		viper.SetDefault("scim.internalSchemaId.user", "urn:ietf:params:scim:schemas:core:2.0:User")
		viper.SetDefault("scim.internalSchemaId.group", "urn:ietf:params:scim:schemas:core:2.0:Group")
		viper.SetDefault("scim.internalSchemaId.resourceType", "urn:ietf:params:scim:schemas:core:2.0:ResourceType")
		viper.SetDefault("scim.internalSchemaId.spConfig", "urn:ietf:params:scim:schemas:core:2.0:ServiceProviderConfig")
		viper.SetDefault("scim.internalSchemaId.schema", "urn:ietf:params:scim:schemas:core:2.0:Schema")
		viper.SetDefault("scim.internalSchemaId.error", "urn:ietf:params:scim:api:messages:2.0:Error")
		viper.SetDefault("scim.internalSchemaId.listResponse", "urn:ietf:params:scim:api:messages:2.0:ListResponse")
		viper.SetDefault("scim.api.userIdUrlParam", "userId")
		viper.SetDefault("scim.api.groupIdUrlParam", "groupId")
		viper.SetDefault("scim.api.schemaIdUrlParam", "schemaId")
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
		internalSchemaRepo := persistence.GetInternalSchemaRepository()
		internalSchemas := viper.GetStringSlice("scim.stock.internal_schema")
		for _, v := range internalSchemas {
			if schema, _, err := helper.LoadSchema(v); err != nil {
				panic(err)
			} else if err = internalSchemaRepo.Create(schema); err != nil {
				panic(err)
			}
		}

		schemaRepo := persistence.GetSchemaRepository()
		for _, path := range viper.GetStringSlice("scim.stock.schema") {
			if schema, _, err := helper.LoadSchema(path); err != nil {
				panic(err)
			} else if err = schemaRepo.Create(schema); err != nil {
				panic(err)
			}
		}

		resourceTypeRepo := persistence.GetResourceTypeRepository()
		for _, path := range viper.GetStringSlice("scim.stock.resource_type") {
			if resource, _, err := helper.LoadResource(path); err != nil {
				panic(err)
			} else if err = resourceTypeRepo.Create(resource); err != nil {
				panic(err)
			}
		}

		serviceProviderConfigRepo := persistence.GetServiceProviderConfigRepository()
		for _, path := range viper.GetStringSlice("scim.stock.sp_config") {
			if resource, _, err := helper.LoadResource(path); err != nil {
				panic(err)
			} else if err = serviceProviderConfigRepo.Create(resource); err != nil {
				panic(err)
			}
		}
	})
	oneServer.Do(func() {
		schemaSrv = &schemaService{}
		resourceTypeSrv = &resourceTypeService{}
		userSrv = &userService{}
		spConfigSrv = &spConfigService{}

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
