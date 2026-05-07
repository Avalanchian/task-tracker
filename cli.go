package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

const StoragePath = "task_store.json"

type FlagStore struct {
	todo, inProgress, done bool
}

type CLI struct {
	Entries  []Task
	Args     []string
	Options  FlagStore
	Commands map[string]*flag.FlagSet
	Writer   io.Writer
	flagset  *flag.FlagSet
}

func NewCLI(w io.Writer) (*CLI, error) {
	entries, err := getEntriesFromJSON(StoragePath)
	if err != nil {
		return fmt.Errorf("Failed creating new CLI, %w", err)
	}

	cli := &CLI{
		Args:     os.Args,
		Entries:  entries,
		Commands: createSubCommands(),
		Writer:   w,
		flagset:  flag.NewFlagSet(),
	}
	cli.setupFlags()

	return cli
}

func getEntriesFromJSON(path string) ([]Task, error) {
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("i/o error opening JSON store, %w", err)
	}

	fdInfo, err = fd.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting file info, %w", err)
	}

	var entries []Task
	if fdInfo.Size() != 0 {
		err = json.NewDecoder(fd).Decode(&entries)
		if err != nil {
			return nil, fmt.Errorf("JSON decoding error, %w")
		}
	}

	return entries, nil
}

func createSubCommands() map[string]*flag.FlagSet {
	addCmd := flag.NewFlagSet("add", flag.ContinueOnError)
	listCmd := flag.NewFlagSet("list", flag.ContinueOnError)
	updateCmd := flag.NewFlagSet("update", flag.ContinueOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ContinueOnError)
	markCmd := flag.NewFlagSet("mark", flag.ContinueOnError)

	commands := map[string]*flag.FlagSet{
		"add":    addCmd,
		"list":   listCmd,
		"update": updateCmd,
		"delete": deleteCmd,
		"mark":   markCmd,
	}

	return commands
}

func (cli *CLI) setupFlags() {
	cli.Commands["list"].BoolVar(&cli.Options.todo, "todo", false, "todo")
	cli.Commands["list"].BoolVar(&cli.Options.inProgress, "inprogress", false, "in progress")
	cli.Commands["list"].BoolVar(&cli.Options.done, "done", false, "done")

	cli.Commands["mark"].BoolVar(&cli.Options.todo, "todo", false, "todo")
	cli.Commands["mark"].BoolVar(&cli.Options.inProgress, "inprogress", false, "in progress")
	cli.Commands["mark"].BoolVar(&cli.Options.done, "done", false, "done")
}
