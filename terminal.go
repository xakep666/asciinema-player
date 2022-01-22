package player

import (
	"io"
)

// Terminal is interface for terminal interaction.
type Terminal interface {
	io.WriteCloser

	// Dimensions returns terminal window size.
	Dimensions() (width, height int)

	// ToRaw puts terminal to raw mode. Implementation must store previous terminal state.
	ToRaw() error

	// Restore puts terminal into state stored in ToRaw.
	Restore() error

	// Control starts "event loop" where Terminal may call methods of PlaybackControl. Method blocks until Close.
	Control(PlaybackControl)
}

// PlaybackControl describes playback control methods for Terminal.
type PlaybackControl interface {
	// Pause pauses playback. If playback already paused it will continue.
	Pause()

	// Stop interrupts playback. Must be called once.
	Stop()

	sealed()
}
