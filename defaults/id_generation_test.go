package defaults

import (
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdGenerationValueDefaulter_Default(t *testing.T) {
	defaulter := &idGenerationValueDefaulter{}
	r := resource.NewResource()
	defaulter.Default(r, nil)
	assert.NotEqual(t, 0, len(r.Attributes["id"].(string)))
}
