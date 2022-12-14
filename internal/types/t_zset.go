package types

import (
	"math"
)

var _ Item = (*ZSet)(nil)

type SerdeZSet struct {
	Keys   []string  `json:"keys"`
	Scores []float64 `json:"scores"`
}

type ZSet struct {
	Inner *SortedSet
}

func NewZSet() *ZSet {
	return &ZSet{Inner: NewSortedSet()}
}

func NewZSetFromSs(value *SortedSet) *ZSet {
	return &ZSet{Inner: value}
}

/// impl Item for ZSet

func (s *ZSet) Value() interface{} {
	return s.Inner
}

func (l ZSet) Type() uint64 {
	return ValueTypeZSet
}

func (l ZSet) TypeFancy() string {
	return ValueTypeFancyZSet
}

func (s ZSet) Len() int {
	return s.Inner.Len()
}

/// Impl ZSet

// Converts a ZSet to a Set.
func (s *ZSet) ToSet() *Set {
	set := NewSetEmpty()

	for key := range s.Inner.Dict {
		set.AddMember(key)
	}

	return set
}

// Union returns a new ZSet that is a union of both sets.
func (s *ZSet) Union(o *ZSet, mode int, weight float64) *ZSet {
	set := NewZSet()

	for key, node := range s.Inner.Dict {
		otherNode := o.Inner.Dict[key]
		if otherNode != nil {
			if mode == 1 { // Min
				set.Inner.AddOrUpdate(key, math.Min(node.Score, otherNode.Score*weight))
			} else if mode == 2 { // Max
				set.Inner.AddOrUpdate(key, math.Max(node.Score, otherNode.Score*weight))
			} else { // NOTE: mode _should_ be 0 here, but this might not always be correct :)
				if (math.IsInf(node.Score, -1) && math.IsInf(otherNode.Score, 1)) ||
					math.IsInf(otherNode.Score, -1) && math.IsInf(node.Score, 1) {
					set.Inner.AddOrUpdate(key, 0)
				} else {
					set.Inner.AddOrUpdate(key, node.Score+(otherNode.Score*weight))
				}
			}
		} else {
			set.Inner.AddOrUpdate(key, node.Score)
		}
	}

	for key, node := range o.Inner.Dict {
		_, exists := set.Inner.Dict[key]
		if !exists {
			set.Inner.AddOrUpdate(key, node.Score*weight)
		}
	}

	return set
}

// Intersect returns a new Set that is a intersection of both sets.
func (s *ZSet) Intersect(o *ZSet, mode int, weight float64) *ZSet {
	set := NewZSet()

	for key, node := range s.Inner.Dict {
		otherNode := o.Inner.Dict[key]
		if otherNode != nil {
			if mode == 1 { // Min
				set.Inner.AddOrUpdate(key, math.Min(node.Score, otherNode.Score*weight))
			} else if mode == 2 { // Max
				set.Inner.AddOrUpdate(key, math.Max(node.Score, otherNode.Score*weight))
			} else { // NOTE: It _should_ be 0 here, but this might not always be correct :)
				if (math.IsInf(node.Score, -1) && math.IsInf(otherNode.Score, 1)) ||
					math.IsInf(otherNode.Score, -1) && math.IsInf(node.Score, 1) {
					set.Inner.AddOrUpdate(key, 0)
				} else {
					set.Inner.AddOrUpdate(key, node.Score+(otherNode.Score*weight))
				}
			}
		}
	}

	return set
}

// Diff returns a new Set that is a union of both sets.
func (s *ZSet) Diff(o *ZSet) *ZSet {
	set := NewZSet()

	for key, node := range s.Inner.Dict {
		otherNode := o.Inner.Dict[key]
		if otherNode == nil {
			set.Inner.AddOrUpdate(key, node.Score)
		}
	}

	return set
}
