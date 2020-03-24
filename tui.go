package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type TUI struct {
	name   string
	desc   string
	author string
	stdout *os.File
	stderr *os.File
	h      int
	w      int
	pane   *TUIPane
	onDraw func(*TUI) int
}

func (t *TUI) GetName() string {
	return t.name
}

func (t *TUI) GetDesc() string {
	return t.desc
}

func (t *TUI) GetAuthor() string {
	return t.author
}

func (t *TUI) GetStdout() *os.File {
	return t.stdout
}

func (t *TUI) GetStderr() *os.File {
	return t.stderr
}

func (t *TUI) GetPane() *TUIPane {
	return t.pane
}

func (t *TUI) GetWidth() int {
	return t.w
}

func (t *TUI) GetHeight() int {
	return t.h
}

func (t *TUI) Run(stdout *os.File, stderr *os.File) int {
	t.stdout = stdout
	t.stderr = stderr

	t.clear()

	done := make(chan bool)
	go t.startMainLoop()
	<-done

	return 0
}

func (t *TUI) SetOnDraw(f func(*TUI) int) {
	t.onDraw = f
}

func (t *TUI) SetPane(p *TUIPane) {
	t.pane = p
}

func (t *TUI) Write(x int, y int, s string) {
	fmt.Fprintf(t.stdout, "\u001b[1000A\u001b[1000D")
	if x > 0 {
		fmt.Fprintf(t.stdout, "\u001b["+strconv.Itoa(x)+"C")
	}
	if y > 0 {
		fmt.Fprintf(t.stdout, "\u001b["+strconv.Itoa(y)+"B")
	}
	fmt.Fprintf(t.stdout, s)
}

func (t *TUI) getSize() (int, int, error) {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	nums := strings.Split(string(out), " ")
	h, err := strconv.Atoi(nums[0])
	if err != nil {
		return 0, 0, err
	}
	w, err := strconv.Atoi(strings.Replace(nums[1], "\n", "", 1))
	if err != nil {
		return 0, 0, err
	}
	return w, h, nil
}

func (t *TUI) refreshSize() bool {
	w, h, err := t.getSize()
	if err != nil {
		return false
	}
	if t.w != w || t.h != h {
		t.w = w
		t.h = h

		t.pane.SetWidth(w)
		t.pane.SetHeight(h)
		return true
	}
	return false
}

func (t *TUI) clear() {
	fmt.Fprintf(t.stdout, "\u001b[2J\u001b[1000A\u001b[1000D")
}

func (t *TUI) startMainLoop() {
	for {
		sizeChanged := t.refreshSize()
		if sizeChanged {
			t.clear()
			if t.onDraw != nil {
				t.onDraw(t)
			}
			t.pane.Draw()
		}
		t.pane.Iterate()
		time.Sleep(time.Millisecond * time.Duration(1000))
	}
}

func NewTUI(n string, d string, a string) *TUI {
	t := &TUI{name: n, desc: d, author: a}
	p := NewTUIPane("main", t)
	t.SetPane(p)
	return t
}
