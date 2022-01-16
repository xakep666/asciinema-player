package player

// FrameType is a type of Frame.
type FrameType string

const (
	// InputFrame contains data sent from stdin of recorded shell.
	InputFrame FrameType = "i"

	// OutputFrame contains data written to stdout of recorded shell.
	OutputFrame FrameType = "o"
)

const FormatVersion = 2

// Header represents asciinema-v2 header (first line). It doesn't include unneeded fields.
type Header struct {
	// Version is a format version. Must be 2.
	Version int `json:"version"`

	// With is a captured terminal width.
	Width int `json:"width"`

	// Height is a captured terminal height.
	Height int `json:"height"`
}

// FrameSource describes frames source.
type FrameSource interface {
	// Header returns asciinema-v2 header.
	Header() Header

	// Next advances to next available frame. It must return false if error occurs or there is no more frames.
	Next() bool

	// Frame returns current frame. It becomes unusable after Next call.
	Frame() Frame

	// Err returns error if it happens during iteration.
	Err() error
}
