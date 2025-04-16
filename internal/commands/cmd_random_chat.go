package commands

import (
	"context"
	"math/rand"
	"strings"

	"github.com/Satr10/go-whatsapp-bot/internal/models"
	"github.com/google/uuid"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

var WaitingQueue = make(map[types.JID]struct{})
var Rand rand.Rand
var ActiveChat = make(map[string]*models.AnonymousChat)

func StartMatching(ctx context.Context, client *whatsmeow.Client, evt *events.Message, args []string) error {
	if _, _, exist := IsUserInChat(evt.Info.Chat); exist {
		_, err := SendTextMessage(ctx, client, evt.Info.Chat, "Anda Sudah Mendapatkan Pasangan Chat")
		return err
	}

	roomID := uuid.NewString()
	if len(WaitingQueue) > 0 {
		// Get random partner from map
		var keys []types.JID
		for k := range WaitingQueue {
			keys = append(keys, k)
		}
		randomIndex := Rand.Intn(len(keys))
		partnerID := keys[randomIndex]

		ActiveChat[roomID] = &models.AnonymousChat{
			Peer1: evt.Info.Chat,
			Peer2: partnerID,

			Active: true,
			RoomID: roomID,
		}
		delete(WaitingQueue, evt.Info.Chat)
		delete(WaitingQueue, partnerID)
		_, err := SendTextMessage(ctx, client, evt.Info.Chat, "Kamu Mendapat Pasangan")
		if err != nil {
			return err
		}

		_, err = SendTextMessage(ctx, client, partnerID, "Kamu Mendapat Pasangan")
		if err != nil {
			return err
		}

	} else {
		WaitingQueue[evt.Info.Chat] = struct{}{}
		_, err := SendTextMessage(ctx, client, evt.Info.Chat, "Menunggu Pasangan")
		return err
	}
	return nil
}

func Leave(ctx context.Context, client *whatsmeow.Client, evt *events.Message, args []string) error {
	if _, exists := WaitingQueue[evt.Info.Chat]; exists {
		delete(WaitingQueue, evt.Info.Chat)
		_, err := SendTextMessage(ctx, client, evt.Info.Chat, "Keluar Dari Antrian Anonymous Chat")
		return err
	}
	if roomID, _, exist := IsUserInChat(evt.Info.Chat); exist {
		_, err := SendTextMessage(ctx, client, ActiveChat[roomID].Peer1, "Removed From Room")
		if err != nil {
			return err
		}
		_, err = SendTextMessage(ctx, client, ActiveChat[roomID].Peer2, "Removed From Room")
		if err != nil {
			return err
		}
		delete(ActiveChat, roomID)

	} else {
		_, err := SendTextMessage(ctx, client, evt.Info.Chat, "Kamu Tidak Masuk Room Dan Tidak Di Antrian")
		return err
	}
	return nil
}

func ForwardMessage(ctx context.Context, client *whatsmeow.Client, evt *events.Message, args []string) error {
	if _, partner, exist := IsUserInChat(evt.Info.Chat); exist {
		_, err := SendTextMessage(ctx, client, partner, strings.Join(args, " "))
		return err
	} else {
		_, err := SendTextMessage(ctx, client, evt.Info.Chat, "Not In Room")
		return err
	}
}

func IsUserInChat(jid types.JID) (string, types.JID, bool) {
	for _, chat := range ActiveChat {
		if chat.Peer1 == jid {
			return chat.RoomID, chat.Peer2, true
		}
		if chat.Peer2 == jid {
			return chat.RoomID, chat.Peer1, true
		}
	}
	return "", types.JID{}, false
}
