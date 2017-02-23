package main

import (
	"errors"
	"testing"
)

// TestShouldSendMessageToPushNotificationsSender tests whether the processor attempts to send all the incomming messages to Pushover connector
func TestShouldSendMessageToPushNotificationsSender(t *testing.T) {

	// **** GIVEN ****

	// The REST API server is initialized and connected to the message handler mock
	pcm := NewPushNotificationsSenderMock()
	processor := NewProcessor(pcm)

	// start the processor
	go processor.Run()

	// **** WHEN ****

	// a push notification is obtained by the process (via IncommingPushNotificationMessageHandler interface method HandleMessage())
	testMessage := PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: ""}
	err, responseCode := processor.HandleMessage(testMessage)

	// **** THEN ****

	// the request should respond correctly
	if err != nil {
		t.Errorf("Handling of the message failed with error %s, response code %d.", err, responseCode)
		return
	}

	// the right message shoud be delivered to the mock
	pcm.AssertMessageAcceptedOnce(t, testMessage)
}

// TestShouldPropagateSuccessOrPermanentFailureResponses tests whether the processor returns the right response codes from the external service calls
func TestShouldPropagateSuccessOrPermanentFailureResponses(t *testing.T) {

	var testcases = []struct {
		id                 string
		responseStatusCode int
	}{
		{"ShouldPropagateSuccess200", 200},
		{"ShouldPropagateError400", 400},
		{"ShouldPropagateError401", 401},
		{"ShouldPropagateError402", 402},
		{"ShouldPropagateError403", 403},
		{"ShouldPropagateError404", 404},
		{"ShouldPropagateError405", 405},
		{"ShouldPropagateError426", 426},
	}
	// **** GIVEN ****

	// The REST API server is initialized and connected to the message handler mock
	pcm := NewPushNotificationsSenderMock()
	processor := NewProcessor(pcm)

	testMessage := PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: ""}

	// start the processor
	go processor.Run()

	for _, tc := range testcases {

		t.Run(tc.id, func(t *testing.T) {

			// **** WHEN ****

			// set the response in the mock
			pcm.ForceResponse(nil, tc.responseStatusCode)

			// a push notification is obtained by the process (via IncommingPushNotificationMessageHandler interface method HandleMessage())
			err, responseCode := processor.HandleMessage(testMessage)

			// **** THEN ****

			// the request should respond correctly
			if err == nil {

				// check the reseponse code
				if responseCode != tc.responseStatusCode {
					t.Errorf("The handling of the message returned response code %d, but expected was %d", responseCode, tc.responseStatusCode)
				}

			} else {

				t.Errorf("Handling of the message failed with error %s, response code %d.", err, responseCode)
			}
		})
	}
}

// TestShouldReturnAcceptedOnTemporaryFailure tests whether the processor returns the HTTP error 202 (Accepted) if the call to the external service fails with a temporary failure
func TestShouldReturn202AcceptedOnTemporaryFailure(t *testing.T) {

	var testcases = []struct {
		id                 string
		responseError      error
		responseStatusCode int
	}{
		{"ShouldReturn202OnPostError", errors.New("post error"), 0},
		{"ShouldReturn202OnInternalServerError", nil, 500},
		{"ShouldReturn202OnGatewayTimeOut504", nil, 504},
		{"ShouldReturn202OnNetworkReadTimeOut598", nil, 598},
		{"ShouldReturn202OnNetworkTimeOut599", nil, 599},
	}
	// **** GIVEN ****

	// The REST API server is initialized and connected to the message handler mock
	pcm := NewPushNotificationsSenderMock()
	processor := NewProcessor(pcm)

	testMessage := PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: ""}

	// start the processor
	go processor.Run()

	for _, tc := range testcases {

		t.Run(tc.id, func(t *testing.T) {

			// **** WHEN ****

			// set the response in the mock
			pcm.ForceResponse(tc.responseError, tc.responseStatusCode)

			// a push notification is obtained by the process (via IncommingPushNotificationMessageHandler interface method HandleMessage())
			err, responseCode := processor.HandleMessage(testMessage)

			// **** THEN ****

			// the request should respond correctly
			if err == nil {

				// check the response code
				if responseCode != 202 {
					t.Errorf("The handling of the message returned response code %d, but expected was 202", responseCode)
				}

			} else {

				t.Errorf("Handling of the message failed with error %s, response code %d.", err, responseCode)
			}
		})
	}
}
