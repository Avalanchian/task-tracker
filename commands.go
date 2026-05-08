package main

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func Add(args []string, entries []Task) ([]Task, string) {
	var ids []uint
	var newID uint

	for _, task := range entries {
		ids = append(ids, task.ID)
	}

	for i := 1; i <= len(entries)+1; i++ {
		if slices.Contains(ids, uint(i)) {
			continue
		}

		entries = append(entries, NewTask(uint(i), args[0]))
		newID = uint(i)
		break
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

func Delete(args []string, entries []Task) ([]Task, string, error) {
	var deleteCount int
	var toDelete []uint
	var remaining []Task

	builder := new(strings.Builder)

	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return nil, "", fmt.Errorf("unable to convert %q to valid id, %w", arg, err)
		}
		toDelete = append(toDelete, uint(id))
	}

	for _, task := range entries {
		if !slices.Contains(toDelete, task.ID) {
			continue
		}

		outStr := fmt.Sprintf("Task deleted successfully (ID: %d)", task.ID)
		builder.WriteString(outStr)
		deleteCount++

		if deleteCount != len(args) {
			builder.WriteString("\n")
		}

		remaining = append(remaining, task)
	}

	return remaining, builder.String(), nil
}

func Mark(args []string, flags FlagStore, entries []Task) ([]Task, string, error) {
	var toMark []uint
	var markFlag TaskStatus
	var defaultBehaviour bool
	var entriesOut []Task

	builder := new(strings.Builder)

	switch {
	case flags.todo:
		markFlag = ToDo
	case flags.inProgress:
		markFlag = InProgress
	case flags.done:
		markFlag = Done
	default:
		defaultBehaviour = true
	}

	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return nil, "", fmt.Errorf("unable to convert %q to vailid id, %w", arg, err)
		}
		toMark = append(toMark, uint(id))
	}

	for i, task := range entries {
		if !slices.Contains(toMark, task.ID) {
			entriesOut = append(entriesOut, task)
			continue
		}

		if defaultBehaviour && task.Status != Done {
			task.Status++
			task.Updated = time.Now()
			outStr := fmt.Sprintf("Successfully marked as %s (ID: %d)", task.Status, task.ID)
			if i != len(entries)-1 {
				outStr += "\n"
			}
			builder.WriteString(outStr)
		}

		if defaultBehaviour {
			entriesOut = append(entriesOut, task)
			continue
		}

		task.Status = markFlag
		task.Updated = time.Now()
		outStr := fmt.Sprintf("Successfully marked as %s (ID: %d)", task.Status, task.ID)
		if i != len(entries)-1 {
			outStr += "\n"
		}
		builder.WriteString(outStr)

		entriesOut = append(entriesOut, task)
	}

	return entriesOut, builder.String(), nil
}

func setDescriptionLength(entries []Task) int {
	var maxDescLength int

	for _, task := range entries {
		descriptionLength := utf8.RuneCountInString(task.Description)

		if descriptionLength >= DescriptionCutoff {
			return DescriptionCutoff
		}

		if descriptionLength > maxDescLength {
			maxDescLength = descriptionLength
		}
	}

	return maxDescLength
}

func truncateDescription(description string, maxLength int) string {
	if utf8.RuneCountInString(description) <= maxLength {
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
		"%4d  %-*s  %-11s  %v  %v\033[0m",
		task.ID,
		maxDescLength,
		truncateDescription(task.Description, DescriptionCutoff),
		task.Status,
		task.Created.Format(time.DateTime),
		task.Updated.Format(time.DateTime),
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
