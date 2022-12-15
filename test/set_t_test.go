package test

import (
	"fmt"
	"testing"

	"github.com/hbina/radish/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestSetGetRandomMembers(t *testing.T) {
	set := types.NewSetEmpty()
	for i := 0; i < 50; i++ {
		set.AddMember(fmt.Sprint(i))
	}
	set2 := types.NewSetEmpty()
	for i := 0; i < 1000; i++ {
		set2.AddMember(*set.GetRandomMember())
		if set.Len() == set2.Len() {
			break
		}
	}
	assert.Equal(t, set.Len(), set2.Len())
}
