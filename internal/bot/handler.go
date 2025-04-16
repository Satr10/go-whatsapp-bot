package bot

import (
	"go.mau.fi/whatsmeow/types/events"
)

// eventHandler is the main event handling function for the bot. It listens to
// events and takes actions accordingly. The events it listens to are:
//
//   - events.Connected: Upon connection, the bot logs a message and could
//     perform additional actions if needed.
//
//   - events.Disconnected: Upon disconnection, the bot logs a message and could
//     perform additional actions if needed.
//
//   - events.Message: When a message is received, it is passed to the command
//     handler for processing.
//
//   - events.Receipt: The bot could handle receipts if it needs to track message
//     status.
//
// - events.Presence: The bot could handle presence updates if needed.
//
//   - events.LoggedOut: Upon logout, the bot logs a message and could perform
//     additional actions if needed.
//
// Add more event types if needed.
func (b *Bot) eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Connected:
		b.logger.Infof("Handler: Connected")
		// Perform actions upon connection if needed

	case *events.Disconnected:
		b.logger.Infof("Handler: Disconnected")
		// Handle disconnection if needed

	case *events.Message:
		// Pass the event to the command handler
		b.cmdHandler.HandleEvent(v)
		// b.logger.Infof("Message: ", v.Message.Conversation)

	case *events.Receipt:
		// Optional: Handle receipts if your bot needs to track message status
		// b.logger.Debugf("Received receipt for %v from %s", v.MessageIDs, v.SourceString())

	case *events.Presence:
		// Optional: Handle presence updates if needed
		// b.logger.Debugf("Presence update from %s: Available=%t", v.From, !v.Unavailable)

	case *events.LoggedOut:
		b.logger.Warnf("Logged out, reason: %s", v.Reason)
		// Handle logout, maybe attempt re-login or shutdown

	// Add more event types if needed
	default:
		// b.logger.Debugf("Ignored event type: %T", evt)
	}
}
