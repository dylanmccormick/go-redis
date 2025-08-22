package main

import (
	"flag"
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
	config := &Config{
		Db:       db,
		LogLevel: INFO,
		Port:     42069,
	}

	pingCommand := flag.NewFlagSet("PING", flag.ExitOnError)

	var verbose bool
	var help bool
	var version bool

	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&help, "h", false, "help message output")
	flag.BoolVar(&version, "version", false, "CLI version")
	pingCommand.BoolVar(&verbose, "verbose", false, "verbose output")
	pingCommand.BoolVar(&verbose, "v", false, "verbose output")

	flag.Parse()

	if verbose {
		fmt.Println("Verbose selected")
	}

	if help {
		fmt.Println("help selected")
	}

	if version {
		fmt.Println("version selected")
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
	var wg sync.WaitGroup
	sc := server.ServerConfig{Port: c.Port, Database: c.Db}

	go server.StartServer(sc, &wg)

	return &sc
}
