package handlers

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	apiclient "github.com/Satr10/go-whatsapp-bot/internal/client"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

// HandleMessage processes incoming message events
func HandleMessage(client *whatsmeow.Client, msg *events.Message) {
	text := msg.Message.GetConversation()
	lowerText := strings.ToLower(text)
	if text == "" {
		return // Ignore empty messages
	}

	// Simple command handling example
	if lowerText == ".ping" {
		err := SendTextMessage(client, msg.Info.ID, msg.Info.Sender, msg.Message, msg.Info.Chat, "PONG", true)
		if err != nil {
			// Log error properly
			fmt.Printf("Error sending pong: %v\n", err)
		}
	}

	// test kirim gambar
	if lowerText == ".gambar" {
		err := SendImage(client, msg.Info.Chat, "testing_folder/e64b254f3c010d51fd20958f365561a4_t.jpg", "TestGambar", true)
		if err != nil {
			// Log error properly
			fmt.Printf("Error sending image: %v\n", err)
		}
	}

	// test reply
	if lowerText == ".reply" {
		err := SendTextMessage(client, msg.Info.ID, msg.Info.Sender, msg.Message, msg.Info.Chat, "halo", true)
		if err != nil {
			fmt.Println("Error Mengirim Specific Reply: ", err)
		}
	}

	if lowerText == ".mp3" {
		err := SendVoiceMessage(client, msg.Info.Chat, "testing_folder/youtube_gBnalcSi138_audio.ogg", false)
		if err != nil {
			fmt.Println("Error kirim mp3: ", err)
		}
	}

	if lowerText == ".quote" {
		quote, err := apiclient.GetQuote()
		if err != nil {
			fmt.Println("Error Mendapatkan GetQuote: ", err)
		}

		message := fmt.Sprintf("✨ *%v*\n_- %v_", quote.Quote, quote.Author)
		err = SendTextMessage(client, msg.Info.ID, msg.Info.Sender, msg.Message, msg.Info.Chat, message, true)
		if err != nil {
			fmt.Println("Error Mengirim Specific Reply: ", err)
		}
	}

	if strings.Contains(lowerText, "fufufafa") {
		quote, err := apiclient.RandomFufufafa()
		if err != nil {
			fmt.Println("Error Mendapatkan GetQuote: ", err)
		}

		message := quote.Content

		err = SendTextMessage(client, msg.Info.ID, msg.Info.Sender, msg.Message, msg.Info.Chat, message, true)
		if err != nil {
			fmt.Println("Error Mengirim Specific Reply: ", err)
		}
	}

}

// SendImage is a helper function to send images to a chat.
// It takes a client, recipient JID, image path, caption, and deleteAfter flag.
// It reads the image at the given path and uploads it to WhatsApp servers.
// It then sends a message with the image to the recipient chat.
// If deleteAfter is true, the image is deleted from the local file system after sending.
// Returns an error if there is a problem reading the image, uploading the image, or sending the message.
func SendImage(client *whatsmeow.Client, recipient types.JID, imagePath string, caption string, deleteAfter bool) error {
	client.SendChatPresence(recipient, types.ChatPresenceComposing, types.ChatPresenceMediaText)
	gambarBytes, err := os.ReadFile(imagePath)
	if err != nil {
		fmt.Println("Error reading image:", err)
		return err
	}
	resp, err := client.Upload(context.Background(), gambarBytes, whatsmeow.MediaImage)
	if err != nil {
		fmt.Println("Error uploading image:", err)
		return err
	}

	imageMsg := &waE2E.ImageMessage{
		Caption:  proto.String(caption),
		Mimetype: proto.String("image/png"), // replace this with the actual mime type
		// you can also optionally add other fields like ContextInfo and JpegThumbnail here

		URL:           &resp.URL,
		DirectPath:    &resp.DirectPath,
		MediaKey:      resp.MediaKey,
		FileEncSHA256: resp.FileEncSHA256,
		FileSHA256:    resp.FileSHA256,
		FileLength:    &resp.FileLength,
	}
	_, err = client.SendMessage(context.Background(), recipient, &waE2E.Message{
		ImageMessage: imageMsg,
	})

	// hapus file setelah kirim, berguna untuk api dll
	if deleteAfter {
		err := os.Remove(imagePath)
		if err != nil {
			fmt.Println("Error Menghapus Gambar: ", err)
			return err
		}
	}

	return err
}

