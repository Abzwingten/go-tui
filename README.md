[![GoDoc](https://godoc.org/github.com/gasiordev/go-tui?status.svg)](https://godoc.org/github.com/gasiordev/go-tui)
[![Build Status](https://travis-ci.org/gasiordev/go-tui.svg?branch=master)](https://travis-ci.org/gasiordev/go-tui)

# go-tui

Package `gasiordev/go-tui` is meant to make creating terminal user interface
easier.

### Install

Ensure you have your
[workspace directory](https://golang.org/doc/code.html#Workspaces) created and
run the following:

```
go get -u github.com/gasiordev/go-tui
```

### Example

```
package main

import (
    "os"
	"github.com/gasiordev/ntree/tui"
)

func getOnTUIDraw(n *NTree) func(*tui.TUI) int {
	fn := func(c *tui.TUI) int {
		return 0
	}
	return fn
}

func getOnTUIPaneDraw(n *NTree, p *tui.TUIPane) func(*tui.TUIPane) int {
	t := tui.NewTUIWidgetSample()
	t.InitPane(p)
	fn := func(x *tui.TUIPane) int {
		return t.Run(x)
	}
	return fn
}

func main() {
	myTUI := tui.NewTUI("My Project", "Its description", "Author")
	myTUI.SetOnDraw(getOnTUIDraw(n))

	p0 := myTUI.GetPane()

	p01, p02 := p0.SplitVertically(-50, tui.UNIT_PERCENT)
	p021, p022 := p02.SplitVertically(-40, tui.UNIT_CHAR)

	p11, p12 := p01.SplitHorizontally(20, tui.UNIT_CHAR)
	p21, p22 := p021.SplitHorizontally(50, tui.UNIT_PERCENT)
	p31, p32 := p022.SplitHorizontally(-35, tui.UNIT_CHAR)

	s1 := tui.NewTUIPaneStyleFrame()
	s2 := tui.NewTUIPaneStyleMargin()

	s3 := &tui.TUIPaneStyle{
		NE: "/", NW: "\\", SE: " ", SW: " ", E: " ", W: " ", N: "_", S: " ",
	}

	p11.SetStyle(s1)
	p12.SetStyle(s1)
	p21.SetStyle(s2)
	p22.SetStyle(s2)
	p31.SetStyle(s3)
	p32.SetStyle(s1)

	p11.SetOnDraw(getOnTUIPaneDraw(n, p11))
	p12.SetOnDraw(getOnTUIPaneDraw(n, p12))
	p21.SetOnDraw(getOnTUIPaneDraw(n, p21))
	p22.SetOnDraw(getOnTUIPaneDraw(n, p22))
	p31.SetOnDraw(getOnTUIPaneDraw(n, p31))
	p32.SetOnDraw(getOnTUIPaneDraw(n, p32))

	p11.SetOnIterate(getOnTUIPaneDraw(n, p11))
	p12.SetOnIterate(getOnTUIPaneDraw(n, p12))
	p21.SetOnIterate(getOnTUIPaneDraw(n, p21))
	p22.SetOnIterate(getOnTUIPaneDraw(n, p22))
	p31.SetOnIterate(getOnTUIPaneDraw(n, p31))
	p32.SetOnIterate(getOnTUIPaneDraw(n, p32))

    myTUI.Run(os.Stdout, os.Stderr)
}
```
