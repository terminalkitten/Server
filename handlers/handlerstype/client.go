package handlerstype

import (
	"github.com/gorilla/websocket"
	"github.com/brokenbydefault/Server/Apps/NanosubApp/nanosubappstypes"
	"github.com/brokenbydefault/Server/Apps/NanoMFAApp/nanomfatypes"
)

type Client struct {
	*websocket.Conn
	*nanosubappstypes.SubClient
	*nanomfaappstypes.MFAClient
}
