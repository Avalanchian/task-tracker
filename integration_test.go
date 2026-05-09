package main

import (
	"bytes"
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
}
