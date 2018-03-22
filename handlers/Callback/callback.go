package Callback

import (
	"net/http"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"github.com/brokenbydefault/Nanofy"
	"github.com/brokenbydefault/Nanollet/Block"
	"bytes"
	"github.com/brokenbydefault/Server/Apps/NanofyApp"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Util"
)

var NV0Address = *Nanofy.NV0.CreateFlagPublicKey()

func Start(w http.ResponseWriter, r *http.Request) {
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	blk, err := jsonparser.GetString(jsn, "block")
	if err != nil {
		return
	}

	flagblock, _ := Block.NewBlockFromJSON([]byte(blk))

	if bytes.Compare(flagblock.Destination, NV0Address) == 0 {
		acc, _ := jsonparser.GetString(jsn, "account")
		flaghash, _ := jsonparser.GetString(jsn, "hash")
		h, _ := Util.UnsafeHexDecode(flaghash)

		go NanofyApp.ReceiveNV0(Wallet.Address(acc), h, flagblock)
	}

}
