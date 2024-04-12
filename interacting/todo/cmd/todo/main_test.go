package main_test

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

var (
	binName  = "todo"
	fileName = ".todo.json"
)

func TestMain(m *testing.M) {
	if os.Getenv("TODO_FILENAME") != "" {
		fileName = os.Getenv("TODO_FILENAME")
	}
	fmt.Printf("Using test filename %s\n", fileName)

	if err := os.Remove(fileName); err != nil {
		var perr *os.PathError
		if errors.As(err, &perr) {
		}
	}
	fmt.Println("Building tool...")

	build := exec.Command("go", "build", "-o", binName)

	if err := build.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "cannnot build tool %s: %s", binName, err)
		os.Exit(1)
	}

	fmt.Println("Running tests...")
	result := m.Run()

	fmt.Println("Cleaning up...")
	if err := os.Remove(binName); err != nil {
		fmt.Fprintf(os.Stderr, "failed to remvoe a file %s", binName)
		os.Exit(1)
	}
	if err := os.Remove(fileName); err != nil {
		fmt.Fprintf(os.Stderr, "failed to remove a file %s", fileName)
		os.Exit(1)
	}

	os.Exit(result)
}

func TestTodoCLI(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cmdPath := filepath.Join(dir, binName)

	task := "test task number 1"
	t.Run("AddNewTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})

	task2 := "test task number 2"
	t.Run("AddNewTaskFromSTDIN", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add")
		cmdStdin, err := cmd.StdinPipe()
		if err != nil {
			t.Fatal(err)
		}
		go func() {
			defer cmdStdin.Close()
			io.WriteString(cmdStdin, task2)
		}()
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("ListTasks", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("  1: %s\n  2: %s\n", task, task2)
		assert.Equal(t, expected, string(out), "expected %q, got %q", expected, string(out))
	})
	t.Run("CompleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-complete", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
		cmd = exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := "X 1: " + task + "\n  2: test task number 2\n"
		assert.Equal(t, expected, string(out), "expected %q, got %q", expected, string(out))
	})
	t.Run("DeleteTask", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-delete", "2")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
		cmd = exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := "X 1: test task number 1\n"
		assert.Equal(t, expected, string(out), "expected %q, got %q", expected, string(out))
	})
	t.Run("DelteTaskAgain", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-delete", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
		cmd = exec.Command(cmdPath, "-list")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := ""
		assert.Equal(t, expected, string(out), "expected %q, got %q", expected, string(out))
	})
	task3 := "test task number 3"
	t.Run("VerboseOutput", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-add", task3)
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
		cmd = exec.Command(cmdPath, "-verbose")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("  1: test task number 3, Created at: %s\n", time.Now().Format(time.DateTime))
		assert.Equal(t, expected, string(out), "expected %q, got %q", expected, string(out))
	})
	t.Run("FilteredOutput", func(t *testing.T) {
		cmd := exec.Command(cmdPath, "-filter")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected := fmt.Sprintf("  1: test task number 3\n")
		assert.Equal(t, expected, string(out), "expected %q, got %q", expected, string(out))
		cmd = exec.Command(cmdPath, "-complete", "1")
		if err := cmd.Run(); err != nil {
			t.Fatal(err)
		}
		cmd = exec.Command(cmdPath, "-filter")
		out, err = cmd.CombinedOutput()
		if err != nil {
			t.Fatal(err)
		}
		expected = ""
		assert.Equal(t, expected, string(out), "expected %q, got %q", expected, string(out))
	})
}
