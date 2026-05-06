package main

import (
	"errors"
	"log"
	"os"
)

const FileStore = "task_store.json"

func main() {
	cli, err := NewCLI(FileStore, os.Args, os.Stderr)
	if errors.Is(err, NoCommandGiven) {
		return
	} else if err != nil {
		log.Printf("Error initializing CLI, %v", err)
		return
	}
	defer cli.Finish()

	err = cli.Act()
	if err != nil {
		log.Printf("Error performing CLI action, %v", err)
		return
	}

	err = cli.Update()
	if err != nil {
		log.Printf("Error updating CLI, %v", err)
		return
	}
}
