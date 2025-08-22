// Package cmd handles all commands for the redis-lite server
package cmd

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/dylanmccormick/go-redis/internal/database"
	"github.com/dylanmccormick/go-redis/internal/resp"
	"github.com/dylanmccormick/go-redis/internal/util"
)

func HandleMessage(db *database.Database, buffer []byte) (string, error) {

	command := bytes.Split(buffer, util.SeparatorBytes)
	fmt.Println(string(command[0]))
	request, _, err := resp.ParseRESP(command)
	if err != nil {
		fmt.Printf("Invalid request found: %s", err)
		return "", err
	}
	return HandleRespRequest(db, request)
}

func HandleRespRequest(db *database.Database, request any) (string, error) {

	switch v := request.(type) {
	case []any:
		return HandleCommand(db, requestToStringArr(v))
		// send to array handlers?
	case string:
		return HandleCommand(db, []string{v})
	case int:
		return "", fmt.Errorf("Not yet implemented. Don't think this should be happening")
	default:
		return "", fmt.Errorf("Unexpected type from resp statement: %T", v)

	}

}

func requestToStringArr(arr []any) ([]string) {

	var output []string

	for _, val := range(arr) {
		switch v := val.(type) {
		case []any:
			fmt.Printf("Got sub array. Figure out how to handle")
			return output

		case string:
			output = append(output, v)

		case int:
			output = append(output, strconv.Itoa(v))

		default:
			fmt.Printf("Got unexpected type: %T\n", v)
			return output
		}

	}
	return output
}

func HandleCommand(db *database.Database, args []string) (string, error) {
	switch strings.ToLower(args[0]) {
	case "ping":
		if len(args) > 2 {
			return "", fmt.Errorf("too many arguments passed to PING command")
		}
		if len(args) > 1 {
			return fmt.Sprintf("\"%s\"", args[1]), nil
		}
		return "PONG", nil
	case "echo":
		if len(args) > 2 {
			return "", fmt.Errorf("too many arguments passed to ECHO command")
		}
		return fmt.Sprintf("\"%s\"", args[1]), nil
	case "set":
		if len(args) != 3 {
			return "", fmt.Errorf("incorrect number of arguments for SET command")
		}
		return handleSet(db, args[1], args[2])
	case "get":
		if len(args) != 2 {
			return "", fmt.Errorf("too many arguments passed to GET command")
		}
		response, err :=  handleGet(db, args[1])
		if err != nil {
			return "", err
		}

		s, ok := response.(string)
		if !ok {
			return "not set up to handle non-string attributes yet", nil
		}

		return s, nil
	case "rpush":
		response, err := handleRPush(db, args[1], args[2])
		if err != nil {
			return "", err
		}
		return response, nil

	case "lpush":
		response, err := handleLPush(db, args[1], args[2])
		if err != nil {
			return "", err
		}
		return response, nil

	case "rpop":
		response, err := handleRPop(db, args[1])
		if err != nil {
			return "", err
		}
		return response, nil

	case "lpop":
		response, err := handleLPop(db, args[1])
		if err != nil {
			return "", err
		}
		return response, nil

	}

	return "", fmt.Errorf("%v is not a valid command", args[0])
}

func handleSet(db *database.Database, key, value string) (string, error) {
	fmt.Printf("Key: %s, Value: %s\n", key, value)

	db.Set(key, value)

	return "OK", nil

}

func handleGet(db *database.Database, key string) (any, error) {

	return db.Get(key)
}


func handleRPush(db *database.Database, key, value string) (string, error) {
	db.RPush(key, value)
	return "OK", nil
}

func handleLPush(db *database.Database, key, value string) (string, error) {
	db.LPush(key, value)
	return "OK", nil
}

func handleRPop(db *database.Database, key string) (string, error) {
	return db.RPop(key)
}

func handleLPop(db *database.Database, key string) (string, error) {
	return db.LPop(key)
}

