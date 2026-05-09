package main

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

// DescriptionCutoff determines the maximum description length displayed when printing
// tasks.
const DescriptionCutoff = 30

// TaskStatus is an Enum representing the state of a task.
type TaskStatus int

const (
	ToDo TaskStatus = iota
	InProgress
	Done
)

var statusText = map[TaskStatus]string{
	ToDo:       "todo",
	InProgress: "in progress",
	Done:       "done",
}

// String implements the Stringer interface.
func (ts TaskStatus) String() string {
	return statusText[ts]
}

// Task represents an entry in the todo list.
type Task struct {
	ID          uint
	Description string
	Status      TaskStatus
	Created     time.Time
	Updated     time.Time
}

// NewTask is mostly a helpful way of creating a new task for the Add function.
func NewTask(id uint, desc string) Task {
	return Task{
		ID:          id,
		Description: desc,
		Status:      ToDo,
		Created:     time.Now(),
		Updated:     time.Now(),
	}
}

// String implements the Stringer interface.
func (t Task) String() string {
	builder := new(strings.Builder)
	maxLen := min(utf8.RuneCountInString(t.Description), DescriptionCutoff)

	switch t.Status {
	case ToDo:
		builder.WriteString("\033[31m")
	case InProgress:
		builder.WriteString("\033[33m")
	case Done:
		builder.WriteString("\033[32m")
	}

	taskStr := fmt.Sprintf(
		"%d  %-*s  %s  %v  %v\033[0m",
		t.ID,
		maxLen,
		truncateString(t.Description, maxLen),
		t.Status,
		t.Created,
		t.Updated,
	)
	builder.WriteString(taskStr)

	return builder.String()
}

func truncateString(s string, maxLen int) string {
	if utf8.RuneCountInString(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
