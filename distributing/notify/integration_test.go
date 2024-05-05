//go:build integration

package notify_test

import (
	"pragprog.com/rggo/distributing/notify"
	"testing"
)

func TestSend(t *testing.T) {
	n := notify.New("test title", "test msg", notify.SeverityNormal)
	err := n.Send()
	if err != nil {
		t.Error(err)
	}
}
