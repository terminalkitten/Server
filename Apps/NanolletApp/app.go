package NanolletApp

import (
	"github.com/gorilla/websocket"
	"github.com/buger/jsonparser"
	"github.com/brokenbydefault/Server/Apps/internal"
	"strings"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
)

var whitelist = [...]string{"account_balance", "accounts_balances", "account_info", "account_history", "pending", "accounts_pending", "process", "block"}

func StartMessaging(m []byte, c *websocket.Conn) {
	m = jsonparser.Delete(m, "app")

	action, err := jsonparser.GetString(m, "action")
	if err != nil || !isActionAllowed(action) {
		internal.ReplyMessage(internal.NOT_ALLOWED, c)
		return
	}

	count, err := jsonparser.GetInt(m, "count")
	if !isCountAllowed(count) {
		internal.ReplyMessage(internal.NOT_ALLOWED, c)
		return
	}


	rc, err := Connectivity.HTTP.SendRequestReader(m)
	if err != nil {
		internal.ReplyMessage(internal.INTERNAL_ERROR, c)
		return
	}

	internal.ReplyMessaging(rc, c)
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

func isCountAllowed(count int64) bool {
	return count <= 30
}