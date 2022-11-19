package redis

import (
	"fmt"
	"strconv"
	"strings"
)

func CreateRespReply(in []byte) (string, []byte) {
	var res strings.Builder

	if len(in) == 0 {
		return "", []byte{}
	} else {
		currByte := in[0]
		if currByte == '+' {
			str, leftover := TakeBytesUntilClrf(in[1:])
			in = leftover
			res.WriteString(string(str))
		} else if currByte == '-' {
			str, leftover := TakeBytesUntilClrf(in[1:])
			in = leftover
			res.WriteString(string(str))
		} else if currByte == ':' {
			str, leftover := TakeBytesUntilClrf(in[1:])
			in = leftover
			strInt64, err := strconv.ParseInt(string(str), 10, 32)

			if err != nil {
				return "", []byte{}
			}

			res.WriteString(fmt.Sprint(int(strInt64)))
		} else if currByte == '$' {
			lenStr, leftover := TakeBytesUntilClrf(in[1:])
			in = leftover

			lenInt64, err := strconv.ParseInt(string(lenStr), 10, 32)

			if err != nil {
				return "", []byte{}
			}

			if lenInt64 < 0 {
				res.WriteString("(nil)")
			} else {
				// TODO: Reuse lenInt for optimization purposes?
				bulkStr, leftover := TakeBytesUntilClrf(in)
				in = leftover

				res.WriteString(string(bulkStr))
			}
		} else if currByte == '*' {
			lenStr, leftover := TakeBytesUntilClrf(in[1:])
			in = leftover

			lenInt64, err := strconv.ParseInt(string(lenStr), 10, 32)

			if err != nil {
				return "", []byte{}
			}

			if lenInt64 < 0 {
				return "", []byte{}
			} else if lenInt64 == 0 {
				res.WriteString("(empty)")
			} else {
				for idx := 0; idx < int(lenInt64) && len(in) != 0; idx++ {
					reply, leftover := CreateRespReply(in)
					in = leftover
					res.WriteString(fmt.Sprintf("%d) \"%s\"", idx+1, EscapeString(reply)))
					if idx+1 < int(lenInt64) && len(in) != 0 {
						res.WriteByte('\n')
					}
				}
			}
		}
	}

	return res.String(), in
}

func ParseRespString(in []byte) (string, []byte, bool) {
	str, leftover := TakeBytesUntilClrf(in)
	return string(str), leftover, true
}

func TakeBytesUntilClrf(in []byte) ([]byte, []byte) {

	// If there's not even enough space for \r\n,
	// its a bad input and we return empty
	if len(in) < 2 {
		return []byte{}, []byte{}
	}

	res := []byte{}

	for idx := 0; idx+2 < len(in) && in[idx] != '\r' && in[idx+1] != '\n'; idx++ {
		res = append(res, in[idx])
	}

	return res, in[len(res)+2:]
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
