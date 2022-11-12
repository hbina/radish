package redis

import "log"

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

func (s *Set) OnDelete(key string, db RedisDb) {
	log.Printf("Deleting %s with key %s from database ID %d\n", s.TypeFancy(), key, db.id)
}

func (s *Set) AddMember(keys ...string) {
	for _, key := range keys {
		s.inner[key] = struct{}{}
	}
}

func (s *Set) RemoveMember(keys ...string) {
	for _, key := range keys {
		delete(s.inner, key)
	}
}

func (s *Set) GetMembers(key string) []string {
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
