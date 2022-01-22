package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"

	player "github.com/xakep666/asciinema-player/v3"
)

type MessageType int

const (
	DimensionsMessage MessageType = iota
	DataMessage
	PlayPauseMessage
	StopMessage
)

type Dimensions struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Message struct {
	Type MessageType `json:"type"`

	Dimensions *Dimensions `json:"dimensions,omitempty"`
	Data       string      `json:"data,omitempty"`
}

type WSTerm struct {
	conn *websocket.Conn
	stop chan struct{}

	dimensions Dimensions
}

func NewWSTerm(conn *websocket.Conn) (*WSTerm, error) {
	var msg Message
	if err := wsjson.Read(context.Background(), conn, &msg); err != nil {
		return nil, fmt.Errorf("dimensions read error: %w", err)
	}

	if msg.Type != DimensionsMessage {
		return nil, fmt.Errorf("first message was not about dimensions")
	}

	return &WSTerm{
		conn: conn,
		stop: make(chan struct{}),

		dimensions: *msg.Dimensions,
	}, nil
}

func (t *WSTerm) Write(p []byte) (n int, err error) {
	return len(p), wsjson.Write(context.Background(), t.conn, Message{
		Type: DataMessage,
		Data: string(p),
	})
}

func (t *WSTerm) Close() error {
	close(t.stop)
	return t.conn.Close(websocket.StatusNormalClosure, "goodbye")
}

func (t *WSTerm) Dimensions() (width, height int) {
	return t.dimensions.Width, t.dimensions.Height
}

func (t *WSTerm) ToRaw() error { return nil }

func (t *WSTerm) Restore() error { return nil }

func (t *WSTerm) Control(control player.PlaybackControl) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-t.stop
		cancel()
	}()

	for {
		var msg Message

		err := wsjson.Read(ctx, t.conn, &msg)
		switch {
		case errors.Is(err, nil):
		case errors.Is(err, context.Canceled):
			return
		default:
			log.Printf("read error: %s", err)
			return
		}

		switch msg.Type {
		case PlayPauseMessage:
			control.Pause()
		case StopMessage:
			control.Stop()
		}
	}
}
