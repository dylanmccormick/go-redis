package resp

import (
	"fmt"
	"strconv"

	"github.com/dylanmccormick/go-redis/internal/util"
)

func ParseRESP(splitCommand [][]byte) (any, int) {

	switch splitCommand[0][0] {

	case util.Star:
		return parseArray(splitCommand)
		// fmt.Println("ARRAY:")
		// length := splitCommand[1][0]
		// idx, err := parseArray(data[1:])
		// if err != nil {
		// 	return -1, err
		// }
		//
		// return idx + 1, nil

	case util.DollarSign:
		return parseBulkString(splitCommand)

	case util.Dash:
		return parseBulkString(splitCommand)

	case util.Colon:
		return parseInt(splitCommand)

	case util.Plus:
		return parseSimpleString(splitCommand)
	}

	panic("no command found")

}

func parseSimpleString(command [][]byte) (string, int) {

	str := string(command[0][1:])
	return str, 0
}

func parseBulkString(command [][]byte) (string, int) {
	stringSize, _ := strconv.Atoi(string(command[0][1:]))

	if stringSize == -1 {
		return "", 0
	}

	return string(command[1]), 1

}

func parseInt(command [][]byte) (int, int) {
	output, err := strconv.Atoi(string(command[0][1:]))
	if err != nil {
		panic(err)
	}

	return output, 0

}

func parseArray(splitCommand [][]byte) ([]any, int) {
	arraySize, _ := strconv.Atoi(string(splitCommand[0][1:]))

	var arr = make([]any, arraySize)
	index := 0
	for i := 1; i <= arraySize; i++ {
		fmt.Printf("%s\n", splitCommand[i])
		val, idx := ParseRESP(splitCommand[i+index:])
		arr[i-1] = val
		index += idx
	}

	fmt.Printf("array: %v", arr)

	return arr, arraySize + index
}

func serializeArray(arr []any) string {
	stringBuilder := fmt.Sprintf("%s%d", string(util.Star), len(arr))
	for _, val := range arr {
		stringBuilder += Serialize(val)
	}

	return stringBuilder
}

func Serialize(val any) string {

	switch v := val.(type) {
	case int:
		fmt.Printf("val is an int: %d\n", v)
		return fmt.Sprintf("%s%d%s", string(util.Colon), v, util.CRLF)
	case string:
		fmt.Printf("val is an string: %s\n", v)
		return fmt.Sprintf("%s%d%s%s", string(util.DollarSign), len(v), v, util.CRLF)
	case []any:
		fmt.Printf("val is an array: %v\n", v)
		return serializeArray((val).([]any))
	}

	return ""
}
