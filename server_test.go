package pushoverbroker_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"

	"github.com/martinjansa/pushoverbroker"
)

// implements the IncommingPushNotificationMessageHandler interface
type MessageHandlerMock struct {
	handleMessageCalled int
	notification        pushoverbroker.PushNotification
}

// NewMessageHandlerMock initializes the mock
func NewMessageHandlerMock() *MessageHandlerMock {
	mh := new(MessageHandlerMock)
	mh.handleMessageCalled = 0
	return mh
}

func (mh *MessageHandlerMock) HandleMessage(message pushoverbroker.PushNotification) {
	mh.handleMessageCalled++
	mh.notification = message
}

func (mh *MessageHandlerMock) AssertMessageAcceptedOnce(t *testing.T, message pushoverbroker.PushNotification) {
	if mh.handleMessageCalled != 1 {
		t.Errorf("1 message expected, %d received.", mh.handleMessageCalled)
	}
	if mh.notification != message {
		t.Error("The received push notification does not match the expected value.")
	}
}

// TestServerShouldAcceptPOST1MessagesJsonWithEmptyMessage is a test function for the REST API call
func TestServerShouldAcceptPOST1MessagesJsonWithEmptyMessage(t *testing.T) {

	// **** GIVEN ****

	// The REST API server is initialized and connected to the message handler mock
	port := 8501
	messageHandlerMock := NewMessageHandlerMock()
	brokerServer := pushoverbroker.NewServer(port, messageHandlerMock)
	go brokerServer.Run()

	// **** WHEN ****

	//a json POST request is sent via the REST API
	testMessage := pushoverbroker.PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: ""}
	jsonData, err := json.Marshal(testMessage)
	if err != nil {
		panic(err)
	}
	url := "http://localhost:" + strconv.Itoa(port) + "/1/messages.json"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	//req.Header.Set("X-Custom-Header", "myvalue")
	//req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Errorf("POST request failed with error %s.", err)
		return
	}
	defer resp.Body.Close()

	// **** THEN ****

	// the request should respond correctly
	if resp.StatusCode != 200 {
		t.Errorf("POST request returned status code %d and status message %s. Expected code 200.", resp.StatusCode, resp.Status)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	t.Logf("POST request response body '%s'.", string(body))

	// the right message shoud be delivered to the mock
	messageHandlerMock.AssertMessageAcceptedOnce(t, testMessage)
}
