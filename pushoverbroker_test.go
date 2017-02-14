package pushoverbroker_test

import (
	"strconv"
	"testing"

	"github.com/martinjansa/pushoverbroker/server"
)

// implements the server.MessageHadler interface
type MessageHandlerMock struct {
	handleMessageCalled  int
	handleMessageRequest string
}

// NewMessageHandlerMock initializes the mock
func NewMessageHandlerMock() *MessageHandlerMock {
	mh := new(MessageHandlerMock)
	mh.handleMessageCalled = 0
	mh.handleMessageRequest = ""
	return mh
}

func (mh *MessageHandlerMock) HandleMessage(message server.Message) {
	mh.handleMessageCalled++
	mh.handleMessageRequest = message.GetRequest()
}

func (mh *MessageHandlerMock) AssertMessageAcceptedOnce(t *testing.T, message server.Message) {
	if mh.handleMessageCalled != 1 {
		t.Errorf("1 message expected, %d received.", mh.handleMessageCalled)
	}
	if mh.handleMessageRequest != message.GetRequest() {
		t.Error("The received request does not match the expected value.")
	}
}

// TestAPIShouldAcceptPOST1MessagesJsonWithEmptyMessage is a test function for the REST API call
func TestAPIShouldAcceptPOST1MessagesJsonWithEmptyMessage(t *testing.T) {
	messageHandlerMock := NewMessageHandlerMock()
	brokerServer := server.NewServer(messageHandlerMock)
	testRequest := server.Message{Request: "test request"}
	port := 8500
	go brokerServer.Martini.RunOnAddr(":" + strconv.Itoa(port))
	messageHandlerMock.AssertMessageAcceptedOnce(t, testRequest)
}
