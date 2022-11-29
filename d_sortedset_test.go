package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedSetGetRangeByIndex(t *testing.T) {
	ss := NewSortedSet[string, int, struct{}]()

	ss.AddOrUpdate("a", 1, struct{}{})
	ss.AddOrUpdate("b", 2, struct{}{})
	ss.AddOrUpdate("c", 3, struct{}{})
	ss.AddOrUpdate("d", 4, struct{}{})
	ss.AddOrUpdate("e", 5, struct{}{})
	ss.AddOrUpdate("f", 6, struct{}{})
	ss.AddOrUpdate("g", 7, struct{}{})

	assert.Equal(t, 7, ss.Len())

	{
		rn := ss.GetRangeByIndex(3, 3, false, false)
		assert.Equal(t, 1, len(rn))
		r := rn[0]
		assert.Equal(t, "d", r.key)
		assert.Equal(t, 4, r.score)
		assert.Equal(t, struct{}{}, r.value)
	}

	{
		rn := ss.GetRangeByIndex(1, 2, false, false)
		assert.Equal(t, 2, len(rn))
		r := rn[0]
		assert.Equal(t, "b", r.key)
		assert.Equal(t, 2, r.score)
		assert.Equal(t, struct{}{}, r.value)
		r = rn[1]
		assert.Equal(t, "c", r.key)
		assert.Equal(t, 3, r.score)
		assert.Equal(t, struct{}{}, r.value)
	}

	{
		res := make([]string, 0)
		rn := ss.GetRangeByIndex(0, -2, true, false)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"g", "f", "e", "d", "c", "b"}, res)
	}

	{
		res := make([]string, 0)
		rn := ss.GetRangeByIndex(-5, 2, false, false)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"c"}, res)
	}

	{
		res := make([]string, 0)
		rn := ss.GetRangeByIndex(0, -2, false, false)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, res)
	}

	{
		res := make([]string, 0)
		rn := ss.GetRangeByIndex(5, -1, false, false)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"f", "g"}, res)
	}

	{
		res := make([]string, 0)
		rn := ss.GetRangeByIndex(0, -5, false, false)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"a", "b", "c"}, res)
	}
}

func TestConvertIndexToRank(t *testing.T) {
	ss := NewSortedSet[string, int, struct{}]()

	ss.AddOrUpdate("a", 1, struct{}{})
	ss.AddOrUpdate("b", 2, struct{}{})
	ss.AddOrUpdate("c", 3, struct{}{})
	ss.AddOrUpdate("d", 4, struct{}{})
	ss.AddOrUpdate("e", 5, struct{}{})
	ss.AddOrUpdate("f", 6, struct{}{})
	ss.AddOrUpdate("g", 7, struct{}{})

	start, end := ss.ConvertIndexToRank(0, 6, false)
	assert.Equal(t, 1, start)
	assert.Equal(t, 7, end)
	start, end = ss.ConvertIndexToRank(0, 6, true)
	assert.Equal(t, 1, start)
	assert.Equal(t, 7, end)

	start, end = ss.ConvertIndexToRank(0, 2, false)
	assert.Equal(t, 1, start)
	assert.Equal(t, 3, end)
	start, end = ss.ConvertIndexToRank(0, 2, true)
	assert.Equal(t, 5, start)
	assert.Equal(t, 7, end)

	start, end = ss.ConvertIndexToRank(2, 5, false)
	assert.Equal(t, 3, start)
	assert.Equal(t, 6, end)
	start, end = ss.ConvertIndexToRank(2, 5, true)
	assert.Equal(t, 2, start)
	assert.Equal(t, 5, end)

	start, end = ss.ConvertIndexToRank(-3, -1, false)
	assert.Equal(t, 5, start)
	assert.Equal(t, 7, end)
	start, end = ss.ConvertIndexToRank(-3, -1, true)
	assert.Equal(t, 1, start)
	assert.Equal(t, 3, end)

	start, end = ss.ConvertIndexToRank(4, 3, false)
	assert.Equal(t, 1, start)
	assert.Equal(t, -1, end)
	start, end = ss.ConvertIndexToRank(4, 3, true)
	assert.Equal(t, 1, start)
	assert.Equal(t, -1, end)
}

func TestSortedSetGetByRank(t *testing.T) {
	ss := NewSortedSet[string, int, struct{}]()

	ss.AddOrUpdate("a", 1, struct{}{})
	ss.AddOrUpdate("b", 2, struct{}{})
	ss.AddOrUpdate("c", 3, struct{}{})
	ss.AddOrUpdate("d", 4, struct{}{})
	ss.AddOrUpdate("e", 5, struct{}{})
	ss.AddOrUpdate("f", 6, struct{}{})
	ss.AddOrUpdate("g", 7, struct{}{})

	assert.Equal(t, 7, ss.Len())

	for i := 0; i < ss.Len(); i++ {
		res := ss.GetByIndex(i, false)
		assert.Equal(t, i+1, res.score)
	}

	for i := -1; i >= -ss.Len(); i-- {
		res := ss.GetByIndex(i, false)
		assert.Equal(t, ss.Len()+i+1, res.score)
	}
}

func TestSortedFindRankOfKey(t *testing.T) {
	ss := NewSortedSet[string, int, struct{}]()

	ss.AddOrUpdate("x", 10, struct{}{})
	ss.AddOrUpdate("y", 20, struct{}{})
	ss.AddOrUpdate("z", 30, struct{}{})

	assert.Equal(t, 3, ss.Len())

	assert.Equal(t, 0, ss.FindRankOfKey("x"))
	assert.Equal(t, 1, ss.FindRankOfKey("y"))
	assert.Equal(t, 2, ss.FindRankOfKey("z"))
}
