package player_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	player "github.com/xakep666/asciinema-player/v3"
)

func TestStreamFrameSource(t *testing.T) {
	cast, err := os.ReadFile(filepath.Join("testdata", "test.cast"))
	if err != nil {
		t.Fatalf("File read failed: %s", err)
	}

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	source, err := player.NewStreamFrameSource(bytes.NewReader(cast))
	if err != nil {
		t.Fatalf("Source create failed: %s", err)
	}

	if err = enc.Encode(source.Header()); err != nil {
		t.Fatalf("Failed to encode header: %s", err)
	}

	for source.Next() {
		frame := source.Frame()
		if err = enc.Encode([]interface{}{frame.Time, frame.Type, string(frame.Data)}); err != nil {
			t.Fatalf("Failed to encode header: %s", err)
		}
	}

	if err = source.Err(); err != nil {
		t.Fatalf("Source error: %s", err)
	}

	if !bytes.Equal(cast, buf.Bytes()) {
		t.Fatalf("Output not equal to input. Output:\n%s", buf.String())
	}
}
