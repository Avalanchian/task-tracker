package main

import (
	"encoding/json"
	"io"
	"os"
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	t.Run("add to empty file", func(t *testing.T) {
		fd, err := os.CreateTemp("", "test_add")
		assertNoIOError(t, err)
		defer fd.Close()

		task := Task{
			Id:          1,
			Description: "a task description",
			Status:      ToDo,
		}

		Add(task, fd)

		var got []Task
		_, err = fd.Seek(0, io.SeekStart)
		assertNoIOError(t, err)

		err = json.NewDecoder(fd).Decode(&got)
		if err != nil {
			t.Fatalf("error during json decoding, %v", err)
		}

		if !reflect.DeepEqual(got[0], task) {
			t.Errorf("got %+v, want %+v", got, task)
		}
	})

	t.Run("add to non-empty file", func(t *testing.T) {
		fd, err := os.CreateTemp("", "test_add")
		assertNoIOError(t, err)
		defer fd.Close()

		task1 := Task{
			Id:          1,
			Description: "a task description",
			Status:      ToDo,
		}
		task2 := Task{
			Id:          2,
			Description: "another description",
			Status:      InProgress,
		}

		tasks := []Task{task1}
		err = json.NewEncoder(fd).Encode(tasks)
		assertNoJsonError(t, err)

		Add(task2, fd)

		tasks = append(tasks, task2)
		var got []Task
		_, err = fd.Seek(0, io.SeekStart)
		assertNoIOError(t, err)

		err = json.NewDecoder(fd).Decode(&got)
		assertNoJsonError(t, err)

		assertTaskListsEqual(t, got, tasks)

	})
}

func assertNoIOError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("error during io operation, %v", err)
	}
}

func assertNoJsonError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("error during json (de-)serialization, %v", err)
	}
}

func assertTaskListsEqual(t testing.TB, got, want []Task) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}
