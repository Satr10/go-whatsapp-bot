package commands

import (
	"context"
	"fmt"
	"strings"

	"go.mau.fi/whatsmeow/types/events"
)

func (h *Handler) HandleImage(evt *events.Message) {
	msgText := ""
	if evt.Message.ImageMessage.Caption != nil {
		msgText = *evt.Message.ImageMessage.Caption
		h.logger.Infof(msgText)
	}
	// Trim whitespace and check for prefix
	trimmedText := strings.TrimSpace(msgText)
	if !strings.HasPrefix(trimmedText, h.prefix) {
		return // Not a command
	}

	// Split into command and arguments
	parts := strings.Fields(trimmedText)
	if len(parts) == 0 {
		return // Empty command
	}

	// mendapatkan command dan argumen
	commandName := strings.ToLower(strings.TrimPrefix(parts[0], h.prefix))
	args := parts[1:]

	h.logger.Infof("Executing command '%s' from %s with args: %v", commandName, evt.Info.Sender, args)

	// Look up and execute command
	cmdFunc, exists := h.registry[commandName]
	if !exists {
		return
	}

	go func() {
		ctx := context.Background()
		err := cmdFunc(ctx, h.client, evt, args)
		if err != nil {
			h.logger.Errorf("Error executing image command '%s': %v", commandName, err)
			SendTextMessage(ctx, h.client, evt.Info.Chat, fmt.Sprintf("Error: %v", err))
		}
	}()
}
