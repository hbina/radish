package util

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func StringifyRespBytes(in []byte) (string, bool) {
	resp, leftover := ParseResp2(in)

	if len(leftover) != 0 || resp == nil {
		return "", false
	}

	inList := false
	_, isArr := resp.(*Resp2Array)

	if isArr {
		inList = true
	}

	str, ok := stringifyRespType(resp, 0, inList)

	if !ok {
		return "", false
	}

	return str, true
}

type Resp2 interface {
	Width() int
	WriteToConn(net.Conn) error
}

var _ Resp2 = &Resp2SimpleString{}
var _ Resp2 = &Resp2ErrorString{}
var _ Resp2 = &Resp2BulkString{}
var _ Resp2 = &Resp2Integer{}
var _ Resp2 = &Resp2NilArray{}
var _ Resp2 = &Resp2Array{}

type Resp2SimpleString struct {
	inner []byte
}

func NewResp2Ss(str string) *Resp2SimpleString {
	return &Resp2SimpleString{
		inner: []byte(str),
	}
}

func (rs *Resp2SimpleString) Width() int {
	return 0
}

func (rs *Resp2SimpleString) WriteToConn(conn net.Conn) error {
	err := WriteAll(conn, []byte("+"))

	if err != nil {
		return err
	}

	err = WriteAll(conn, rs.inner)

	if err != nil {
		return err
	}

	err = WriteAll(conn, []byte("\r\n"))

	if err != nil {
		return err
	}

	return nil
}

type Resp2ErrorString struct {
	inner []byte
}

func NewResp2Es(str string) *Resp2ErrorString {
	return &Resp2ErrorString{
		inner: []byte(str),
	}
}

func (rs *Resp2ErrorString) Width() int {
	return 0
}

func (rs *Resp2ErrorString) WriteToConn(conn net.Conn) error {
	err := WriteAll(conn, []byte("-"))

	if err != nil {
		return err
	}

	err = WriteAll(conn, rs.inner)

	if err != nil {
		return err
	}

	err = WriteAll(conn, []byte("\r\n"))

	if err != nil {
		return err
	}

	return nil
}

type Resp2BulkString struct {
	inner []byte
}

func NewResp2Bs(str string) *Resp2BulkString {
	return &Resp2BulkString{
		inner: []byte(str),
	}
}

func (rs *Resp2BulkString) Width() int {
	return 0
}

func (rs *Resp2BulkString) WriteToConn(conn net.Conn) error {
	err := WriteAll(conn, []byte(fmt.Sprintf("$%d\r\n", len(rs.inner))))

	if err != nil {
		return err
	}

	err = WriteAll(conn, rs.inner)

	if err != nil {
		return err
	}

	err = WriteAll(conn, []byte("\r\n"))

	if err != nil {
		return err
	}

	return nil
}

type Resp2NilBulk struct {
}

func (rs *Resp2NilBulk) Width() int {
	return 0
}

func (rs *Resp2NilBulk) WriteToConn(conn net.Conn) error {
	err := WriteAll(conn, []byte("$-1\r\n"))

	if err != nil {
		return err
	}

	return nil
}

type Resp2Integer struct {
	inner int
}

func NewResp2I(value int) *Resp2Integer {
	return &Resp2Integer{
		inner: value,
	}
}

func (rs *Resp2Integer) Width() int {
	return 0
}

func (rs *Resp2Integer) WriteToConn(conn net.Conn) error {
	err := WriteAll(conn, []byte(fmt.Sprintf(":%d\r\n", rs.inner)))

	if err != nil {
		return err
	}

	return nil
}

type Resp2Array struct {
	inner []Resp2
}

func NewResp2Arr() *Resp2Array {
	return &Resp2Array{
		inner: make([]Resp2, 0),
	}
}

func (rs *Resp2Array) Width() int {
	ourWidth := 0
	currLen := len(rs.inner)

	for currLen != 0 {
		ourWidth += 1
		currLen /= 10
	}

	return ourWidth
}

