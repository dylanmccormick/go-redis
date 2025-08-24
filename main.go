package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/dylanmccormick/go-redis/internal/cmd"
	"github.com/dylanmccormick/go-redis/internal/database"
	"github.com/dylanmccormick/go-redis/internal/server"
)

type Config struct {
	Db       *database.Database
	LogLevel LogLevel
	Port     int
}

type LogLevel int

const (
	SILENT LogLevel = iota
	INFO
	ERROR
	DEBUG
)

func main() {
	db := database.InitializeDB()
	cmd.Load(db)
	defer cmd.Save(db)
	config := &Config{
		Db:       db,
		LogLevel: INFO,
		Port:     6379,
	}
	firstArg := ""
	if len(os.Args) > 1 {
		firstArg = os.Args[1]
	}

	switch firstArg {
	case "start":
		sc := config.serveInteractive()
		fmt.Println("interactive shell later")
		sc.Shell()
	case "":
		fmt.Println("Starting server")
		config.serve()
	default:
		response, err := cmd.HandleCommand(config.Db, os.Args[1:])
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Println(response)
	}
}

func (c *Config) serve() *server.ServerConfig {
	var wg sync.WaitGroup
	wg.Add(1)
	sc := server.ServerConfig{Port: c.Port, Database: c.Db}

	go server.StartServer(sc, &wg)

	wg.Wait()
	return &sc
}

func (c *Config) serveInteractive() *server.ServerConfig {
	// TODO: Figure out how to make this the same method as the other one. I think I'm using wait groups incorrectly or something. 
	// Or this needs to have a wait group and be run as a goroutine. 
	var wg sync.WaitGroup
	sc := server.ServerConfig{Port: c.Port, Database: c.Db}

	go server.StartServer(sc, &wg)

	return &sc
}
