package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	player "github.com/xakep666/asciinema-player/v3"
)

func errExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	maxWait  time.Duration
	speed    float64
	filePath string
)

func init() {
	flag.DurationVar(&maxWait, "maxWait", 2*time.Second, "maximum time between frames")
	flag.Float64Var(&speed, "speed", 1, "speed adjustment: <1 - increase, >1 - decrease")
	flag.StringVar(&filePath, "f", "", "path to asciinema v2 file")
	flag.Parse()
}

func main() {
	if filePath == "" {
		fmt.Println("Please specify file\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(filePath)
	errExit(err)
	defer file.Close()

	source, err := player.NewStreamFrameSource(file)
	errExit(err)

	term, err := player.NewOSTerminal()
	errExit(err)
	defer term.Close()

	p, err := player.NewPlayer(source, term, player.WithSpeed(speed), player.WithMaxWait(maxWait), player.WithIgnoreSizeCheck())
	errExit(err)

	err = p.Start()
	if err != nil {
		fmt.Println("Start failed:", err)
	}
}
