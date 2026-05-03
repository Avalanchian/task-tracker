package main

import "os"

const FileStore = "task_store.json"

func main() {
	cli := NewCLI(FileStore, os.Args, os.Stderr)
	defer cli.Finish()

	cli.Act()
	cli.Update()
}
