//go:build !windows
// +build !windows

package player

func enableVT100(fd int) error { return nil }
