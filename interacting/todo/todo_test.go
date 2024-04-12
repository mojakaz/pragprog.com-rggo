package todo_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"pragprog.com/rggo/interacting/todo"
	"testing"
)

func TestList_Add(t *testing.T) {
	l := todo.List{}
	task := "New"
	l.Add(task)

	if l[0].Task != task {
		t.Errorf("expected %s, got %s", task, l[0].Task)
	}
}

func TestList_Complete(t *testing.T) {
	l := todo.List{}
	task := "New ToDo Item"
	task2 := "Another ToDo Item"
	l.Add(task)
	l.Add(task2)
	assert.Equal(t, task, l[0].Task, "expected %s, got %s", task, l[0].Task)
	assert.False(t, l[0].Done, "new task should not be completed")
	err := l.Complete(1)
	require.NoError(t, err, "expected no error")
	assert.True(t, l[0].Done, "new task should be completed")

	assert.Equal(t, task2, l[1].Task)
	assert.False(t, l[1].Done)
	err = l.Complete(2)
	require.NoError(t, err)
	assert.True(t, l[1].Done)
}

func TestList_Delete(t *testing.T) {
	l := todo.List{}
	tasks := []string{
		"New Task1",
		"New Task2",
		"New Task3",
	}
	for _, value := range tasks {
		l.Add(value)
	}
	assert.Equal(t, tasks[0], l[0].Task, "expected %q, got %q", tasks[0], l[0].Task)
	err := l.Delete(2)
	require.NoError(t, err, "expected no error")
	assert.Equal(t, len(l), 2, "expected list length %d, got %d", 2, len(l))
	assert.Equal(t, tasks[2], l[1].Task, "expected %q, got %q", tasks[2], l[1].Task)
}

func TestList_SaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}
	task := "New Task"
	l1.Add(task)
	assert.Equal(t, l1[0].Task, task, "expected %q, got %q", task, l1[0].Task)
	tf, err := os.CreateTemp("", "")
	require.NoError(t, err, "failed to create a temp file: %w", err)
	defer func() {
		if err := os.Remove(tf.Name()); err != nil {
			panic(err)
		}
	}()
	err = l1.Save(tf.Name())
	require.NoError(t, err, "failed to save the list to file: %w", err)
	err = l2.Get(tf.Name())
	require.NoError(t, err, "failed to get the list from file: %w", err)
	assert.Equal(t, l1[0].Task, l2[0].Task, "expected task %q match task %q", l1[0].Task, l2[0].Task)
}
