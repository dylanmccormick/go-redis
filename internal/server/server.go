// Package server creates a redis-lite server
package server

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/dylanmccormick/go-redis/internal/cmd"
	"github.com/dylanmccormick/go-redis/internal/database"
	"github.com/dylanmccormick/go-redis/internal/util"
)

type ServerConfig struct {
	Port     int
	Database *database.Database
}

func StartServer(c ServerConfig, wg *sync.WaitGroup) {
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
		log.Printf("Received request")

		go c.handleRequest(conn)
	}
}

func (c *ServerConfig) handleRequest(conn net.Conn) {
	scanner := bufio.NewReader(conn)
	defer conn.Close()
	log.Printf("Handling connection")

	for {
		buffer := make([]byte, 4096)
		i, err := scanner.Read(buffer)
		fmt.Printf("Read %d bytes\n", i)
		if err != nil {
			if err == io.EOF {
				fmt.Println("End of connection closed gracefully")
				return
			} else {
				fmt.Printf("Unknown error from TCP request: %s", err)
			}
			return
		}
		buffer = util.ClearZeros(buffer)
		log.Printf("clean read into buffer %#v\n", string(buffer))
		response, err := cmd.HandleMessage(c.Database, buffer)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Println(response)
		fmt.Fprintf(conn, "+%s\r\n", response)
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
		resp, err := cmd.HandleCommand(c.Database, strArr)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		fmt.Println(resp)
	}
}
