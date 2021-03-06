package main

import "testing"

// PushNotificationsSenderMock implements the PushNotificationsSender.PushNotificationsSender interface
type PushNotificationsSenderMock struct {
	responseErr         error
	responseCode        int
	limits              *Limits
	responseBody        string
	handleMessageCalled int
	notification        PushNotification
}

// NewPushNotificationsSenderMock initializes the mock
func NewPushNotificationsSenderMock() *PushNotificationsSenderMock {
	pcm := new(PushNotificationsSenderMock)
	pcm.responseErr = nil
	pcm.responseCode = 200
	pcm.limits = nil
	pcm.handleMessageCalled = 0
	pcm.responseBody = ""
	return pcm
}

// ForceResponse configures the response to be returned from the PostPushNotificationMessage() call
func (pcm *PushNotificationsSenderMock) ForceResponse(responseErr error, reseponseCode int, limits *Limits, responseBody string) {
	pcm.handleMessageCalled = 0
	pcm.responseErr = responseErr
	pcm.responseCode = reseponseCode
	pcm.responseBody = responseBody
	pcm.limits = limits
}

// PostPushNotificationMessage receives the push notification message and returns the predefined error and response code
func (pcm *PushNotificationsSenderMock) PostPushNotificationMessage(response *PushNotificationHandlingResponse, message PushNotification) error {
	pcm.handleMessageCalled++
	pcm.notification = message
	response.responseCode = pcm.responseCode
	response.limits = pcm.limits
	response.jsonResponseBody = pcm.responseBody
	return pcm.responseErr
}

// AssertMessageAcceptedOnce checks that the message was accepted
func (pcm *PushNotificationsSenderMock) AssertMessageAcceptedOnce(t *testing.T, message PushNotification) {
	if pcm.handleMessageCalled != 1 {
		t.Errorf("1 message expected, %d received.", pcm.handleMessageCalled)
	}
	if pcm.notification != message {
		t.Error("The received push notification does not match the expected value.")
	}
}
