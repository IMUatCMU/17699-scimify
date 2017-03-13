package processor

import (
	"github.com/go-scim/scimify/resource"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpdateMetaProcessor_Process(t *testing.T) {
	r := resource.NewResourceFromMap(map[string]interface{}{
		"id": "4FD76312-456B-4233-A357-16EA035637E2",
		"meta": map[string]interface{}{
			"version": "W\\/\"a330bc54f0671c9\"",
		},
	})

	processor := &updateMetaProcessor{}
	err := processor.Process(r, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, r.Attributes["meta"].(map[string]interface{})["version"].(string))
	assert.True(t, len(r.Attributes["meta"].(map[string]interface{})["lastModified"].(string)) > 0)
}
