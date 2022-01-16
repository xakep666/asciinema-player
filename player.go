package player

import (
	"fmt"
	"time"
)

var (
	ErrUnexpectedVersion = fmt.Errorf("unexpected asciicast version")
	ErrSmallTerminal     = fmt.Errorf("terminal too small for frames")
)

type Player struct {
	frameSource FrameSource
	terminal    Terminal
	options     options

	pause chan struct{}
	stop  chan struct{}
}

func NewPlayer(frameSource FrameSource, terminal Terminal, opts ...Option) (*Player, error) {
	defaultOptions := options{
		maxWait: 0,
		speed:   1,
	}

	for _, o := range opts {
		o(&defaultOptions)
	}

	termWidth, termHeight := terminal.Dimensions()
	hdr := frameSource.Header()

	if hdr.Version != FormatVersion {
		return nil, ErrUnexpectedVersion
	}

	if !defaultOptions.ignoreSizeCheck && (hdr.Height > termHeight || hdr.Width > termWidth) {
		return nil, ErrSmallTerminal
	}

	p := &Player{
		frameSource: frameSource,
		terminal:    terminal,
		options:     defaultOptions,
		pause:       make(chan struct{}),
		stop:        make(chan struct{}),
	}

	go terminal.Control(p)

	return p, nil
}

// Start starts playback. Method blocks until Stop call.
func (p *Player) Start() (err error) {
	if err = p.terminal.ToRaw(); err != nil {
		return fmt.Errorf("put terminal to raw mode failed: %w", err)
	}

	defer func() {
		if restoreErr := p.terminal.Restore(); restoreErr != nil {
			err = fmt.Errorf("restore terminal failed: %w", restoreErr)
		}
	}()

	timer := time.NewTimer(1)
	<-timer.C // wait for first tick

	prevFrameTime := 0.
	for {
		if !p.frameSource.Next() {
			return p.frameSource.Err()
		}

		frame := p.frameSource.Frame()
		if frame.Type != OutputFrame {
			continue
		}

		timer.Stop()
		timer.Reset(p.nextFrameDelay(frame, prevFrameTime))
		prevFrameTime = frame.Time

		select {
		case <-timer.C:
			// play
		case <-p.pause:
			select {
			case <-p.pause:
				// play
			case <-p.stop:
				return nil
			}
		case <-p.stop:
			return nil
		}

		if _, err = p.terminal.Write(frame.Data); err != nil {
			return fmt.Errorf("frame write failed: %w", err)
		}
	}
}

func (p *Player) nextFrameDelay(frame Frame, prevFrameTime float64) time.Duration {
	delay := time.Duration((frame.Time - prevFrameTime) / p.options.speed * float64(time.Second))
	if p.options.maxWait > 0 && delay > p.options.maxWait {
		return p.options.maxWait
	}

	return delay
}

// Pause pauses playback. If playback already paused it will continue.
func (p *Player) Pause() {
	p.pause <- struct{}{}
}

// Stop interrupts playback. Must be called once.
func (p *Player) Stop() {
	p.stop <- struct{}{}
}

func (p *Player) sealed() {}
