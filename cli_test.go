package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestNewCLI(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "testdata.json")
	fd, err := os.Create(path)
	checkIOError(err)
	fd.Close()

	out := new(bytes.Buffer)
	_ = NewCLI(path, []string{"task-tracker"}, out)

	if len(out.String()) == 0 {
		t.Errorf("expected error not present, %q", out.String())
	}
}
