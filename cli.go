package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

type SubCommands map[string]*flag.FlagSet

type CLI struct {
	File    *os.File
	Entries TaskList
	OutMsg  string
	Actions SubCommands
	Args    []string
	fs      *flag.FlagSet
}

func NewCLI(path string, args []string, w io.Writer) *CLI {
	cli := new(CLI)
	cli.Args = args

	cli.fs = flag.NewFlagSet("task-tracker", flag.ExitOnError)
	cli.setUsageFunc(w)

	if len(cli.Args) < 2 {
		cli.fs.Usage()
	}

	cli.createSubCommands()
	cli.constructEntriesFromFile(path)

	return cli
}

func (cli *CLI) Act() {
	switch cli.Args[1] {
	case "add":
		addSet := cli.Actions["add"]
		addSet.Parse(cli.Args[2:])
		cli.Entries, cli.OutMsg = Add(addSet.Args(), cli.Entries)
	case "list":
		listSet := cli.Actions["list"]
		listSet.Parse(cli.Args[2:])
		cli.OutMsg = List(cli.Entries)
	case "delete":
		deleteSet := cli.Actions["delete"]
		deleteSet.Parse(cli.Args[2:])
		cli.Entries, cli.OutMsg = Delete(deleteSet.Args(), cli.Entries)
	case "update":
		updateSet := cli.Actions["update"]
		updateSet.Parse(cli.Args[2:])
		if updateSet.NArg() != 2 {
			updateSet.Usage()
		}
		cli.Entries, cli.OutMsg = Update(updateSet.Args(), cli.Entries)
	default:
		cli.fs.Usage()
	}
}

func (cli *CLI) Update() {
	err := cli.File.Truncate(0)
	checkIOError(err)
	_, err = cli.File.Seek(0, io.SeekStart)
	checkIOError(err)

	err = json.NewEncoder(cli.File).Encode(cli.Entries)
	checkJSONError(err)
}

func (cli *CLI) Finish() {
	fmt.Println(cli.OutMsg)
	cli.File.Close()
}

func (cli *CLI) createSubCommands() {
	cli.Actions = make(SubCommands)

	cli.Actions["add"] = flag.NewFlagSet("add", flag.ExitOnError)
	cli.Actions["list"] = flag.NewFlagSet("list", flag.ExitOnError)
	cli.Actions["delete"] = flag.NewFlagSet("delete", flag.ExitOnError)
	cli.Actions["update"] = flag.NewFlagSet("update", flag.ExitOnError)
}

func (cli *CLI) constructEntriesFromFile(path string) {
	var err error

	cli.File, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	checkIOError(err)

	fileInfo, err := cli.File.Stat()
	checkIOError(err)

	if fileInfo.Size() != 0 {
		err = json.NewDecoder(cli.File).Decode(&cli.Entries)
		checkJSONError(err)
	}
}

func (cli *CLI) setUsageFunc(w io.Writer) {
	cli.fs.SetOutput(w)
	cli.fs.Usage = func() {
		fmt.Fprintf(w, "Usage of %s:", cli.Args[0])
	}
}

func checkIOError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "i/o error %v", err)
	}
}

func checkConversionError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "string conversion error %v", err)
	}
}

func checkJSONError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "json error %v", err)
	}
}
