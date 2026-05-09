package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

		cli, err := NewCLI(args, buf, fd.Name())
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

		cli, err := NewCLI(args, buf, fd.Name())
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

		if len(cli.Entries) != 4 {
			t.Errorf("failed to add, got %d, want 4", len(cli.Entries))
		}
	})
}

func createTestEntries(t testing.TB, num int) (entries []Task) {
	for i := range num {
		desc := fmt.Sprintf("description of task %d", i+1)
		task := NewTask(uint(i+1), desc)
		entries = append(entries, task)
	}

	return
}
