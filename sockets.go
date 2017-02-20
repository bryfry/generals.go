package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type GeneralsIO struct {
	*websocket.Conn
	Updates chan<- Update
}

const (
	ENGINEIO_OPEN    = 0
	ENGINEIO_CLOSE   = 1
	ENGINEIO_PING    = 2
	ENGINEIO_PONG    = 3
	ENGINEIO_MESSAGE = 4
	SOCKETIO_CONNECT = 0
	SOCKETIO_EVENT   = 2
)

// Packet > Message > Event

type Packet struct {
	Type        int
	MessageType int
	Payload     []byte
}

func (g *GeneralsIO) Connect() (err error) {
	dialer := &websocket.Dialer{}
	dialer.EnableCompression = false
	generals_url := "ws://botws.generals.io/socket.io/?EIO=3&transport=websocket"
	conn, _, err := dialer.Dial(generals_url, nil)
	if err != nil {
		return err
	}
	g.Conn = conn
	return err
}

func (g *GeneralsIO) Emit(msg ...interface{}) error {
	w, err := g.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("42"))
	j := json.NewEncoder(w)
	return j.Encode(msg)
}

func (g *GeneralsIO) Ping(t int) {
	for range time.Tick(time.Duration(t) * time.Second) {
		l.Debug("Ping")
		w, err := g.NextWriter(websocket.TextMessage)
		if err != nil {
			l.Fatal(err)
		}
		_, err = w.Write([]byte("2ping"))
		if err != nil {
			l.Fatal(err)
		}
	}
}

func DecodePacket(b []byte) (p Packet, err error) {
	// TODO: b > 0
	t, err := strconv.Atoi(string(b[0]))
	if err != nil {
		return p, err
	}

	switch t {
	case ENGINEIO_OPEN, ENGINEIO_CLOSE, ENGINEIO_PING, ENGINEIO_PONG:
		p.Type = t
		p.Payload = b[1:]
	case ENGINEIO_MESSAGE:
		p.Type = t
		mt, err := strconv.Atoi(string(b[1]))
		if err != nil {
			return p, err
		}
		p.MessageType = mt
		p.Payload = b[2:]
	}
	return
}

func (g *GeneralsIO) RecievePackets() {
	for {
		_, packet, err := g.ReadMessage()
		if err != nil {
			l.Fatal(err)
		}
		p, err := DecodePacket(packet)
		if err != nil {
			l.Fatal(err)
		}
		switch p.Type {
		case ENGINEIO_PONG:
			l.Debug("Pong")
		case ENGINEIO_OPEN:
			l.Info("Open", zap.String("packet", string(p.Payload)))
		case ENGINEIO_MESSAGE:
			switch p.MessageType {
			case SOCKETIO_CONNECT:
				l.Debug("Connected")
			case SOCKETIO_EVENT:
				// generals.go
				err = DecodeEvent(p.Payload, g.Updates)
				if err != nil {
					l.Error("Unhandled Event", zap.String("packet", string(packet)), zap.Error(err))
				}
			default:
				l.Error("Unhandled Message", zap.String("packet", string(packet)))
			}
		default:
			l.Error("Unhandled Packet", zap.String("packet", string(packet)))
		}

	}

}
