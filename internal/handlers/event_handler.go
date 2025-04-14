package handlers

import (
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

// EventHandler holds dependencies needed by handlers (like the client to send replies)
type EventHandler struct {
	client *whatsmeow.Client
	log    waLog.Logger // Or your preferred logger
	// Add other dependencies like database access if needed
}

func NewEventHandler(client *whatsmeow.Client) *EventHandler {
	return &EventHandler{
		client: client,
		log:    waLog.Stdout("EventHandler", "DEBUG", true), // Configure properly
	}
}

// Dispatch routes incoming events to specific handler functions
func (eh *EventHandler) Dispatch(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		eh.log.Infof("Received message from %s: %s", v.Info.Sender, v.Message.GetConversation())
		HandleMessage(eh.client, v) // Call the specific message handler
	case *events.Receipt:
		// Handle receipts if needed
		eh.log.Debugf("Received receipt: %s", v.MessageIDs)
	case *events.Presence:
		// Handle presence updates if needed
		eh.log.Infof("Presence update: %s is %t", v.From, v.Unavailable)
	default:
		// eh.log.Debugf("Received unknown event type: %T", evt)
	}
}
