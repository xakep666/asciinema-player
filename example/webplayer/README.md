# webplayer

A simple web player for casts stored on server.

Usage:
```
Usage of webplayer:
  -listen-addr string
        HTTP server listen address
  -rootdir string
        Root directory for casts (default ".")
```

I.e. run with `go run -v . -rootdir=../.. -listen-addr=:8080` and go to http://localhost:8080/
You will see some available casts, just click on it to start playback.
