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

package types

import (
	"math"
	"math/rand"
)

type Score int // the type of score

const SKIPLIST_MAXLEVEL = 32 /* Should be enough for 2^64 elements */
const SKIPLIST_P = 0.25      /* Skiplist P = 1/4 */

type SortedSetLevel struct {
	// Pointer to the next node
	forward *SortedSetNode
	// The number of node we skip if we go to the forward node
	span int
}

// Node in skip list
type SortedSetNode struct {
	Key      string  // unique key of this node
	Score    float64 // score to determine the order of this node in the set
	backward *SortedSetNode
	level    []SortedSetLevel
}

type SortedSet struct {
	header *SortedSetNode
	tail   *SortedSetNode
	level  int
	Dict   map[string]*SortedSetNode
}

func createNode(level int, key string, score float64) *SortedSetNode {
	node := SortedSetNode{
		Key:   key,
		Score: score,
		level: make([]SortedSetLevel, level),
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
func (ss *SortedSet) insertNode(score float64, key string) *SortedSetNode {
	var update [SKIPLIST_MAXLEVEL]*SortedSetNode
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
			(x.level[i].forward.Score < score ||
				(x.level[i].forward.Score == score && // score is the same but the key is different
					x.level[i].forward.Key < key)) {
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
			update[i].level[i].span = ss.Len()
		}
		ss.level = level
	}

	x = createNode(level, key, score)

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
	return x
}

/* Internal function used by delete, DeleteByScore and DeleteByRank */
func (ss *SortedSet) deleteNode(x *SortedSetNode, update [SKIPLIST_MAXLEVEL]*SortedSetNode) {
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
	delete(ss.Dict, x.Key)
}

/* Delete an element with matching score/key from the skiplist. */
func (ss *SortedSet) delete(score float64, key string) bool {
	var update [SKIPLIST_MAXLEVEL]*SortedSetNode

	x := ss.header
	for i := ss.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.Score < score ||
				(x.level[i].forward.Score == score &&
					x.level[i].forward.Key < key)) {
			x = x.level[i].forward
		}
		update[i] = x
	}

	/* We may have multiple elements with the same score, what we need
	 * is to find the element with both the right score and object. */
	x = x.level[0].forward
	if x != nil && score == x.Score && x.Key == key {
		ss.deleteNode(x, update)
		delete(ss.Dict, key)
		return true
	}
	return false /* not found */
}

// NewSortedSet returns a new empty sorted set
func NewSortedSet() *SortedSet {
	sortedSet := SortedSet{
		level: 1,
		Dict:  make(map[string]*SortedSetNode),
	}

	var key string
	var score float64
	sortedSet.header = createNode(SKIPLIST_MAXLEVEL, key, score)
	return &sortedSet
}

// Get the number of elements
func (ss *SortedSet) Len() int {
	return len(ss.Dict)
}

// PeekMin returns the element with the lowest score if it exists.
// Otherwise it returns nil.
//
// Time complexity: O(1)
func (ss *SortedSet) PeekMin() *SortedSetNode {
	return ss.header.level[0].forward
}

// PopMin returns the element with the lowest score if it exists and
// removes it.
// Otherwise it returns nil.
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet) PopMin() *SortedSetNode {
	x := ss.header.level[0].forward
	if x != nil {
		ss.Remove(x.Key)
	}
	return x
}

// PeekMax returns the element with the highest score if it exists.
//
// Time Complexity : O(1)
func (ss *SortedSet) PeekMax() *SortedSetNode {
	return ss.tail
}

// PopMin returns the element with the highest score if it exists and
// removes it.
// Otherwise it returns nil.
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet) PopMax() *SortedSetNode {
	x := ss.tail
	if x != nil {
		ss.Remove(x.Key)
	}
	return x
}

// Add an element into the sorted set with specific key / value / score.
// If the element is added, this method returns true; otherwise false means updated.
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet) AddOrUpdate(key string, score float64) bool {
	var newNode *SortedSetNode = nil

	found := ss.Dict[key]
	if found != nil {
		// score does not change, only update value
		if found.Score != score { // score changes, delete and re-insert
			ss.delete(found.Score, found.Key)
			newNode = ss.insertNode(score, key)
		}
	} else {
		newNode = ss.insertNode(score, key)
	}

	if newNode != nil {
		ss.Dict[key] = newNode
	}
	return found == nil
}

