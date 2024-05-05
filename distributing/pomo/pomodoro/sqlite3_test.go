//go:build sqlite3
// +build sqlite3

package pomodoro_test

import (
	"os"
	"pragprog.com/rggo/interactiveTools/pomo/pomodoro"
	"pragprog.com/rggo/interactiveTools/pomo/pomodoro/repository"
	"testing"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()
	tf, err := os.CreateTemp("", "pomo")
	if err != nil {
		t.Fatal(err)
	}
	tf.Close()
	dbRepo, err := repository.NewSQLite3Repo(tf.Name())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Using SQLite repository with tempfile %s", tf.Name())
	return dbRepo, func() {
		os.Remove(tf.Name())
	}
}
