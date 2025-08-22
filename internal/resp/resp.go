// Package resp for serializing and deserializing RESP commands
package resp

import (
	"fmt"
	"strconv"

	"github.com/dylanmccormick/go-redis/internal/util"
)

func ParseRESP(splitCommand [][]byte) (any, int, error) {

	switch splitCommand[0][0] {

	case util.Star:
		// []any
		return parseArray(splitCommand)

	case util.DollarSign:
		// string
		return parseBulkString(splitCommand)

	case util.Dash:
		// string
		return parseBulkString(splitCommand)

	case util.Colon:
		// int
		return parseInt(splitCommand)

	case util.Plus:
		// string
		return parseSimpleString(splitCommand)
	}

	panic("no command found")

}

func parseSimpleString(command [][]byte) (string, int, error) {

	str := string(command[0][1:])
	return str, 0, nil
}

func parseBulkString(command [][]byte) (string, int, error) {
	stringSize, _ := strconv.Atoi(string(command[0][1:]))

	if stringSize == -1 {
		return "", 0, nil
	}

	str := string(command[1])
	if stringSize != len(str) {
		return "", -1, fmt.Errorf("string length %d does not match declared length %d", len(str), stringSize)
	}

	return string(command[1]), 1, nil

}

func parseInt(command [][]byte) (int, int, error) {
	output, err := strconv.Atoi(string(command[0][1:]))
	if err != nil {
		panic(err)
	}

	return output, 0, nil
}

func parseArray(splitCommand [][]byte) ([]any, int, error) {
	arraySize, _ := strconv.Atoi(string(splitCommand[0][1:]))

	if arraySize == -1 {
		return []any{nil}, 0, nil
	}

	var arr = make([]any, arraySize)
	index := 0

	// add the number of extra CRLFs that have been parsed by each command
	// eg INT: ":1\r\n\" has 0 EXTRA CRLFs
	// eg BULK STRING "$5\r\nhello\r\n" has 1 EXTRA CRLF
	// eg ARRAY "*2\r\n*2\r\n:1\r\n:3\r\n$5\r\nhello\r\n" Nested array adds 2 EXTRA CRLFs
	for i := 1; i <= arraySize; i++ {

		if len(splitCommand) <= i+index || len(splitCommand[i+index:][0]) == 0 {
			return []any{nil}, -1, fmt.Errorf("mismatch between array size and length of command")
		}

		val, idx, err := ParseRESP(splitCommand[i+index:])
		if err != nil {
			return []any{nil}, -1, err
		}
		arr[i-1] = val
		index += idx
	}

	return arr, arraySize + index, nil
}

func serializeArray(arr []any) (string, error) {
	stringBuilder := fmt.Sprintf("%s%d%s", string(util.Star), len(arr), util.CRLF)
	for _, val := range arr {
		str, err := Serialize(val)
		if err != nil {
			return "", err
		}
		stringBuilder += str
	}

	return stringBuilder, nil
}

func Serialize(val any) (string, error) {

	switch v := val.(type) {
	case int:
		return fmt.Sprintf("%s%d%s", string(util.Colon), v, util.CRLF), nil
	case string:
		return fmt.Sprintf("%s%d%s%s%s", string(util.DollarSign), len(v), util.CRLF, v, util.CRLF), nil
	case []any:
		return serializeArray((val).([]any))
	case nil:
		return fmt.Sprintf("%s%d%s", string(util.DollarSign), -1, util.CRLF), nil
	default:
		return "", fmt.Errorf("unexpected type %v", v)
	}

}
