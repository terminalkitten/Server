package NanosubApp

import (
	"github.com/buger/jsonparser"
	"github.com/brokenbydefault/Server/Apps/internal"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Nanollet/RPC"
	"encoding/json"
	"bytes"
	"github.com/brokenbydefault/Server/handlers/handlerstype"
	"github.com/Inkeliz/blakEd25519"
)

var whitelist = []string{"subscribe", "unsubscribe"}

func StartMessaging(m []byte, c *handlerstype.Client) {
	action, err := jsonparser.GetString(m, "action")
	if err != nil || !internal.IsActionAllowed(action, whitelist) {
		internal.ReplyMessage(internal.NOT_ALLOWED, c)
		return
	}

	switch action {
	case "subscribe":
		j := RPCClient.SubscribeRequest{}
		json.Unmarshal(m, &j)

		Subscribe(j.PublicKey, c)
	case "unsubscribe":
		Unsubscribe(c)
	}

}

func Subscribe(pk Wallet.PublicKey, c *handlerstype.Client) {
	var pkb [32]byte
	copy(pkb[:], pk)

	// If the connection already subscribed to this public-key or invalid
	if bytes.Equal(c.SubClient.PK, pk) || len(pk) != blakEd25519.PublicKeySize {
		return
	}

	if c.SubClient.PK != nil {
		Unsubscribe(c)
	}

	if subs, ok := subscribers[pkb]; ok {
		subscribers[pkb] = append(subs, c)
	} else {
		subscribers[pkb] = []*handlerstype.Client{c}
	}

	c.SubClient.PK = pk

	m := RPCClient.Subscription{
		PublicKey: c.SubClient.PK,
	}
	internal.ReplyJSON(&m, c)
}

func Unsubscribe(c *handlerstype.Client) {
	var pkb [32]byte

	if c.SubClient.PK == nil {
		return
	}

	copy(pkb[:], c.SubClient.PK)

	if subs, ok := subscribers[pkb]; ok {
		if len(subs) >= 2 {
			cons := make([]*handlerstype.Client, len(subs)-1)

			var index int
			for _, con := range subs {
				if c != con {
					cons[index] = con
				}
				index++
			}

			subscribers[pkb] = cons
		} else {
			delete(subscribers, pkb)
		}
	}

	c.SubClient.PK = nil

	m := RPCClient.Subscription{
		PublicKey: c.SubClient.PK,
	}
	internal.ReplyJSON(&m, c)
}

var subscribers = make(map[[32]byte][]*handlerstype.Client)

func ReceiveCallback(origin Wallet.PublicKey, dest Wallet.PublicKey, amm *Numbers.RawAmount, hash Block.BlockHash, block []byte) {
	var sub []*handlerstype.Client
	var searchOrigin [32]byte
	var searchDest [32]byte

	copy(searchOrigin[:], origin)
	sub = append(sub, subscribers[searchOrigin]...)

	copy(searchDest[:], dest)
	sub = append(sub, subscribers[searchDest]...)

	// If no one is listening to these addresses: we will stop here.
	if len(sub) == 0 {
		return
	}

	m := RPCClient.CallbackResponse{
		Hash:        hash,
		Origin:      origin,
		Destination: dest,
		Amount:      amm,
		Block:       block,
		DefaultResponse: RPCClient.DefaultResponse{
			Error: "",
		},
	}

	j, err := json.Marshal(&m)
	if err != nil {
		return
	}

	for _, con := range sub {
		internal.ReplyMessage(j, con)
	}

}
