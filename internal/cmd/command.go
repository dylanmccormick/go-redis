// Package cmd handles all commands for the redis-lite server
package cmd

import (
	"bytes"
	"fmt"
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
	response := fmt.Sprintf("%#v", request)

	return response, nil
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


