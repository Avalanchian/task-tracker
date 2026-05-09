package main

import (
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
