package util

import (
	"fmt"
	"strconv"
	"strings"
)

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

func StringifyRespBytes(in []byte) (string, bool, []byte) {
	resp, leftover := ConvertBytesToRespType(in)

	if resp == nil {
		return "", false, []byte{}
	}

	inList := false
	_, isArr := resp.(*RespArray)
	_, isMap := resp.(*RespMap)

	if isArr || isMap {
		inList = true
	}

	str, ok := stringifyRespType(resp, 0, inList)

	if !ok {
		return "", false, []byte{}
	}

	return str, true, leftover
}

func stringifyRespType(res Resp, width int, inList bool) (string, bool) {
	if res == nil {
		return "", false
	} else if rs, ok := res.(*RespBulkString); ok {
		return fmt.Sprintf("\"%s\"", string(rs.inner)), true
	} else if rs, ok := res.(*RespSimpleString); ok && inList {
		return fmt.Sprintf("\"%s\"", string(rs.inner)), true
	} else if rs, ok := res.(*RespErrorString); ok && inList {
		return fmt.Sprintf("\"%s\"", string(rs.inner)), true
	} else if rs, ok := res.(*RespSimpleString); ok {
		return string(rs.inner), true
	} else if rs, ok := res.(*RespErrorString); ok {
		return string(rs.inner), true
	} else if rs, ok := res.(*RespInteger); ok {
		return fmt.Sprintf("(integer) %d", rs.inner), true
	} else if rs, ok := res.(*RespFloat); ok {
		return fmt.Sprintf("(double) %f", rs.inner), true
	} else if _, ok := res.(*RespNilBulk); ok {
		return "(nil)", true
	} else if _, ok := res.(*RespNilArray); ok {
		return "(nil)", true
	} else if _, ok := res.(*RespNil); ok {
		return "(nil)", true
	} else if rs, ok := res.(*RespArray); ok {
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
	} else if rs, ok := res.(*RespMap); ok {
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
		} else if len(arr)%2 != 0 {
			return FormatErr, true
		} else {
			for i := 0; i < len(arr); i += 2 {
				first, ok := stringifyRespType(arr[i], res.Width()+width, true)

				if !ok {
					return "", false
				}

				second, ok := stringifyRespType(arr[i+1], res.Width()+width, true)

				if !ok {
					return "", false
				}

				if i == 0 {
					str.WriteString(fmt.Sprintf("%s => %s\n", first, second))
				} else if i == len(arr)-2 {
					str.WriteString(fmt.Sprintf("%s%s => %s", padding.String(), first, second))
				} else {
					str.WriteString(fmt.Sprintf("%s%s => %s\n", padding.String(), first, second))
				}
			}
			return str.String(), true
		}
	}

	return "", false
}

