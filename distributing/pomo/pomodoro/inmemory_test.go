//go:build inmemory
// +build inmemory

package pomodoro_test

import (
	"pragprog.com/rggo/interactiveTools/pomo/pomodoro"
	"pragprog.com/rggo/interactiveTools/pomo/pomodoro/repository"
	"testing"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()
	return repository.NewInMemoryRepo(), func() {}
}
