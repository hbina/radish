// Copyright (c) 2016, Jerry.Wang
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, ss
//  list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//  ss list of conditions and the following disclaimer in the documentation
//  and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

// Plus some modifications from the original package.

package redis

import (
	"math/rand"

	"golang.org/x/exp/constraints"
)

type Score int // the type of score

const SKIPLIST_MAXLEVEL = 32 /* Should be enough for 2^64 elements */
const SKIPLIST_P = 0.25      /* Skiplist P = 1/4 */

type SortedSetLevel[K constraints.Ordered, S constraints.Ordered, V any] struct {
	// Pointer to the next node
	forward *SortedSetNode[K, S, V]
	// The number of node we skip if we go to the forward node
	span int
}

// Node in skip list
type SortedSetNode[K constraints.Ordered, S constraints.Ordered, V any] struct {
	key      K // unique key of this node
	value    V // associated data
	score    S // score to determine the order of this node in the set
	backward *SortedSetNode[K, S, V]
	level    []SortedSetLevel[K, S, V]
}

// Get the key of the node
func (ss *SortedSetNode[K, S, V]) Key() K {
	return ss.key
}

// Get the score of the node
func (ss *SortedSetNode[K, S, V]) Score() S {
	return ss.score
}

// Get the score of the node
func (ss *SortedSetNode[K, S, V]) Value() V {
	return ss.value
}

type SortedSet[K constraints.Ordered, S constraints.Ordered, V any] struct {
	header *SortedSetNode[K, S, V]
	tail   *SortedSetNode[K, S, V]
	length int
	level  int
	dict   map[K]*SortedSetNode[K, S, V]
}

func createNode[K constraints.Ordered, S constraints.Ordered, V any](level int, key K, score S, value V) *SortedSetNode[K, S, V] {
	node := SortedSetNode[K, S, V]{
		score: score,
		key:   key,
		value: value,
		level: make([]SortedSetLevel[K, S, V], level),
	}
	return &node
}

// randomLevel returns the level that this node should go up to.
func randomLevel() int {
	level := 1

	for rand.Intn(1/SKIPLIST_P) == 1 {
		level += 1

		if SKIPLIST_MAXLEVEL == level {
			break
		}
	}

	return level
}

// insertNode inserts a new node with the given score, key and value
func (ss *SortedSet[K, S, V]) insertNode(score S, key K, value V) *SortedSetNode[K, S, V] {
	var update [SKIPLIST_MAXLEVEL]*SortedSetNode[K, S, V]
	var rank [SKIPLIST_MAXLEVEL]int

	x := ss.header

	// We start from the top levels to the bottom
	for i := ss.level - 1; i >= 0; i-- {
		// Store rank that is crossed to reach the insert position
		if ss.level-1 == i {
			// The highest level rank is 0
			rank[i] = 0
		} else {
			// The rank at this level continues from the previous level
			rank[i] = rank[i+1]
		}

		// We are moving forward x as far as we can in our current level
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score && // score is the same but the key is different
					x.level[i].forward.key < key)) {
			rank[i] += x.level[i].span
			x = x.level[i].forward
		}

		// This is the furthest node that the current score can go at this level
		update[i] = x
	}

	level := randomLevel()

	// The random level for this node is higher than our set level
	// So we must update rank and update to match this
	if level > ss.level {
		for i := ss.level; i < level; i++ {
			rank[i] = 0
			update[i] = ss.header
			update[i].level[i].span = ss.length
		}
		ss.level = level
	}

	x = createNode(level, key, score, value)

	// Update span and forward metadata of current nodes and nodes to be updated
	for i := 0; i < level; i++ {
		// Rewire the pointers around
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x

		// Divide the span that was originally covered by the previous node from:
		// previous_node <----> next_node
		// to:
		// previous_node <----> new_node <------> next_node
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	// Increment span for the untouched levels because we insterted a new node
	// right in front of us
	for i := level; i < ss.level; i++ {
		update[i].level[i].span++
	}

	if update[0] == ss.header {
		// This is a special when the new node must be put exactly in front of the header
		x.backward = nil
	} else {
		// Otherwise, backward is just the previous node which can be always be obtained from update[0]
		x.backward = update[0]
	}

	if x.level[0].forward != nil {
		// If there is a node in front of the current node, then that node's backward must point to us
		x.level[0].forward.backward = x
	} else {
		// Otherwise, we are the furthermost node at level 0, we are the new tail
		ss.tail = x
	}
	ss.length++
	return x
}

/* Internal function used by delete, DeleteByScore and DeleteByRank */
func (ss *SortedSet[K, S, V]) deleteNode(x *SortedSetNode[K, S, V], update [SKIPLIST_MAXLEVEL]*SortedSetNode[K, S, V]) {
	for i := 0; i < ss.level; i++ {
		if update[i].level[i].forward == x {
			update[i].level[i].span += x.level[i].span - 1
			update[i].level[i].forward = x.level[i].forward
		} else {
			update[i].level[i].span -= 1
		}
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x.backward
	} else {
		ss.tail = x.backward
	}
	for ss.level > 1 && ss.header.level[ss.level-1].forward == nil {
		ss.level--
	}
	ss.length--
	delete(ss.dict, x.key)
}

