package main

import "log"
import "net/http"
import "fmt"

// Processor handles the incomming messages, is responsible for the queing, persinstence and repeated attempts to deliver
type Processor struct {
	PushNotificationsSender PushNotificationsSender
	LimitsCounter           LimitsCounter
}

// NewProcessor creates a new instance of the Processor
func NewProcessor(PushNotificationsSender PushNotificationsSender, LimitsCounter LimitsCounter) *Processor {
	p := new(Processor)
	p.PushNotificationsSender = PushNotificationsSender
	p.LimitsCounter = LimitsCounter
	return p
}

// HandleMessage receives a message to be processed (see IncommingPushNotificationMessageHandler interface)
func (p *Processor) HandleMessage(response *PushNotificationHandlingResponse, message PushNotification) error {

	// simple forward of the received message to the Pushover connector and return the result
	responseErr := p.PushNotificationsSender.PostPushNotificationMessage(response, message)

	acceptRequestToQueue := false

	// if the call succeeded
	if responseErr == nil {

		// if the response represents a temporary error and we should enqueue the message and try later
		switch {
		case response.responseCode >= 100 && response.responseCode < 300: // success codes
			// store the currnt limits into the cache
			p.LimitsCounter.SetLimits(message.GetToken(), response.limits)
			break

		case
			response.responseCode == 500, // Internal Server Error
			response.responseCode == 504, // Gateway Timeout
			response.responseCode == 598, // Network Read Timeout
			response.responseCode == 599: // Network Timeout
			acceptRequestToQueue = true
			break

		default: // all other failures
			// always generate a status=0 response
			response.jsonResponseBody = "{\"status\": 0 }"
			break

		}

	} else {

		// if the posting failed we assume the sender works fine (should be checked by the production tests), but connection cannot be made temporarily
		log.Printf("PushNotificationsSender.PostPushNotificationMessage failed with error %s.", responseErr.Error())

		acceptRequestToQueue = true
	}

	// if the message should be accepted to queue
	if acceptRequestToQueue {

		// decrement the limits for the current message
		err := p.LimitsCounter.DecrementLimits(message.GetToken())

		// if succeeded
		if err == nil {

			// TODO: queing of the message and trying later should be done here!

			// return HTTP error 202 (Accepted)
			responseErr = nil
			response.responseCode = http.StatusAccepted
			response.limits, _ = p.LimitsCounter.GetLimits(message.GetToken())
			response.jsonResponseBody = "{\"status\": 1 }"

		} else {
			// return the not permited reponse
			response.responseCode = http.StatusForbidden
			response.limits, _ = p.LimitsCounter.GetLimits(message.GetToken())
			response.jsonResponseBody = fmt.Sprintf("{\"status\": 0, \"errors\": [\"%s\"] }", err.Error())
		}
	}

	return responseErr
}

// Run starts the message processing loop
func (p *Processor) Run() {
	log.Print("Not Process::Run has not been implemented yet")
}
