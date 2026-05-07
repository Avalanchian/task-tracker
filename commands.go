package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func Add(args []string, entries TaskList) (TaskList, string) {
	builder := new(strings.Builder)

	var currentIDs []int
	for _, task := range entries {
		currentIDs = append(currentIDs, task.Id)
	}

	for i, desc := range args {
		for id := 1; id <= len(entries)+1; id++ {
			if slices.Contains(currentIDs, id) {
				continue
			}

			task := Task{
				Id:          id,
				Description: desc,
				Status:      ToDo,
			}
			entries = append(entries, task)
			currentIDs = append(currentIDs, id)
			msg := fmt.Sprintf("Task added successfully (ID: %d)", id)

			if i != len(args)-1 {
				msg += "\n"
			}

			builder.WriteString(msg)
			break
		}
	}
	return entries, builder.String()
}

func List(entries TaskList) string {
	return entries.String()
}

func Delete(args []string, entries TaskList) (TaskList, string, error) {
	builder := new(strings.Builder)
	builder.WriteString("Deleted:\n")

	for i, arg := range args {
		id, err := strconv.Atoi(arg)
		if err = checkConversionError(err); err != nil {
			return nil, "", fmt.Errorf("error in Delete command, %w", err)
		}
		for j, task := range entries {
			if task.Id == id {
				entries = slices.Delete(entries, j, j+1)
				builder.WriteString(task.String())
				break
			}
		}
		if i != len(args)-1 {
			builder.WriteString("\n")
		}
	}
	return entries, builder.String(), nil
}

func Update(args []string, entries TaskList) (TaskList, string, error) {
	builder := new(strings.Builder)
	builder.WriteString("Updated:\n")

	id, err := strconv.Atoi(args[0])
	if err = checkConversionError(err); err != nil {
		return nil, "", fmt.Errorf("error in Update command, %w", err)
	}
	for i, task := range entries {
		if task.Id == id {
			entries[i].Description = args[1]
			builder.WriteString(entries[i].String())
			break
		}
	}
	return entries, builder.String(), nil
}

func Mark(args []string, entries TaskList) (TaskList, string, error) {
	builder := new(strings.Builder)
	builder.WriteString("Marked:\n")

	for i, index := range args {
		id, err := strconv.Atoi(index)
		if err = checkConversionError(err); err != nil {
			return nil, "", fmt.Errorf("error in Mark command, %w", err)
		}

		for j, task := range entries {
			if task.Id == id && (task.Status == ToDo || task.Status == InProgress) {
				entries[j].Status += 1
				builder.WriteString(entries[j].String())
				break
			}
		}
		if i != len(args)-1 {
			builder.WriteString("\n")
		}
	}
	return entries, builder.String(), nil
}