/* Delete an element with matching score/key from the skiplist. */
func (ss *SortedSet[K, S, V]) delete(score S, key K) bool {
	var update [SKIPLIST_MAXLEVEL]*SortedSetNode[K, S, V]

	x := ss.header
	for i := ss.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score &&
					x.level[i].forward.key < key)) {
			x = x.level[i].forward
		}
		update[i] = x
	}
	/* We may have multiple elements with the same score, what we need
	 * is to find the element with both the right score and object. */
	x = x.level[0].forward
	if x != nil && score == x.score && x.key == key {
		ss.deleteNode(x, update)
		return true
	}
	return false /* not found */
}

// NewSortedSet returns a new empty sorted set
func NewSortedSet[K constraints.Ordered, S constraints.Ordered, V any]() *SortedSet[K, S, V] {
	sortedSet := SortedSet[K, S, V]{
		level: 1,
		dict:  make(map[K]*SortedSetNode[K, S, V]),
	}

	var key K
	var score S
	var value V
	sortedSet.header = createNode(SKIPLIST_MAXLEVEL, key, score, value)
	return &sortedSet
}

// Get the number of elements
func (ss *SortedSet[K, S, V]) Len() int {
	return ss.length
}

// PeekMin returns the element with the lowest score if it exists.
// Otherwise it returns nil.
//
// Time complexity: O(1)
func (ss *SortedSet[K, S, V]) PeekMin() *SortedSetNode[K, S, V] {
	return ss.header.level[0].forward
}

// PopMin returns the element with the lowest score if it exists and
// removes it.
// Otherwise it returns nil.
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet[K, S, V]) PopMin() *SortedSetNode[K, S, V] {
	x := ss.header.level[0].forward
	if x != nil {
		ss.Remove(x.key)
	}
	return x
}

// PeekMax returns the element with the highest score if it exists.
//
// Time Complexity : O(1)
func (ss *SortedSet[K, S, V]) PeekMax() *SortedSetNode[K, S, V] {
	return ss.tail
}

// PopMin returns the element with the highest score if it exists and
// removes it.
// Otherwise it returns nil.
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet[K, S, V]) PopMax() *SortedSetNode[K, S, V] {
	x := ss.tail
	if x != nil {
		ss.Remove(x.key)
	}
	return x
}

// Add an element into the sorted set with specific key / value / score.
// If the element is added, this method returns true; otherwise false means updated.
//
// Time complexity: O(log(N))
func (ss *SortedSet[K, S, V]) AddOrUpdate(key K, score S, value V) bool {
	var newNode *SortedSetNode[K, S, V] = nil

	found := ss.dict[key]
	if found != nil {
		// score does not change, only update value
		if found.score == score {
			found.value = value
		} else { // score changes, delete and re-insert
			ss.delete(found.score, found.key)
			newNode = ss.insertNode(score, key, value)
		}
	} else {
		newNode = ss.insertNode(score, key, value)
	}

	if newNode != nil {
		ss.dict[key] = newNode
	}
	return found == nil
}

// Delete element specified by key
//
// Time complexity: O(log(N))
func (ss *SortedSet[K, S, V]) Remove(key K) *SortedSetNode[K, S, V] {
	found := ss.dict[key]
	if found != nil {
		ss.delete(found.score, found.key)
		return found
	}
	return nil
}

// TODO: Add reverse, offset
type GetByScoreRangeOptions struct {
	Limit        int  // limit the max nodes to return
	ExcludeStart bool // exclude start value, so it search in interval (start, end] or (start, end)
	ExcludeEnd   bool // exclude end value, so it search in interval [start, end) or (start, end)
}

