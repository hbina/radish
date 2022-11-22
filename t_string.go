package redis

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

func (l *String) Type() uint64 {
	return ValueTypeString
}

func (l *String) TypeFancy() string {
	return ValueTypeFancyString
}

func (s *String) Len() int {
	return len(s.inner)
}

func (s *String) Get(idx int) byte {
	v := s.inner[idx]
	return v
}
