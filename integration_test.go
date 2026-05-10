package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"testing"
)

func TestCLIAdd(t *testing.T) {
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

	t.Run("use lowest index", func(t *testing.T) {
		buf := new(bytes.Buffer)
		args := []string{"task-tracker", "add", "a test description"}

		testEntries := createTestEntries(t, 5)
		testEntries = slices.Delete(testEntries, 2, 3)

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

		runCLI(t, args, buf, fd.Name())

		if !strings.Contains(buf.String(), "3") {
			t.Errorf("did not use lowest id, got %q (want ID 3)", buf.String())
		}
	})
}

func TestCLIList(t *testing.T) {
	testEntries := createTestEntries(t, 6)
	testEntries[3].Status = InProgress
	testEntries[4].Status = InProgress
	testEntries[5].Status = Done

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

	t.Run("list all", func(t *testing.T) {
		buf := new(bytes.Buffer)
		args := []string{"task-tracker", "list"}

		runCLI(t, args, buf, fd.Name())

		if strings.Count(buf.String(), "\n") != 6 {
			t.Errorf("incorrect list output, got %q", buf.String())
		}
	})

	t.Run("list ToDo", func(t *testing.T) {
		buf := new(bytes.Buffer)
		args := []string{"task-tracker", "list", "-todo"}

		runCLI(t, args, buf, fd.Name())

		if strings.Count(buf.String(), "\n") != 3 {
			t.Errorf("incorrect list output, got %q", buf.String())
		}
	})

	t.Run("list InProgress", func(t *testing.T) {
		buf := new(bytes.Buffer)
		args := []string{"task-tracker", "list", "-inprogress"}

		runCLI(t, args, buf, fd.Name())

		if strings.Count(buf.String(), "\n") != 2 {
			t.Errorf("incorrect list output, got %q", buf.String())
		}
	})

	t.Run("list Done", func(t *testing.T) {
		buf := new(bytes.Buffer)
		args := []string{"task-tracker", "list", "-done"}

		runCLI(t, args, buf, fd.Name())

		if strings.Count(buf.String(), "\n") != 1 {
			t.Errorf("incorrect list output, got %q", buf.String())
		}
	})
}

func TestCLIDelete(t *testing.T) {
	tmpDir := t.TempDir()

	buf := new(bytes.Buffer)
	args := []string{"task-tracker", "delete", "2"}

	testEntries := createTestEntries(t, 3)
	fd := setupTempFileStore(t, tmpDir, testEntries)

	cli := runCLI(t, args, buf, fd.Name())

	if len(cli.Entries) != 2 {
		t.Errorf("entry not removed, got %v", cli.Entries)
	}

	if !strings.Contains(buf.String(), "2") {
		t.Errorf("incorrect ID, got %q, want 2", buf.String())
	}
}

func TestCLIMark(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("no options (todo -> in progress)", func(t *testing.T) {
		buf := new(bytes.Buffer)
		args := []string{"task-tracker", "mark", "2"}

		testEntries := createTestEntries(t, 3)
		fd := setupTempFileStore(t, tmpDir, testEntries)

		cli := runCLI(t, args, buf, fd.Name())

		got := cli.Entries[1].Status
		want := InProgress

		if got != want {
			t.Errorf("incorrect status, got %q, want %q", got, want)
		}
	})

	t.Run("no options (in progress -> done)", func(t *testing.T) {
		buf := new(bytes.Buffer)
		args := []string{"task-tracker", "mark", "2"}

		testEntries := createTestEntries(t, 3)
		testEntries[1].Status = InProgress
		fd := setupTempFileStore(t, tmpDir, testEntries)

		cli := runCLI(t, args, buf, fd.Name())

		got := cli.Entries[1].Status
		want := Done

		if got != want {
			t.Errorf("incorrect status, got %q, want %q", got, want)
		}
	})

	t.Run("todo -> done", func(t *testing.T) {
		buf := new(bytes.Buffer)
		args := []string{"task-tracker", "mark", "-done", "2"}

		testEntries := createTestEntries(t, 3)
		fd := setupTempFileStore(t, tmpDir, testEntries)

		cli := runCLI(t, args, buf, fd.Name())

		got := cli.Entries[1].Status
		want := Done

		if got != want {
			t.Errorf("incorrect status, got %q, want %q", got, want)
		}
	})

	t.Run("in progress -> todo", func(t *testing.T) {
		buf := new(bytes.Buffer)
		args := []string{"task-tracker", "mark", "-todo", "2"}

		testEntries := createTestEntries(t, 3)
		testEntries[1].Status = InProgress
		fd := setupTempFileStore(t, tmpDir, testEntries)

		cli := runCLI(t, args, buf, fd.Name())

		got := cli.Entries[1].Status
		want := ToDo

		if got != want {
			t.Errorf("incorrect status, got %q, want %q", got, want)
		}
	})
}

func TestCLIUpdate(t *testing.T) {
	tmpDir := t.TempDir()

	buf := new(bytes.Buffer)
	args := []string{"task-tracker", "update", "2", "a new description"}

	testEntries := createTestEntries(t, 3)
	fd := setupTempFileStore(t, tmpDir, testEntries)

	cli := runCLI(t, args, buf, fd.Name())

	got := cli.Entries[1].Description
	want := "a new description"

	if got != want {
		t.Errorf("incorrect status, got %q, want %q", got, want)
	}
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

func setupTempFileStore(t testing.TB, dir string, tasks []Task) *os.File {
	t.Helper()

	fd, err := os.CreateTemp(dir, "")
	if err != nil {
		t.Fatalf("failed to create temp file store, %v", err)
	}

	err = json.NewEncoder(fd).Encode(tasks)
	if err != nil {
		fd.Close()
		t.Fatalf("failed to encode entries as json, %v", err)
	}
	fd.Close()
	return fd
}
