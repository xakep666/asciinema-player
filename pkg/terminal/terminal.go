package terminal

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/creack/termios/raw"
	"github.com/kr/pty"
	"github.com/nsf/termbox-go"
	"golang.org/x/crypto/ssh/terminal"
)

type Terminal interface {
	Size() (int, int, error)
	Write([]byte) error
	ToRaw() error
	Reset() error
	TimeoutEvent(timeout time.Duration) (termbox.Event, error)

	io.Closer
}

var ErrEventTimeout = errors.New("event timeout")

type Pty struct {
	Stdin     *os.File
	Stdout    *os.File
	prevState *raw.Termios
}

func NewPty() (*Pty, error) {
	if err := termbox.Init(); err != nil {
		return nil, err
	}
	return &Pty{Stdin: os.Stdin, Stdout: os.Stdout}, nil
}

func (p *Pty) Size() (int, int, error) {
	return pty.Getsize(p.Stdout)
}

func (p *Pty) Write(data []byte) error {
	_, err := p.Stdout.Write(data)
	if err != nil {
		return err
	}

	err = p.Stdout.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (p *Pty) ToRaw() error {
	fd := p.Stdin.Fd()
	var err error
	if terminal.IsTerminal(int(fd)) {
		p.prevState, err = raw.MakeRaw(fd)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pty) Reset() error {
	if p.prevState == nil {
		return nil
	}
	return raw.TcSetAttr(p.Stdin.Fd(), p.prevState)
}

func (p *Pty) TimeoutEvent(timeout time.Duration) (termbox.Event, error) {
	ev := make(chan termbox.Event)

	go func() {
		ev <- termbox.PollEvent()
	}()

	after := time.After(timeout)

	select {
	case <-after:
		return termbox.Event{}, ErrEventTimeout
	case event := <-ev:
		return event, nil
	}

}

func (p *Pty) Close() error {
	p.Reset()
	termbox.Close()
	return nil
}
