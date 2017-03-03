package main

import (
	"github.com/go-scim/scimify/service"
	"net/http"
)

//var once sync.Once

//func bootstrapData() {
//	once.Do(func() {
//		if common, _, err := helper.LoadSchema("./schemas/common_schema.json"); err != nil {
//			panic(err)
//		} else if user, _, err := helper.LoadSchema("./schemas/user_schema.json"); err != nil {
//			panic(err)
//		} else {
//			repository := persistence.GetSchemaRepository()
//
//			finalUserSchema := &resource.Schema{
//				Id:user.Id,
//				Name:user.Name,
//				Schemas:user.Schemas,
//				Description:user.Description,
//				Attributes:make([]*resource.Attribute, 0),
//			}
//			finalUserSchema.MergeWith(common, user)
//
//			// We need to create a unified interface for ScimObject (getId, getData)
//			repository.Create()
//		}
//	})
//}

func main() {
	http.ListenAndServe(":8080", service.Mux())
}
