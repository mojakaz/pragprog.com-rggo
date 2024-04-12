package app

import (
	"context"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/gauge"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum4k/termdash/widgets/text"
	"time"
)

type widgets struct {
	//donTimer       *donut.Donut
	gauTimer *gauge.Gauge
	disType  *segmentdisplay.SegmentDisplay
	txtInfo  *text.Text
	txtTimer *text.Text
	//updateDonTimer chan []int
	updateGauTimer     chan []int
	updateTxtInfo      chan string
	updateTxtTimer     chan string
	updateDurationType chan time.Duration
}

func (w *widgets) update(timer []int, durationType time.Duration, txtInfo, txtTimer string, redrawCh chan<- bool) {
	if txtInfo != "" {
		w.updateTxtInfo <- txtInfo
	}
	if durationType >= 0 {
		w.updateDurationType <- durationType
	}
	if txtTimer != "" {
		w.updateTxtTimer <- txtTimer
	}
	if len(timer) > 0 {
		//w.updateDonTimer <- timer
		w.updateGauTimer <- timer
	}
	redrawCh <- true
}

func newWidgets(ctx context.Context, errorCh chan<- error) (*widgets, error) {
	w := &widgets{}
	var err error
	//w.updateDonTimer = make(chan []int)
	w.updateGauTimer = make(chan []int)
	w.updateDurationType = make(chan time.Duration)
	w.updateTxtInfo = make(chan string)
	w.updateTxtTimer = make(chan string)
	//w.donTimer, err = newDonut(ctx, w.updateDonTimer, errorCh)
	w.gauTimer, err = newGauge(ctx, w.updateGauTimer, errorCh)
	if err != nil {
		return nil, err
	}
	w.disType, err = newSegmentDisplay(ctx, w.updateDurationType, errorCh)
	if err != nil {
		return nil, err
	}
	w.txtInfo, err = newText(ctx, w.updateTxtInfo, errorCh)
	if err != nil {
		return nil, err
	}
	w.txtTimer, err = newText(ctx, w.updateTxtTimer, errorCh)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func newText(ctx context.Context, updateText <-chan string, errorCh chan<- error) (*text.Text, error) {
	txt, err := text.New()
	if err != nil {
		return nil, err
	}
	// Goroutine to update Text
	go func() {
		for {
			select {
			case t := <-updateText:
				txt.Reset()
				errorCh <- txt.Write(t)
			case <-ctx.Done():
				return
			}
		}
	}()
	return txt, nil
}

func newDonut(ctx context.Context, donUpdater <-chan []int, errorCh chan<- error) (*donut.Donut, error) {
	don, err := donut.New(
		donut.Clockwise(),
		donut.CellOpts(cell.FgColor(cell.ColorBlue)),
	)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case d := <-donUpdater:
				if d[0] <= d[1] {
					errorCh <- don.Absolute(d[0], d[1])
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return don, nil
}

func newSegmentDisplay(ctx context.Context, updateDuration <-chan time.Duration, errorCh chan<- error) (*segmentdisplay.SegmentDisplay, error) {
	sd, err := segmentdisplay.New()
	if err != nil {
		return nil, err
	}
	// Goroutine to update SegmentDisplay
	go func() {
		for {
			select {
			case d := <-updateDuration:
				errorCh <- sd.Write([]*segmentdisplay.TextChunk{
					segmentdisplay.NewChunk(d.String()),
				})
			case <-ctx.Done():
				return
			}
		}
	}()
	return sd, nil
}

func newGauge(ctx context.Context, updateGauge <-chan []int, errorCh chan<- error) (*gauge.Gauge, error) {
	gau, err := gauge.New(
		gauge.Height(5),
		gauge.Border(linestyle.Light, cell.FgColor(cell.ColorGreen)),
		gauge.EmptyTextColor(cell.ColorGray),
		gauge.FilledTextColor(cell.ColorCyan),
		gauge.HideTextProgress(),
	)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case g := <-updateGauge:
				if g[0] <= g[1] {
					errorCh <- gau.Absolute(g[0], g[1])
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return gau, nil
}
