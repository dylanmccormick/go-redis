package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"

	resp "github.com/dylanmccormick/go-redis/internal/resp"
	"github.com/dylanmccormick/go-redis/internal/util"
)

func main() {

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

	switch os.Args[1] {
	case "PING":
		pingCommand.Parse(os.Args[2:])

		fmt.Println("PONG")
	default:
		// Ok here we're going to split into command and arguments

		input := os.Args[1]
		fmt.Printf("%#v\n", input)
		splitCommand := bytes.Split([]byte(input), util.SeparatorBytes)
		array, size := resp.ParseRESP(splitCommand)
		fmt.Printf("Command: %v\t Size: %v\n", array, size)
		v := resp.Serialize(array)
		fmt.Printf("%#v", v)

	}

}
