package models

import "go.mau.fi/whatsmeow/types"

type AnonymousChat struct {
	// ID used if want to use with a dedicated db

	Peer1 types.JID
	Peer2 types.JID

	Active bool
	RoomID string
}

type AnonymousChatUserID struct {
	userID types.JID
}