// Delete element specified by key
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet) Remove(key string) *SortedSetNode {
	// Check the dict first so we don't have to iterate the nodes
	found := ss.Dict[key]
	if found != nil {
		ss.delete(found.Score, found.Key)
		return found
	}
	return nil
}

// Delete element specified by rank
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet) RemoveByRank(rank int) *SortedSetNode {
	node, _ := ss.findNodeByRank(rank)
	if node != nil {
		ss.delete(node.Score, node.Key)
		return node
	}
	return nil
}

// TODO: Add reverse, offset
type GetRangeOptions struct {
	Reverse        bool // Start iterating from the back
	Offset         int  // How many nodes to skip
	Limit          int  // limit the max nodes to return
	StartExclusive bool // exclude start value, so it search in interval (start, end] or (start, end)
	StopExclusive  bool // exclude end value, so it search in interval [start, end) or (start, end)
}

func DefaultRangeOptions() GetRangeOptions {
	return GetRangeOptions{
		Reverse:        false,
		Offset:         0,
		Limit:          math.MaxInt,
		StartExclusive: false,
		StopExclusive:  false,
	}
}

// GetRangeByScore returns an array of nodes that satisfy the given score range.
//
// Time complexity: O(log(N))
func (ss *SortedSet) GetRangeByScore(start float64, end float64, options GetRangeOptions) []*SortedSetNode {
	if options.Reverse {
		start, end = end, start
	}

	startNode, startRank := ss.findNodeByScore(start, true)

	if startNode == nil {
		return []*SortedSetNode{}
	}

	if options.StartExclusive && startNode.Score == start {
		startRank += 1
	}

	endNode, endRank := ss.findNodeByScore(end, false)

	if options.StopExclusive && endNode.Score == end {
		endRank -= 1
	}

	nodes := ss.GetRangeByRank(startRank, endRank, options)

	return nodes
}

// GetRangeByLex returns an array of nodes that satisfy the given score range.
//
// Time complexity: O(log(N))
func (ss *SortedSet) GetRangeByLex(start string, end string, options GetRangeOptions) []*SortedSetNode {
	if options.Reverse {
		start, end = end, start
		options.StartExclusive, options.StopExclusive = options.StopExclusive, options.StartExclusive
	}

	startNode, startRank := ss.FindNodeByLex(start)

	if startNode == nil {
		return []*SortedSetNode{}
	}

	if options.StartExclusive && startNode.Key == start {
		startRank += 1
	}

	endNode, endRank := ss.FindNodeByLex(end)

	if (options.StopExclusive && endNode.Key == end) || endNode.Key > end {
		endRank -= 1
	}

	if start == "+" {
		startRank = ss.Len()
	} else if start == "-" {
		startRank = 1
	}

	if end == "+" {
		endRank = ss.Len()
	} else if end == "-" {
		endRank = 1
	}

	nodes := ss.GetRangeByRank(startRank, endRank, options)

	return nodes
}

// SanitizeIndex sanitizes the given 0-based range.
// Returns (1,-1) if its an empty range.
// TODO: This is such a hack. Reimplement so we don't even need this whole mess.
func (ss *SortedSet) SanitizeIndex(start int, end int, reverse bool) (int, int) {
	// If start is negative, calculate the absolute value
	if start < 0 {
		start += ss.Len()
	}

	// If end is negative, calculate the absolute value
	if end < 0 {
		end += ss.Len()
	}

	// If absolute index does not make sense, return 1,-1
	if start > end || start >= ss.Len() {
		return 1, -1
	}

	// Convert 0-based to 1-based index
	start++
	end++

	// Constraint index to minimum
	if start < 1 {
		start = 1
	}

	// Constraint index to max
	if end > ss.Len() {
		end = ss.Len()
	}

	// If reverse, calculate the opposite
	if reverse {
		start, end = ss.Len()-end+1, ss.Len()-start+1
	}

	return start, end
}

