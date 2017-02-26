package defaults

import (
	"context"
	"github.com/go-scim/scimify/resource"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetaGenerationValueDefaulter_Default(t *testing.T) {
	viper.Set("server.rootPath", "http://foo.com/v2/")

	ctx := context.Background()
	ctx = context.WithValue(ctx, resource.CK_ResourceType, "User")
	ctx = context.WithValue(ctx, resource.CK_ResourceTypeURI, "/User")

	r := resource.NewResourceFromMap(map[string]interface{}{
		"id": "bar",
	})

	defaulter := &metaGenerationValueDefaulter{}
	defaulter.Default(r, ctx)

	meta, ok := r.Attributes["meta"].(map[string]interface{})
	assert.True(t, ok)
	assert.NotNil(t, meta)
	assert.Equal(t, "User", meta["resourceType"].(string))
	assert.Equal(t, "http://foo.com/v2/User/bar", meta["location"].(string))
	assert.Equal(t, "W\"1", meta["version"].(string))
	assert.NotEmpty(t, meta["created"].(string))
	assert.NotEmpty(t, meta["lastModified"].(string))
}
