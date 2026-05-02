package main

import "fmt"

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
