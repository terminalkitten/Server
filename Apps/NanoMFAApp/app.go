package NanoMFAApp

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Server/handlers/handlerstype"
	"bytes"
	"github.com/brokenbydefault/Server/Apps/internal"
	"github.com/buger/jsonparser"
	"encoding/json"
	"github.com/brokenbydefault/MFA/mfatypes"
	"github.com/Inkeliz/blakEd25519"
)

var whitelist = []string{"subscribe", "unsubscribe", "send"}

func StartMessaging(m []byte, c *handlerstype.Client) {
	action, err := jsonparser.GetString(m, "action")
	if err != nil || !internal.IsActionAllowed(action, whitelist) {
		internal.ReplyMessage(internal.NOT_ALLOWED, c)
		return
	}

	switch action {
	case "subscribe":
		j := mfatypes.SubscribeRequest{}
		json.Unmarshal(m, &j)

		Subscribe(j.PublicKey, c)
	case "send":
		j := mfatypes.EnvelopeRequest{}
		json.Unmarshal(m, &j)

		Send(j.PublicKey, j.Envelope, c)
	case "unsubscribe":
		Unsubscribe(c)
	}

}

func Subscribe(pk Wallet.PublicKey, c *handlerstype.Client) {
	var pkb [32]byte
	copy(pkb[:], pk)

	// If the connection already subscribed to this public-key or invalid
	if len(pk) != blakEd25519.PublicKeySize  || bytes.Equal(c.MFAClient.PK, pk) {
		return
	}

	if c.MFAClient.PK != nil {
		Unsubscribe(c)
	}

	if _, ok := subscribers[pkb]; !ok {
		subscribers[pkb] = c
		c.MFAClient.PK = pk
	}

	m := mfatypes.Subscription{
		PublicKey: c.MFAClient.PK,
	}
	internal.ReplyJSON(&m, c)
}

func Unsubscribe(c *handlerstype.Client) {
	var pkb [32]byte

	if c.MFAClient.PK != nil {
		copy(pkb[:], c.MFAClient.PK)

		delete(subscribers, pkb)
		c.MFAClient.PK = nil
	}

	m := mfatypes.Subscription{
		PublicKey: c.MFAClient.PK,
	}
	internal.ReplyJSON(&m, c)
}

var subscribers = make(map[[32]byte]*handlerstype.Client)

func Send(pk Wallet.PublicKey, capsule []byte, c *handlerstype.Client) {
	var pkb [32]byte
	copy(pkb[:], pk)

	receiver, ok := subscribers[pkb]
	if !ok {
		internal.ReplyMessage(internal.RECEIVER_NOT_FOUND, c)
		return
	}

	j := mfatypes.CallbackResponse{
		Envelope: capsule,
	}
	err := internal.ReplyJSON(&j, receiver)
	if err != nil {
		internal.ReplyMessage(internal.RECEIVER_NOT_FOUND, c)
		return
	}

	s := mfatypes.Subscription{
		PublicKey: pk,
	}
	internal.ReplyJSON(&s, c)
}
