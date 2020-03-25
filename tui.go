package tui

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// TUI is main interface definition. It has a name, description, author
// (which are not used anywhere yet), current terminal width and height,
// pointer to main pane, pointer to a function that is triggered when
// interface is being drawn (that happens when app is started and when
// terminal size is changed), and finally pointers to standard output
// and standard error File instances.
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

// GetName returns TUI name
func (t *TUI) GetName() string {
	return t.name
}

// GetDesc returns TUI description
func (t *TUI) GetDesc() string {
	return t.desc
}

// GetAuthor returns TUI author
func (t *TUI) GetAuthor() string {
	return t.author
}

// GetStdout returns stdout property
func (t *TUI) GetStdout() *os.File {
	return t.stdout
}

// GetStderr returns stderr property
func (t *TUI) GetStderr() *os.File {
	return t.stderr
}

// GetPane returns initial/first terminal pane
func (t *TUI) GetPane() *TUIPane {
	return t.pane
}

// GetWidth returns cached terminal width
func (t *TUI) GetWidth() int {
	return t.w
}

// GetHeight returns cached terminal height
func (t *TUI) GetHeight() int {
	return t.h
}

// Run clears the terminal and starts program's main loop
func (t *TUI) Run(stdout *os.File, stderr *os.File) int {
	t.stdout = stdout
	t.stderr = stderr

	t.clear()

	done := make(chan bool)
	go t.startMainLoop()
	<-done

	return 0
}

// SetOnDraw attaches function that will be triggered when interface is being
// drawn (what happens on initialisation and terminal resize)
func (t *TUI) SetOnDraw(f func(*TUI) int) {
	t.onDraw = f
}

// SetPane sets the main terminal pane
func (t *TUI) SetPane(p *TUIPane) {
	t.pane = p
}

// Write prints out on the terminal window at a specified position
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

// getSize gets terminal size by calling stty command
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

// refreshSize gets terminal size and caches it
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

// clear clears terminal window
func (t *TUI) clear() {
	fmt.Fprintf(t.stdout, "\u001b[2J\u001b[1000A\u001b[1000D")
}

// startMainLoop initialises program's main loop, controls the terminal size,
// ensures panes are correctly drawn and calls methods attached to their
// onIterate property
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

// NewTUI creates new instance of TUI and returns it
func NewTUI(n string, d string, a string) *TUI {
	t := &TUI{name: n, desc: d, author: a}
	p := NewTUIPane("main", t)
	t.SetPane(p)
	return t
}
