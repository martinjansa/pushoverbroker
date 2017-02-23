package main

import "testing"

// PushNotificationsSenderMock implements the PushNotificationsSender.PushNotificationsSender interface
type PushNotificationsSenderMock struct {
	responseErr         error
	responseCode        int
	limits              *Limits
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
	return pcm
}

// ForceResponse configures the response to be returned from the PostPushNotificationMessage() call
func (pcm *PushNotificationsSenderMock) ForceResponse(responseErr error, reseponseCode int, limits *Limits) {
	pcm.handleMessageCalled = 0
	pcm.responseErr = responseErr
	pcm.responseCode = reseponseCode
	pcm.limits = limits
}

// PostPushNotificationMessage receives the push notification message and returns the predefined error and response code
func (pcm *PushNotificationsSenderMock) PostPushNotificationMessage(message PushNotification) (error, int, *Limits) {
	pcm.handleMessageCalled++
	pcm.notification = message
	return pcm.responseErr, pcm.responseCode, pcm.limits
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
