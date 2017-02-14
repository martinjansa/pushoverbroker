package pushoverbroker_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
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

	// GIVEN REST API server is initialized and connected to the message handler mock
	messageHandlerMock := NewMessageHandlerMock()
	brokerServer := server.NewServer(messageHandlerMock)
	testRequest := server.Message{Request: "test request"}
	portStr := "8500"
	go brokerServer.Martini.RunOnAddr(":" + portStr)

	// WHEN a json POST request is sent via the REST API
	url := "http://localhost:" + portStr + "/1/messages.json"
	var jsonStr = []byte(testRequest.GetRequest())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	//req.Header.Set("X-Custom-Header", "myvalue")
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("POST request failed with error %s.", err)
	}
	defer resp.Body.Close()

	// CHECK

	// the request should respond correctly
	if resp.StatusCode != 200 {
		t.Errorf("POST request returned status code %d and status message %s. Expected code 200.", resp.StatusCode, resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	t.Logf("POST request response body '%s'.", string(body))

	// the right message shoud be delivered to the mock
	messageHandlerMock.AssertMessageAcceptedOnce(t, testRequest)
}
