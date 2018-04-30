package Callback

import (
	"net/http"
	"github.com/buger/jsonparser"
	"io/ioutil"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Server/Apps/NanofyApp"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Server/Apps/NanosubApp"
	"github.com/brokenbydefault/Nanollet/Util"
	"github.com/brokenbydefault/Nanollet/Numbers"
)

type Receiver interface {
	ReceiveCallback(origin Wallet.PublicKey, dest Wallet.PublicKey, hash Block.BlockHash, block []byte)
}

func Start(_ http.ResponseWriter, r *http.Request) {
	jsn, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	var (
		addr        Wallet.PublicKey
		hash        = make(Block.BlockHash, 32)
		destination Wallet.PublicKey
		amount      *Numbers.RawAmount
		block       []byte
		isSend      bool
	)

	paths := [][]string{
		{"account"},
		{"hash"},
		{"block"},
		{"amount"},
		{"is_send"},
	}
	jsonparser.EachKey(jsn, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
		switch idx {
		case 0:
			addr, _ = Wallet.Address(value).GetPublicKey()
		case 1:
			hash, _ = Util.UnsafeHexDecode(string(value))
		case 2:
			// We need to get string then get back to json because the callback is really bad, the json is a string!
			blk, _ := jsonparser.ParseString(value)
			block = []byte(blk)
		case 3:
			amount, _ = Numbers.NewRawFromString(string(value))
		case 4:
			isSend, _ = jsonparser.ParseBoolean(value)
		}
	}, paths...)

	if isSend {
		dest, _ := jsonparser.GetString(block, "link")
		destination, _ = Util.UnsafeHexDecode(dest)
	} else {
		dest, _ := jsonparser.GetString(block, "destination")
		destination, _ = Wallet.Address(dest).GetPublicKey()
	}

	go NanofyApp.ReceiveCallback(addr, destination, amount, hash, block)
	go NanosubApp.ReceiveCallback(addr, destination, amount, hash, block)
}
