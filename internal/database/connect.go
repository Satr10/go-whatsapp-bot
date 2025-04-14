package database

import (
	"github.com/Satr10/go-whatsapp-bot/internal/config"
	// for sqlite
	// _ "github.com/mattn/go-sqlite3"

	// for postgres
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func ConnectContainer() (*sqlstore.Container, waLog.Logger, error) {
	dbLog := waLog.Stdout("Database", config.Config("LOG_LEVEL"), true)

	container, err := sqlstore.New(config.Config("DB_TYPE"), config.Config("DB_CONNECTION_STRING"), dbLog)
	if err != nil {
		panic(err)
	}
	return container, dbLog, err
}
