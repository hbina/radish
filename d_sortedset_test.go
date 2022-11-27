package redis

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkOrder(t *testing.T, nodes []*SortedSetNode[string, int64, string], expectedOrder []string) {
	if len(expectedOrder) != len(nodes) {
		t.Errorf("nodes does not contain %d elements", len(expectedOrder))
	}
	for i := 0; i < len(expectedOrder); i++ {
		if nodes[i].Key() != expectedOrder[i] {
			t.Errorf("nodes[%d] is %q, but the expected key is %q", i, nodes[i].Key(), expectedOrder[i])
		}

	}
}

func checkIterByRankRange(t *testing.T, ss *SortedSet[string, int64, string], start int, end int, expectedOrder []string) {
	var keys []string

	// check nil callback should do nothing
	ss.IterRangeByRank(start, end, false, nil)

	ss.IterRangeByRank(start, end, false, func(key string, _ string) bool {
		keys = append(keys, key)
		return true
	})
	if len(expectedOrder) != len(keys) {
		t.Errorf("keys does not contain %d elements", len(expectedOrder))
	}
	for i := 0; i < len(expectedOrder); i++ {
		if keys[i] != expectedOrder[i] {
			t.Errorf("keys[%d] is %q, but the expected key is %q", i, keys[i], expectedOrder[i])
		}
	}

	// check return early
	if len(expectedOrder) < 1 {
		return
	}
	// reset data
	keys = []string{}
	var i int
	ss.IterRangeByRank(start, end, false, func(key string, _ string) bool {
		keys = append(keys, key)
		i++
		// return early
		return i < len(expectedOrder)-1
	})
	if len(expectedOrder)-1 != len(keys) {
		t.Errorf("keys does not contain %d elements", len(expectedOrder)-1)
	}
	for i := 0; i < len(expectedOrder)-1; i++ {
		if keys[i] != expectedOrder[i] {
			t.Errorf("keys[%d] is %q, but the expected key is %q", i, keys[i], expectedOrder[i])
		}
	}

}

func checkRankRangeIterAndOrder(t *testing.T, sortedset *SortedSet[string, int64, string], start int, end int, remove bool, expectedOrder []string) {
	checkIterByRankRange(t, sortedset, start, end, expectedOrder)
	nodes := sortedset.GetRangeByRank(start, end, false, remove)
	checkOrder(t, nodes, expectedOrder)
}

func TestCase1(t *testing.T) {
	sortedset := NewSortedSet[string, int64, string]()

	sortedset.AddOrUpdate("a", 89, "Kelly")
	sortedset.AddOrUpdate("b", 100, "Staley")
	sortedset.AddOrUpdate("c", 100, "Jordon")
	sortedset.AddOrUpdate("d", -321, "Park")
	sortedset.AddOrUpdate("e", 101, "Albert")
	sortedset.AddOrUpdate("f", 99, "Lyman")
	sortedset.AddOrUpdate("g", 99, "Singleton")
	sortedset.AddOrUpdate("h", 70, "Audrey")

	sortedset.AddOrUpdate("e", 99, "ntrnrt")

	sortedset.Remove("b")

	node := sortedset.GetByRank(3, false)
	if node == nil || node.Key() != "a" {
		t.Error("GetByRank() does not return expected value `a`")
	}

	node = sortedset.GetByRank(-3, false)
	if node == nil || node.Key() != "f" {
		t.Error("GetByRank() does not return expected value `f`")
	}

	// get all nodes since the first one to last one
	checkRankRangeIterAndOrder(t, sortedset, 1, -1, false, []string{"d", "h", "a", "e", "f", "g", "c"})

	// get & remove the 2nd/3rd nodes in reserve order
	checkRankRangeIterAndOrder(t, sortedset, -2, -3, true, []string{"g", "f"})

	// get all nodes since the last one to first one
	checkRankRangeIterAndOrder(t, sortedset, -1, 1, false, []string{"c", "e", "a", "h", "d"})

}

