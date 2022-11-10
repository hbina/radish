package redis

import "log"

var _ Item = (*Set)(nil)

type Set struct {
	value map[string]struct{}
}

func NewSetFromMap(value map[string]struct{}) *Set {
	return &Set{value: value}
}

func NewSetEmpty() *Set {
	return &Set{value: map[string]struct{}{}}
}

func (s Set) Value() interface{} {
	return s.value
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
	s.value[key] = struct{}{}
}

func (s *Set) GetMembers(key string) []string {
	r := make([]string, 0, len(s.value))
	for k := range s.value {
		r = append(r, k)
	}
	return r
}
