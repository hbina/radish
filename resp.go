package redis

import (
	"fmt"
	"strconv"
	"strings"
)

func StringifyRespBytes(in []byte) (string, []byte) {
	res, leftover := convertBytesToRespType(in)
	return stringifyRespType(res, 0), leftover
}

type RespType interface {
	TypeId() string
	Value() interface{}
	RequiredWidth() int
}

type RespString struct {
	inner []byte
}

func (rs *RespString) TypeId() string {
	return "string"
}

func (rs *RespString) Value() interface{} {
	return string(rs.inner)
}

func (rs *RespString) RequiredWidth() int {
	return 0
}

type RespInteger struct {
	inner int
}

func (rs *RespInteger) TypeId() string {
	return "integer"
}

func (rs *RespInteger) Value() interface{} {
	return rs.inner
}

func (rs *RespInteger) RequiredWidth() int {
	return 0
}

type RespNil struct {
}

func (rs *RespNil) TypeId() string {
	return "nil"
}

func (rs *RespNil) Value() interface{} {
	return "(nil)"
}

func (rs *RespNil) RequiredWidth() int {
	return 0
}

type RespArray struct {
	inner []RespType
}

func (rs *RespArray) TypeId() string {
	return "array"
}

func (rs *RespArray) Value() interface{} {
	return rs.inner
}

// TODO: Do we actually need this?
func (rs *RespArray) RequiredWidth() int {
	ourWidth := 0
	currLen := len(rs.inner)

	for currLen != 0 {
		ourWidth += 1
		currLen /= 10
	}

	return ourWidth
}

func stringifyRespType(res RespType, width int) string {
	if res == nil {
		return ""
	} else if res.TypeId() == "string" {
		return fmt.Sprintf("\"%s\"", res.Value().(string))
	} else if res.TypeId() == "integer" {
		return fmt.Sprintf("(integer) %d", res.Value().(int))
	} else if res.TypeId() == "nil" {
		return res.Value().(string)
	} else if res.TypeId() == "array" {
		var str strings.Builder
		arr := res.Value().([]RespType)

		if width > 0 {
			width += 2
		}

		var padding strings.Builder
		for i := 0; i < width; i++ {
			padding.WriteByte(' ')
		}

		if len(arr) == 0 {
			return "(empty)"
		} else {
			for i, v := range arr {
				if i == 0 {
					str.WriteString(fmt.Sprintf("%d) %s\n", i+1, stringifyRespType(v, res.RequiredWidth()+width)))
				} else if i == len(arr)-1 {
					str.WriteString(fmt.Sprintf("%s%d) %s", padding.String(), i+1, stringifyRespType(v, res.RequiredWidth()+width)))
				} else {
					str.WriteString(fmt.Sprintf("%s%d) %s\n", padding.String(), i+1, stringifyRespType(v, res.RequiredWidth()+width)))
				}
			}
			return str.String()
		}
	}

	return ""
}

func convertBytesToRespType(input []byte) (RespType, []byte) {
	// We need at least 1 byte for the first redis type byte
	if len(input) > 0 {
		currByte := input[0]
		if currByte == '+' || currByte == '-' {
			str, leftover, ok := TakeBytesUntilClrf(input[1:])

			if !ok {
				return nil, input
			}

			rs := RespString{
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

			rs := RespInteger{
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
				rs := RespNil{}

				return &rs, leftover
			} else {
				if int(lenInt64)+2 > len(leftover) {
					return nil, input
				} else {
					if leftover[lenInt64] == '\r' && leftover[lenInt64+1] == '\n' {

						rs := RespString{
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
				return nil, input
			} else if lenInt64 == 0 {
				rs := RespArray{
					inner: make([]RespType, 0),
				}

				return &rs, leftover
			} else {
				replies := make([]RespType, 0, lenInt64)
				for idx := 0; idx < int(lenInt64) && len(nextInput) != 0; idx++ {
					reply, leftover := convertBytesToRespType(nextInput)

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

				rs := RespArray{
					inner: replies,
				}

				return &rs, nextInput
			}
		}
	}

	return nil, input
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

func SplitStringIntoArgs(s string) ([]string, bool) {
	res := []string{}
	var currStr strings.Builder
	inQuote := false

	for idx := 0; idx < len(s); idx++ {
		currChar := s[idx]
		hasNext := len(s) > idx+1
		if currChar == '\\' {
			// we are escaping something
			if hasNext {
				nextChar := s[idx]
				currStr.WriteByte(currChar)
				currStr.WriteByte(currChar)
				currStr.WriteByte(nextChar)
				idx++
			} else {
				return res, false
			}
		} else if currChar == '"' {
			inQuote = !inQuote
		} else if currChar == ' ' && !inQuote {
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
