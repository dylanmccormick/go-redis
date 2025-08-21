// Package server creates a redis-lite server
package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/dylanmccormick/go-redis/internal/cmd"
	"github.com/dylanmccormick/go-redis/internal/database"
	"github.com/dylanmccormick/go-redis/internal/util"
)

type ServerConfig struct {
	Port int
	Database *database.Database
}

func StartServer(c ServerConfig) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(c.Port))
	log.Printf("Server is listening on port :%s\n", strconv.Itoa(c.Port))
	if err != nil {
		log.Fatalf("A listener error occurred: %s", err)
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("A listener error occurred: %s", err)
		}

		go c.handleRequest(conn)
	}
}

func(c *ServerConfig) handleRequest(conn net.Conn) {
	scanner := bufio.NewReader(conn)
	defer conn.Close()
	log.Printf("Handling connection")

	for {
		buffer := make([]byte, 4096)
		i, err := scanner.Read(buffer)
		fmt.Printf("Read %d bytes\n", i)
		if err != nil {
			log.Fatalf("Error reading from request")
		}
		buffer = util.ClearZeros(buffer)
		log.Printf("clean read into buffer %#v\n", string(buffer))
		response, err := cmd.HandleMessage(c.Database, buffer)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Println(response)
	}
}

func (c *ServerConfig) Shell() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("localhost:%s> ", strconv.Itoa(c.Port))
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Printf("read input: %s\n", input)
		strArr := strings.Split(strings.Trim(input, "\n"), " ")
		resp , err:= cmd.HandleCommand(c.Database, strArr)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		fmt.Println(resp)
	}
}
