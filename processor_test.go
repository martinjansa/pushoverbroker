package main

import (
	"encoding/json"
	"errors"
	"strings"
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

	var response = PushNotificationHandlingResponse{}
	err := processor.HandleMessage(&response, testMessage)

	// **** THEN ****

	// the request should respond correctly
	if err != nil {
		t.Errorf("Handling of the message failed with error %s, response code %d.", err, response.responseCode)
		return
	}

	// the right message shoud be delivered to the mock
	pcm.AssertMessageAcceptedOnce(t, testMessage)
}

// TestShouldPropagateSuccessOrPermanentFailureResponses tests whether the processor returns the right response codes from the external service calls
func TestShouldPropagateSuccessOrPermanentFailureResponses(t *testing.T) {

	var testcases = []struct {
		id                         string
		responseStatusCode         int
		responseLimits             *Limits
		responseBody               string
		expectedResponseBodyStatus int
	}{
		{"ShouldPropagateSuccess200", 200, &Limits{limit: 1000, remaining: 500, reset: 123456789}, "{\"status\": 1}", 1},
		{"ShouldPropagateError400", 400, nil, "{\"status\": 0}", 1},
		{"ShouldPropagateError401", 401, nil, "", 1},
		{"ShouldPropagateError402", 402, nil, "", 1},
		{"ShouldPropagateError403", 403, nil, "", 1},
		{"ShouldPropagateError404", 404, nil, "", 1},
		{"ShouldPropagateError405", 405, nil, "", 1},
		{"ShouldPropagateError426", 426, nil, "", 1},
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
			pcm.ForceResponse(nil, tc.responseStatusCode, tc.responseLimits, tc.responseBody)

			// a push notification is obtained by the process (via IncommingPushNotificationMessageHandler interface method HandleMessage())
			var response = PushNotificationHandlingResponse{}
			err := processor.HandleMessage(&response, testMessage)

			// **** THEN ****

			// the request should respond correctly
			if err == nil {

				// check the reseponse code
				if response.responseCode != tc.responseStatusCode {
					t.Errorf("The handling of the message returned response code %d, but expected was %d", response.responseCode, tc.responseStatusCode)
				}

			} else {

				t.Errorf("Handling of the message failed with error %s, response code %d.", err, response.responseCode)
			}

			// if the limits match the expected limits
			if response.limits != tc.responseLimits {

				t.Errorf("Returned limits %s don't match the expected value %s.", response.limits, tc.responseLimits)
			}

			// get the content of the body
			type ResponseJSONBodyContent struct {
				Status  int    `json:"status"`
				Request string `json:"request"`
			}
			var responseJSONBodyContent = ResponseJSONBodyContent{0, ""}
			err = json.NewDecoder(strings.NewReader(response.jsonResponseBody)).Decode(&responseJSONBodyContent)
			if err != nil {
				t.Errorf("POST request returned JSON \"%s\", which failed to decode with error %s.", response.jsonResponseBody, err.Error())
				return
			}

			// check the status value
			if responseJSONBodyContent.Status != tc.expectedResponseBodyStatus {
				t.Errorf("POST request returned JSON with status %d, expected status was %d.", responseJSONBodyContent.Status, tc.expectedResponseBodyStatus)
				return
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
		responseLimits     *Limits
	}{
		{"ShouldReturn202OnPostError", errors.New("post error"), 0, nil},
		{"ShouldReturn202OnInternalServerError", nil, 500, nil},
		{"ShouldReturn202OnGatewayTimeOut504", nil, 504, nil},
		{"ShouldReturn202OnNetworkReadTimeOut598", nil, 598, nil},
		{"ShouldReturn202OnNetworkTimeOut599", nil, 599, nil},
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
			pcm.ForceResponse(tc.responseError, tc.responseStatusCode, tc.responseLimits, "{\"status\": 1}")

			// a push notification is obtained by the process (via IncommingPushNotificationMessageHandler interface method HandleMessage())
			var response = PushNotificationHandlingResponse{}
			err := processor.HandleMessage(&response, testMessage)

			// **** THEN ****

			// the request should respond correctly
			if err == nil {

				// check the response code
				if response.responseCode != 202 {
					t.Errorf("The handling of the message returned response code %d, but expected was 202", response.responseCode)
				}

			} else {

				t.Errorf("Handling of the message failed with error %s, response code %d.", err, response.responseCode)
			}

			// if the limits match the expected limits
			if response.limits != tc.responseLimits {

				t.Errorf("Returned limits %s don't match the expected value %s.", response.limits, tc.responseLimits)
			}
		})
	}
}
