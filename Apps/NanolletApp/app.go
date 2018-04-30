package NanolletApp

import (
	"github.com/buger/jsonparser"
	"github.com/brokenbydefault/Server/Apps/internal"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"github.com/brokenbydefault/Server/handlers/handlerstype"
)

var whitelist = [...]string{"account_balance", "accounts_balances", "account_info", "account_history", "pending", "accounts_pending", "process", "block"}

func StartMessaging(m []byte, c *handlerstype.Client) {
	m = jsonparser.Delete(m, "app")

	action, err := jsonparser.GetString(m, "action")
	if err != nil || !internal.IsActionAllowed(action, whitelist[:]) {
		internal.ReplyMessage(internal.NOT_ALLOWED, c)
		return
	}

	count, err := jsonparser.GetInt(m, "count")
	if err == nil && (count <= 0 || count > 1000) {
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
