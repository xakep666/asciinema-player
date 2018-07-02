package asciicast

import (
	"math"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/xakep666/asciinema-player/pkg/terminal"
)

// Player is an interface for playing asciicasts.
type Player interface {
	Play(asciicast *Asciicast, maxWait, speed float64) error
}

// TerminalPlayer is a asciicast to terminal player.
type TerminalPlayer struct {
	Terminal terminal.Terminal
}

// Play plays provided asciicast.
// On *nix systems it firstly puts terminal to raw mode (will be restored after finish).
// Playing is actually putting frame data to terminal without escaping.
// Player can be interrupted by hitting Ctrl-C.
// Player can be paused and unpaused by hitting space key.
// If player paused you can switch to next frame by pressing tab key.
func (p *TerminalPlayer) Play(asciicast *Asciicast, maxWait time.Duration, speed float64) error {
	p.Terminal.ToRaw()
	defer p.Terminal.Reset()

	stdout := asciicast.Frames.
		Filter(IsOutputFrame).
		ToRelativeTime().
		CapRelativeTime(math.Min(asciicast.Header.IdleTimeLimit, maxWait.Seconds())).
		ToAbsoluteTime().
		AdjustSpeed(speed)

	baseTime := time.Now()
	ctrlC := false
	paused := false
	var pauseTime time.Time

	for _, frame := range stdout {
		delay := frame.Duration() - time.Now().Sub(baseTime)

	btnLoop:
		for !ctrlC && delay > 0 {
			if paused {
			pauseLoop:
				for {
					event, err := p.Terminal.TimeoutEvent(time.Second)
					switch err {
					case nil, terminal.ErrEventTimeout:
						// pass
					default:
						return err
					}

					if event.Type != termbox.EventKey {
						continue
					}
					switch event.Key {
					case termbox.KeyCtrlC:
						ctrlC = true
						break pauseLoop
					case termbox.KeySpace:
						paused = false
						baseTime = baseTime.Add(time.Now().Sub(pauseTime))
						break pauseLoop
					case termbox.KeyTab:
						delay = 0
						pauseTime = time.Now()
						baseTime = pauseTime.Add(-frame.Duration())
						break pauseLoop
					}
				}
			} else {
				event, err := p.Terminal.TimeoutEvent(delay)
				switch err {
				case nil:
					// pass
				case terminal.ErrEventTimeout:
					break btnLoop
				default:
					return err
				}

				if event.Type != termbox.EventKey {
					continue
				}
				switch event.Key {
				case termbox.KeyCtrlC:
					ctrlC = true
					break btnLoop
				case termbox.KeySpace:
					paused = true
					pauseTime = time.Now()
					slept := frame.Duration() - pauseTime.Sub(baseTime)
					delay -= slept
				}
			}

		}
		if ctrlC {
			break
		}
		if _, err := p.Terminal.Write(frame.Data); err != nil {
			return err
		}
	}

	return nil
}
