package NanofyApp

import (
	"github.com/brokenbydefault/Nanollet/Wallet"
	"github.com/brokenbydefault/Nanollet/RPC"
	"github.com/brokenbydefault/Nanollet/RPC/Connectivity"
	"github.com/brokenbydefault/Nanofy"
	"github.com/brokenbydefault/Server/SQL"
	"github.com/brokenbydefault/Nanollet/Block"
	"log"
)

func ReceiveNV0(acc Wallet.Address, flaghash []byte, flagblock Block.UniversalBlock) {
	// The FlagHash is explicit needed in the function to save us from Blake2 computation.

	addr, err := acc.GetPublicKey()
	if err != nil {
		return
	}

	sigblock, err := RPCClient.GetBlockByHash(Connectivity.HTTP, flagblock.Previous)
	if err != nil {
		log.Print(err)
		return
	}

	if !Nanofy.VerifyBlock(&addr, flagblock, sigblock) {
		log.Print(err)
		return
	}

	stmt, err := SQL.Connection.Prepare("INSERT history(`PubKey`, `FileKey`, `FlagBlock`, `SigBlock`) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Print(err)
		return
	}

	_, err = stmt.Exec([]byte(addr), []byte(sigblock.Destination), []byte(flaghash), []byte(flagblock.Previous))
	if err != nil {
		log.Print(err)
	}
}
