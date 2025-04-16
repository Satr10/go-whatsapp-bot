package bot

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Satr10/go-whatsapp-bot/internal/commands"
	"github.com/Satr10/go-whatsapp-bot/internal/config"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Bot struct {
	client      *whatsmeow.Client
	logger      waLog.Logger
	cmdHandler  *commands.Handler
	upTimeSince time.Time
}

func New(driver string, dbPath string, logger waLog.Logger) (*Bot, error) {

	dbLog := waLog.Stdout("Database", "WARN", true)
	container, err := sqlstore.New("pgx", config.Config("DB_CONNECTION_STRING"), dbLog)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get device from database: %w", err)
	}

	clientLog := waLog.Stdout("Client", "WARN", true)
	client := whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID == nil {
		// No ID stored, new login
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				// Render the QR code here
				// e.g. qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
				// or just manually `echo 2@... | qrencode -t ansiutf8` in a terminal
				fmt.Println("QR code:", evt.Code)
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		// Already logged in, just connect
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}
	botInstance := &Bot{
		client:      client,
		logger:      logger,
		cmdHandler:  commands.NewHandler(client, logger),
		upTimeSince: time.Now(),
	}

	client.SendPresence(types.PresenceAvailable)
	client.AddEventHandler(botInstance.eventHandler)

	return botInstance, err
}

func (b *Bot) Client() *whatsmeow.Client {
	return b.client
}

func (b *Bot) Disconnect() {
	b.logger.Infof("Disconnecting...")
	b.client.Disconnect()

}
