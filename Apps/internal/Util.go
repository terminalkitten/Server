package internal

import (
	"github.com/gorilla/websocket"
	"io"
	"encoding/json"
)

var NOT_ALLOWED = []byte(`{"error":"method not allowed"}`)
var INTERNAL_ERROR = []byte(`{"error":"internal error"}`)
var SENDER_ERROR = []byte(`{"error":"data is invalid"}`)

func ReplyMessaging(m io.ReadCloser, c *websocket.Conn) error {
	w, err := c.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, m)
	if err != nil {
		return err
	}
	w.Close()
	return nil
}

func ReplyMessage(m []byte, c *websocket.Conn) error {
	return c.WriteMessage(websocket.TextMessage, m)
}

func ReplyJSON(m interface{}, c *websocket.Conn) error {
	j, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	return ReplyMessage(j, c)
}
