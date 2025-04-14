package bot

import (
	"context"
	"fmt"

	"github.com/Satr10/go-whatsapp-bot/internal/auth"
	"github.com/Satr10/go-whatsapp-bot/internal/config"
	"github.com/Satr10/go-whatsapp-bot/internal/handlers"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type Bot struct {
	client  *whatsmeow.Client
	store   *sqlstore.Container
	log     waLog.Logger
	handler *handlers.EventHandler // Use a dedicated handler structure
}

// NewBot creates a new Bot instance
func NewBot(store *sqlstore.Container) (*Bot, error) {
	deviceStore, err := store.GetFirstDevice() // Or logic to select/create device
	if err != nil {
		return nil, fmt.Errorf("failed to get device store: %w", err)
	}

	// Configure logger (can be passed in or created here)
	clientLog := waLog.Stdout("Client", config.Config("LOG_LEVEL"), true) // Use log level from config

	client := whatsmeow.NewClient(deviceStore, clientLog)

	// Create the event handler instance (passing dependencies if needed)
	eventHandler := handlers.NewEventHandler(client) // Pass client for sending replies

	b := &Bot{
		client:  client,
		store:   store,
		log:     clientLog, // Can have a separate bot logger too
		handler: eventHandler,
	}

	// Register the main event dispatcher
	b.client.AddEventHandler(b.eventDispatcher)

	return b, nil
}

// eventDispatcher routes events to the appropriate handler
func (b *Bot) eventDispatcher(evt interface{}) {
	// Delegate to the handler structure
	b.handler.Dispatch(evt)
}

// Start connects the client and handles authentication
func (b *Bot) Start(ctx context.Context) error {
	var err error
	if b.client.Store.ID == nil {
		// No ID stored, new login required
		err = auth.PerformLogin(ctx, b.client) // Implement performLogin in auth.go
	} else {
		// Already logged in, just connect
		err = b.client.Connect()
	}

	if err != nil {
		return fmt.Errorf("failed to connect/login: %w", err)
	}

	b.log.Infof("WhatsApp client connected successfully")

	// send online
	err = b.client.SendPresence(types.PresenceAvailable)
	if err != nil {
		fmt.Println("Error Sending Online", err)
	}

	// Keep running until context is cancelled (handled by main)
	<-ctx.Done()
	return nil
}

// Stop disconnects the client
func (b *Bot) Stop() {
	b.log.Infof("Disconnecting WhatsApp client...")
	b.client.Disconnect()
}
