package redis

import "log"

var _ Item = (*String)(nil)

type String struct {
	value *string
}

func NewString(value *string) *String {
	return &String{value: value}
}

func (s *String) Value() interface{} {
	return s.value
}

func (l *String) Type() uint64 {
	return ValueTypeString
}

func (l *String) TypeFancy() string {
	return ValueTypeFancyString
}

func (s *String) OnDelete(key *string, db *RedisDb) {
	log.Printf("Deleting string with key %s from database ID %d\n", *key, db.id)
}
