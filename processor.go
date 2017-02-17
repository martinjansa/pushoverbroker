package pushoverbroker

import (
	"errors"
)

// Processor handles the incomming messages, is responsible for the queing, persinstence and repeated attempts to deliver
type Processor struct {
	pushoverConnector PushoverConnector
}

// NewPushoverProcessor creates a new instance of the Processor
func NewPushoverProcessor(pushoverConnector PushoverConnector) *Processor {
	p := new(Processor)
	p.pushoverConnector = pushoverConnector

	return p
}

// HandleMessage receives a message to be processed (see IncommingPushNotificationMessageHandler interface)
func (p *Processor) HandleMessage(message PushNotification) error {
	return errors.New("Not implemented")
}

// Run starts the message processing loop
func (p *Processor) Run() error {
	return errors.New("Not implemented")
}
