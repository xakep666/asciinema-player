package terminal

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/nsf/termbox-go"
	"golang.org/x/crypto/ssh/terminal"
)

var pty *Pty

func TestMain(t *testing.M) {
	if !terminal.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Println("Not running terminal tests because stdin is not terminal")
		fmt.Println("To run terminal tests build a test binary (\"go build -o test ./pkg/terminal\") and run it (\"./test -test.v\")")
		return
	}
	var err error
	pty, err = NewPty()
	if err != nil {
		fmt.Printf("PTY construct error: %v\n", err)
		return
	}
	defer func() {
		if err := pty.Close(); err != nil {
			fmt.Printf("PTY close error: %v\n", err)
		}
	}()
	t.Run()
}

func TestTerminalSize(t *testing.T) {
	curWidth, curHeight := termbox.Size()
	t.Logf("current terminal size: width=%d, height=%d", curWidth, curHeight)

	pty, err := NewPty()
	if err != nil {
		t.Errorf("PTY construct error: %v", err)
		return
	}
	defer func() {
		if err := pty.Close(); err != nil {
			t.Errorf("PTY close error: %v", err)
		}
	}()

	width, height, err := pty.Size()
	if err != nil {
		t.Errorf("get size error: %v", err)
		return
	}

	if width != curWidth || height != curHeight {
		t.Errorf("unexpected sizes: width=%d, height=%d", width, height)
		return
	}
}

func TestTerminalRaw(t *testing.T) {
	oldState, err := terminal.GetState(int(os.Stdin.Fd()))
	if err != nil {
		t.Errorf("get terminal state failed: %v", err)
		return
	}

	if err := pty.ToRaw(); err != nil {
		t.Errorf("putting terminal to raw mode failed: %v", err)
		return
	}
	if err := pty.Reset(); err != nil {
		t.Errorf("restoring terminal failed: %v", err)
		return
	}

	newState, err := terminal.GetState(int(os.Stdin.Fd()))
	if err != nil {
		t.Errorf("get terminal state failed: %v", err)
		return
	}

	if !reflect.DeepEqual(oldState, newState) {
		t.Errorf("new and old terminal state differs: new: %#v\n, old: %#v", newState, oldState)
		return
	}
}

func TestTerminalEventPoll(t *testing.T) {
	event, err := pty.TimeoutEvent(time.Second)
	switch err {
	case nil:
		// pass
	case ErrEventTimeout:
		t.Logf("event timeout")
		return
	default:
		t.Errorf("unexpected poll error: %v", err)
		return
	}
	t.Logf("event caught: %#v", event)
}

func TestTerminalWrite(t *testing.T) {
	testStr := "terminal_test\n"
	n, err := pty.Write([]byte(testStr))
	if err != nil {
		t.Errorf("unexpected write error: %v", err)
		return
	}
	if n != len(testStr) {
		t.Errorf("wanted to write %d bytes, wrote %d bytes", len(testStr), n)
		return
	}
}
