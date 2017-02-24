package worker

import (
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/validation"
	"testing"
)

type validateWorkerTest struct {
	name         string
	resourcePath string
}

func BenchmarkValidateWorker(b *testing.B) {
	worker := GetCreationValidatorWorker()

	schema, _, err := helper.LoadSchema("../test_data/test_user_schema_all.json")
	if err != nil {
		b.Fatal(err)
	}

	resource, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	input := &ValidationInput{
		Resource: resource,
		Context: &validation.ValidatorContext{
			Data: map[string]interface{}{
				validation.Schema: schema,
			},
		},
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			worker.Do(input)
		}
	})
}
