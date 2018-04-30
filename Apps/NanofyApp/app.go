package NanofyApp

import (
	"github.com/buger/jsonparser"
	"github.com/brokenbydefault/Server/Apps/internal"
	"encoding/json"
	"github.com/brokenbydefault/Server/SQL"
	"github.com/brokenbydefault/Nanofy/nanofytypes"
	"database/sql"
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/Block"
	"github.com/brokenbydefault/Nanofy"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Server/handlers/handlerstype"
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"log"
)

type NanofyVersion struct {
	PublicKey Wallet.PublicKey
	Address   Wallet.Address
}

var NV0 = Nanofy.NewNanofierVersion0()
var whitelist = [...]string{"file"}
var supportedVersions = []Nanofy.Nanofier{
	Nanofy.NewNanofierVersion0(),
	Nanofy.NewNanofierVersion1(),
}

func Start() {
	go Recover()
}

func StartMessaging(m []byte, c *handlerstype.Client) {
	action, err := jsonparser.GetString(m, "action")
	if err != nil || !internal.IsActionAllowed(action, whitelist[:]) {
		internal.ReplyMessage(internal.NOT_ALLOWED, c)
		return
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

func ReceiveCallback(origin Wallet.PublicKey, dest Wallet.PublicKey, amm *Numbers.RawAmount, hash Block.BlockHash, block []byte) {

	for _, version := range supportedVersions {
		if dest.CreateAddress() == version.Address() {
			if blk, err := Block.NewBlockFromJSON(block); err == nil {
				Store(origin, hash, blk)
			} else {
				log.Println(err)
			}
		}
	}
}

func Store(pk Wallet.PublicKey, flaghash []byte, flagblock Block.UniversalBlock) {
	// The FlagHash is explicit needed in the function to save us from Blake2 computation.

	sigblock, err := RPCClient.GetBlockByHash(Connectivity.HTTP, flagblock.Previous)
	if err != nil {
		return
	}

	prevblock, err := RPCClient.GetBlockByHash(Connectivity.HTTP, sigblock.Previous)
	if err != nil {
		return
	}

	nanofier, err := Nanofy.NewNanofierFromFlagBlock(&flagblock)
	if err != nil {
		// unsupported version
		return
	}

	if !nanofier.VerifyBlock(&pk, &flagblock, &sigblock, &prevblock) {
		// invalid block
		return
	}

	stmt, err := SQL.Connection.Prepare("INSERT history(`PubKey`, `FileKey`, `FlagBlock`, `SigBlock`) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		return
	}

	destination, _ := sigblock.GetTarget()
	dest, _ := destination.GetPublicKey()

	_, err = stmt.Exec([]byte(pk), []byte(dest), []byte(flaghash), []byte(flagblock.Previous))
	if err != nil {
		log.Print(err)
	}
}
