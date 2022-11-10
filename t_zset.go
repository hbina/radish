package redis

import (
	"log"

	"github.com/zavitax/sortedset-go"
)

var _ Item = (*ZSet)(nil)

type ZSet struct {
	value sortedset.SortedSet[string, float64, struct{}]
}

func NewZSetEmpty() *ZSet {
	return &ZSet{value: *sortedset.New[string, float64, struct{}]()}
}

func NewZSet(value sortedset.SortedSet[string, float64, struct{}]) *ZSet {
	return &ZSet{value: value}
}

func (s *ZSet) Value() interface{} {
	return s.value
}

func (l ZSet) Type() uint64 {
	return ValueTypeZSet
}

func (l ZSet) TypeFancy() string {
	return ValueTypeFancyZSet
}

func (s ZSet) OnDelete(key string, db RedisDb) {
	log.Printf("Deleting %s with key %s from database ID %d\n", s.TypeFancy(), key, db.id)
}