// findNodeByRank returns the node with the requested rank
//
// Time complexity: O(log(N)) with high probability.
func (ss *SortedSet) findNodeByRank(start int) (*SortedSetNode, int) {
	node := ss.header
	nodeRank := 0

	for i := ss.level - 1; i >= 0; i-- {
		for node.level[i].forward != nil &&
			nodeRank+node.level[i].span < start {
			nodeRank += node.level[i].span
			node = node.level[i].forward
		}
	}

	// Move forward once because startRank is the last node that succeeds <start
	if node != nil {
		nodeRank += node.level[0].span
		node = node.level[0].forward
	}

	return node, nodeRank
}

// findNodeByScore returns the node just before the requested score and its rank.
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet) findNodeByScore(score float64, forward bool) (*SortedSetNode, int) {
	node := ss.header
	nodeRank := 0

	for i := ss.level - 1; i >= 0; i-- {
		for node.level[i].forward != nil &&
			node.level[i].forward.Score < score {
			nodeRank += node.level[i].span
			node = node.level[i].forward
		}
	}

	// Move forward once because node is the last node that succeeds <score
	if node != nil && (forward || (node.level[0].forward != nil && node.level[0].forward.Score == score)) {
		nodeRank += node.level[0].span
		node = node.level[0].forward
	}

	return node, nodeRank
}

// FindNodeByLex returns the node with the requested key
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet) FindNodeByLex(key string) (*SortedSetNode, int) {
	node := ss.header
	nodeRank := 0

	for i := ss.level - 1; i >= 0; i-- {
		for node.level[i].forward != nil &&
			node.level[i].forward.Key < key {
			nodeRank += node.level[i].span
			node = node.level[i].forward
		}
	}

	// Move forward once because startRank is the last node that succeeds <start
	if node != nil {
		nodeRank += node.level[0].span
		node = node.level[0].forward
	}

	return node, nodeRank
}

// GetRangeByIndex returns array of nodes within specific index range [start, end].
// The given start and end must be a valid rank-based index which can be obtained from 'SanitizeRank'.
//
// If start is greater than end, the returned array is in reserved order.
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet) GetRangeByIndex(start int, end int, options GetRangeOptions) []*SortedSetNode {
	start, end = ss.SanitizeIndex(start, end, options.Reverse)
	return ss.GetRangeByRank(start, end, options)
}

// GetRangeByRank returns array of nodes within specific rank range [start, end].
// If start is greater than end, returns an empty array.
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet) GetRangeByRank(start int, end int, options GetRangeOptions) []*SortedSetNode {

	if start > end {
		return []*SortedSetNode{}
	}

	if options.Reverse {
		node, nodeRank := ss.findNodeByRank(end)

		nodes := make([]*SortedSetNode, 0)
		for node != nil && options.Limit != 0 && nodeRank >= start {
			if options.Offset == 0 {
				options.Limit--

				nodes = append(nodes, node)
				node = node.backward
			} else {
				options.Offset--

				node = node.backward
			}
			nodeRank--
		}

		return nodes
	} else {
		node, nodeRank := ss.findNodeByRank(start)

		nodes := make([]*SortedSetNode, 0)
		for node != nil && options.Limit != 0 && nodeRank <= end {
			if options.Offset == 0 {
				options.Limit--

				nodes = append(nodes, node)
				node = node.level[0].forward
			} else {
				options.Offset--

				node = node.level[0].forward
			}
			nodeRank++
		}

		return nodes
	}
}

// Get node by index.
//
// If remove is true, the returned nodes are removed
// If node is not found at specific rank, nil is returned
//
// Time complexity: O(log(N))
func (ss *SortedSet) GetByIndex(rank int, remove bool) *SortedSetNode {
	nodes := ss.GetRangeByIndex(rank, rank, DefaultRangeOptions())
	if len(nodes) == 1 {
		return nodes[0]
	}
	return nil
}

// Get node by key
//
// If node is not found, nil is returned
// Time complexity: O(1)
func (ss *SortedSet) GetByKey(key string) *SortedSetNode {
	return ss.Dict[key]
}

// Find the rank of the node specified by key
// Note that the rank is 1-based integer. Rank 1 means the first node
//
// If the node is not found, 0 is returned. Otherwise rank(>0) is returned
//
// Time complexity: O(log(N)) with high probability
func (ss *SortedSet) FindRankOfKey(key string) int {
	node, rank := ss.FindNodeByLex(key)
	if node == nil || node == ss.header || node.Key != key {
		return 0
	}
	return rank
}
