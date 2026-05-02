package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type TaskStatus int

const (
	ToDo TaskStatus = iota
	InProgress
	Complete
)

var StatusStrings = map[TaskStatus]string{
	ToDo:       "todo",
	InProgress: "in progress",
	Complete:   "complete",
}

func (ts TaskStatus) String() string {
	return StatusStrings[ts]
}

type Task struct {
	Id          int
	Description string
	Status      TaskStatus
}

func Add(task Task, file *os.File) {
	tasks := decodeTaskList(file)
	tasks = append(tasks, task)
	encodeTaskList(tasks, file)
}

func encodeTaskList(tasks []Task, file *os.File) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Errorf("error seeking file, %v", err)
	}

	err = json.NewEncoder(file).Encode(tasks)
	if err != nil {
		fmt.Errorf("error during json encoding, %v", err)
	}
}

func decodeTaskList(file *os.File) []Task {
	var tasks []Task

	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Errorf("error seeking file, %v", err)
	}

	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		fmt.Errorf("error during json decoding, %v", err)
	}

	return tasks
}
