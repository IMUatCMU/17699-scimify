package defaults

import (
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetaUpdateValueDefaulter_Default(t *testing.T) {
	r := resource.NewResourceFromMap(map[string]interface{}{
		"meta": map[string]interface{}{
			"version": "W\"1",
		},
	})

	defaulter := &metaUpdateValueDefaulter{}
	defaulter.Default(r, nil)
	assert.Equal(t, "W\"2", r.Attributes["meta"].(map[string]interface{})["version"].(string))
	assert.True(t, len(r.Attributes["meta"].(map[string]interface{})["lastModified"].(string)) > 0)
}
