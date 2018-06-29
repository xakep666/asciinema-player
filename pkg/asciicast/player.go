package asciicast

import (
	"math"
	"time"

	termbox "github.com/nsf/termbox-go"
	"github.com/xakep666/asciinema-player/pkg/terminal"
)

type Player interface {
	Play(asciicast *Asciicast, maxWait, speed float64) error
}

type TerminalPlayer struct {
	Terminal terminal.Terminal
}

func NewTerminalPlayer() (TerminalPlayer, error) {
	term, err := terminal.NewPty()
	if err != nil {
		return TerminalPlayer{}, err
	}
	return TerminalPlayer{Terminal: term}, nil
}

func (p *TerminalPlayer) Play(asciicast *Asciicast, maxWait, speed float64) error {
	p.Terminal.ToRaw()
	defer p.Terminal.Reset()

	stdout := asciicast.Frames.
		Filter(IsOutputFrame).
		ToRelativeTime().
		CapRelativeTime(math.Min(asciicast.Header.IdleTimeLimit, maxWait)).
		ToAbsoluteTime().
		AdjustSpeed(speed)

	baseTime := time.Now()
	ctrlC := false
	paused := false
	var pauseTime time.Time

	for _, frame := range stdout {
		delay := frame.Time - time.Now().Sub(baseTime).Seconds()

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
						baseTime.Add(time.Now().Sub(pauseTime))
						break pauseLoop
					case 0x2e: // period (dot)
						delay = 0
						pauseTime = time.Now()
						baseTime = time.Unix(pauseTime.Unix()-int64(frame.Time*float64(time.Second)), 0)
						break pauseLoop
					}
				}
			} else {
				event, err := p.Terminal.TimeoutEvent(time.Duration(float64(time.Second) * delay))
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
					slept := frame.Time - pauseTime.Sub(baseTime).Seconds()
					delay -= slept
				}
			}

		}
		if ctrlC {
			break
		}
		p.Terminal.Write(frame.Data)
	}

	return nil
}
