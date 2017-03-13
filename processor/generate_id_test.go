package processor

import (
	"testing"
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
)

func TestGenerateIdProcessor_Process(t *testing.T) {
	processor := &generateIdProcessor{}
	r := resource.NewResource()
	err := processor.Process(r, nil)
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(r.Attributes["id"].(string)))
}
