package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

func TestCLI(t *testing.T) {
	t.Run("add to empty list", func(t *testing.T) {
		buf := new(bytes.Buffer)
		args := []string{"task-tracker", "add", "a test description"}

		tmpDir := t.TempDir()
		fd, err := os.CreateTemp(tmpDir, "")
		if err != nil {
			t.Fatalf("failed to create temp file store, %v", err)
		}
		fd.Close()

		cli := runCLI(t, args, buf, fd.Name())

		if len(cli.Entries) != 1 {
			t.Errorf("failed to add, got %d, want 1", len(cli.Entries))
		}
	})

	t.Run("add to populated list", func(t *testing.T) {
		buf := new(bytes.Buffer)
		args := []string{"task-tracker", "add", "a test description"}

		testEntries := createTestEntries(t, 3)

		tmpDir := t.TempDir()
		fd, err := os.CreateTemp(tmpDir, "")
		if err != nil {
			t.Fatalf("failed to create temp file store, %v", err)
		}

		err = json.NewEncoder(fd).Encode(testEntries)
		if err != nil {
			fd.Close()
			t.Fatalf("failed to encode entries as json, %v", err)
		}
		fd.Close()

		cli := runCLI(t, args, buf, fd.Name())

		if len(cli.Entries) != 4 {
			t.Errorf("failed to add, got %d, want 4", len(cli.Entries))
		}
	})
}

func runCLI(t testing.TB, args []string, w io.Writer, path string) *CLI {
	t.Helper()

	cli, err := NewCLI(args, w, path)
	if err != nil {
		t.Fatalf("failed to create CLI, %v", err)
	}

	err = cli.Run()
	if err != nil {
		t.Fatalf("error while running CLI, %v", err)
	}

	err = cli.Save()
	if err != nil {
		t.Fatalf("error wile saving data, %v", err)
	}

	return cli
}

func createTestEntries(t testing.TB, num int) (entries []Task) {
	t.Helper()

	for i := range num {
		desc := fmt.Sprintf("description of task %d", i+1)
		task := NewTask(uint(i+1), desc)
		entries = append(entries, task)
	}

	return
}
