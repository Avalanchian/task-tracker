package main

import (
	"fmt"
	"strings"
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
	return fmt.Sprintf("%4d %30s %11s", t.Id, t.Description, t.Status)
}

type TaskList []Task

func (l *TaskList) String() string {
	builder := new(strings.Builder)
	maxLen := 0

	for _, task := range *l {
		if len(task.Description) > maxLen {
			maxLen = len(task.Description)
		}
	}

	for i, task := range *l {
		taskStr := fmt.Sprintf("%4d %-*s %s", task.Id, maxLen, task.Description, task.Status)
		if i != len(*l)-1 {
			taskStr += "\n"
		}
		builder.WriteString(taskStr)
	}

	return builder.String()
}
