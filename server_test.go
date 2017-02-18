package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
)

// implements the IncommingPushNotificationMessageHandler interface
type MessageHandlerMock struct {
	handleMessageCalled int
	notification        PushNotification
}

// NewMessageHandlerMock initializes the mock
func NewMessageHandlerMock() *MessageHandlerMock {
	mh := new(MessageHandlerMock)
	mh.handleMessageCalled = 0
	return mh
}

func (mh *MessageHandlerMock) HandleMessage(message PushNotification) error {
	mh.handleMessageCalled++
	mh.notification = message
	return nil
}

func (mh *MessageHandlerMock) AssertMessageAcceptedOnce(t *testing.T, message PushNotification) {
	if mh.handleMessageCalled != 1 {
		t.Errorf("1 message expected, %d received.", mh.handleMessageCalled)
	}
	if mh.notification != message {
		t.Error("The received push notification does not match the expected value.")
	}
}

// TestServerShouldAcceptPOST1MessagesJsonWithEmptyMessage is a test function for the REST API call
func ImplTestServerShouldAcceptPOST1MessagesJSON(t *testing.T, port int, message string) {

	// **** GIVEN ****

	// The REST API server is initialized and connected to the message handler mock
	messageHandlerMock := NewMessageHandlerMock()
	brokerServer := NewServer(port, messageHandlerMock)
	go brokerServer.Run()

	// **** WHEN ****

	//a json POST request is sent via the REST API
	expectedMessage := PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: message}

	// encode message into the URL form values
	form := url.Values{}
	form.Set("token", "<dummy token>")
	form.Set("user", "<dummy user>")
	form.Set("message", message)
	formStr := form.Encode()

	// Prepare the POST request with form data
	urlStr := "http://localhost:" + strconv.Itoa(port) + "/1/messages.json"
	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(formStr))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formStr)))

	// port the request
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
	messageHandlerMock.AssertMessageAcceptedOnce(t, expectedMessage)
}

// TestServerShouldAcceptPOST1MessagesJsonWithEmptyMessage is a test function for the REST API call
func TestServerShouldAcceptPOST1MessagesJsonWithEmptyMessage(t *testing.T) {

	ImplTestServerShouldAcceptPOST1MessagesJSON(t, 8502, "")
}

// TestServerShouldAcceptPOST1MessagesJsonWithMessage is a test function for the REST API call
func TestServerShouldAcceptPOST1MessagesJsonWithMessage(t *testing.T) {

	ImplTestServerShouldAcceptPOST1MessagesJSON(t, 8503, "<dummy message>")
}
