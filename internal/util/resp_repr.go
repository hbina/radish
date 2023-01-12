package util

import (
	"fmt"
	"net"
)

type Resp interface {
	Width() int
	Write(net.Conn, bool)
}

var _ Resp = &RespSimpleString{}
var _ Resp = &RespErrorString{}
var _ Resp = &RespBulkString{}
var _ Resp = &RespInteger{}
var _ Resp = &RespNil{}
var _ Resp = &RespArray{}
var _ Resp = &RespMap{}

type RespSimpleString struct {
	inner []byte
}

func (rs *RespSimpleString) Width() int {
	return 0
}

func (rs *RespSimpleString) Write(c net.Conn, useV3 bool) {
	c.Write([]byte(fmt.Sprintf("+%s\r\n", rs.inner)))
}

type RespErrorString struct {
	inner []byte
}

func (rs *RespErrorString) Width() int {
	return 0
}

func (rs *RespErrorString) Write(c net.Conn, useV3 bool) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", rs.inner)))
}

type RespBulkString struct {
	inner []byte
}

func (rs *RespBulkString) Width() int {
	return 0
}

func (rs *RespBulkString) Write(c net.Conn, useV3 bool) {
	c.Write([]byte(fmt.Sprintf("$%d\r\n%s\r\n", len(rs.inner), rs.inner)))
}

type RespInteger struct {
	inner int
}

func (rs *RespInteger) Width() int {
	return 0
}

func (rs *RespInteger) Write(c net.Conn, useV3 bool) {
	c.Write([]byte(fmt.Sprintf(":%d\r\n", rs.inner)))
}

type RespNilKind = int

const (
	RespNilKindBulk RespNilKind = iota
	RespNilKindArray
)

// RespNil represents both bulk nil and array nil
// kind = 0 => bulk
// kind = 1 => array
type RespNil struct {
	kind RespNilKind
}

func (rs *RespNil) Width() int {
	return 0
}

func (rs *RespNil) Write(c net.Conn, useV3 bool) {
	if useV3 {
		c.Write([]byte("_\r\n"))
	} else {
		if rs.kind == 0 {
			c.Write([]byte("$-1\r\n"))
		} else {
			c.Write([]byte("*-1\r\n"))
		}
	}
}

func NewRespNil(kind RespNilKind) *RespNil {
	return &RespNil{
		kind: kind,
	}
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

func (rs *RespArray) Write(c net.Conn, useV3 bool) {
	c.Write([]byte(fmt.Sprintf("*%d\r\n", len(rs.inner))))
	for _, r := range rs.inner {
		r.Write(c, useV3)
	}
}

type RespMap struct {
	inner map[string]Resp
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

func (rs *RespMap) Write(c net.Conn, useV3 bool) {
	c.Write([]byte(fmt.Sprintf("%%%d\r\n", len(rs.inner))))
	for k, r := range rs.inner {
		c.Write
		r.Write(c, useV3)
	}
}
