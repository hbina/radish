package types

import (
	"encoding/json"
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

	for key := range s.inner.Dict {
		set.AddMember(key)
	}

	return set
}

// Union returns a new ZSet that is a union of both sets.
func (s *ZSet) Union(o *ZSet, mode int, weight float64) *ZSet {
	set := NewZSet()

	for key, node := range s.inner.Dict {
		otherNode := o.inner.Dict[key]
		if otherNode != nil {
			if mode == 1 { // Min
				set.inner.AddOrUpdate(key, math.Min(node.Score, otherNode.Score*weight))
			} else if mode == 2 { // Max
				set.inner.AddOrUpdate(key, math.Max(node.Score, otherNode.Score*weight))
			} else { // NOTE: mode _should_ be 0 here, but this might not always be correct :)
				if (math.IsInf(node.Score, -1) && math.IsInf(otherNode.Score, 1)) ||
					math.IsInf(otherNode.Score, -1) && math.IsInf(node.Score, 1) {
					set.inner.AddOrUpdate(key, 0)
				} else {
					set.inner.AddOrUpdate(key, node.Score+(otherNode.Score*weight))
				}
			}
		} else {
			set.inner.AddOrUpdate(key, node.Score)
		}
	}

	for key, node := range o.inner.Dict {
		_, exists := set.inner.Dict[key]
		if !exists {
			if weight == 0 {
				set.inner.AddOrUpdate(key, 0)
			} else {
				set.inner.AddOrUpdate(key, node.Score*weight)
			}
		}
	}

	return set
}

// Intersect returns a new Set that is a intersection of both sets.
func (s *ZSet) Intersect(o *ZSet, mode int, weight float64) *ZSet {
	set := NewZSet()

	for key, node := range s.inner.Dict {
		otherNode := o.inner.Dict[key]
		if otherNode != nil {
			if mode == 1 { // Min
				set.inner.AddOrUpdate(key, math.Min(node.Score, otherNode.Score*weight))
			} else if mode == 2 { // Max
				set.inner.AddOrUpdate(key, math.Max(node.Score, otherNode.Score*weight))
			} else { // NOTE: It _should_ be 0 here, but this might not always be correct :)
				if (math.IsInf(node.Score, -1) && math.IsInf(otherNode.Score, 1)) ||
					math.IsInf(otherNode.Score, -1) && math.IsInf(node.Score, 1) {
					set.inner.AddOrUpdate(key, 0)
				} else {
					set.inner.AddOrUpdate(key, node.Score+(otherNode.Score*weight))
				}
			}
		}
	}

	return set
}

// Diff returns a new Set that is a union of both sets.
func (s *ZSet) Diff(o *ZSet) *ZSet {
	set := NewZSet()

	for key, node := range s.inner.Dict {
		otherNode := o.inner.Dict[key]
		if otherNode == nil {
			set.inner.AddOrUpdate(key, node.Score)
		}
	}

	return set
}

func (s *ZSet) Marshal() ([]byte, error) {
	keys := make([]string, 0, s.Len())
	scores := make([]float64, 0, s.Len())

	for _, z := range s.inner.Dict {
		keys = append(keys, z.Key)
		scores = append(scores, z.Score)
	}

	sz := SerdeZSet{
		Keys:   keys,
		Scores: scores,
	}

	str, err := json.Marshal(sz)
	return str, err
}

func ZSetUnmarshal(data []byte) (*ZSet, bool) {
	var sz SerdeZSet
	err := json.Unmarshal(data, &sz)

	if err != nil {
		return nil, false
	}

	ss := NewSortedSet()

	for idx := range sz.Keys {
		ss.AddOrUpdate(sz.Keys[idx], sz.Scores[idx])
	}

	return NewZSetFromSs(ss), true
}

func (ss *ZSet) GetByKey(key string) *SortedSetNode {
	return ss.inner.GetByKey(key)
}

func (ss *ZSet) AddOrUpdate(key string, score float64) bool {
	return ss.inner.AddOrUpdate(key, score)
}

func (ss *ZSet) FindNodeByLex(key string) (*SortedSetNode, int) {
	return ss.inner.FindNodeByLex(key)
}

func (ss *ZSet) Remove(key string) *SortedSetNode {
	return ss.inner.Remove(key)
}

func (ss *ZSet) GetRangeByRank(start int, end int, options GetRangeOptions) []*SortedSetNode {
	return ss.inner.GetRangeByRank(start, end, options)
}
