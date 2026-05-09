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

// Add appends a new Task to entries with a Description given in args. It returns a
// new task list and an output string to log the successful addition.
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

// List returns a string of all entries requested via flags with appropriate color
// and formatting, with tasks delineated by newlines (\n).
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

// Delete returns a new task list from entries, without those tasks with IDs provided
// by args. If any of args cannot be converted to uint, then Delete returns a
// *strconv.NumError.
func Delete(args []string, entries []Task) ([]Task, string, error) {
	var deleteCount int
	var toDelete []uint
	var remaining []Task

	builder := new(strings.Builder)

	toDelete, err := convertArgIDs(args)
	if err != nil {
		return nil, "", fmt.Errorf("could not convert arg IDs for deletion, %w", err)
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

// Mark returns a new task list from entries. The tasks with IDs provided by args are
// updated based on the flags provided. If no flags are provided, default behaviour is to
// change ToDo to InProgress and InProgress to Done. If any of args cannot be converted
// to uint, Mark returns a *strconv.NumError.
func Mark(args []string, flags FlagStore, entries []Task) ([]Task, string, error) {
	var markCount int
	var entriesOut []Task

	builder := new(strings.Builder)

	markFlag, defaultBehaviour := setMarkBehaviour(flags)
	toMark, err := convertArgIDs(args)
	if err != nil {
		return nil, "", fmt.Errorf("could not convert arg IDs for marking, %w", err)
	}

	for _, task := range entries {
		if !slices.Contains(toMark, task.ID) {
			entriesOut = append(entriesOut, task)
			continue
		}

		if defaultBehaviour && task.Status != Done {
			task.Status++
			markCount++
			task.Updated = time.Now()
			outStr := fmt.Sprintf("Successfully marked as %s (ID: %d)", task.Status, task.ID)
			if markCount != len(args) {
				outStr += "\n"
			}
			builder.WriteString(outStr)
		}

		if defaultBehaviour {
			entriesOut = append(entriesOut, task)
			continue
		}

		task.Status = markFlag
		markCount++
		task.Updated = time.Now()
		outStr := fmt.Sprintf("Successfully marked as %s (ID: %d)", task.Status, task.ID)
		if markCount != len(args) {
			outStr += "\n"
		}
		builder.WriteString(outStr)

		entriesOut = append(entriesOut, task)
	}

	return entriesOut, builder.String(), nil
}

// Update returns a new task list where the task with ID given by args[0] has a newr
// Description, given by args[1]. If args[0] cannot be converted to uint, Update returns
// a *strconv.NumError.
func Update(args []string, entries []Task) ([]Task, string, error) {
	var entriesOut []Task

	builder := new(strings.Builder)

	toUpdate, err := convertArgIDs(args[:1])
	if err != nil {
		return nil, "", fmt.Errorf("Update failed during ID conversion, %w", err)
	}

	for _, task := range entries {
		if !slices.Contains(toUpdate, task.ID) {
			entriesOut = append(entriesOut, task)
			continue
		}

		task.Description = args[1]
		task.Updated = time.Now()
		outStr := fmt.Sprintf("Successfully updated description (ID: %d)", task.ID)
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

func setMarkBehaviour(flags FlagStore) (TaskStatus, bool) {
	switch {
	case flags.todo:
		return ToDo, false
	case flags.inProgress:
		return InProgress, false
	case flags.done:
		return Done, false
	default:
		return ToDo, true
	}
}

func convertArgIDs(args []string) ([]uint, error) {
	var convertedArgs []uint
	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return nil, fmt.Errorf("unable to convert %q to vailid id, %w", arg, err)
		}
		convertedArgs = append(convertedArgs, uint(id))
	}

	return convertedArgs, nil
}
