# togif

Renders asciicast to gif image using `fyne` gui toolkit.

Usage:
```
Usage of togif:
  -f string
        path to asciinema v2 file
  -maxWait duration
        maximum time between frames (default 2s)
  -o string
        path to output gif
  -speed float
        speed adjustment: <1 - increase, >1 - decrease (default 1)

```

Required parameters are `-f` and `-o`.

Example gif rendered from `app-demo.cast` with default settings:
![app-demo](demo.gif)
