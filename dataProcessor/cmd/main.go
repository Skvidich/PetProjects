package main

import (
	"dataProcessor/cmd/cli"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("must set path to ini file as only argument")
	}

	cli.RunCLI(os.Args[1])
}
