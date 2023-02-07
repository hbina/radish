package test

import (
	"math"
	"testing"

	"github.com/hbina/radish/internal/types"
	"github.com/stretchr/testify/assert"
)

func TestSortedSet1(t *testing.T) {
	ss := types.NewSortedSet()

	ss.AddOrUpdate("a", 1)
	ss.AddOrUpdate("b", 2)
	ss.AddOrUpdate("c", 3)
	ss.AddOrUpdate("d", 4)
	ss.AddOrUpdate("e", 5)
	ss.AddOrUpdate("f", 6)
	ss.AddOrUpdate("g", 7)

	assert.Equal(t, 7, ss.Len())

	{
		rn := ss.GetRangeByIndex(3, 3, types.DefaultRangeOptions())
		assert.Equal(t, 1, len(rn))

		res := make([]string, 0)
		for _, r := range rn {
			res = append(res, r.Key)
		}

		assert.Equal(t, []string{"d"}, res)
	}

	{
		rn := ss.GetRangeByIndex(1, 2, types.DefaultRangeOptions())

		res := make([]string, 0)
		for _, r := range rn {
			res = append(res, r.Key)
		}

		assert.Equal(t, []string{"b", "c"}, res)
	}

	{
		options := types.DefaultRangeOptions()
		options.Reverse = true
		rn := ss.GetRangeByIndex(0, -2, options)

		res := make([]string, 0)
		for _, r := range rn {
			res = append(res, r.Key)
		}

		assert.Equal(t, []string{"g", "f", "e", "d", "c", "b"}, res)
	}

	{
		rn := ss.GetRangeByIndex(-5, 2, types.DefaultRangeOptions())

		res := make([]string, 0)
		for _, r := range rn {
			res = append(res, r.Key)
		}

		assert.Equal(t, []string{"c"}, res)
	}

	{
		rn := ss.GetRangeByIndex(0, -2, types.DefaultRangeOptions())

		res := make([]string, 0)
		for _, r := range rn {
			res = append(res, r.Key)
		}

		assert.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, res)
	}

	{
		rn := ss.GetRangeByIndex(5, -1, types.DefaultRangeOptions())

		res := make([]string, 0)
		for _, r := range rn {
			res = append(res, r.Key)
		}

		assert.Equal(t, []string{"f", "g"}, res)
	}

	{
		rn := ss.GetRangeByIndex(0, -5, types.DefaultRangeOptions())

		res := make([]string, 0)
		for _, r := range rn {
			res = append(res, r.Key)
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
		assert.Equal(t, float64(i+1), res.Score)
	}

	for i := -1; i >= -ss.Len(); i-- {
		res := ss.GetByIndex(i, false)
		assert.Equal(t, float64(ss.Len()+i+1), res.Score)
	}
}

func TestSortedSet2(t *testing.T) {
	ss := types.NewSortedSet()

	ss.AddOrUpdate("a", 1)
	ss.AddOrUpdate("b", 2)
	ss.AddOrUpdate("c", 3)
	ss.AddOrUpdate("d", 4)

	{
		rn := ss.GetRangeByIndex(-5, 2, types.DefaultRangeOptions())
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"a", "b", "c"}, res)
	}
	{
		options := types.DefaultRangeOptions()
		options.Reverse = true
		rn := ss.GetRangeByIndex(1, 5, options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"c", "b", "a"}, res)
	}
	{
		a := ss.FindRankOfKey("A")
		assert.Equal(t, 0, a)
		e := ss.FindRankOfKey("e")
		assert.Equal(t, 0, e)
	}
}

