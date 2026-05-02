package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Task struct {
	Id          int
	Description string
	Status      string
}

func Add(task Task, file *os.File) {
	_, err := file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Errorf("error seeking file, %v", err)
	}

	var tasks []Task
	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		fmt.Errorf("error during json decoding, %v", err)
	}

	tasks = append(tasks, task)

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		fmt.Errorf("error seeking file, %v", err)
	}
	err = json.NewEncoder(file).Encode(tasks)
	if err != nil {
		fmt.Errorf("error during json encoding, %v", err)
	}
}