func (rs *Resp2Array) WriteToConn(conn net.Conn) error {
	err := WriteAll(conn, []byte(fmt.Sprintf("*%d\r\n", len(rs.inner))))

	if err != nil {
		return err
	}

	for _, r := range rs.inner {
		err = r.WriteToConn(conn)

		if err != nil {
			return err
		}
	}

	return nil
}

type Resp2NilArray struct {
}

func (rs *Resp2NilArray) Width() int {
	return 0
}

func (rs *Resp2NilArray) WriteToConn(conn net.Conn) error {
	err := WriteAll(conn, []byte("*-1\r\n"))

	if err != nil {
		return err
	}

	return nil
}

func stringifyRespType(res Resp2, width int, inList bool) (string, bool) {
	if res == nil {
		return "", false
	} else if rs, ok := res.(*Resp2BulkString); ok {
		return fmt.Sprintf("\"%s\"", string(rs.inner)), true
	} else if rs, ok := res.(*Resp2SimpleString); ok && inList {
		return fmt.Sprintf("\"%s\"", string(rs.inner)), true
	} else if rs, ok := res.(*Resp2ErrorString); ok {
		return string(rs.inner), true
	} else if rs, ok := res.(*Resp2Integer); ok {
		return fmt.Sprintf("(integer) %d", rs.inner), true
	} else if _, ok := res.(*Resp2NilArray); ok {
		return "(nil)", true
	} else if rs, ok := res.(*Resp2Array); ok {
		var str strings.Builder
		arr := rs.inner

		if width > 0 {
			width += 2
		}

		var padding strings.Builder
		for i := 0; i < width; i++ {
			padding.WriteByte(' ')
		}

		if len(arr) == 0 {
			return "(empty)", true
		} else {
			for i, v := range arr {
				s, ok := stringifyRespType(v, res.Width()+width, true)

				if !ok {
					return "", false
				}

				if i == 0 {
					str.WriteString(fmt.Sprintf("%d) %s\n", i+1, s))
				} else if i == len(arr)-1 {
					str.WriteString(fmt.Sprintf("%s%d) %s", padding.String(), i+1, s))
				} else {
					str.WriteString(fmt.Sprintf("%s%d) %s\n", padding.String(), i+1, s))
				}
			}
			return str.String(), true
		}
	}

	return "", false
}

func ParseResp2(input []byte) (Resp2, []byte) {
	// We need at least 1 byte for the first redis type byte
	if len(input) > 0 {
		currByte := input[0]
		if currByte == '+' {
			str, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok {
				return nil, input
			}

			rs := Resp2SimpleString{
				inner: str,
			}

			return &rs, leftover
		} else if currByte == '-' {
			str, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok {
				return nil, input
			}

			rs := Resp2ErrorString{
				inner: str,
			}

			return &rs, leftover
		} else if currByte == ':' {
			str, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok {
				return nil, input
			}

			valInt64, err := strconv.ParseInt(string(str), 10, 32)

			if err != nil {
				return nil, input
			}

			rs := Resp2Integer{
				inner: int(valInt64),
			}

			return &rs, leftover

		} else if currByte == '$' {
			lenStr, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok {
				return nil, input
			}

			lenInt64, err := strconv.ParseInt(string(lenStr), 10, 32)

			if err != nil {
				return nil, input
			}

			if lenInt64 < 0 {
				rs := Resp2NilBulk{}

				return &rs, leftover
			} else {
				if int(lenInt64)+2 > len(leftover) {
					return nil, input
				} else {
					if leftover[lenInt64] == '\r' && leftover[lenInt64+1] == '\n' {

						rs := Resp2BulkString{
							inner: leftover[:lenInt64],
						}

						return &rs, leftover[lenInt64+2:]
					} else {
						return nil, input
					}
				}
			}
		} else if currByte == '*' {
			lenStr, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok {
				return nil, input
			}

			lenInt64, err := strconv.ParseInt(string(lenStr), 10, 32)

			if err != nil {
				return nil, input
			}

			// We parsed the length of the array, now we march forward
			nextInput := leftover

			if lenInt64 < 0 {
				rs := Resp2NilBulk{}

				return &rs, leftover
			} else if lenInt64 == 0 {
				rs := Resp2Array{
					inner: make([]Resp2, 0),
				}

				return &rs, leftover
			} else {
				replies := make([]Resp2, 0, lenInt64)
				for idx := 0; idx < int(lenInt64) && len(nextInput) != 0; idx++ {
					reply, leftover := ParseResp2(nextInput)

					// If any of the elements are bad or we can't make progress, just bail
					if reply == nil || len(leftover) == len(nextInput) {
						return nil, input
					}

					nextInput = leftover
					replies = append(replies, reply)
				}

				if len(replies) != int(lenInt64) {
					return nil, input
				}

				rs := Resp2Array{
					inner: replies,
				}

				return &rs, nextInput
			}
		} else {
			str, leftover, ok := TakeBytesUntilClrf(input)

			if !ok {
				return nil, input
			}

			rs := Resp2SimpleString{
				inner: str,
			}

			return &rs, leftover
		}
	}

	return nil, input
}

