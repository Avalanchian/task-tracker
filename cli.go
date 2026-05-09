package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

// StoragePath determines the file path of the persistent JSON storage for the application.
const StoragePath = "task_store.json"

// FlagStore holds bools that are set when flags are passed.
type FlagStore struct {
	todo, inProgress, done bool
}

// CLI stores and coordinates the data associated with the application.
type CLI struct {
	Entries  []Task
	Args     []string
	Options  FlagStore
	Commands map[string]*flag.FlagSet
	Writer   io.Writer
	OutStr   string
	flagset  *flag.FlagSet
}

// NewCLI creates and initializes a CLI, including the retrieval of tasks from the
// JSON store.
func NewCLI(w io.Writer) (*CLI, error) {
	entries, err := getEntriesFromJSON(StoragePath)
	if err != nil {
		return nil, fmt.Errorf("Failed creating new CLI, %w", err)
	}

	cli := &CLI{
		Args:     os.Args,
		Entries:  entries,
		Commands: createSubCommands(),
		Writer:   w,
		flagset:  flag.NewFlagSet("task-tracker", flag.ContinueOnError),
	}
	cli.setupFlags()

	return cli, nil
}

// Run handles the provided sub-command and selects the appropriate command function to
// transform the entry list. All sub-commands can produce parsing errors, but some
// commands, such as delete, mark, and update can return *strconv.NumError.
func (cli *CLI) Run() error {
	var err error

	if len(cli.Args) < 2 {
		return err
	}

	switch cli.Args[1] {
	case "add":
		err = cli.Commands["add"].Parse(cli.Args[2:])
		if err != nil {
			return fmt.Errorf("error parsing add command, %w", err)
		}
		cli.Entries, cli.OutStr = Add(cli.Commands["add"].Args(), cli.Entries)
	case "list":
		err = cli.Commands["list"].Parse(cli.Args[2:])
		if err != nil {
			return fmt.Errorf("error parsing list command, %w", err)
		}
		cli.OutStr = List(cli.Options, cli.Entries)
	case "delete":
		err = cli.Commands["delete"].Parse(cli.Args[2:])
		if err != nil {
			return fmt.Errorf("error parsing delete command, %w", err)
		}
		cli.Entries, cli.OutStr, err = Delete(cli.Commands["delete"].Args(), cli.Entries)
		if err != nil {
			return fmt.Errorf("could not complete delete command, %w", err)
		}
	case "mark":
		err = cli.Commands["mark"].Parse(cli.Args[2:])
		if err != nil {
			return fmt.Errorf("error parsing mark command, %w", err)
		}
		cli.Entries, cli.OutStr, err = Mark(cli.Commands["mark"].Args(), cli.Options, cli.Entries)
		if err != nil {
			return fmt.Errorf("could not complete mark command, %w", err)
		}
	case "update":
		err = cli.Commands["update"].Parse(cli.Args[2:])
		if err != nil {
			return fmt.Errorf("error parsing update command, %w", err)
		}
		cli.Entries, cli.OutStr, err = Update(cli.Commands["update"].Args(), cli.Entries)
		if err != nil {
			return fmt.Errorf("could not complete update command, %w", err)
		}
	}

	fmt.Fprintf(cli.Writer, cli.OutStr)
	if len(cli.OutStr) != 0 {
		fmt.Fprintf(cli.Writer, "\n")
	}

	return err
}

// Save writes cli.Entries as a JSON string in persistent storage. Initially it saves to a
// swap file and then calles os.Rename to overwrite, providing some assurance against data
// loss.
func (cli *CLI) Save() error {
	fd, err := os.Create(StoragePath + "swap")
	if err != nil {
		return fmt.Errorf("could not create temporary storage, %w", err)
	}

	err = json.NewEncoder(fd).Encode(cli.Entries)
	if err != nil {
		fd.Close()
		return fmt.Errorf("error encoding entries as JSON, %w", err)
	}

	fd.Close()
	err = os.Rename(fd.Name(), StoragePath)
	if err != nil {
		return fmt.Errorf("error renaming temporary storage, %w", err)
	}

	return nil
}

func getEntriesFromJSON(path string) ([]Task, error) {
	fd, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("i/o error opening JSON store, %w", err)
	}
	defer fd.Close()

	fdInfo, err := fd.Stat()
	if err != nil {
		return nil, fmt.Errorf("error getting file info, %w", err)
	}

	var entries []Task
	if fdInfo.Size() != 0 {
		err = json.NewDecoder(fd).Decode(&entries)
		if err != nil {
			return nil, fmt.Errorf("JSON decoding error, %w", err)
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
