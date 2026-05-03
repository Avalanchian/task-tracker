package main

const FileStore = "task_store.json"

func main() {
	cli := NewCLI(FileStore)
	defer cli.Finish()

	cli.Act()
	cli.Update()
}
