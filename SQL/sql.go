package SQL

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"github.com/brokenbydefault/Server/Config"
	"log"
)

var Connection *sql.DB

func Start() error {
	db, err := sql.Open("mysql", Config.Config["DB_USER"]+":"+Config.Config["DB_PASS"]+"@/nanofy?charset=utf8mb4")
	if err != nil {
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}

	db.SetConnMaxLifetime(time.Second * 3)

	Connection = db
	return nil
}
