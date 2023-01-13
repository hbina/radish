package util

type Resp interface {
	Width() int
}

var _ Resp = &RespSimpleString{}
var _ Resp = &RespErrorString{}
var _ Resp = &RespBulkString{}
var _ Resp = &RespInteger{}
var _ Resp = &RespNilBulk{}
var _ Resp = &RespArray{}
var _ Resp = &RespNilArray{}
var _ Resp = &RespMap{}
var _ Resp = &RespNil{}

type RespSimpleString struct {
	inner []byte
}

func (rs *RespSimpleString) Width() int {
	return 0
}

func NewRss(inner string) *RespSimpleString {
	return &RespSimpleString{
		inner: []byte(inner),
	}
}

type RespErrorString struct {
	inner []byte
}

func (rs *RespErrorString) Width() int {
	return 0
}

type RespBulkString struct {
	inner []byte
}

func (rs *RespBulkString) Width() int {
	return 0
}

type RespNilBulk struct {
}

func (rs *RespNilBulk) Width() int {
	return 0
}

type RespInteger struct {
	inner int
}

func (rs *RespInteger) Width() int {
	return 0
}

type RespFloat struct {
	inner float64
}

func (rs *RespFloat) Width() int {
	return 0
}

type RespArray struct {
	inner []Resp
}

func (rs *RespArray) Width() int {
	ourWidth := 0
	currLen := len(rs.inner)

	for currLen != 0 {
		ourWidth += 1
		currLen /= 10
	}

	return ourWidth
}

type RespNilArray struct {
}

func (rs *RespNilArray) Width() int {
	return 0
}

type RespMap struct {
	inner []Resp
}

func (rs *RespMap) Width() int {
	ourWidth := 0
	currLen := len(rs.inner)

	for currLen != 0 {
		ourWidth += 1
		currLen /= 10
	}

	return ourWidth
}

type RespNil struct {
}

func (rs *RespNil) Width() int {
	return 0
}
