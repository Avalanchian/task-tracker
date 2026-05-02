package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func Add(args []string, entries TaskList) (TaskList, string) {
	for _, desc := range args {
		task := Task{
			Id:          len(entries) + 1,
			Description: desc,
			Status:      ToDo,
		}
		entries = append(entries, task)
	}
	return entries, fmt.Sprintf("Task added successfully (ID: %d)", len(entries))
}

func List(entries TaskList) string {
	var taskStrings []string
	for _, task := range entries {
		taskStrings = append(taskStrings, task.String())
	}
	return strings.Join(taskStrings, "\n")
}

func Delete(args []string, entries TaskList) (TaskList, string) {
	builder := new(strings.Builder)
	builder.WriteString("Deleted:\n")

	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		checkConversionError(err)
		for i, task := range entries {
			if task.Id == id {
				entries = slices.Delete(entries, i, i+1)
				builder.WriteString(task.String() + "\n")
				break
			}
		}
	}
	return entries, builder.String()
}

func Update(args []string, entries TaskList) (TaskList, string) {
	builder := new(strings.Builder)
	builder.WriteString("Updated:\n")

	id, err := strconv.Atoi(args[0])
	checkConversionError(err)
	for _, task := range entries {
		if task.Id == id {
			task.Description = args[1]
			builder.WriteString(task.String() + "\n")
			break
		}
	}
	return entries, builder.String()
}
