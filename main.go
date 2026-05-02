package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
)

const FileStore = "task_store.json"

type TaskStatus int

const (
	ToDo TaskStatus = iota
	InProgress
	Done
)

var StatusStrings = map[TaskStatus]string{
	ToDo:       "todo",
	InProgress: "in progress",
	Done:       "done",
}

func (t TaskStatus) String() string {
	return StatusStrings[t]
}

type Task struct {
	Id          int
	Description string
	Status      TaskStatus
}

func (t Task) String() string {
	return fmt.Sprintf("%4d %40s %11s", t.Id, t.Description, t.Status)
}

func main() {
	// Set subcommands for CLI
	cmdAdd := flag.NewFlagSet("add", flag.ExitOnError)
	cmdList := flag.NewFlagSet("list", flag.ExitOnError)
	cmdDelete := flag.NewFlagSet("delete", flag.ExitOnError)
	cmdUpdate := flag.NewFlagSet("update", flag.ExitOnError)

	// declare container for tasks
	var entries []Task

	// Open and check size of file
	fd, err := os.OpenFile(FileStore, os.O_RDWR|os.O_CREATE, 0666)
	checkIOError(err)
	defer fd.Close()

	fdInfo, err := fd.Stat()
	checkIOError(err)

	// Skip decoding for empty files.
	if fdInfo.Size() != 0 {
		err = json.NewDecoder(fd).Decode(&entries)
		checkJSONError(err)
	}

	if len(os.Args) < 2 {
		showHelpAndExit()
	}

	// input dependent action
	switch os.Args[1] {
	case "add":
		cmdAdd.Parse(os.Args[2:])
		for _, desc := range cmdAdd.Args() {
			task := Task{
				Id:          len(entries),
				Description: desc,
				Status:      ToDo,
			}
			entries = append(entries, task)
		}
		fmt.Printf("Task Added Successfully (ID: %d)\n", len(entries))
	case "list":
		cmdList.Parse(os.Args[2:])
		for _, t := range entries {
			fmt.Printf("%s\n", t)
		}
	case "delete":
		cmdDelete.Parse(os.Args[2:])
		id, err := strconv.Atoi(cmdDelete.Arg(0))
		checkConversionError(err)
		for i, task := range entries {
			if task.Id == id {
				entries = slices.Delete(entries, i, i+1)
				break
			}
		}
	case "update":
		cmdUpdate.Parse(os.Args[2:])
		id, err := strconv.Atoi(cmdUpdate.Arg(0))
		checkConversionError(err)
		for _, task := range entries {
			if task.Id == id {
				task.Description = cmdUpdate.Arg(1)
				break
			}
		}
	default:
		showHelpAndExit()
	}

	// truncate to avoid invalid data on underwrites
	err = fd.Truncate(0)
	checkIOError(err)

	// save and quit
	err = json.NewEncoder(fd).Encode(entries)
	checkJSONError(err)
	os.Exit(0)
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
