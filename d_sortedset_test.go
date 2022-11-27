package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedSet(t *testing.T) {
	ss := NewSortedSet[string, int, struct{}]()

	ss.AddOrUpdate("a", 1, struct{}{})
	ss.AddOrUpdate("b", 2, struct{}{})
	ss.AddOrUpdate("c", 3, struct{}{})
	ss.AddOrUpdate("d", 4, struct{}{})
	ss.AddOrUpdate("e", 5, struct{}{})
	ss.AddOrUpdate("f", 6, struct{}{})
	ss.AddOrUpdate("g", 7, struct{}{})

	assert.Equal(t, 7, ss.Len())

	rn := ss.GetRangeByRank(4, 4, false, true)
	assert.Equal(t, 1, len(rn))
	assert.Equal(t, 6, ss.Len())

	r := rn[0]
	assert.Equal(t, "d", r.key)
	assert.Equal(t, 4, r.score)
	assert.Equal(t, struct{}{}, r.value)

	{
		rn := ss.GetRangeByRank(1, 2, false, false)
		assert.Equal(t, 2, len(rn))
	}
}

func TestSortedSetGetRangeByRank(t *testing.T) {
	ss := NewSortedSet[string, int, struct{}]()

	ss.AddOrUpdate("a", 1, struct{}{})
	ss.AddOrUpdate("b", 2, struct{}{})
	ss.AddOrUpdate("c", 3, struct{}{})
	ss.AddOrUpdate("d", 4, struct{}{})
	assert.Equal(t, 4, ss.Len())

	{
		res := make([]string, 0)
		rn := ss.GetRangeByIndex(-5, 2, false, false)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"a", "b", "c"}, res)
	}

	{
		res := make([]string, 0)
		rn := ss.GetRangeByIndex(0, -2, false, false)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"a", "b", "c"}, res)
	}

	{
		res := make([]string, 0)
		rn := ss.GetRangeByIndex(5, -1, false, false)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{}, res)
	}

	{
		res := make([]string, 0)
		rn := ss.GetRangeByIndex(0, -5, false, false)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{}, res)
	}
}

func TestRecalibrateIndex(t *testing.T) {
	ss := NewSortedSet[string, int, struct{}]()

	ss.AddOrUpdate("a", 1, struct{}{})
	ss.AddOrUpdate("b", 2, struct{}{})
	ss.AddOrUpdate("c", 3, struct{}{})
	ss.AddOrUpdate("d", 4, struct{}{})
	ss.AddOrUpdate("e", 5, struct{}{})
	ss.AddOrUpdate("f", 6, struct{}{})
	ss.AddOrUpdate("g", 7, struct{}{})

	start, end := ss.RecalibrateRank(1, 7, false)
	assert.Equal(t, 1, start)
	assert.Equal(t, 7, end)
	start, end = ss.RecalibrateRank(1, 7, true)
	assert.Equal(t, 1, start)
	assert.Equal(t, 7, end)

	start, end = ss.RecalibrateRank(1, 3, false)
	assert.Equal(t, 1, start)
	assert.Equal(t, 3, end)
	start, end = ss.RecalibrateRank(1, 3, true)
	assert.Equal(t, 5, start)
	assert.Equal(t, 7, end)

	start, end = ss.RecalibrateRank(2, 5, false)
	assert.Equal(t, 2, start)
	assert.Equal(t, 5, end)
	start, end = ss.RecalibrateRank(2, 5, true)
	assert.Equal(t, 3, start)
	assert.Equal(t, 6, end)

	start, end = ss.RecalibrateRank(-3, -1, false)
	assert.Equal(t, 5, start)
	assert.Equal(t, 7, end)
	start, end = ss.RecalibrateRank(-3, -1, true)
	assert.Equal(t, 1, start)
	assert.Equal(t, 3, end)

	start, end = ss.RecalibrateRank(4, 3, false)
	assert.Equal(t, 4, start)
	assert.Equal(t, 3, end)
	start, end = ss.RecalibrateRank(4, 3, true)
	assert.Equal(t, 5, start)
	assert.Equal(t, 4, end)
}