func TestCase2(t *testing.T) {

	// create a new set
	sortedset := NewSortedSet[string, int64, string]()

	// fill in new node
	sortedset.AddOrUpdate("a", 89, "Kelly")
	sortedset.AddOrUpdate("b", 100, "Staley")
	sortedset.AddOrUpdate("c", 100, "Jordon")
	sortedset.AddOrUpdate("d", -321, "Park")
	sortedset.AddOrUpdate("e", 101, "Albert")
	sortedset.AddOrUpdate("f", 99, "Lyman")
	sortedset.AddOrUpdate("g", 99, "Singleton")
	sortedset.AddOrUpdate("h", 70, "Audrey")

	// update an existing node
	sortedset.AddOrUpdate("e", 99, "ntrnrt")

	// remove node
	sortedset.Remove("b")

	nodes := sortedset.GetRangeByScore(-500, 500, nil)
	checkOrder(t, nodes, []string{"d", "h", "a", "e", "f", "g", "c"})

	nodes = sortedset.GetRangeByScore(500, -500, nil)
	//t.Logf("%v", nodes)
	checkOrder(t, nodes, []string{"c", "g", "f", "e", "a", "h", "d"})

	nodes = sortedset.GetRangeByScore(600, 500, nil)
	checkOrder(t, nodes, []string{})

	nodes = sortedset.GetRangeByScore(500, 600, nil)
	checkOrder(t, nodes, []string{})

	rank := sortedset.FindRankOfKey("f")
	if rank != 5 {
		t.Error("FindRank() does not return expected value `5`")
	}

	rank = sortedset.FindRankOfKey("d")
	if rank != 1 {
		t.Error("FindRank() does not return expected value `1`")
	}

	nodes = sortedset.GetRangeByScore(99, 100, nil)
	checkOrder(t, nodes, []string{"e", "f", "g", "c"})

	nodes = sortedset.GetRangeByScore(90, 50, nil)
	checkOrder(t, nodes, []string{"a", "h"})

	nodes = sortedset.GetRangeByScore(99, 100, &GetByScoreRangeOptions{
		ExcludeStart: true,
	})
	checkOrder(t, nodes, []string{"c"})

	nodes = sortedset.GetRangeByScore(100, 99, &GetByScoreRangeOptions{
		ExcludeStart: true,
	})
	checkOrder(t, nodes, []string{"g", "f", "e"})

	nodes = sortedset.GetRangeByScore(99, 100, &GetByScoreRangeOptions{
		ExcludeEnd: true,
	})
	checkOrder(t, nodes, []string{"e", "f", "g"})

	nodes = sortedset.GetRangeByScore(100, 99, &GetByScoreRangeOptions{
		ExcludeEnd: true,
	})
	checkOrder(t, nodes, []string{"c"})

	nodes = sortedset.GetRangeByScore(50, 100, &GetByScoreRangeOptions{
		Limit: 2,
	})
	checkOrder(t, nodes, []string{"h", "a"})

	nodes = sortedset.GetRangeByScore(100, 50, &GetByScoreRangeOptions{
		Limit: 2,
	})
	checkOrder(t, nodes, []string{"c", "g"})

	minNode := sortedset.PeekMin()
	if minNode == nil || minNode.Key() != "d" {
		t.Error("PeekMin() does not return expected value `d`")
	}

	minNode = sortedset.PopMin()
	if minNode == nil || minNode.Key() != "d" {
		t.Error("PopMin() does not return expected value `d`")
	}

	nodes = sortedset.GetRangeByScore(-500, 500, nil)
	checkOrder(t, nodes, []string{"h", "a", "e", "f", "g", "c"})

	maxNode := sortedset.PeekMax()
	if maxNode == nil || maxNode.Key() != "c" {
		t.Error("PeekMax() does not return expected value `c`")
	}

	maxNode = sortedset.PopMax()
	if maxNode == nil || maxNode.Key() != "c" {
		t.Error("PopMax() does not return expected value `c`")
	}

	nodes = sortedset.GetRangeByScore(500, -500, nil)
	checkOrder(t, nodes, []string{"g", "f", "e", "a", "h"})
}

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

	r := rn[0]
	assert.Equal(t, "d", r.key)
	assert.Equal(t, 4, r.score)
	assert.Equal(t, struct{}{}, r.value)
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
