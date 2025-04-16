package commands

import (
	"context"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)

// PingCommand handles the !ping command.
func PingCommand(ctx context.Context, client *whatsmeow.Client, evt *events.Message, args []string) error {
	_, err := SendTextMessage(ctx, client, evt.Info.Chat, "pong!")
	return err
}
