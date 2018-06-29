# asciinema-player
[![Go Report Card](https://goreportcard.com/badge/github.com/xakep666/asciinema-player)](https://goreportcard.com/report/github.com/xakep666/asciinema-player) [![GoDoc](https://godoc.org/github.com/xakep666/asciinema-player/pkg/asciicast?status.svg)](https://godoc.org/github.com/xakep666/asciinema-player/pkg/asciicast) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

asciinema-player is a library and cli-app to play terminal sessions recorded by asciinema (http://github.com/asciinema/asciinema)

## Prerequisites
* Golang >= 1.10 or Vgo

## Installation
Library:
```bash
go get -v -u github.com/xakep666/pkg/asciicast
```

App:
```bash
go get -v -u github.com/xakep666/cmd/asciinema-player
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
    parsed, err := parser.Parse(file)
	if err != nil {
	    return err
	}

	tp, err := asciicast.NewTerminalPlayer()
	if err != nil {
        return err
    }

	err = tp.Play(parsed, maxWait, speed)
	if err != nil {
        return err
    }
```
Library usage example is app actually.

## License
Asciinema-player project is licensed under the terms of the MIT license. Please see LICENSE in this repository for more details.