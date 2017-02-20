package main

import "testing"

// PushNotificationsSenderMock implements the PushNotificationsSender.PushNotificationsSender interface
type PushNotificationsSenderMock struct {
	response            error
	handleMessageCalled int
	notification        PushNotification
}

// NewPushNotificationsSenderMock initializes the mock
func NewPushNotificationsSenderMock() *PushNotificationsSenderMock {
	pcm := new(PushNotificationsSenderMock)
	pcm.response = nil
	pcm.handleMessageCalled = 0
	return pcm
}

// PostPushNotificationMessage receives the push notification message and returns the predefined response
func (pcm *PushNotificationsSenderMock) PostPushNotificationMessage(message PushNotification) error {
	pcm.handleMessageCalled++
	pcm.notification = message
	return pcm.response
}

// ForceResponse configures the response to be returned from the PostPushNotificationMessage() call
func (pcm *PushNotificationsSenderMock) ForceResponse(response error) {
	pcm.response = response
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
