package redis

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetGetRandomMembers(t *testing.T) {
	set := NewSetEmpty()
	for i := 0; i < 50; i++ {
		set.AddMember(fmt.Sprint(i))
	}
	set2 := NewSetEmpty()
	for i := 0; i < 1000; i++ {
		set2.AddMember(*set.GetRandomMember())
		if set.Len() == set2.Len() {
			break
		}
	}
	assert.Equal(t, set.Len(), set2.Len())
}
