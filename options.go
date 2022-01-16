package player

import (
	"time"
)

type options struct {
	maxWait         time.Duration
	speed           float64
	ignoreSizeCheck bool
}

// Option for Player.
type Option func(*options)

// WithMaxWait sets minimal delay between frames. Zero or negative value are ignored.
func WithMaxWait(t time.Duration) Option {
	return func(o *options) {
		if t > 0 {
			o.maxWait = t
		}
	}
}

// WithSpeed sets playback speed.
// Values greater than 1 speeds up playback.
// Values between 0 and 1 slows down playback.
// Negative values are ignored.
func WithSpeed(speed float64) Option {
	return func(o *options) {
		if speed > 0 {
			o.speed = speed
		}
	}
}

//WithIgnoreSizeCheck turns off check that terminal can fit frames.
func WithIgnoreSizeCheck() Option {
	return func(o *options) {
		o.ignoreSizeCheck = true
	}
}
