package auth

import (
	"context"
	"os"

	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
)

func PerformLogin(ctx context.Context, c *whatsmeow.Client) (err error) {
	qrChan, err := c.GetQRChannel(ctx)
	if err != nil {
		return err
	}
	err = c.Connect()
	if err != nil {
		return err
	}
	for evt := range qrChan {
		if evt.Event == "code" {
			qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
		}
	}
	return err

}
