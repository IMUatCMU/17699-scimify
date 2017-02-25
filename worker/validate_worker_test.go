package worker

import (
	"context"
	"github.com/go-scim/scimify/helper"
	"github.com/go-scim/scimify/resource"
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

	r, _, err := helper.LoadResource("../test_data/single_test_user_david.json")
	if err != nil {
		b.Fatal(err)
	}

	input := &ValidationInput{
		Resource: r,
		Option:   validation.ValidationOptions{ReadOnlyIsMandatory: false, UnassignedImmutableIsIgnored: false},
		Context:  context.WithValue(context.Background(), resource.CK_Schema, schema),
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			worker.Do(input)
		}
	})
}
