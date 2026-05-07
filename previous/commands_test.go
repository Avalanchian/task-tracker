package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestAdd(t *testing.T) {
	for i := range 3 {
		t.Run(fmt.Sprintf("list with %d tasks", i), func(t *testing.T) {
			entries := setupTempEntries(i)
			args := []string{"add test task"}

			entries, _ = Add(args, entries)

			if len(entries) != i+1 {
				t.Errorf("entry not added %d", len(entries))
			}
		})
	}
}

func TestList(t *testing.T) {
	t.Run("returns string of entries", func(t *testing.T) {
		num := 3
		entries := setupTempEntries(num)

		outMsg := List(entries)

		got := strings.Count(outMsg, "\n")
		if got != num-1 {
			t.Errorf("incorrect list printing, %q", outMsg)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("removes an entry from the list", func(t *testing.T) {
		num := 3
		entries := setupTempEntries(num)

		entries, _, _ = Delete([]string{"2"}, entries)

		if len(entries) != num-1 {
			t.Errorf("delete malfunction, %+v", entries)
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("changes description of entry", func(t *testing.T) {
		num := 3
		entries := setupTempEntries(num)

		entries, _, _ = Update([]string{"2", "new desc of task 2"}, entries)

		if entries[1].Description != "new desc of task 2" {
			t.Errorf("update malfunction, got %q, want %q", entries[1].Description, "new desc of task 2")
		}
	})
}

<<<<<<< HEAD
func TestMark(t *testing.T) {
	t.Run("no options", func(t *testing.T) {
		num := 3
		entries := setupTempEntries(num)
		entries[2].Status = InProgress

		entries, _, _ = Mark([]string{"2", "3"}, Flags{}, entries)

		if entries[1].Status != InProgress {
			t.Errorf("got %v, want %v", entries[1].Status, InProgress)
		}
		if entries[2].Status != Done {
			t.Errorf("got %v, want %v", entries[2].Status, Done)
		}
	})

	t.Run("done flag", func(t *testing.T) {
		num := 3
		entries := setupTempEntries(num)
		flags := Flags{
			toDo:       false,
			inProgress: true,
			done:       false,
		}

		entries, _, _ = Mark([]string{"2", "3"}, flags, entries)

		if entries[1].Status != InProgress {
			t.Errorf("got %v, want %v", entries[1].Status, InProgress)
		}
		if entries[2].Status != InProgress {
			t.Errorf("got %v, want %v", entries[2].Status, InProgress)
		}
	})
}

=======
>>>>>>> parent of f6a6fa6 (Added mark command and improved string formats.)
func setupTempEntries(n int) TaskList {
	var entries TaskList

	for i := range n {
		task := Task{
			Id:          i + 1,
			Description: fmt.Sprintf("description of task %d", i+1),
			Status:      ToDo,
		}
		entries = append(entries, task)
	}

	return entries
}
