# asciinema-player
[![Build Status](https://github.com/xakep666/asciinema-player/actions/workflows/testing.yml/badge.svg)](https://github.com/xakep666/asciinema-player/actions/workflows/testing.yml)
[![codecov](https://codecov.io/gh/xakep666/asciinema-player/branch/master/graph/badge.svg)](https://codecov.io/gh/xakep666/asciinema-player)
[![Go Report Card](https://goreportcard.com/badge/github.com/xakep666/asciinema-player)](https://goreportcard.com/report/github.com/xakep666/asciinema-player)
[![GoDev](https://pkg.go.dev/badge/github.com/xakep666/asciinema-player/pkg/asciicast)](https://godoc.org/github.com/xakep666/asciinema-player/pkg/asciicast)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

asciinema-player is a library and cli-app to play terminal sessions recorded by asciinema (http://github.com/asciinema/asciinema)

## Prerequisites
* Golang >= 1.17

## Installation
Library:
```bash
go get -v -u github.com/xakep666/asciinema-player
```

App:
```bash
go get -v -u github.com/xakep666/asciinema-player/cmd/asciinema-player
```

## Usage
### App
```
$ ./asciinema-player --help
  Usage of ./asciinema-player:
    -f string
          path to asciinema v2 file
    -maxWait duration
          maximum time between frames (default 2s)
    -speed float
          speed adjustment: <1 - increase, >1 - decrease (default 1)
```
For example you can play test session `./asciinema-player -f test.cast`

[![asciicast](https://asciinema.org/a/189343.png)](https://asciinema.org/a/189343)

### Library
```go
frameSource, err := player.NewStreamFrameSource(reader)
if err != nil {
    return err
}

term, err := player.NewOSTerminal()
if err != nil {
    return err
}

defer term.Close()

player, err := player.NewPlayer(frameSource, terminal)
if err != nil {
    return err
}

err = player.Play()
if err != nil {
    return err
}
```
Library usage example is app, actually.

## License
Asciinema-player project is licensed under the terms of the MIT license. Please see LICENSE in this repository for more details.