func ConvertBytesToRespType(input []byte) (Resp, []byte) {
	// We need at least 1 byte for the first redis type byte
	if len(input) > 0 {
		currByte := input[0]
		if currByte == '_' {
			str, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok || len(str) != 0 {
				return nil, []byte{}
			}

			return &RespNil{}, leftover
		} else if currByte == '+' {
			str, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok {
				return nil, []byte{}
			}

			return &RespSimpleString{
				inner: str,
			}, leftover
		} else if currByte == '-' {
			str, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok {
				return nil, []byte{}
			}

			return &RespErrorString{
				inner: str,
			}, leftover
		} else if currByte == ':' {
			str, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok {
				return nil, []byte{}
			}

			valInt64, err := strconv.ParseInt(string(str), 10, 32)

			if err != nil {
				return nil, []byte{}
			}

			return &RespInteger{
				inner: int(valInt64),
			}, leftover
		} else if currByte == ',' {
			str, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok {
				return nil, []byte{}
			}

			valFloat64, err := strconv.ParseFloat(string(str), 64)

			if err != nil {
				return nil, []byte{}
			}

			return &RespFloat{
				inner: valFloat64,
			}, leftover
		} else if currByte == '$' {
			lenStr, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok {
				return nil, []byte{}
			}

			lenInt64, err := strconv.ParseInt(string(lenStr), 10, 32)

			if err != nil {
				return nil, []byte{}
			}

			if lenInt64 < 0 {
				return &RespNilBulk{}, leftover
			} else {
				if int(lenInt64)+2 > len(leftover) {
					return nil, []byte{}
				} else {
					if leftover[lenInt64] == '\r' && leftover[lenInt64+1] == '\n' {
						return &RespBulkString{
							inner: leftover[:lenInt64],
						}, leftover[lenInt64+2:]
					} else {
						return nil, []byte{}
					}
				}
			}
		} else if currByte == '*' {
			lenStr, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok {
				return nil, []byte{}
			}

			lenInt64, err := strconv.ParseInt(string(lenStr), 10, 32)

			if err != nil {
				return nil, []byte{}
			}

			// We parsed the length of the array, now we march forward
			nextInput := leftover

			if lenInt64 < 0 {
				return &RespNilArray{}, leftover
			} else if lenInt64 == 0 {
				return &RespArray{
					inner: make([]Resp, 0),
				}, leftover
			} else {
				replies := make([]Resp, 0, lenInt64)
				for idx := 0; idx < int(lenInt64) && len(nextInput) != 0; idx++ {
					reply, leftover := ConvertBytesToRespType(nextInput)

					// If any of the elements are bad or we can't make progress, just bail
					if reply == nil || len(leftover) == len(nextInput) {
						return nil, []byte{}
					}

					nextInput = leftover
					replies = append(replies, reply)
				}

				if len(replies) != int(lenInt64) {
					return nil, []byte{}
				}

				return &RespArray{
					inner: replies,
				}, nextInput
			}
		} else if currByte == '%' {
			lenStr, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok {
				return nil, []byte{}
			}

			lenInt64, err := strconv.ParseInt(string(lenStr), 10, 64)

			if err != nil {
				return nil, []byte{}
			}

			// We parsed the length of the array, now we march forward
			nextInput := leftover

			if lenInt64 < 0 {
				return &RespNilArray{}, leftover
			} else if lenInt64 == 0 {
				return &RespMap{
					inner: make([]Resp, 0),
				}, leftover
			} else {
				lenInt64 = lenInt64 * 2
				replies := make([]Resp, 0, lenInt64)
				for idx := 0; idx < int(lenInt64) && len(nextInput) != 0; idx++ {
					reply, leftover := ConvertBytesToRespType(nextInput)

					// If any of the elements are bad or we can't make progress, just bail
					if reply == nil || len(leftover) == len(nextInput) {
						return nil, []byte{}
					}

					nextInput = leftover
					replies = append(replies, reply)
				}

				if len(replies) != int(lenInt64) {
					return nil, []byte{}
				}

				return &RespMap{
					inner: replies,
				}, nextInput
			}
		} else {
			str, leftover, ok := TakeBytesUntilClrf(input)

			if !ok {
				return nil, []byte{}
			}

			return &RespSimpleString{
				inner: str,
			}, leftover
		}
	}

	return nil, []byte{}
}

func TakeBytesUntilClrf(in []byte) ([]byte, []byte, bool) {
	if len(in) == 0 {
		return []byte{}, []byte{}, false
	}

	idx := 0
	// We don't have to check for escapes here because we check for both CRLF
	for len(in) > idx+1 && !(in[idx] == '\r' && in[idx+1] == '\n') {
		idx++
	}

	if len(in) > idx+1 && in[idx] == '\r' && in[idx+1] == '\n' {
		return in[:idx], in[idx+2:], true
	} else {
		return in, []byte{}, false
	}
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

func ConvertRespToArgs(resp Resp) [][]byte {
	if resp == nil {
		return [][]byte{}
	}

	arr, ok := resp.(*RespArray)

	if ok {

		args := make([][]byte, 0)

		for _, r := range arr.inner {
			str, ok := r.(*RespSimpleString)

			if ok {
				args = append(args, []byte(str.inner))
				continue
			}

			bulkStr, ok := r.(*RespBulkString)

			if ok {
				args = append(args, []byte(bulkStr.inner))
				continue
			}

			return [][]byte{}
		}

		return args
	}

	str, ok := resp.(*RespSimpleString)

	if ok {
		return [][]byte{[]byte(str.inner)}
	}

	bulkStr, ok := resp.(*RespBulkString)

	if ok {
		return [][]byte{[]byte(bulkStr.inner)}
	}

	return [][]byte{}
}
