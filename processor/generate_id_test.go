package processor

import (
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateIdProcessor_Process(t *testing.T) {
	processor := &generateIdProcessor{}
	r := resource.NewResource()
	err := processor.Process(&ProcessorContext{Resource: r})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(r.Attributes["id"].(string)))
}
