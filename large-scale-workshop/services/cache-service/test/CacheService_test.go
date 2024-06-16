package test

import (
	"testing"

	service "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/cache-service/service"
	"github.com/stretchr/testify/assert"
)

func TestCacheService(t *testing.T) {
	cs, err := service.NewCacheService("TestNode", 1099, "")
	assert.NoError(t, err)

	cs.Set("key1", "value1")
	value, err := cs.Get("key1")
	assert.NoError(t, err)
	assert.Equal(t, "value1", value)

	cs.Delete("key1")
	value, err = cs.Get("key1")
	assert.Error(t, err)

	cs.Set("key2", "value2")
	value, err = cs.Get("key2")
	assert.Equal(t, "value2", value)

}