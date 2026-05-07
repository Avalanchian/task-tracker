package main

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
)

const DescriptionCutoff = 30

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

func (ts TaskStatus) String() string {
	return statusText[ts]
}

type Task struct {
	ID          uint
	Description string
	Status      TaskStatus
	Created     time.Time
	Updated     time.Time
}

func NewTask(id uint, desc string) *Task {
	return &Task{
		ID:          id,
		Description: desc,
		Status:      ToDo,
		Created:     time.Now(),
		Updated:     time.Now(),
	}
}

func (t *Task) String() string {
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
