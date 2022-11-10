package redis

import (
	"log"

	"github.com/wangjia184/sortedset"
)

var _ Item = (*ZSet)(nil)

type ZSet struct {
	value sortedset.SortedSet
}

func NewZSetEmpty() *ZSet {
	return &ZSet{value: *sortedset.New()}
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
