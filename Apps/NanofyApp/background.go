package NanofyApp

import (
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"github.com/brokenbydefault/Server/SQL"
	"database/sql"
	"time"
	"log"
)

// RecoverOldNV0 will request the node all pending transactions from the NV0 address, it'll retrieve all blocks that was
// a signature. It's need because the callback can crash or don't be notified, or even if you are starting your own
// Nanofy server today.
func Recover() {

	stmt, err := SQL.Connection.Prepare(`SELECT "1" FROM history WHERE FlagBlock = ?`)
	if err != nil {
		log.Print(err)
		time.Sleep(1 * time.Second)
		Recover()
	}

	for range time.Tick(30 * time.Second) {
		for _, version := range supportedVersions {
			blocks, err := RPCClient.GetAccountPending(Connectivity.HTTP, 1<<31, version.Amount(), version.Address())
			if err != nil {
				log.Print(err)
				continue
			}

			for _, pend := range blocks {
				var i string

				err := stmt.QueryRow([]byte(pend.Hash)).Scan(&i)
				if err == sql.ErrNoRows {

					blk, err := RPCClient.GetBlockByHash(Connectivity.HTTP, pend.Hash)
					if err != nil {
						log.Print(err)
						continue
					}

					pk, _ := pend.Source.GetPublicKey()
					go Store(pk, pend.Hash, blk)

				} else if err != nil {
					log.Print(err)
				}

			}
		}
	}

}
