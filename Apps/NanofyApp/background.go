package NanofyApp

import (
	"github.com/brokenbydefault/Nanofy"
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"github.com/brokenbydefault/Nanollet/Numbers"
	"github.com/brokenbydefault/Server/SQL"
	"database/sql"
	"time"
	"log"
)

var NV0Address = Nanofy.NV0.CreateFlagPublicKey().CreateAddress()
var amm, _ = Numbers.NewRawFromString("1")

func TrackNV0() {

	stmt, err := SQL.Connection.Prepare(`SELECT "1" FROM history WHERE FlagBlock = ?`)
	if err != nil {
		log.Print(err)
		time.Sleep(1 * time.Second)
		TrackNV0()
	}

	for range time.Tick(30 * time.Second) {
		blocks, err := RPCClient.GetAccountPending(Connectivity.HTTP, 1<<31, amm, NV0Address)
		if err != nil {
			log.Print(err)
			continue
		}

		for _, b := range blocks {

			var i string
			err := stmt.QueryRow([]byte(b.Hash)).Scan(&i)
			if err == sql.ErrNoRows {

				blk, err := RPCClient.GetBlockByHash(Connectivity.HTTP, b.Hash)
				if err != nil {
					log.Print(err)
					continue
				}

				go ReceiveNV0(b.Source, b.Hash, blk)
			}
			if err != nil {
				log.Print(err)
			}

		}
	}

}
