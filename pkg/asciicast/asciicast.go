package asciicast

// Asciicast represents asciicast v2 file.
// First line of file is Header JSON-object. It consists meta information: version, environment variables, terminal size, etc.
// Next lines is JSON-arrays (one per line) with frame delay, frame type and frame data (escaped string).
type Asciicast struct {
	Header Header
	Frames Frames
}
