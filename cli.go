package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

var NoCommandGiven = errors.New("No command given")

type SubCommands map[string]*flag.FlagSet

type CLI struct {
	File    *os.File
	Entries TaskList
	OutMsg  string
	Actions SubCommands
	Args    []string
	fs      *flag.FlagSet
	writer  io.Writer
}

func NewCLI(path string, args []string, w io.Writer) (*CLI, error) {
	cli := new(CLI)
	cli.Args = args

	cli.fs = flag.NewFlagSet("task-tracker", flag.ExitOnError)
	cli.writer = w
	cli.setUsageFunc(cli.writer)

	if len(cli.Args) < 2 {
		cli.fs.Usage()
		return nil, NoCommandGiven
	}

	cli.createSubCommands()
	err := cli.constructEntriesFromFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed constructing entries from file, %w", err)
	}

	return cli, nil
}

func (cli *CLI) Act() error {
	var err error

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
		cli.Entries, cli.OutMsg, err = Delete(deleteSet.Args(), cli.Entries)
		if err != nil {
			return fmt.Errorf("error deleting %s, %w", deleteSet.Arg(0), err)
		}
	case "update":
		updateSet := cli.Actions["update"]
		updateSet.Parse(cli.Args[2:])
		if updateSet.NArg() != 2 {
			updateSet.Usage()
		}
		cli.Entries, cli.OutMsg, err = Update(updateSet.Args(), cli.Entries)
		if err != nil {
			return fmt.Errorf("error updating %s, %w", updateSet.Arg(0), err)
		}
	default:
		cli.fs.Usage()
	}
	return nil
}

func (cli *CLI) Update() error {
	err := cli.File.Truncate(0)
	if err = checkIOError(err); err != nil {
		return fmt.Errorf("error truncating file %w", err)
	}
	_, err = cli.File.Seek(0, io.SeekStart)
	if err = checkIOError(err); err != nil {
		return fmt.Errorf("error seeking to file start %w", err)
	}

	err = json.NewEncoder(cli.File).Encode(cli.Entries)
	if err = checkJSONError(err); err != nil {
		return fmt.Errorf("error encoding tasks as JSON %w", err)
	}
	return nil
}

func (cli *CLI) Finish() {
	fmt.Fprintln(cli.writer, cli.OutMsg)
	cli.File.Close()
}

func (cli *CLI) createSubCommands() {
	cli.Actions = make(SubCommands)

	cli.Actions["add"] = flag.NewFlagSet("add", flag.ExitOnError)
	cli.Actions["list"] = flag.NewFlagSet("list", flag.ExitOnError)
	cli.Actions["delete"] = flag.NewFlagSet("delete", flag.ExitOnError)
	cli.Actions["update"] = flag.NewFlagSet("update", flag.ExitOnError)
}

func (cli *CLI) constructEntriesFromFile(path string) error {
	var err error

	cli.File, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err = checkIOError(err); err != nil {
		return fmt.Errorf("error opening file %s, %w", path, err)
	}

	fileInfo, err := cli.File.Stat()
	if err = checkIOError(err); err != nil {
		cli.File.Close()
		return fmt.Errorf("error getting file info, %w", err)
	}

	if fileInfo.Size() != 0 {
		err = json.NewDecoder(cli.File).Decode(&cli.Entries)
		if err = checkJSONError(err); err != nil {
			cli.File.Close()
			return fmt.Errorf("error decoding JSON entries, %w", err)
		}
	}
	return nil
}

func (cli *CLI) setUsageFunc(w io.Writer) {
	cli.fs.SetOutput(w)
	cli.fs.Usage = func() {
		fmt.Fprintf(w, "Usage of %s:", cli.Args[0])
	}
}

func checkIOError(err error) error {
	if err != nil {
		return fmt.Errorf("i/o error %w", err)
	}
	return nil
}

func checkConversionError(err error) error {
	if err != nil {
		return fmt.Errorf("string conversion error %w", err)
	}
	return nil
}

func checkJSONError(err error) error {
	if err != nil {
		return fmt.Errorf("json error %w", err)
	}
	return nil
}
