package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Satr10/go-whatsapp-bot/internal/bot"
	"github.com/Satr10/go-whatsapp-bot/internal/commands"
	"github.com/Satr10/go-whatsapp-bot/internal/config"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}

func main() {

	// initialize random for stuff
	commands.Rand = *rand.New(rand.NewSource(time.Now().UnixMicro()))

	logger := waLog.Stdout("Main", "INFO", true)
	dbPath := config.Config("DB_CONNECTION_STRING")
	dbDriver := config.Config("DB_TYPE")
	botInstance, err := bot.New(dbDriver, dbPath, logger)
	if err != nil {
		logger.Errorf("Failed to create bot: %v", err)
		os.Exit(1)
	}

	// Listen to Ctrl+C (you can also do something else that prevents the program from exiting)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	logger.Infof("Shutting down...")
	botInstance.Disconnect()
}
