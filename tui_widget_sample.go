package tui

import (
	"time"
)

type TUIWidgetSample struct {
}

func (w *TUIWidgetSample) InitPane(p *TUIPane) {
	p.SetMinWidth(5)
	p.SetMinHeight(3)
}

func (w *TUIWidgetSample) Run(p *TUIPane) int {
	t := time.Now()
	p.Write(0, 0, t.Format("15:04:05"), false)
	return 1
}

func NewTUIWidgetSample() *TUIWidgetSample {
	w := &TUIWidgetSample{}
	return w
}
