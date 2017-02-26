package defaults

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBumpVersion(t *testing.T) {
	v, err := BumpVersion("W\"1")
	assert.Nil(t, err)
	assert.Equal(t, "W\"2", v)
}
