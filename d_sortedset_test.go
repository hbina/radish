package redis

import (
	"math"
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
		rn := ss.GetRangeByIndex(3, 3, DefaultRangeOptions())
		assert.Equal(t, 1, len(rn))

		res := make([]string, 0)
		for _, r := range rn {
			res = append(res, r.Key())
		}

		assert.Equal(t, []string{"d"}, res)
	}

	{
		rn := ss.GetRangeByIndex(1, 2, DefaultRangeOptions())

		res := make([]string, 0)
		for _, r := range rn {
			res = append(res, r.Key())
		}

		assert.Equal(t, []string{"b", "c"}, res)
	}

	{
		options := DefaultRangeOptions()
		options.reverse = true
		rn := ss.GetRangeByIndex(0, -2, options)

		res := make([]string, 0)
		for _, r := range rn {
			res = append(res, r.Key())
		}

		assert.Equal(t, []string{"g", "f", "e", "d", "c", "b"}, res)
	}

	{
		rn := ss.GetRangeByIndex(-5, 2, DefaultRangeOptions())

		res := make([]string, 0)
		for _, r := range rn {
			res = append(res, r.Key())
		}

		assert.Equal(t, []string{"c"}, res)
	}

	{
		rn := ss.GetRangeByIndex(0, -2, DefaultRangeOptions())

		res := make([]string, 0)
		for _, r := range rn {
			res = append(res, r.Key())
		}

		assert.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, res)
	}

	{
		rn := ss.GetRangeByIndex(5, -1, DefaultRangeOptions())

		res := make([]string, 0)
		for _, r := range rn {
			res = append(res, r.Key())
		}

		assert.Equal(t, []string{"f", "g"}, res)
	}

	{
		rn := ss.GetRangeByIndex(0, -5, DefaultRangeOptions())

		res := make([]string, 0)
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
		rn := ss.GetRangeByIndex(-5, 2, DefaultRangeOptions())
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"a", "b", "c"}, res)
	}
	{
		options := DefaultRangeOptions()
		options.reverse = true
		rn := ss.GetRangeByIndex(1, 5, options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"c", "b", "a"}, res)
	}
}

func TestSortedSet3(t *testing.T) {
	ss := NewSortedSet[string, float64, struct{}]()

	ss.AddOrUpdate("a", math.Inf(-1), struct{}{})
	ss.AddOrUpdate("b", 1, struct{}{})
	ss.AddOrUpdate("c", 2, struct{}{})
	ss.AddOrUpdate("d", 3, struct{}{})
	ss.AddOrUpdate("e", 4, struct{}{})
	ss.AddOrUpdate("f", 5, struct{}{})
	ss.AddOrUpdate("g", math.Inf(1), struct{}{})

	{
		options := DefaultRangeOptions()
		rn := ss.GetRangeByScore(math.Inf(-1), 2, options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"a", "b", "c"}, res)
	}
	{
		options := DefaultRangeOptions()
		rn := ss.GetRangeByScore(0, 3, options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"b", "c", "d"}, res)
	}
	{
		options := DefaultRangeOptions()
		rn := ss.GetRangeByScore(3, 6, options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"d", "e", "f"}, res)
	}
	{
		options := DefaultRangeOptions()
		options.reverse = true
		rn := ss.GetRangeByScore(2, math.Inf(-1), options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"c", "b", "a"}, res)
	}
	{
		options := DefaultRangeOptions()
		options.offset = 2
		options.limit = 3
		rn := ss.GetRangeByScore(0, 10, options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"d", "e", "f"}, res)
	}
	{
		options := DefaultRangeOptions()
		options.offset = 2
		options.limit = 10
		rn := ss.GetRangeByScore(0, 10, options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"d", "e", "f"}, res)
	}
	{
		options := DefaultRangeOptions()
		options.reverse = true
		options.offset = 0
		options.limit = 2
		rn := ss.GetRangeByScore(10, 0, options)
		assert.Equal(t, 2, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key())
		}
		assert.Equal(t, []string{"f", "e"}, res)
	}
}

func TestSortedSet4(t *testing.T) {
	ss := NewSortedSet[string, float64, struct{}]()

	ss.AddOrUpdate("b", 1, struct{}{})
	ss.AddOrUpdate("c", 2, struct{}{})
	ss.AddOrUpdate("d", 3, struct{}{})
	ss.AddOrUpdate("e", 4, struct{}{})
	ss.AddOrUpdate("f", 5, struct{}{})

	{
		options := DefaultRangeOptions()
		rn := ss.GetRangeByScore(6, math.Inf(1), options)
		assert.Equal(t, 0, len(rn))
	}
}
