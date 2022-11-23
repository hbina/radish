package redis

import (
	"math/rand"
)

var _ Item = (*Set)(nil)

type Set struct {
	inner map[string]struct{}
}

func NewSetFromMap(value map[string]struct{}) *Set {
	return &Set{inner: value}
}

func NewSetEmpty() *Set {
	return &Set{inner: map[string]struct{}{}}
}

func (s *Set) Value() interface{} {
	return s.inner
}

func (l *Set) Type() uint64 {
	return ValueTypeSet
}

func (l *Set) TypeFancy() string {
	return ValueTypeFancySet
}

func (s *Set) AddMember(keys ...string) {
	for _, key := range keys {
		s.inner[key] = struct{}{}
	}
}

// RemoveMember removes the given member from the set.
// Returns true if the key exists. False otherwise.
func (s *Set) RemoveMember(key string) bool {
	_, exists := s.inner[key]
	delete(s.inner, key)
	return exists
}

func (s *Set) GetMembers() []string {
	r := make([]string, 0, len(s.inner))
	for k := range s.inner {
		r = append(r, k)
	}
	return r
}

func (s *Set) Exists(key string) bool {
	_, exists := s.inner[key]
	return exists
}

func (s *Set) Len() int {
	return len(s.inner)
}

// Pop removes a random key from the set.
func (s *Set) Pop() *string {
	member := s.GetRandomMember()
	if member != nil {
		s.RemoveMember(*member)
		return member
	}
	return nil
}

// GetRandomMeber returns a random member from the set.
func (s *Set) GetRandomMember() *string {
	if s.Len() > 0 {
		keys := s.GetMembers()
		idx := rand.Intn(len(keys))
		key := keys[idx]
		return &key
	}
	return nil
}

// Intersect returns a new Set that is an intersection of both sets.
// TODO: Better intersection algorithm?
func (s *Set) Intersect(o *Set) *Set {
	set := NewSetEmpty()

	// loop over smaller set
	if s.Len() < o.Len() {
		for elem := range s.inner {
			if o.Exists(elem) {
				set.AddMember(elem)
			}
		}
	} else {
		for elem := range o.inner {
			if s.Exists(elem) {
				set.AddMember(elem)
			}
		}
	}

	return set
}

// Union returns a new Set that is a union of both sets.
func (s *Set) Union(o *Set) *Set {
	set := NewSetEmpty()

	for elem := range s.inner {
		set.AddMember(elem)
	}

	for elem := range o.inner {
		set.AddMember(elem)
	}

	return set
}

// Diff returns a new Set that is a diff of both sets.
func (s *Set) Diff(o *Set) *Set {
	set := NewSetEmpty()

	for elem := range s.inner {
		if !o.Exists(elem) {
			set.AddMember(elem)
		}
	}

	return set
}

// TODO: For now we only store strings so this should be enough.
func (s *Set) ForEachF(f func(a string)) {
	for k := range s.inner {
		f(k)
	}
}
