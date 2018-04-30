package internal

import (
	"github.com/gorilla/websocket"
	"io"
	"encoding/json"
	"strings"
	"github.com/brokenbydefault/Server/handlers/handlerstype"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"bytes"
	"github.com/brokenbydefault/Nanollet/RPC"
)

var NOT_ALLOWED = []byte(`{"error":"method not allowed"}`)
var INTERNAL_ERROR = []byte(`{"error":"internal error"}`)
var SENDER_ERROR = []byte(`{"error":"data is invalid"}`)
var RECEIVER_NOT_FOUND = []byte(`{"error":"receiver not found"}`)

func ReplyMessaging(m io.ReadCloser, c *handlerstype.Client) error {
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

func ReplyMessage(m []byte, c *handlerstype.Client) error {
	return c.WriteMessage(websocket.TextMessage, m)
}

func ReplyJSON(m interface{}, c *handlerstype.Client) error {
	j, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	return ReplyMessage(j, c)
}

func IsActionAllowed(action string, whitelist []string) bool {
	action = strings.ToLower(action)
	for _, ac := range whitelist {
		if action == ac {
			return true
		}
	}
	return false
}

func Subscribe(pk Wallet.PublicKey, c *handlerstype.Client, list map[[32]byte][]*handlerstype.Client) {
	var pkb [32]byte
	copy(pkb[:], pk)

	// If the connection already subscribed to this public-key we return
	if bytes.Equal(c.SubClient.PK, pk) {
		return
	}

	if c.SubClient.PK != nil {
		Unsubscribe(c, list)
	}

	if subs, ok := list[pkb]; ok {
		list[pkb] = append(subs, c)
	} else {
		list[pkb] = []*handlerstype.Client{c}
	}

	c.SubClient.PK = pk

	m := RPCClient.Subscription{
		PublicKey: c.SubClient.PK,
	}
	ReplyJSON(&m, c)
}

func Unsubscribe(c *handlerstype.Client, list map[[32]byte][]*handlerstype.Client) {
	var pkb [32]byte
	copy(pkb[:], c.SubClient.PK)

	if c.SubClient.PK == nil {
		return
	}

	if subs, ok := list[pkb]; ok {
		if len(subs) >= 2 {
			cons := make([]*handlerstype.Client, len(subs)-1)

			var index int
			for _, con := range subs {
				if c != con {
					cons[index] = con
				}
				index++
			}

			list[pkb] = cons
		} else {
			delete(list, pkb)
		}
	}

	c.SubClient.PK = nil

	m := RPCClient.Subscription{
		PublicKey: c.SubClient.PK,
	}
	ReplyJSON(&m, c)
}