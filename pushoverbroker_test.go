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

// implements the pushoverconnector.PushoverConnector interface
type PushoverConnectorMock struct {
	response            error
	handleMessageCalled int
	notification        pushoverbroker.PushNotification
}

// NewPushoverConnectorMock initializes the mock
func NewPushoverConnectorMock() *PushoverConnectorMock {
	pcm := new(PushoverConnectorMock)
	pcm.response = nil
	pcm.handleMessageCalled = 0
	return pcm
}

func (pcm *PushoverConnectorMock) PostPushNotificationMessage(message pushoverbroker.PushNotification) error {
	pcm.handleMessageCalled++
	pcm.notification = message
	return pcm.response
}

func (pcm *PushoverConnectorMock) AssertMessageAcceptedOnce(t *testing.T, message pushoverbroker.PushNotification) {
	if pcm.handleMessageCalled != 1 {
		t.Errorf("1 message expected, %d received.", pcm.handleMessageCalled)
	}
	if pcm.notification != message {
		t.Error("The received push notification does not match the expected value.")
	}
}

// TestShouldForwardEmptyMessage is a test function for the REST API call
func TestShouldForwardEmptyMessage(t *testing.T) {

	// **** GIVEN ****

	// The REST API server is initialized and connected to the message handler mock
	pcm := NewPushoverConnectorMock()
	broker := pushoverbroker.NewPushoverBroker(pcm)

	port := 8500
	go broker.Run(port)

	// **** WHEN ****

	//a json POST request is sent via the REST API
	url := "http://localhost:" + strconv.Itoa(port) + "/1/messages.json"
	testMessage := pushoverbroker.PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: ""}
	jsonData, err := json.Marshal(testMessage)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
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
	pcm.AssertMessageAcceptedOnce(t, testMessage)
}
