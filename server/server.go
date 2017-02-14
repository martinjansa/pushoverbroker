package server

import "github.com/go-martini/martini"

// Message represents a message with json request that is passed to the REST API
type Message struct {
	Request string
}

// GetRequest  returns the request from the message
func (m *Message) GetRequest() string {
	return m.Request
}

// MessageHandler handles message accepted by the REST API
type MessageHandler interface {
	HandleMessage(message Message)
}

// Server is the REST API server.
type Server *martini.ClassicMartini

// NewServer creates a new server. Accepts the messageHandler that will handle all the received messages
func NewServer(messageHandler MessageHandler) Server {
	m := Server(martini.Classic())

	return m
}
