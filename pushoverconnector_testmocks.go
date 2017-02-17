package pushoverbroker

import "testing"

// PushoverConnectorMock implements the pushoverconnector.PushoverConnector interface
type PushoverConnectorMock struct {
	response            error
	handleMessageCalled int
	notification        PushNotification
}

// NewPushoverConnectorMock initializes the mock
func NewPushoverConnectorMock() *PushoverConnectorMock {
	pcm := new(PushoverConnectorMock)
	pcm.response = nil
	pcm.handleMessageCalled = 0
	return pcm
}

// PostPushNotificationMessage receives the push notification message and returns the predefined response
func (pcm *PushoverConnectorMock) PostPushNotificationMessage(message PushNotification) error {
	pcm.handleMessageCalled++
	pcm.notification = message
	return pcm.response
}

// AssertMessageAcceptedOnce checks that the message was accepted
func (pcm *PushoverConnectorMock) AssertMessageAcceptedOnce(t *testing.T, message PushNotification) {
	if pcm.handleMessageCalled != 1 {
		t.Errorf("1 message expected, %d received.", pcm.handleMessageCalled)
	}
	if pcm.notification != message {
		t.Error("The received push notification does not match the expected value.")
	}
}
