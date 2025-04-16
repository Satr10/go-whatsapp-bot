package commands

import (
	"context"
	"strings"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

// CommandFunc defines the signature for command handler functions.
type CommandFunc func(ctx context.Context, client *whatsmeow.Client, evt *events.Message, args []string) error

// Handler manages command registration and execution.
type Handler struct {
	client   *whatsmeow.Client
	registry map[string]CommandFunc
	logger   waLog.Logger
	prefix   string
}

// NewHandler creates a new command handler.
func NewHandler(client *whatsmeow.Client, logger waLog.Logger) *Handler {
	h := &Handler{
		client:   client,
		registry: make(map[string]CommandFunc),
		logger:   logger,
		prefix:   ".",
	}
	// Register commands here
	h.RegisterCommand("ping", PingCommand) // Example command
	// Add more commands: h.RegisterCommand("help", HelpCommand)
	return h
}

// RegisterCommand adds a new command to the registry.
func (h *Handler) RegisterCommand(name string, handlerFunc CommandFunc) {
	h.registry[strings.ToLower(name)] = handlerFunc
	h.logger.Infof("Registered command: %s%s", h.prefix, name)
}

// HandleEvent processes incoming message events to check for commands.
func (h *Handler) HandleEvent(evt *events.Message) {
	// ignore message from self
	if evt.Info.IsFromMe {
		return
	}

	// mendapatkan pesan dari extended message
	msgText := ""
	if evt.Message.GetConversation() != "" {
		msgText = evt.Message.GetConversation()
	} else if evt.Message.ExtendedTextMessage != nil && evt.Message.ExtendedTextMessage.Text != nil {
		msgText = evt.Message.ExtendedTextMessage.GetText()
	} else {
		// Add support for other message types if needed (e.g., captions)
		return
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

	// Look up command in registry
	cmdFunc, exists := h.registry[commandName]
	if !exists {
		// Optionally send "unknown command" message
		// h.logger.Infof("Unknown command received: %s", commandName)
		// _, _ = h.client.SendMessage(context.Background(), evt.Info.Chat, &waProto.Message{Conversation: proto.String("Unknown command.")})
		return
	}

	h.logger.Infof("Executing command '%s' from %s with args: %v", commandName, evt.Info.Sender, args)

	// Execute command in a goroutine to avoid blocking the event handler
	go func() {
		ctx := context.Background()
		err := cmdFunc(ctx, h.client, evt, args)
		if err != nil {
			h.logger.Errorf("Error executing command '%s': %v", commandName, err)
			// Optionally send error message back to user
			// SendTextMessage(ctx, evt.Info.Chat, fmt.Sprintf("Error: %v", err))
		}
	}()
}

// SendTextMessage is a utility function within the commands package
// Note: This could also live in internal/commands/utils.go
func SendTextMessage(ctx context.Context, client *whatsmeow.Client, recipient types.JID, text string) (whatsmeow.SendResponse, error) {
	msg := &waE2E.Message{Conversation: proto.String(text)}
	return client.SendMessage(ctx, recipient, msg)
}
