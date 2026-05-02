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

		entries, _ = Delete([]string{"2"}, entries)

		if len(entries) != num-1 {
			t.Errorf("delete malfunction, %+v", entries)
		}
	})
}

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
