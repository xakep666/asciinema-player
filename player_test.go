package player_test

import (
	"bytes"
	"encoding/base64"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	player "github.com/xakep666/asciinema-player/v3"
)

type bufferTerminal struct {
	bytes.Buffer
	Width, Height int
}

func (b *bufferTerminal) Close() error { return nil }

func (b *bufferTerminal) Dimensions() (width, height int) { return b.Width, b.Height }

func (b *bufferTerminal) ToRaw() error { return nil }

func (b *bufferTerminal) Restore() error { return nil }

func (b *bufferTerminal) Control(control player.PlaybackControl) {}

func TestPlayer(t *testing.T) {
	cast, err := os.ReadFile(filepath.Join("testdata", "test.cast"))
	if err != nil {
		t.Fatalf("Cast read failed: %s", err)
	}

	source, err := player.NewStreamFrameSource(bytes.NewReader(cast))
	if err != nil {
		t.Fatalf("Source create failed: %s", err)
	}

	term := &bufferTerminal{Width: 100, Height: 100}

	p, err := player.NewPlayer(source, term, player.WithMaxWait(1))
	if err != nil {
		t.Fatalf("Player setup failed: %s", err)
	}

	if err = p.Start(); err != nil {
		t.Fatalf("Play failed: %s", err)
	}

	out, err := os.ReadFile(filepath.Join("testdata", "terminal_out.bin"))
	if err != nil {
		t.Fatalf("Terminal out read failed: %s", err)
	}

	if !bytes.Equal(out, term.Bytes()) {
		t.Fatalf("Output mismatch. Run `echo \"%s\" | base64 -d > %s` to update expected output.",
			base64.StdEncoding.EncodeToString(term.Bytes()), filepath.Join("testdata", "terminal_out.bin"),
		)
	}
}

func TestPlayer_PausePlay(t *testing.T) {
	cast, err := os.ReadFile(filepath.Join("testdata", "test.cast"))
	if err != nil {
		t.Fatalf("Cast read failed: %s", err)
	}

	source, err := player.NewStreamFrameSource(bytes.NewReader(cast))
	if err != nil {
		t.Fatalf("Source create failed: %s", err)
	}

	term := &bufferTerminal{Width: 100, Height: 100}

	p, err := player.NewPlayer(source, term, player.WithMaxWait(100*time.Millisecond))
	if err != nil {
		t.Fatalf("Player setup failed: %s", err)
	}

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		defer wg.Done()

		p.Pause()
		p.Pause()
	}()

	if err = p.Start(); err != nil {
		t.Fatalf("Play failed: %s", err)
	}

	wg.Wait()

	out, err := os.ReadFile(filepath.Join("testdata", "terminal_out.bin"))
	if err != nil {
		t.Fatalf("Terminal out read failed: %s", err)
	}

	if !bytes.Equal(out, term.Bytes()) {
		t.Fatalf("Output mismatch. Run `echo \"%s\" | base64 -d > %s` to update expected output.",
			base64.StdEncoding.EncodeToString(term.Bytes()), filepath.Join("testdata", "terminal_out.bin"),
		)
	}
}

func TestPlayer_Stop(t *testing.T) {
	cast, err := os.ReadFile(filepath.Join("testdata", "test.cast"))
	if err != nil {
		t.Fatalf("Cast read failed: %s", err)
	}

	source, err := player.NewStreamFrameSource(bytes.NewReader(cast))
	if err != nil {
		t.Fatalf("Source create failed: %s", err)
	}

	term := &bufferTerminal{Width: 100, Height: 100}

	p, err := player.NewPlayer(source, term, player.WithMaxWait(100*time.Millisecond))
	if err != nil {
		t.Fatalf("Player setup failed: %s", err)
	}

	var wg sync.WaitGroup

	go func() {
		wg.Add(1)
		defer wg.Done()

		time.Sleep(300 * time.Millisecond)
		p.Stop()
	}()

	if err = p.Start(); err != nil {
		t.Fatalf("Play failed: %s", err)
	}

	out, err := os.ReadFile(filepath.Join("testdata", "terminal_out.bin"))
	if err != nil {
		t.Fatalf("Terminal out read failed: %s", err)
	}

	if !bytes.HasPrefix(out, term.Bytes()) || bytes.Equal(out, term.Bytes()) {
		t.Fatalf("Terminal out must be prefix of expected output but not equal to it")
	}
}

func TestPlayer_DisableSizeCheck(t *testing.T) {
	cast, err := os.ReadFile(filepath.Join("testdata", "test.cast"))
	if err != nil {
		t.Fatalf("Cast read failed: %s", err)
	}

	source, err := player.NewStreamFrameSource(bytes.NewReader(cast))
	if err != nil {
		t.Fatalf("Source create failed: %s", err)
	}

	term := &bufferTerminal{Width: 1, Height: 1}

	_, err = player.NewPlayer(source, term)
	if !errors.Is(err, player.ErrSmallTerminal) {
		t.Fatalf("Unexpected error returned: %s, expected ErrSmallTerminal", err)
	}

	_, err = player.NewPlayer(source, term, player.WithIgnoreSizeCheck())
	if err != nil {
		t.Fatalf("Unexpected error returned: %s, expected nil", err)
	}
}
