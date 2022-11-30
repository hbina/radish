package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortedSet1(t *testing.T) {
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
		rn := ss.GetRangeByIndex(3, 3, false)
		assert.Equal(t, 1, len(rn))
		r := rn[0]
		assert.Equal(t, "d", r.key)
		assert.Equal(t, 4, r.score)
		assert.Equal(t, struct{}{}, r.value)
	}

	{
		rn := ss.GetRangeByIndex(1, 2, false)
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
		rn := ss.GetRangeByIndex(0, -2, true)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"g", "f", "e", "d", "c", "b"}, res)
	}

	{
		res := make([]string, 0)
		rn := ss.GetRangeByIndex(-5, 2, false)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"c"}, res)
	}

	{
		res := make([]string, 0)
		rn := ss.GetRangeByIndex(0, -2, false)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, res)
	}

	{
		res := make([]string, 0)
		rn := ss.GetRangeByIndex(5, -1, false)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"f", "g"}, res)
	}

	{
		res := make([]string, 0)
		rn := ss.GetRangeByIndex(0, -5, false)
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"a", "b", "c"}, res)
	}

	start, end := ss.SanitizeIndex(0, 6, false)
	assert.Equal(t, 1, start)
	assert.Equal(t, 7, end)
	start, end = ss.SanitizeIndex(0, 6, true)
	assert.Equal(t, 1, start)
	assert.Equal(t, 7, end)

	start, end = ss.SanitizeIndex(0, 2, false)
	assert.Equal(t, 1, start)
	assert.Equal(t, 3, end)
	start, end = ss.SanitizeIndex(0, 2, true)
	assert.Equal(t, 5, start)
	assert.Equal(t, 7, end)

	start, end = ss.SanitizeIndex(2, 5, false)
	assert.Equal(t, 3, start)
	assert.Equal(t, 6, end)
	start, end = ss.SanitizeIndex(2, 5, true)
	assert.Equal(t, 2, start)
	assert.Equal(t, 5, end)

	start, end = ss.SanitizeIndex(-3, -1, false)
	assert.Equal(t, 5, start)
	assert.Equal(t, 7, end)
	start, end = ss.SanitizeIndex(-3, -1, true)
	assert.Equal(t, 1, start)
	assert.Equal(t, 3, end)

	start, end = ss.SanitizeIndex(4, 3, false)
	assert.Equal(t, 1, start)
	assert.Equal(t, -1, end)
	start, end = ss.SanitizeIndex(4, 3, true)
	assert.Equal(t, 1, start)
	assert.Equal(t, -1, end)

	for i := 0; i < ss.Len(); i++ {
		res := ss.GetByIndex(i, false)
		assert.Equal(t, i+1, res.score)
	}

	for i := -1; i >= -ss.Len(); i-- {
		res := ss.GetByIndex(i, false)
		assert.Equal(t, ss.Len()+i+1, res.score)
	}
}

func TestSortedSet2(t *testing.T) {
	ss := NewSortedSet[string, int, struct{}]()

	ss.AddOrUpdate("a", 1, struct{}{})
	ss.AddOrUpdate("b", 2, struct{}{})
	ss.AddOrUpdate("c", 3, struct{}{})
	ss.AddOrUpdate("d", 4, struct{}{})

	{
		rn := ss.GetRangeByIndex(-5, 2, false)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"a", "b", "c"}, res)
	}
	{
		rn := ss.GetRangeByIndex(1, 5, true)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"c", "b", "a"}, res)
	}
}
