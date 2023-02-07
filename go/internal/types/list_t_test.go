package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListForEachF(t *testing.T) {
	list := NewList()
	for j := 0; j < 5; j++ {
		list.LPush(fmt.Sprint(j))
	}
	c := 4
	list.ForEachF(func(a string) {
		assert.Equal(t, a, fmt.Sprint(c))
		c--
	})
}
