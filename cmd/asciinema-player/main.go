package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/xakep666/asciinema-player/pkg/asciicast"
	"github.com/xakep666/asciinema-player/pkg/parser"
	"github.com/xakep666/asciinema-player/pkg/terminal"
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

	parsed, err := parser.Parse(file)
	errExit(err)

	term, err := terminal.NewPty()
	errExit(err)

	tp := &asciicast.TerminalPlayer{Terminal: term}

	err = tp.Play(parsed, maxWait, speed)
	errExit(err)
}
