package redis

var _ Item = (*ZSet)(nil)

type ZSet struct {
	inner SortedSet
}

func NewZSet() *ZSet {
	return &ZSet{inner: *New()}
}

func NewZSetFromSs(value SortedSet) *ZSet {
	return &ZSet{inner: value}
}

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
	return s.inner.GetCount()
}
