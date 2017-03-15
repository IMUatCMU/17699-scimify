package processor

import (
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateMetaProcessor_Process(t *testing.T) {
	viper.Set("server.rootPath", "http://foo.com/v2/")

	r := resource.NewResourceFromMap(map[string]interface{}{
		"id": "bar",
	})

	processor := &generateMetaProcessor{
		ResourceType:    "User",
		ResourceTypeUri: "/User",
	}
	err := processor.Process(&ProcessorContext{Resource: r})
	assert.Nil(t, err)

	meta, ok := r.Attributes["meta"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotNil(t, meta)
	assert.Equal(t, "User", meta["resourceType"].(string))
	assert.Equal(t, "http://foo.com/v2/User/bar", meta["location"].(string))
	assert.NotEmpty(t, meta["version"].(string))
	assert.NotEmpty(t, meta["created"].(string))
	assert.NotEmpty(t, meta["lastModified"].(string))
}
