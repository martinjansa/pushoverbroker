package pushoverbroker_test

import (
	"testing"

	"github.com/martinjansa/pushoverbroker"
)

// TestShouldSendMessageToPushoverConnector tests whether the processor attempts to send all the incomming messages to Pushover connector
func TestShouldSendMessageToPushoverConnector(t *testing.T) {

	// **** GIVEN ****

	// The REST API server is initialized and connected to the message handler mock
	pcm := pushoverbroker.NewPushoverConnectorMock()
	processor := pushoverbroker.NewPushoverProcessor(pcm)

	go processor.Run()

	// **** WHEN ****

	// a push notification is obtained by the process (via IncommingPushNotificationMessageHandler interface method HandleMessage())
	testMessage := pushoverbroker.PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: ""}
	err := processor.HandleMessage(testMessage)

	// **** THEN ****

	// the request should respond correctly
	if err != nil {
		t.Errorf("Handling of the message failed with error %s.", err)
		return
	}

	// the right message shoud be delivered to the mock
	pcm.AssertMessageAcceptedOnce(t, testMessage)
}
