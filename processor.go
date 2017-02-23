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
func (p *Processor) HandleMessage(message PushNotification) (error, int, *Limits) {

	// simple forward of the received message to the Pushover connector and return the result
	responseErr, reseponseCode, limits := p.PushNotificationsSender.PostPushNotificationMessage(message)

	acceptRequestToQueue := false

	// if the call succeeded
	if responseErr == nil {

		// if the response represents a temporary error and we should enqueue the message and try later
		switch reseponseCode {
		case
			500, // Internal Server Error
			504, // Gateway Timeout
			598, // Network Read Timeout
			599: // Network Timeout
			acceptRequestToQueue = true
		}

	} else {

		// if the posting failed we assume the sender works fine (should be checked by the production tests), but connection cannot be made temporarily
		log.Printf("PushNotificationsSender.PostPushNotificationMessage failed with error %s.", responseErr.Error())

		acceptRequestToQueue = true
	}

	// if the message should be accepted to queue
	if acceptRequestToQueue {

		// TODO: quing of the message and trying later should be done here!

		// return HTTP error 202 (Accepted)
		return nil, 202, nil
	}

	return responseErr, reseponseCode, limits
}

// Run starts the message processing loop
func (p *Processor) Run() {
	log.Print("Not Process::Run has not been implemented yet")
}
