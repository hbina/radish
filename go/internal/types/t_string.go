package types

import "encoding/json"

var _ Item = (*String)(nil)

type String struct {
	inner string
}

// impl Item for String

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

// impl String

func (s *String) Get(idx int) byte {
	v := s.inner[idx]
	return v
}

func (s *String) Marshal() ([]byte, error) {
	str, err := json.Marshal(s.inner)
	return str, err
}

func StringUnmarshal(data []byte) (*String, bool) {
	var set string
	err := json.Unmarshal(data, &set)

	if err != nil {
		return nil, false
	}

	return NewString(set), true
}

func (s *String) AsBytes() []byte {
	return []byte(s.inner)
}

func (s *String) AsString() string {
	return s.inner
}

func (s *String) SubString(start, end int) string {
	return s.inner[start:end]
}

func (s *String) Reverse() String {
	n := len(s.inner)
	runes := make([]rune, n)
	for _, rune := range s.inner {
		n--
		runes[n] = rune
	}
	return *NewString(string(runes[n:]))
}
