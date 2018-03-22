package Websocket

import (
	"github.com/gorilla/websocket"
	"github.com/buger/jsonparser"
	"github.com/brokenbydefault/Server/Apps/NanolletApp"
	"github.com/brokenbydefault/Server/Apps/NanofyApp"
	"log"
)

func readMessage(m []byte, c *websocket.Conn) {
	app, err := jsonparser.GetString(m, "app")
	if err != nil {
		return
	}

	switch app {
	case "nanollet":
		NanolletApp.StartMessaging(m, c)
	case "nanofy":
		NanofyApp.StartMessaging(m, c)
		//case "nanonitor":
	default:
		c.WriteMessage(0, []byte("invalid"))
	}
}

func listen(c *websocket.Conn) {
	defer c.Close()
	c.SetReadLimit(1 << 13)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Print(err)
			}
			break
		}
		readMessage(message, c)
	}
}
