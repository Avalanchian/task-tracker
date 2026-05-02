package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

const FileStore = "task_store.json"

func main() {
	if len(os.Args) < 2 {
		showHelpAndExit()
	}

	// Set subcommands for CLI
	cmdAdd := flag.NewFlagSet("add", flag.ExitOnError)
	cmdList := flag.NewFlagSet("list", flag.ExitOnError)
	cmdDelete := flag.NewFlagSet("delete", flag.ExitOnError)
	cmdUpdate := flag.NewFlagSet("update", flag.ExitOnError)

	// Open and check size of file
	fd, err := os.OpenFile(FileStore, os.O_RDWR|os.O_CREATE, 0666)
	checkIOError(err)
	defer fd.Close()

	fdInfo, err := fd.Stat()
	checkIOError(err)

	// declare container for tasks
	var entries TaskList

	// Skip decoding for empty files.
	if fdInfo.Size() != 0 {
		err = json.NewDecoder(fd).Decode(&entries)
		checkJSONError(err)
	}

	// input dependent action
	var outMsg string
	switch os.Args[1] {
	case "add":
		cmdAdd.Parse(os.Args[2:])
		entries, outMsg = Add(cmdAdd.Args(), entries)
	case "list":
		cmdList.Parse(os.Args[2:])
		outMsg = List(entries)
	case "delete":
		cmdDelete.Parse(os.Args[2:])
		entries, outMsg = Delete(cmdDelete.Args(), entries)
	case "update":
		cmdUpdate.Parse(os.Args[2:])
		if cmdUpdate.NArg() != 2 {
			cmdUpdate.Usage()
			os.Exit(2)
		}
		entries, outMsg = Update(cmdUpdate.Args(), entries)
	default:
		showHelpAndExit()
	}

	// truncate to avoid invalid data on underwrites
	err = fd.Truncate(0)
	checkIOError(err)
	_, err = fd.Seek(0, io.SeekStart)
	checkIOError(err)

	// save and quit
	err = json.NewEncoder(fd).Encode(entries)
	checkJSONError(err)

	fmt.Println(outMsg)
}

func showHelpAndExit() {
	flag.Usage()
	os.Exit(2)
}

func checkIOError(err error) {
	if err != nil {
		fmt.Errorf("i/o error %v", err)
	}
}

func checkConversionError(err error) {
	if err != nil {
		fmt.Errorf("string conversion error %v", err)
	}
}

func checkJSONError(err error) {
	if err != nil {
		fmt.Errorf("json error %v", err)
	}
}
