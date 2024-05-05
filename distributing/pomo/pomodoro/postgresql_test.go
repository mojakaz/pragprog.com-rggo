package pomodoro_test

import (
	"pragprog.com/rggo/interactiveTools/pomo/pomodoro"
	"pragprog.com/rggo/interactiveTools/pomo/pomodoro/repository"
	"testing"
)

func getRepo(t *testing.T, dbName string) (pomodoro.Repository, func()) {
	t.Helper()
	dbRepo, err := repository.NewPostgresRepo(dbName)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Using PostgreSQL repository")
	return dbRepo, func() {}
}
