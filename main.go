package main

import (
	"errors"
	"flag"
	"log"
	"os"
)

func main() {
	cli, err := NewCLI(os.Args, os.Stderr, StoragePath)
	if err != nil {
		log.Printf("failed to create new CLI, %v\n", err)
		return
	}

	err = cli.Run()
	if errors.Is(err, flag.ErrHelp) {
		cli.flagset.Usage()
		return
	}
	if err != nil {
		log.Printf("error while running CLI, %v\n", err)
		return
	}

	err = cli.Save()
	if err != nil {
		log.Printf("error while saving data, %v\n", err)
		return
	}
}
