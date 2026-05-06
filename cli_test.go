package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"testing"
)

func TestNewCLI(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "testdata.json")
	fd, err := os.Create(path)
	checkIOError(err)
	fd.Close()

	out := new(bytes.Buffer)
	_, err = NewCLI(path, []string{"task-tracker"}, out)

	if !errors.Is(err, NoCommandGiven) {
		t.Errorf("expected error not present, %v", err)
	}
}

func TestAct(t *testing.T) {
	cases := [][]string{
		{"task-tracker", "add", "a description"},
		{"task-tracker", "list", ""},
		{"task-tracker", "delete", "4"},
		{"task-tracker", "update", "3", "new description"},
	}

	for _, args := range cases {
		t.Run(fmt.Sprintf("test %s parser", args[1]), func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "testdata.json")
			fd, err := os.Create(path)
			checkIOError(err)
			tasks := setupTempEntries(5)
			err = json.NewEncoder(fd).Encode(tasks)
			checkJSONError(err)
			fd.Close()

			out := new(bytes.Buffer)
			cli, _ := NewCLI(path, args, out)

			cli.Act()

			got := cli.Actions[args[1]].Args()
			want := args[2:]

			if !slices.Equal(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}

func TestCLIUpdate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "testdata.json")
	fd, err := os.Create(path)
	checkIOError(err)
	storedTasks := setupTempEntries(5)
	err = json.NewEncoder(fd).Encode(storedTasks)
	checkJSONError(err)
	fd.Close()

	out := new(bytes.Buffer)
	cli, _ := NewCLI(path, []string{"task-tracker", "delete", "3"}, out)

	cli.Entries = setupTempEntries(2)
	cli.Update()

	var got TaskList
	newFile, err := os.Open(path)
	checkIOError(err)
	defer newFile.Close()
	err = json.NewDecoder(newFile).Decode(&got)
	checkJSONError(err)

	if !reflect.DeepEqual(got, cli.Entries) {
		t.Errorf("got %+v, want %+v", got, cli.Entries)
	}
}
