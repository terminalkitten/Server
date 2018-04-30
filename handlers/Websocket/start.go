package Websocket

import (
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/buger/jsonparser"
	"github.com/brokenbydefault/Server/Apps/NanolletApp"
	"github.com/brokenbydefault/Server/Apps/NanofyApp"
	"log"
	"github.com/brokenbydefault/Server/Apps/NanosubApp"
	"github.com/brokenbydefault/Server/handlers/handlerstype"
	"github.com/brokenbydefault/Server/Apps/NanosubApp/nanosubappstypes"
	"github.com/brokenbydefault/Server/Apps/NanoMFAApp"
	"github.com/brokenbydefault/Server/Apps/NanoMFAApp/nanomfatypes"
)

var SocketUpgrade = &websocket.Upgrader{
	ReadBufferSize:  1 << 14,
	WriteBufferSize: 1 << 18,

	// It's need to everyone can use the API, even on other website. ;)
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Start(w http.ResponseWriter, r *http.Request) {
	conn, err := SocketUpgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	c := &handlerstype.Client{
		Conn:      conn,
		SubClient: new(nanosubappstypes.SubClient),
		MFAClient: new(nanomfaappstypes.MFAClient),
	}
	listening(c)
}

func listening(c *handlerstype.Client) {
	defer c.Close()
	c.SetReadLimit(1 << 13)

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Print(err)
			}

			NanosubApp.Unsubscribe(c)
			NanoMFAApp.Unsubscribe(c)
			break
		}
		redirect(c, message)
	}
}

func redirect(c *handlerstype.Client, m []byte) {
	app, err := jsonparser.GetString(m, "app")
	if err != nil {
		return
	}

	switch app {
	case "nanollet":
		NanolletApp.StartMessaging(m, c)
	case "nanofy":
		NanofyApp.StartMessaging(m, c)
	case "nanosub":
		NanosubApp.StartMessaging(m, c)
	case "nanomfa":
		NanoMFAApp.StartMessaging(m, c)
	default:
		c.WriteMessage(0, []byte(`{"error":"invalid"}`))
	}
}
