package redis

import "log"

var _ Item = (*String)(nil)

type String struct {
	inner string
}

func NewString(value string) *String {
	return &String{inner: value}
}

func (s *String) Value() interface{} {
	return s.inner
}

func (l String) Type() uint64 {
	return ValueTypeString
}

func (l String) TypeFancy() string {
	return ValueTypeFancyString
}

func (s String) OnDelete(key string, db RedisDb) {
	log.Printf("Deleting %s with key %s from database ID %d\n", s.TypeFancy(), key, db.id)
}
