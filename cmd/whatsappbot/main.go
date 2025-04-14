package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Satr10/go-whatsapp-bot/internal/bot"
	"github.com/Satr10/go-whatsapp-bot/internal/database"
)

func main() {
	container, _, err := database.ConnectContainer()
	if err != nil {
		panic(err)
	}

	// 3. Create Bot Instance
	// Pass dependencies (config, store) to the bot constructor
	appBot, err := bot.NewBot(container)
	if err != nil {
		log.Fatalf("Failed to create bot instance: %v", err)
	}

	// 4. Start the Bot (Connects and starts listening)
	// Use context for potentially cancellable startup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := appBot.Start(ctx)
		if err != nil {
			log.Printf("Bot execution failed: %v", err)
			cancel() // Trigger shutdown if bot fails to start/run
		}
	}()

	// 5. Wait for termination signal
	fmt.Println("Bot is running. Press Ctrl+C to exit.")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		fmt.Println("Received termination signal, shutting down...")
	case <-ctx.Done(): // Handle context cancellation (e.g., from bot error)
		fmt.Println("Bot context cancelled, shutting down...")
	}

	// 6. Graceful Shutdown
	appBot.Stop() // Implement Stop in internal/bot
	fmt.Println("Bot stopped.")
}
