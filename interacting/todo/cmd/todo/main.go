package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"pragprog.com/rggo/interacting/todo"
	"strings"
	"time"
)

var todoFileName = ".todo.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"%s tool. Developed for The Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2020\n")
		fmt.Fprintln(flag.CommandLine.Output(), "To add new task, use -add flag followed by task name. You can also provide task name from STDIN.")
		fmt.Fprintln(flag.CommandLine.Output(), "The default filename to save the to-do tasks is .todo.json.\nTo change the filename, specify the new filename with environment variable TODO_FILENAME.\ni.e. `export TODO_FILENAME=<new filename>`")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}

	// Parsing command line flags
	add := flag.Bool("add", false, "Add task to the to-do list")
	list := flag.Bool("list", false, "List all tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	del := flag.Int("delete", 0, "Item to be deleted")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	filter := flag.Bool("filter", false, "Prevent displaying completed itmes")

	flag.Parse()

	if os.Getenv("TODO_FILENAME") != "" {
		todoFileName = os.Getenv("TODO_FILENAME")
	}

	// Define an items list
	l := &todo.List{}

	// Use the Get method to read to-do items from file
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Decide what to do based on the number of arguments provided
	switch {
	// For no extra arguments, print the list
	case *list:
		// List current to-do items
		fmt.Print(l)
	// Enable verbose output
	case *verbose:
		// List current to-do items with time of creation and completion
		var formatted string
		for index, value := range *l {
			prefix := "  "
			suffix := fmt.Sprintf("Created at: %s\n", value.CreatedAt.Format(time.DateTime))
			if value.Done {
				prefix = "X "
				suffix = fmt.Sprintf("Created at: %s, Completed at: %s\n", value.CreatedAt.Format(time.DateTime), value.CompletedAt.Format(time.DateTime))
			}
			formatted += fmt.Sprintf("%s%d: %s, %s", prefix, index+1, value.Task, suffix)
		}
		fmt.Print(formatted)
	// Display uncompleted tasks only
	case *filter:
		var formatted string
		for index, value := range *l {
			if value.Done {
				continue
			}
			formatted += fmt.Sprintf("  %d: %s\n", index+1, value.Task)
		}
		fmt.Print(formatted)
	// Complete the given item
	case *complete > 0:
		// Complete the given item
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *add:
		// When any arguments (excluding flags) are provided,
		// they will be used as the new task
		t, err := getTask(os.Stdin, flag.Args()...)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		s := strings.Split(t, "\n")
		for _, task := range s[:len(s)-1] {
			l.Add(task)
		}
		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case *del > 0:
		// Delete the given item
		if err := l.Delete(*del); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Save the new list
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		// Invalid flag provided
		flag.Usage()
		os.Exit(1)
	}
}

// getTask function decides where to get the description for a new
// task from: arguments or STDIN
func getTask(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}
	s := bufio.NewScanner(r)
	joined := ""
	fmt.Println("To add new tasks, type task name. To add another task, push enter and type another task name. To finish, push enter again.")
	for s.Scan() {
		if err := s.Err(); err != nil {
			return "", err
		}
		if len(s.Text()) == 0 {
			break
		}
		joined += s.Text() + "\n"
	}
	if len(joined) == 0 {
		return "", fmt.Errorf("task cannot be blank")
	}
	return joined, nil
}
