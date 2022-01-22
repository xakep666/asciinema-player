package main

import (
	player "github.com/xakep666/asciinema-player/v3"
	"io/fs"
	"net/http"
	"nhooyr.io/websocket"
)

type TermHandler struct {
	FS      fs.FS
	FileSet map[string]struct{}
}

func (h *TermHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("file")
	if _, ok := h.FileSet[fileName]; !ok {
		http.Error(w, "Requested file was not found", http.StatusNotFound)
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{})
	if err != nil {
		return
	}

	go h.handleConn(fileName, conn)
}

func (h *TermHandler) handleConn(fileName string, conn *websocket.Conn) {
	term, err := NewWSTerm(conn)
	if err != nil {
		conn.Close(websocket.StatusProtocolError, "failed to initiate terminal:"+err.Error())
		return
	}

	defer term.Close()

	file, err := h.FS.Open(fileName)
	if err != nil {
		conn.Close(websocket.StatusProtocolError, "file open failed:"+err.Error())
		return
	}

	defer file.Close()

	src, err := player.NewStreamFrameSource(file)
	if err != nil {
		conn.Close(websocket.StatusProtocolError, "frame source create failed:"+err.Error())
		return
	}

	p, err := player.NewPlayer(src, term)
	if err != nil {
		conn.Close(websocket.StatusProtocolError, "player create failed:"+err.Error())
		return
	}

	if err = p.Start(); err != nil {
		conn.Close(websocket.StatusProtocolError, "playback error:"+err.Error())
	}
}
