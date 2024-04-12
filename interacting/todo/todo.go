package todo

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type List []item

// String implements Stringer interface
func (l *List) String() string {
	var formatted string
	for index, value := range *l {
		prefix := "  "
		if value.Done {
			prefix = "X "
		}
		formatted += fmt.Sprintf("%s%d: %s\n", prefix, index+1, value.Task)
	}
	return formatted
}

// Add creates a new todo item and appends it to the list
func (l *List) Add(task string) {
	t := item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}
	*l = append(*l, t)
}

// Complete marks a todo item as completed by
// setting Done = true and CompletedAt to the current time
func (l *List) Complete(i int) error {
	ls := *l
	if i <= 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}

	// adjusting index for 0 based index
	ls[i-1].Done = true
	ls[i-1].CompletedAt = time.Now()

	return nil
}

// Delete deletes a todo item from the list
func (l *List) Delete(i int) error {
	ls := *l
	if i < 0 || i > len(ls) {
		return fmt.Errorf("item %d does not exist", i)
	}
	*l = append(ls[:i-1], ls[i:]...)
	return nil
}

// Save encodes the List as JSON and saves it
// using the provided file name
func (l *List) Save(name string) error {
	js, err := json.Marshal(l)
	if err != nil {
		return fmt.Errorf("unable to encode the List as JSON: %w", err)
	}
	return os.WriteFile(name, js, 0644)
}

// Get opens the provided file name, decodes
// the JSON data and parses it into a List
func (l *List) Get(name string) error {
	file, err := os.ReadFile(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("unable to open the file name %s: %w", name, err)
	}

	if len(file) == 0 {
		return nil
	}
	return json.Unmarshal(file, l)
}
