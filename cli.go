package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"

	"github.com/adrg/xdg"
)

// StoragePath determines the file path of the persistent JSON storage for the application.
// It is set on application initialisation.
var StoragePath string

func init() {
	err := os.MkdirAll(filepath.Join(xdg.DataHome, "task-tracker"), 0777)
	if err != nil {
		log.Fatalf("Could not create storage directory, %v", err)
	}
	StoragePath = filepath.Join(xdg.DataHome, "task-tracker", "tasks.json")
}

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
	Store    string
	OutStr   string
	flagset  *flag.FlagSet
}

// NewCLI creates and initializes a CLI, including the retrieval of tasks from the
// JSON store.
func NewCLI(args []string, w io.Writer, store string) (*CLI, error) {
	entries, err := getEntriesFromJSON(store)
	if err != nil {
		return nil, fmt.Errorf("Failed creating new CLI, %w", err)
	}

	cli := &CLI{
		Args:     args,
		Entries:  entries,
		Commands: createSubCommands(w),
		Writer:   w,
		Store:    store,
		flagset:  flag.NewFlagSet(args[0], flag.ContinueOnError),
	}
	cli.setupFlags()
	cli.flagset.Usage = func() {
		fmt.Fprintf(cli.Writer, "Usage: %s <command> [options] args\n\n", args[0])
		fmt.Fprintf(cli.Writer, "Commands:\n")

		var commandNames []string
		for key, _ := range cli.Commands {
			commandNames = append(commandNames, key)
		}

		slices.Sort(commandNames)
		for _, name := range commandNames {
			fmt.Fprintf(cli.Writer, "\t%s\n", name)
		}
	}

	return cli, nil
}

// Run handles the provided sub-command and selects the appropriate command function to
// transform the entry list. All sub-commands can produce parsing errors, but some
// commands, such as delete, mark, and update can return *strconv.NumError.
func (cli *CLI) Run() error {
	var err error

	if len(cli.Args) < 2 {
		return flag.ErrHelp
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
	default:
		return flag.ErrHelp
	}

	fmt.Fprintf(cli.Writer, "%s", cli.OutStr)
	if len(cli.OutStr) != 0 {
		fmt.Fprintf(cli.Writer, "\n")
	}

	return err
}

// Save writes cli.Entries as a JSON string in persistent storage. Initially it saves to a
// swap file and then calles os.Rename to overwrite, providing some assurance against data
// loss.
func (cli *CLI) Save() error {
	fd, err := os.Create(cli.Store + ".swp")
	if err != nil {
		return fmt.Errorf("could not create temporary storage, %w", err)
	}

	err = json.NewEncoder(fd).Encode(cli.Entries)
	if err != nil {
		fd.Close()
		return fmt.Errorf("error encoding entries as JSON, %w", err)
	}

	fd.Close()
	err = os.Rename(fd.Name(), cli.Store)
	if err != nil {
		return fmt.Errorf("error renaming temporary storage, %w", err)
	}

	return nil
}

func getEntriesFromJSON(path string) ([]Task, error) {
	fd, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
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

func createSubCommands(w io.Writer) map[string]*flag.FlagSet {
	commands := map[string]*flag.FlagSet{
		"add":    createAddCommand(w),
		"list":   createListCommand(w),
		"update": createUpdateCommand(w),
		"delete": createDeleteCommand(w),
		"mark":   createMarkCommand(w),
	}

	return commands
}

func (cli *CLI) setupFlags() {
	cli.Commands["list"].BoolVar(&cli.Options.todo, "todo", false, "list only outstanding tasks")
	cli.Commands["list"].BoolVar(&cli.Options.inProgress, "inprogress", false, "list tasks that are currently in progress")
	cli.Commands["list"].BoolVar(&cli.Options.done, "done", false, "list completed tasks")

	cli.Commands["mark"].BoolVar(&cli.Options.todo, "todo", false, "mark a task with the given ID as todo")
	cli.Commands["mark"].BoolVar(&cli.Options.inProgress, "inprogress", false, "mark a task with the given ID as in progress")
	cli.Commands["mark"].BoolVar(&cli.Options.done, "done", false, "mark a task with the given ID as done")
}

func createAddCommand(w io.Writer) (command *flag.FlagSet) {
	command = flag.NewFlagSet("add", flag.ExitOnError)
	command.Usage = func() {
		fmt.Fprintf(w, "Usage: todo add <description>\n\n")
		fmt.Fprintf(w, "Adds a new task with given description\n")
	}
	return
}

func createListCommand(w io.Writer) (command *flag.FlagSet) {
	command = flag.NewFlagSet("list", flag.ExitOnError)
	command.Usage = func() {
		fmt.Fprintf(w, "Usage: todo list [options]\n\n")
		fmt.Fprintf(w, "Lists all tasks, filtered by an option (if provided).\n\n")
		fmt.Fprintf(w, "Options:\n")
		command.PrintDefaults()
	}
	return
}

func createMarkCommand(w io.Writer) (command *flag.FlagSet) {
	command = flag.NewFlagSet("mark", flag.ExitOnError)
	command.Usage = func() {
		fmt.Fprintf(w, "Usage: todo mark [options] <ID>\n\n")
		fmt.Fprintf(w, "Marks a task with the given ID as in progress/done.\n")
		fmt.Fprintf(w, "A specific status can be chosen by supplying an option.\n\n")
		fmt.Fprintf(w, "Options:\n")
		command.PrintDefaults()
	}
	return
}

func createDeleteCommand(w io.Writer) (command *flag.FlagSet) {
	command = flag.NewFlagSet("delete", flag.ExitOnError)
	command.Usage = func() {
		fmt.Fprintf(w, "Usage: todo delete <ID>\n\n")
		fmt.Fprintf(w, "deletes the task associated with the provided ID\n")
	}
	return
}

func createUpdateCommand(w io.Writer) (command *flag.FlagSet) {
	command = flag.NewFlagSet("update", flag.ExitOnError)
	command.Usage = func() {
		fmt.Fprintf(w, "Usage: todo update <ID> <description>\n\n")
		fmt.Fprintf(w, "Updates a task with the given ID with a new description.\n")
	}
	return
}
