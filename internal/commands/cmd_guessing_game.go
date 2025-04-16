package commands

import (
	"context"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type DataTebakKata struct {
	Kata string
	chat types.JID
}

var TebakKataAktifMap = make(map[types.JID]*DataTebakKata)

// dapatkan chat jid lalu simpan
func StartTebakKata(ctx context.Context, client *whatsmeow.Client, evt *events.Message, args []string) error {
	if _, exists := TebakKataAktifMap[evt.Info.Chat]; exists {
		_, err := SendTextMessage(ctx, client, evt.Info.Chat, "Anda sudah memulai game tebak kata")
		return err
	}

	// delete the data if it went unanswered after a set of time
	time.AfterFunc(30*time.Second, func() {
		if _, exists := TebakKataAktifMap[evt.Info.Chat]; exists {
			_, _ = SendTextMessage(ctx, client, evt.Info.Chat, "Deleted")
			delete(TebakKataAktifMap, evt.Info.Chat)
		}
	})

	TebakKataAktifMap[evt.Info.Chat] = &DataTebakKata{Kata: "Mouse", chat: evt.Info.Chat}
	_, err := SendTextMessage(ctx, client, evt.Info.Chat, "Jawaban: Mouse")
	return err
}

func TebakKata(ctx context.Context, client *whatsmeow.Client, evt *events.Message, args []string) error {
	if data, exists := TebakKataAktifMap[evt.Info.Chat]; exists {
		if args[0] == data.Kata {
			delete(TebakKataAktifMap, evt.Info.Chat)
			_, err := SendTextMessage(ctx, client, evt.Info.Chat, "BERHASIL")
			return err
		} else {
			_, err := SendTextMessage(ctx, client, evt.Info.Chat, "SALAH")
			return err
		}
	} else {
		_, err := SendTextMessage(ctx, client, evt.Info.Chat, "Tidak ada sesi tebak kata di chat ini")
		return err
	}

}
