package main

import (
	"flag"
	"fmt"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/software"
	"github.com/fyne-io/terminal"

	player "github.com/xakep666/asciinema-player/v3"
)

var (
	filePath   = flag.String("f", "", "path to asciinema v2 file")
	maxWait    = flag.Duration("maxWait", 2*time.Second, "maximum time between frames")
	speed      = flag.Float64("speed", 1, "speed adjustment: <1 - increase, >1 - decrease")
	outputPath = flag.String("o", "", "path to output gif")
)

func main() {
	flag.Parse()

	if *filePath == "" {
		fmt.Println("Please specify asciicast file\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *outputPath == "" {
		fmt.Println("Please specify output file\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Println("Failed to open file", err)
		os.Exit(1)
	}

	src, err := player.NewStreamFrameSource(file)
	if err != nil {
		fmt.Println("Failed to create frame source", err)
		os.Exit(1)
	}

	cnv := software.NewCanvas()
	fyneTerm := terminal.New()
	pr, pw := io.Pipe()

	cnv.SetContent(fyneTerm)
	cnv.Resize(windowSize(src.Header()))
	go fyneTerm.RunWithConnection(nil, pr)

	var (
		images        []*image.Paletted
		delays        []int // 1/100 of second
		prevFrameTime = 0.
	)

	for src.Next() {
		frame := src.Frame()
		if frame.Type != player.OutputFrame {
			continue
		}

		_, err = pw.Write(frame.Data)
		if err != nil {
			fmt.Println("Failed to write to terminal", err)
			os.Exit(1)
		}

		img := cnv.Capture()
		paletted := image.NewPaletted(img.Bounds(), palette.WebSafe)
		draw.Draw(paletted, img.Bounds(), img, image.Point{}, draw.Src)

		images = append(images, paletted)
		delays = append(delays, calcDelay(frame, prevFrameTime, *maxWait, *speed))
		prevFrameTime = frame.Time
	}

	pw.Close()

	outFile, err := os.OpenFile("demo.gif", os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("Output file create failed", err)
		os.Exit(1)
	}

	defer outFile.Close()

	err = gif.EncodeAll(outFile, &gif.GIF{
		Image: images,
		Delay: delays,
	})
	if err != nil {
		fmt.Println("Output gif encode failed", err)
		os.Exit(1)
	}
}

func windowSize(header player.Header) fyne.Size {
	width, height := float32(10), float32(10)

	return fyne.NewSize(width*float32(header.Width), height*float32(header.Height))
}

func calcDelay(frame player.Frame, prevFrameTime float64, maxWait time.Duration, speed float64) int {
	delay := frame.Time - prevFrameTime
	if speed > 0 {
		delay /= speed
	}

	if maxWait > 0 && delay > maxWait.Seconds() {
		return int(100 * maxWait.Seconds())
	}

	return int(100 * delay)
}
