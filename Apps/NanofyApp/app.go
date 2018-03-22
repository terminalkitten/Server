package NanofyApp

import (
	"github.com/gorilla/websocket"
	"github.com/buger/jsonparser"
	"github.com/brokenbydefault/Server/Apps/internal"
	"strings"
	"encoding/json"
	"github.com/brokenbydefault/Server/SQL"
	"github.com/brokenbydefault/Nanofy/nanofytypes"
	"database/sql"
)

var whitelist = [...]string{"file"}

func Start() {
	go TrackNV0()
}

func StartMessaging(m []byte, c *websocket.Conn) {
	action, err := jsonparser.GetString(m, "action")
	if err != nil || !isActionAllowed(action) {
		internal.ReplyMessage(internal.NOT_ALLOWED, c)
	}

	if action == "file" {
		req := nanofytypes.RequestByFile{}
		resp := nanofytypes.Response{}

		err = json.Unmarshal(m, &req)
		if err != nil {
			resp.Error = string(internal.INTERNAL_ERROR)
			internal.ReplyJSON(m, c)
			return
		}

		var PubKey, FlagHash, SigHash []byte
		err := SQL.Connection.QueryRow("SELECT PubKey, FlagBlock, SigBlock FROM history WHERE PubKey = ? AND FileKey = ?", []byte(req.PubKey), []byte(req.FileKey)).Scan(&PubKey, &FlagHash, &SigHash)
		if err != nil && err != sql.ErrNoRows {
			resp.Error = string(internal.INTERNAL_ERROR)
			internal.ReplyJSON(m, c)
			return
		}

		m := nanofytypes.Response{
			Exist:    err != sql.ErrNoRows,
			PubKey:   PubKey,
			FlagHash: FlagHash,
			SigHash:  SigHash,
		}

		internal.ReplyJSON(&m, c)
	}

}

func isActionAllowed(action string) bool {
	action = strings.ToLower(action)
	for _, ac := range whitelist {
		if action == ac {
			return true
		}
	}
	return false
}
