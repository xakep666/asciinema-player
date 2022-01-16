package player

import (
	"errors"

	"golang.org/x/sys/windows"
)

func enableVT100(fd int) error {
	var consoleMode uint32

	err := windows.GetConsoleMode(windows.Handle(fd), &consoleMode)
	if err != nil {
		return err
	}

	consoleMode |= windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING

	err = windows.SetConsoleMode(windows.Handle(fd), consoleMode)
	switch {
	case errors.Is(err, nil):
		return nil
	case errors.Is(err, windows.ERROR_INVALID_PARAMETER): // vt100 not supported
		return nil
	default:
		return err
	}
}
