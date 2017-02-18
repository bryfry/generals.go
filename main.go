package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type SocketIO struct {
	*websocket.Conn
}

func (s *SocketIO) Emit(msg ...interface{}) error {
	w, err := s.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("42"))
	j := json.NewEncoder(w)
	return j.Encode(msg)
}
func (s *SocketIO) Ping() error {
	w, err := s.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("2ping"))
	return err
}

func main() {
	generals_url := "ws://botws.generals.io/socket.io/?EIO=3&transport=websocket"
	user_id := "k_00003"
	username := "k_bot3"

	dialer := &websocket.Dialer{}
	dialer.EnableCompression = false

	conn, _, err := dialer.Dial(generals_url, nil)
	c := &SocketIO{conn}
	if err != nil {
		log.Fatal(err)
	}
	err = c.Emit("set_username", user_id, username)
	if err != nil {
		log.Fatal(err)
	}
	err = c.Emit("join_private", "test", user_id)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for range time.Tick(5 * time.Second) {
			c.Ping()
		}
	}()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Got: ", string(message))
		dec := json.NewDecoder(bytes.NewBuffer(message))
		var msgType int
		dec.Decode(&msgType)

		if msgType == 42 {
			var raw json.RawMessage
			dec.Decode(&raw)
			eventname := ""
			data := []interface{}{&eventname}
			json.Unmarshal(raw, &data)
			//if f, ok := c.events[eventname]; ok {
			//	f(raw)
			//}
		}
	}
}
