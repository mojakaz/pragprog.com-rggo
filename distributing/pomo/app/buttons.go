package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/button"
	"pragprog.com/rggo/interactiveTools/pomo/pomodoro"
)

type buttonSet struct {
	btStart *button.Button
	btPause *button.Button
}

func newButtonSet(ctx context.Context, config *pomodoro.IntervalConfig, w *widgets, s *summary,
	redrawCh chan<- bool, errorCh chan<- error) (*buttonSet, error) {
	startInterval := func() {
		i, err := pomodoro.GetInterval(config)
		errorCh <- err
		start := func(i pomodoro.Interval) {
			message := "Take a break"
			if i.Category == pomodoro.CategoryPomodoro {
				message = "Focus on your task"
			}
			w.update([]int{}, i.ActualDuration, message, "", redrawCh)
			sendNotification(message)
		}
		end := func(pomodoro.Interval) {
			w.update([]int{}, 0, "Nothing running...", "", redrawCh)
			s.update(redrawCh)
			message := fmt.Sprintf("%s finished!", i.Category)
			sendNotification(message)
		}
		periodic := func(i pomodoro.Interval) {
			w.update(
				[]int{int(i.ActualDuration), int(i.PlannedDuration)}, i.ActualDuration, "",
				fmt.Sprint(i.PlannedDuration-i.ActualDuration),
				redrawCh,
			)
		}
		errorCh <- i.Start(ctx, config, start, periodic, end)
	}
	pauseInterval := func() {
		i, err := pomodoro.GetInterval(config)
		if err != nil {
			errorCh <- err
			return
		}
		if err := i.Pause(config); err != nil {
			if errors.Is(err, pomodoro.ErrIntervalNotRunning) {
				return
			}
			errorCh <- err
			return
		}
		w.update([]int{}, i.ActualDuration, "Paused... press start to continue", "", redrawCh)
	}
	btStart, err := button.New("(s)tart", func() error {
		go startInterval()
		return nil
	},
		button.GlobalKey('s'),
		button.WidthFor("(p)ause"),
		button.Height(2),
	)
	if err != nil {
		return nil, err
	}
	btPause, err := button.New("(p)ause", func() error {
		go pauseInterval()
		return nil
	},
		button.FillColor(cell.ColorNumber(220)),
		button.GlobalKey('p'),
		button.Height(2),
	)
	if err != nil {
		return nil, err
	}
	return &buttonSet{btStart, btPause}, nil
}
