package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const timeFormat = "Jan/02 @15:04"

var (
	ErrConnection      = errors.New("Connection error")
	ErrNotFound        = errors.New("Not found")
	ErrInvalidResponse = errors.New("Invalid server response")
	ErrInvalid         = errors.New("Invalid data")
	ErrNotNumber       = errors.New("Not a number")
)

type item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}
type response struct {
	Results      []item `json:"results"`
	Date         int    `json:"date"`
	TotalResults int    `json:"total_results"`
}

func newClient(timeout time.Duration) *http.Client {
	c := &http.Client{
		Timeout: timeout,
	}
	return c
}

func getItems(url string, timeout time.Duration) ([]item, error) {
	r, err := newClient(timeout).Get(url)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnection, err)
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, fmt.Errorf("Cannot read body: %w", err)
		}
		err = ErrInvalidResponse
		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s", err, msg)
	}
	var resp response
	if err := json.NewDecoder(r.Body).Decode(&resp); err != nil {
		return nil, err
	}
	if resp.TotalResults == 0 {
		return nil, fmt.Errorf("%w: No results found", ErrNotFound)
	}
	return resp.Results, nil
}

func getAll(apiRoot string, timeout time.Duration) ([]item, error) {
	u := fmt.Sprintf("%s/todo", apiRoot)
	return getItems(u, timeout)
}

func getOne(apiRoot string, id int, timeout time.Duration) (item, error) {
	u := fmt.Sprintf("%s/todo/%d", apiRoot, id)
	items, err := getItems(u, timeout)
	if err != nil {
		return item{}, err
	}
	if len(items) != 1 {
		return item{}, fmt.Errorf("%w: Invalid results", ErrInvalid)
	}
	return items[0], nil
}

func sendRequest(url, method, contentType string, timeout time.Duration, expStatus int, body io.Reader) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	r, err := newClient(timeout).Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	if r.StatusCode != expStatus {
		msg, err := io.ReadAll(r.Body)
		if err != nil {
			return fmt.Errorf("Cannot read body: %w", err)
		}
		err = ErrInvalidResponse
		if r.StatusCode == http.StatusNotFound {
			err = ErrNotFound
		}
		return fmt.Errorf("%w: %s", err, msg)
	}
	return nil
}

func addItem(apiRoot, task string, timeout time.Duration) error {
	// Define the Add endpoint URL
	u := fmt.Sprintf("%s/todo", apiRoot)
	item := struct {
		Task string `json:"task"`
	}{
		Task: task,
	}
	var body bytes.Buffer
	if err := json.NewEncoder(&body).Encode(item); err != nil {
		return err
	}
	return sendRequest(u, http.MethodPost, "application/json", timeout, http.StatusCreated, &body)
}

func completeItem(apiRoot string, id int, timeout time.Duration) error {
	u := fmt.Sprintf("%s/todo/%d?complete", apiRoot, id)
	return sendRequest(u, http.MethodPatch, "", timeout, http.StatusNoContent, nil)
}

func deleteItem(apiRoot string, id int, timeout time.Duration) error {
	u := fmt.Sprintf("%s/todo/%d", apiRoot, id)
	return sendRequest(u, http.MethodDelete, "", timeout, http.StatusNoContent, nil)
}