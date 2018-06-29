package asciicast

type Header struct {
	Env           map[string]string `json:"env"`
	Width         uint              `json:"width"`
	Height        uint              `json:"height"`
	Timestamp     uint64            `json:"timestamp"` // unix time
	Version       uint              `json:"version"`
	IdleTimeLimit float64           `json:"idle_time_limit"`
}
