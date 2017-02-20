package main

import "log"

// Processor handles the incomming messages, is responsible for the queing, persinstence and repeated attempts to deliver
type Processor struct {
	PushNotificationsSender PushNotificationsSender
}

// NewProcessor creates a new instance of the Processor
func NewProcessor(PushNotificationsSender PushNotificationsSender) *Processor {
	p := new(Processor)
	p.PushNotificationsSender = PushNotificationsSender

	return p
}

// HandleMessage receives a message to be processed (see IncommingPushNotificationMessageHandler interface)
func (p *Processor) HandleMessage(message PushNotification) (error, int) {

	// simple forward of the received message to the Pushover connector and return the result
	responseErr, reseponseCode := p.PushNotificationsSender.PostPushNotificationMessage(message)

	return responseErr, reseponseCode
}

// Run starts the message processing loop
func (p *Processor) Run() {
	log.Print("Not Process::Run has not been implemented yet")
}
