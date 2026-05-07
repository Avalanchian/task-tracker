package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type TaskStatus int

const (
	ToDo TaskStatus = iota
	InProgress
	Done
)

var StatusStrings = map[TaskStatus]string{
	ToDo:       "todo",
	InProgress: "in progress",
	Done:       "done",
}

func (t TaskStatus) String() string {
	return StatusStrings[t]
}

type Task struct {
	Id          int        `json:"ID"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
}

func (t Task) String() string {
	builder := new(strings.Builder)
	maxLen := utf8.RuneCountInString(t.Description)
	cutoff := 30
	if maxLen >= cutoff {
		maxLen = cutoff
	}

	switch t.Status {
	case ToDo:
		builder.WriteString("\033[31m")
	case InProgress:
		builder.WriteString("\033[33m")
	case Done:
		builder.WriteString("\033[32m")
	}

	builder.WriteString(
		fmt.Sprintf(
			"%4d  %-*s  %s\033[0m",
			t.Id,
			maxLen,
			truncateString(t.Description, 30),
			t.Status,
		),
	)
	return builder.String()
}

type TaskList []Task

func (l *TaskList) String() string {
	builder := new(strings.Builder)
	cutoff := 40
	maxLen := 0

	for _, task := range *l {
		if maxLen >= cutoff {
			maxLen = cutoff
			break
		}

		if len(task.Description) > maxLen {
			maxLen = len(task.Description)
		}
	}

	for i, task := range *l {
		var taskStr string

		switch task.Status {
		case ToDo:
			taskStr = fmt.Sprintf(
				"\033[31m%4d %-*s %s\033[0m",
				task.Id,
				maxLen,
				truncateString(task.Description, 40),
				task.Status,
			)
		case InProgress:
			taskStr = fmt.Sprintf(
				"\033[33m%4d %-*s %s\033[0m",
				task.Id,
				maxLen,
				truncateString(task.Description, 40),
				task.Status,
			)
		case Done:
			taskStr = fmt.Sprintf(
				"\033[32m%4d %-*s %s\033[0m",
				task.Id,
				maxLen,
				truncateString(task.Description, 40),
				task.Status,
			)
		}

		if i != len(*l)-1 {
			taskStr += "\n"
		}
		builder.WriteString(taskStr)
	}

	return builder.String()
}

func truncateString(str string, maxLen int) string {
	if utf8.RuneCountInString(str) >= maxLen {
		return string([]rune(str)[:maxLen-3]) + "..."
	}
	return str
}
