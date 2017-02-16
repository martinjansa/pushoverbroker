package pushoverbroker

import (
	"github.com/go-martini/martini"
)

// MessageHandler handles message accepted by the REST API
type MessageHandler interface {
	HandleMessage(message PushNotification)
}

// Server is the REST API server.
type Server *martini.ClassicMartini

// NewServer creates a new server. Accepts the messageHandler that will handle all the received messages
func NewServer(messageHandler MessageHandler) Server {
	m := Server(martini.Classic())

	return m
}
