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

func (s Set) Value() interface{} {
	return s.inner
}

func (l Set) Type() uint64 {
	return ValueTypeSet
}

func (l Set) TypeFancy() string {
	return ValueTypeFancySet
}

func (s Set) OnDelete(key string, db RedisDb) {
	log.Printf("Deleting %s with key %s from database ID %d\n", s.TypeFancy(), key, db.id)
}

func (s *Set) AddMember(key string) {
	s.inner[key] = struct{}{}
}

func (s *Set) GetMembers(key string) []string {
	r := make([]string, 0, len(s.inner))
	for k := range s.inner {
		r = append(r, k)
	}
	return r
}
