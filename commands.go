package main

import (
	"fmt"
	"io"
	"slices"
	"strings"
	"unicode/utf8"
)

func Add(args []string, entries []Task) ([]Task, string) {
	var ids []uint
	var newID uint

	for _, task := range entries {
		ids = append(ids, task.ID)
	}

	for i := 1; i <= len(entries); i++ {
		if slices.Contains(ids, uint(i)) {
			continue
		}
		entries = append(entries, NewTask(uint(i), args[0]))
		newID = uint(i)
	}

	outStr := fmt.Sprintf("Task added successfully (ID: %d)", newID)
	return entries, outStr
}

func List(flags FlagStore, entries []Task) string {
	builder := new(strings.Builder)
	toList := getRequestedTasks(flags, entries)
	maxDescLength := setDescriptionLength(toList)

	for i, task := range toList {
		setTaskColorByStatus(task.Status, builder)
		formatStringForListOutput(task, maxDescLength, builder)

		if i != len(toList)-1 {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func setDescriptionLength(entries []Task) int {
	var maxDescLength int

	for _, task := range entries {
		descriptionLength := utf8.RuneCountInString(task.Description)

		if descriptionLength > DescriptionCutoff {
			maxDescLength = DescriptionCutoff
			break
		}

		if descriptionLength > maxDescLength {
			maxDescLength = descriptionLength
		}
	}

	return maxDescLength
}

func truncateDescription(description string, maxLength int) string {
	if utf8.RuneCountInString(description) < maxLength {
		return description
	}

	return string([]rune(description)[:maxLength-3]) + "..."
}

func getRequestedTasks(flags FlagStore, entries []Task) (toList []Task) {
	switch {
	case flags.todo:
		toList = append(toList, getTasksByStatus(entries, ToDo)...)
	case flags.inProgress:
		toList = append(toList, getTasksByStatus(entries, InProgress)...)
	case flags.done:
		toList = append(toList, getTasksByStatus(entries, Done)...)
	case !(flags.todo || flags.inProgress || flags.done):
		toList = append(toList, entries...)
	}

	return
}

func setTaskColorByStatus(status TaskStatus, w io.StringWriter) {
	switch status {
	case ToDo:
		w.WriteString("\033[31m")
	case InProgress:
		w.WriteString("\033[33m")
	case Done:
		w.WriteString("\033[32m")
	}
}

func formatStringForListOutput(task Task, maxDescLength int, w io.StringWriter) {
	taskString := fmt.Sprintf(
		"%4d  %-*s  %11s  %v  %v\033[0m",
		task.ID,
		maxDescLength,
		truncateDescription(task.Description, maxDescLength),
		task.Status,
		task.Created,
		task.Updated,
	)
	w.WriteString(taskString)
}

func getTasksByStatus(entries []Task, status TaskStatus) []Task {
	var foundTasks []Task

	for _, task := range entries {
		if task.Status != status {
			continue
		}
		foundTasks = append(foundTasks, task)
	}
	return foundTasks
}