// SendTextMessage sends a text message to a recipient chat.
// If isReply is true, a reply message is sent with the quoted message.
// Otherwise, a regular conversation message is sent.
// Returns an error if there is a problem sending the message.
func SendTextMessage(client *whatsmeow.Client, StanzaID types.MessageID, sender types.JID, quotedMessage *waE2E.Message, recipient types.JID, textToSend string, isReply bool) error {

	client.SendChatPresence(recipient, types.ChatPresenceComposing, types.ChatPresenceMediaText)
	defer client.SendChatPresence(recipient, types.ChatPresencePaused, types.ChatPresenceMediaText)

	msgToSend := &waE2E.Message{}
	if isReply {
		msgToSend.ExtendedTextMessage = &waE2E.ExtendedTextMessage{
			Text: proto.String(textToSend),
			ContextInfo: &waE2E.ContextInfo{
				StanzaID:      proto.String(StanzaID),
				Participant:   proto.String(sender.String()),
				QuotedMessage: quotedMessage,
			},
		}
	} else {
		msgToSend.Conversation = proto.String(textToSend)
	}
	_, err := client.SendMessage(context.Background(), recipient, msgToSend)
	return err
}

// SendVoiceMessage is a helper function to send voice messages to a chat.
// It takes a client, recipient JID, audio path, and deleteAfter flag.
// It reads the audio at the given path, uploads it to WhatsApp servers,
// and sends a message with the audio to the recipient chat.
// If deleteAfter is true, the audio is deleted from the local file system after sending.
// Returns an error if there is a problem reading the audio, uploading the audio, or sending the message.
func SendVoiceMessage(client *whatsmeow.Client, recipient types.JID, audioPath string, deleteAfter bool) error {
	client.SendChatPresence(recipient, types.ChatPresenceComposing, types.ChatPresenceMediaAudio)
	audioBytes, err := os.ReadFile(audioPath)
	if err != nil {

		return err
	}
	fileForDuration, err := os.Open(audioPath)
	if err != nil {

		return err
	}
	defer fileForDuration.Close()

	//get audio duration
	duration, err := getOggDurationMs(fileForDuration)
	if err != nil {

		return err
	}

	// turn duration to second
	seconds := duration / 1000

	// mType := mimetype.Detect(audioBytes)
	// fmt.Println(mType.String())

	resp, err := client.Upload(context.Background(), audioBytes, whatsmeow.MediaAudio)
	if err != nil {
		fmt.Println("Error uploading audio:", err)
		return err
	}
	audioMessage := &waE2E.Message{
		AudioMessage: &waE2E.AudioMessage{
			URL:           proto.String(resp.URL),
			DirectPath:    proto.String(resp.DirectPath),
			MediaKey:      resp.MediaKey,
			Mimetype:      proto.String("audio/ogg; codecs=opus"),
			FileEncSHA256: resp.FileEncSHA256,
			FileSHA256:    resp.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(audioBytes))),
			Seconds:       proto.Uint32(uint32(seconds)),
			PTT:           proto.Bool(true),
			ContextInfo:   &waE2E.ContextInfo{},
		},
	}
	_, err = client.SendMessage(context.Background(), recipient, audioMessage)
	if err != nil {
		fmt.Println("Error Mengirim MP3: ", err)
	}
	return err
}

// HELPER FUNCTION

func getOggDurationMs(reader io.Reader) (int64, error) {
	// For simplicity, we read the entire Ogg file into a byte slice
	data, err := io.ReadAll(reader)
	if err != nil {
		return 0, fmt.Errorf("error reading Ogg file: %w", err)
	}

	// Search for the "OggS" signature and calculate the length
	var length int64
	for i := len(data) - 14; i >= 0 && length == 0; i-- {
		if data[i] == 'O' && data[i+1] == 'g' && data[i+2] == 'g' && data[i+3] == 'S' {
			length = int64(readLittleEndianInt(data[i+6 : i+14]))
		}
	}

	// Search for the "vorbis" signature and calculate the rate
	var rate int64
	for i := 0; i < len(data)-14 && rate == 0; i++ {
		if data[i] == 'v' && data[i+1] == 'o' && data[i+2] == 'r' && data[i+3] == 'b' && data[i+4] == 'i' && data[i+5] == 's' {
			rate = int64(readLittleEndianInt(data[i+11 : i+15]))
		}
	}

	if length == 0 || rate == 0 {
		return 0, fmt.Errorf("could not find necessary information in Ogg file")
	}

	durationMs := length * 1000 / rate
	return durationMs, nil
}

func readLittleEndianInt(data []byte) int64 {
	return int64(uint32(data[0]) | uint32(data[1])<<8 | uint32(data[2])<<16 | uint32(data[3])<<24)
}