func TestSortedSet3(t *testing.T) {
	ss := types.NewSortedSet()

	ss.AddOrUpdate("a", math.Inf(-1))
	ss.AddOrUpdate("b", 1)
	ss.AddOrUpdate("c", 2)
	ss.AddOrUpdate("d", 3)
	ss.AddOrUpdate("e", 4)
	ss.AddOrUpdate("f", 5)
	ss.AddOrUpdate("g", math.Inf(1))

	{
		options := types.DefaultRangeOptions()
		rn := ss.GetRangeByScore(math.Inf(-1), 2, options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"a", "b", "c"}, res)
	}
	{
		options := types.DefaultRangeOptions()
		rn := ss.GetRangeByScore(0, 3, options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"b", "c", "d"}, res)
	}
	{
		options := types.DefaultRangeOptions()
		rn := ss.GetRangeByScore(3, 6, options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"d", "e", "f"}, res)
	}
	{
		options := types.DefaultRangeOptions()
		options.Reverse = true
		rn := ss.GetRangeByScore(2, math.Inf(-1), options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"c", "b", "a"}, res)
	}
	{
		options := types.DefaultRangeOptions()
		options.Offset = 2
		options.Limit = 3
		rn := ss.GetRangeByScore(0, 10, options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"d", "e", "f"}, res)
	}
	{
		options := types.DefaultRangeOptions()
		options.Offset = 2
		options.Limit = 10
		rn := ss.GetRangeByScore(0, 10, options)
		assert.Equal(t, 3, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"d", "e", "f"}, res)
	}
	{
		options := types.DefaultRangeOptions()
		options.Reverse = true
		options.Offset = 0
		options.Limit = 2
		rn := ss.GetRangeByScore(10, 0, options)
		assert.Equal(t, 2, len(rn))

		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"f", "e"}, res)
	}
}

func TestSortedSet4(t *testing.T) {
	ss := types.NewSortedSet()

	ss.AddOrUpdate("b", 1)
	ss.AddOrUpdate("c", 2)
	ss.AddOrUpdate("d", 3)
	ss.AddOrUpdate("e", 4)
	ss.AddOrUpdate("f", 5)

	{
		options := types.DefaultRangeOptions()
		rn := ss.GetRangeByScore(6, math.Inf(1), options)
		assert.Equal(t, 0, len(rn))
	}

	{
		options := types.DefaultRangeOptions()
		options.StopExclusive = true
		rn := ss.GetRangeByScore(2, 2, options)
		assert.Equal(t, 0, len(rn))
	}
}

// ZADD key 0 alpha 0 bar 0 cool 0 down 0 elephant 0 foo 0 great 0 hill 0 omega
func TestSortedSet5(t *testing.T) {
	ss := types.NewSortedSet()

	ss.AddOrUpdate("alpha", 0)
	ss.AddOrUpdate("bar", 0)
	ss.AddOrUpdate("cool", 0)
	ss.AddOrUpdate("down", 0)
	ss.AddOrUpdate("elephant", 0)
	ss.AddOrUpdate("foo", 0)
	ss.AddOrUpdate("great", 0)
	ss.AddOrUpdate("hill", 0)
	ss.AddOrUpdate("omega", 0)

	{
		options := types.DefaultRangeOptions()
		rn := ss.GetRangeByLex("-", "cool", options)
		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"alpha", "bar", "cool"}, res)
	}
	{
		options := types.DefaultRangeOptions()
		rn := ss.GetRangeByLex("g", "+", options)
		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"great", "hill", "omega"}, res)
	}
	{
		options := types.DefaultRangeOptions()
		options.Reverse = true
		options.StopExclusive = true
		rn := ss.GetRangeByLex("+", "d", options)
		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"omega", "hill", "great", "foo", "elephant", "down"}, res)
	}
	{
		options := types.DefaultRangeOptions()
		rn := ss.GetRangeByLex("ele", "h", options)
		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"elephant", "foo", "great"}, res)
	}
	{
		options := types.DefaultRangeOptions()
		options.Reverse = true
		options.StartExclusive = true
		rn := ss.GetRangeByLex("cool", "-", options)
		res := make([]string, 0, len(rn))
		for _, r := range rn {
			res = append(res, r.Key)
		}
		assert.Equal(t, []string{"bar", "alpha"}, res)

	}
}
