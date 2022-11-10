package redis

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestKeyExpirer(t *testing.T) {
	c := CreateTestClient()

	s, err := c.Set("a", "v", 53*time.Millisecond).Result()
	assert.Equal(t, "OK", s)
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)

	s, err = c.Get("a").Result()
	assert.NotEqual(t, "v", s)
	assert.Error(t, err)
}
