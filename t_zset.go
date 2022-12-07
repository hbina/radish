package redis

import (
	"math"
)

var _ Item = (*ZSet)(nil)

type SerdeZSet struct {
	Keys   []string  `json:"keys"`
	Scores []float64 `json:"scores"`
}

type ZSet struct {
	inner *SortedSet
}

func NewZSet() *ZSet {
	return &ZSet{inner: NewSortedSet()}
}

func NewZSetFromSs(value *SortedSet) *ZSet {
	return &ZSet{inner: value}
}

/// impl Item for ZSet

func (s *ZSet) Value() interface{} {
	return s.inner
}

func (l ZSet) Type() uint64 {
	return ValueTypeZSet
}

func (l ZSet) TypeFancy() string {
	return ValueTypeFancyZSet
}

func (s ZSet) Len() int {
	return s.inner.Len()
}

/// Impl ZSet

// Converts a ZSet to a Set.
func (s *ZSet) ToSet() *Set {
	set := NewSetEmpty()

	for key := range s.inner.dict {
		set.AddMember(key)
	}

	return set
}

// Union returns a new ZSet that is a union of both sets.
func (s *ZSet) Union(o *ZSet, mode int, weight float64) *ZSet {
	set := NewZSet()

	for key, node := range s.inner.dict {
		otherNode := o.inner.dict[key]
		if otherNode != nil {
			if mode == 1 { // Min
				set.inner.AddOrUpdate(key, math.Min(node.score, otherNode.score*weight))
			} else if mode == 2 { // Max
				set.inner.AddOrUpdate(key, math.Max(node.score, otherNode.score*weight))
			} else { // NOTE: mode _should_ be 0 here, but this might not always be correct :)
				if (math.IsInf(node.score, -1) && math.IsInf(otherNode.score, 1)) ||
					math.IsInf(otherNode.score, -1) && math.IsInf(node.score, 1) {
					set.inner.AddOrUpdate(key, 0)
				} else {
					set.inner.AddOrUpdate(key, node.score+(otherNode.score*weight))
				}
			}
		} else {
			set.inner.AddOrUpdate(key, node.score)
		}
	}

	for key, node := range o.inner.dict {
		_, exists := set.inner.dict[key]
		if !exists {
			set.inner.AddOrUpdate(key, node.score*weight)
		}
	}

	return set
}

// Intersect returns a new Set that is a intersection of both sets.
func (s *ZSet) Intersect(o *ZSet, mode int, weight float64) *ZSet {
	set := NewZSet()

	for key, node := range s.inner.dict {
		otherNode := o.inner.dict[key]
		if otherNode != nil {
			if mode == 1 { // Min
				set.inner.AddOrUpdate(key, math.Min(node.score, otherNode.score*weight))
			} else if mode == 2 { // Max
				set.inner.AddOrUpdate(key, math.Max(node.score, otherNode.score*weight))
			} else { // NOTE: It _should_ be 0 here, but this might not always be correct :)
				if (math.IsInf(node.score, -1) && math.IsInf(otherNode.score, 1)) ||
					math.IsInf(otherNode.score, -1) && math.IsInf(node.score, 1) {
					set.inner.AddOrUpdate(key, 0)
				} else {
					set.inner.AddOrUpdate(key, node.score+(otherNode.score*weight))
				}
			}
		}
	}

	return set
}

// Diff returns a new Set that is a union of both sets.
func (s *ZSet) Diff(o *ZSet) *ZSet {
	set := NewZSet()

	for key, node := range s.inner.dict {
		otherNode := o.inner.dict[key]
		if otherNode == nil {
			set.inner.AddOrUpdate(key, node.score)
		}
	}

	return set
}
