package redis

import (
	"fmt"
	"strconv"
	"strings"
)

func CreateRespReply(in []byte) (string, error) {
	sds, leftover := implCreateRespReply(in, 0)

	if len(leftover) != 0 {
		return "", fmt.Errorf("CreateRespReply have %d leftovers ", len(leftover))
	}

	maxDepth := 0
	for _, sd := range sds {
		if sd.depth > maxDepth {
			maxDepth = sd.depth
		}
	}

	depthCount := make([]int, maxDepth)

	// Build the strings
	for _, sd := range sds {
		r := ""
		for i := 0; i < sd.depth; i++ {
			r += fmt.Sprintf("%d) ", depthCount[i])
			depthCount[i]++
		}
		r += sd.inner
		fmt.Println(r)
	}

	return "", nil
}

type SD struct {
	inner string
	depth int
}

func getSdErr(depth int) ([]SD, []byte) {
	return []SD{{
		inner: "",
		depth: depth,
	}}, []byte{}
}

func getSdFromStr(inner string, depth int) []SD {
	return []SD{{
		inner: inner,
		depth: depth,
	}}
}

func implCreateRespReply(in []byte, depth int) ([]SD, []byte) {
	if len(in) == 0 {
		return getSdErr(depth)
	} else {
		currByte := in[0]
		if currByte == '+' {
			// Simple strings
			str, leftover := TakeBytesUntilClrf(in[1:])
			in = leftover
			res := fmt.Sprintf("\"%s\"", string(str))
			return getSdFromStr(res, depth), in
		} else if currByte == '-' {
			// Error strings
			str, leftover := TakeBytesUntilClrf(in[1:])
			in = leftover
			res := fmt.Sprintf("\"%s\"", string(str))
			return getSdFromStr(res, depth), in
		} else if currByte == ':' {
			// Integers
			str, leftover := TakeBytesUntilClrf(in[1:])
			in = leftover
			strInt64, err := strconv.ParseInt(string(str), 10, 32)

			if err != nil {
				return getSdErr(depth)
			}

			res := fmt.Sprintf("(integer) %d", strInt64)
			return getSdFromStr(res, depth), in
		} else if currByte == '$' {
			// Bulk strings
			lenStr, leftover := TakeBytesUntilClrf(in[1:])
			lenInt64, err := strconv.ParseInt(string(lenStr), 10, 32)
			in = leftover

			if err != nil {
				return getSdErr(depth)
			}

			if lenInt64 < 0 {
				return getSdFromStr("(nil)", depth), in
			} else {
				// TODO: Reuse lenInt for optimization purposes?
				bulkStr, leftover := TakeBytesUntilClrf(in)
				in = leftover

				res := fmt.Sprintf("\"%s\"", string(bulkStr))
				return getSdFromStr(res, depth), in
			}
		} else if currByte == '*' {
			// Arrays
			fmt.Println("depth", depth, " from ", EscapeString(string(in)))
			lenStr, leftover := TakeBytesUntilClrf(in[1:])
			len64, err := strconv.ParseInt(string(lenStr), 10, 32)
			in = leftover

			if err != nil {
				return getSdErr(depth)
			}

			if len64 < 0 {
				return getSdErr(depth)
			} else if len64 == 0 {
				return getSdFromStr("(empty)", depth), in
			} else {
				res := make([]SD, 0)
				for idx := 0; idx < int(len64) && len(in) != 0; idx++ {
					replies, leftover := implCreateRespReply(in, depth+1)
					res = append(res, replies...)
					in = leftover
				}
				return res, in
			}
		}
	}

	return getSdErr(depth)
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