const (
	quoteModeNone = iota
	quoteModeSingle
	quoteModeDouble
)

func SplitStringIntoArgs(s string) ([]string, bool) {
	res := []string{}
	var currStr strings.Builder
	inQuote := quoteModeNone

	for idx := 0; idx < len(s); idx++ {
		currChar := s[idx]
		hasNext := len(s) > idx+1
		if currChar == '\\' {
			// we are escaping something
			if hasNext {
				nextChar := s[idx+1]
				currStr.WriteByte(currChar)
				currStr.WriteByte(currChar)
				currStr.WriteByte(nextChar)
				idx++
			} else {
				return res, false
			}
		} else if currChar == '"' {
			if inQuote == quoteModeDouble {
				inQuote = quoteModeNone
			} else if inQuote == quoteModeSingle {
				currStr.WriteByte('"')
			} else {
				inQuote = quoteModeDouble
			}
		} else if currChar == '\'' {
			if inQuote == quoteModeSingle {
				inQuote = quoteModeNone
			} else if inQuote == quoteModeDouble {
				currStr.WriteByte('\'')
			} else {
				inQuote = quoteModeSingle
			}
		} else if currChar == ' ' && inQuote == quoteModeNone {
			res = append(res, currStr.String())
			currStr.Reset()
		} else {
			currStr.WriteByte(currChar)
		}
	}
	res = append(res, currStr.String())

	return res, true
}

func ConvertCommandArgToResp(args []string) string {
	var str strings.Builder

	str.WriteString(fmt.Sprintf("*%d\r\n", len(args)))

	for _, arg := range args {
		str.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}

	return str.String()
}

func ConvertRespToArgs(resp Resp2) [][]byte {
	if resp == nil {
		return [][]byte{}
	}

	arr, ok := resp.(*Resp2Array)

	if ok {

		args := make([][]byte, 0)

		for _, r := range arr.inner {
			str, ok := r.(*Resp2SimpleString)

			if ok {
				args = append(args, []byte(str.inner))
				continue
			}

			bulkStr, ok := r.(*Resp2BulkString)

			if ok {
				args = append(args, []byte(bulkStr.inner))
				continue
			}

			return [][]byte{}
		}

		return args
	}

	str, ok := resp.(*Resp2SimpleString)

	if ok {
		return [][]byte{[]byte(str.inner)}
	}

	bulkStr, ok := resp.(*Resp2BulkString)

	if ok {
		return [][]byte{[]byte(bulkStr.inner)}
	}

	return [][]byte{}
}
