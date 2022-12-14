package types

var _ Item = (*String)(nil)

type String struct {
	Inner string
}

func NewString(value string) *String {
	return &String{Inner: value}
}

func (s *String) Value() interface{} {
	return s.Inner
}

func (l *String) Type() uint64 {
	return ValueTypeString
}

func (l *String) TypeFancy() string {
	return ValueTypeFancyString
}

func (s *String) Len() int {
	return len(s.Inner)
}

func (s *String) Get(idx int) byte {
	v := s.Inner[idx]
	return v
}
