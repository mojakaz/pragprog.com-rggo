//go:build !inmemory && !containers && !sqlite3

package cmd

import (
	"github.com/spf13/viper"
	"pragprog.com/rggo/interactiveTools/pomo/pomodoro"
	"pragprog.com/rggo/interactiveTools/pomo/pomodoro/repository"
)

func getRepo() (pomodoro.Repository, error) {
	repo, err := repository.NewPostgresRepo(viper.GetString("dbname"))
	if err != nil {
		return nil, err
	}
	return repo, nil
}