// Get the nodes whose score within the specific range
//
// If options is nil, it searchs in interval [start, end] without any limit by default
//
// Time complexity: O(log(N))
func (ss *SortedSet[K, S, V]) GetRangeByScore(start S, end S, options *GetByScoreRangeOptions) []*SortedSetNode[K, S, V] {

	// prepare parameters
	var limit int = int((^uint(0)) >> 1)
	if options != nil && options.Limit > 0 {
		limit = options.Limit
	}

	excludeStart := options != nil && options.ExcludeStart
	excludeEnd := options != nil && options.ExcludeEnd
	reverse := start > end

	if reverse {
		start, end = end, start
		excludeStart, excludeEnd = excludeEnd, excludeStart
	}

	var nodes []*SortedSetNode[K, S, V]

	//determine if out of range
	if ss.length == 0 {
		return nodes
	}

	if reverse { // search from end to start
		x := ss.header

		if excludeEnd {
			for i := ss.level - 1; i >= 0; i-- {
				for x.level[i].forward != nil &&
					x.level[i].forward.score < end {
					x = x.level[i].forward
				}
			}
		} else {
			for i := ss.level - 1; i >= 0; i-- {
				for x.level[i].forward != nil &&
					x.level[i].forward.score <= end {
					x = x.level[i].forward
				}
			}
		}

		for x != nil && limit > 0 {
			if excludeStart {
				if x.score <= start {
					break
				}
			} else {
				if x.score < start {
					break
				}
			}

			next := x.backward

			nodes = append(nodes, x)
			limit--

			x = next
		}
	} else {
		// search from start to end
		x := ss.header

		if excludeStart {
			for i := ss.level - 1; i >= 0; i-- {
				for x.level[i].forward != nil &&
					x.level[i].forward.score <= start {
					x = x.level[i].forward
				}
			}
		} else {
			for i := ss.level - 1; i >= 0; i-- {
				for x.level[i].forward != nil &&
					x.level[i].forward.score < start {
					x = x.level[i].forward
				}
			}
		}

		/* Current node is the last with score < or <= start. */
		x = x.level[0].forward

		for x != nil && limit > 0 {
			if excludeEnd {
				if x.score >= end {
					break
				}
			} else {
				if x.score > end {
					break
				}
			}

			next := x.level[0].forward

			nodes = append(nodes, x)
			limit--

			x = next
		}
	}

	return nodes
}

// ConvertIndexToRank sanitizes the given rank-based range.
// Returns (1,-1) if its an empty range.
func (ss *SortedSet[K, S, V]) ConvertIndexToRank(start int, end int, reverse bool) (int, int) {
	if start < 0 {
		start += ss.Len()
	}
	if end < 0 {
		end += ss.Len()
	}
	if start < 0 {
		start = 0
	}
	if start > end || start >= ss.Len() {
		return 1, -1
	}
	start++
	end++
	if reverse {
		start, end = ss.Len()-end+1, ss.Len()-start+1
	}
	return start, end
}

// findNodeByIter returns the node at the given rank (which is 1-based index).
// Can also remove the node is asked.
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet[K, S, V]) findNodeByRank(start int, remove bool) (traversed int, x *SortedSetNode[K, S, V], update [SKIPLIST_MAXLEVEL]*SortedSetNode[K, S, V]) {
	x = ss.header
	for i := ss.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			traversed+x.level[i].span < start {
			traversed += x.level[i].span
			x = x.level[i].forward
		}
		if remove {
			update[i] = x
		} else {
			if traversed+1 == start {
				break
			}
		}
	}
	return
}

// GetRangeByIndex returns array of nodes within specific index range [start, end].
// The given start and end must be a valid rank-based index which can be obtained from 'SanitizeRank'.
//
// If start is greater than end, the returned array is in reserved order.
// If remove is true, the returned nodes are removed
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet[K, S, V]) GetRangeByIndex(start int, end int, reverse bool, remove bool) []*SortedSetNode[K, S, V] {
	start, end = ss.ConvertIndexToRank(start, end, reverse)

	if start > end {
		return []*SortedSetNode[K, S, V]{}
	}

	var nodes []*SortedSetNode[K, S, V]

	traversed, x, update := ss.findNodeByRank(start, remove)
	traversed++

	x = x.level[0].forward
	for x != nil && traversed <= end {
		next := x.level[0].forward

		nodes = append(nodes, x)

		if remove {
			ss.deleteNode(x, update)
		}

		traversed++
		x = next
	}

	if reverse {
		for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
			nodes[i], nodes[j] = nodes[j], nodes[i]
		}
	}

	return nodes
}

// Get node by index.
//
// If remove is true, the returned nodes are removed
// If node is not found at specific rank, nil is returned
//
// Time complexity: O(log(N))
func (ss *SortedSet[K, S, V]) GetByIndex(rank int, remove bool) *SortedSetNode[K, S, V] {
	nodes := ss.GetRangeByIndex(rank, rank, false, remove)
	if len(nodes) == 1 {
		return nodes[0]
	}
	return nil
}

// Get node by key
//
// If node is not found, nil is returned
// Time complexity: O(1)
func (ss *SortedSet[K, S, V]) GetByKey(key K) *SortedSetNode[K, S, V] {
	return ss.dict[key]
}

// Find the rank of the node specified by key
// Note that the rank is 1-based integer. Rank 1 means the first node
//
// If the node is not found, 0 is returned. Otherwise rank(> 0) is returned
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet[K, S, V]) FindRankOfKey(key K) int {
	var rank int = 0
	node := ss.dict[key]
	if node != nil {
		x := ss.header
		for i := ss.level - 1; i >= 0; i-- {
			for x.level[i].forward != nil &&
				(x.level[i].forward.score < node.score ||
					(x.level[i].forward.score == node.score &&
						x.level[i].forward.key <= node.key)) {
				rank += x.level[i].span
				x = x.level[i].forward
			}

			if x.key == key {
				return rank
			}
		}
	}
	return 0
}
