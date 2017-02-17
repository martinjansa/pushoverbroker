package pushoverbroker

import "log"

// Processor handles the incomming messages, is responsible for the queing, persinstence and repeated attempts to deliver
type Processor struct {
	PushoverConnector PushoverConnector
}

// NewProcessor creates a new instance of the Processor
func NewProcessor(pushoverConnector PushoverConnector) *Processor {
	p := new(Processor)
	p.PushoverConnector = pushoverConnector

	return p
}

// HandleMessage receives a message to be processed (see IncommingPushNotificationMessageHandler interface)
func (p *Processor) HandleMessage(message PushNotification) error {

	// simple forward of the received message to the Pushover connector and return the result
	return p.PushoverConnector.PostPushNotificationMessage(message)
}

// Run starts the message processing loop
func (p *Processor) Run() {
	log.Print("Not Process::Run has not been implemented yet")
}
