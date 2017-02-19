package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"testing"
)

// TestAPI1MessageJSONShouldAcceptEmptyMessageViaTLSAndForwardToPushSender is a test function for the REST API call
func TestAPI1MessageJSONShouldAcceptEmptyMessageViaTLSAndForwardToPushSender(t *testing.T) {

	// **** GIVEN ****

	// get the certificate files path
	certFilePath := path.Join(path.Dir(os.Args[0]), "private", "server.cert.pem")
	keyFilePath := path.Join(path.Dir(os.Args[0]), "private", "server.key.pem")

	// The REST API server is initialized and connected to the message handler mock
	pcm := NewPushNotificationsSenderMock()
	port := 8501
	broker := NewPushoverBroker(port, certFilePath, keyFilePath, pcm)

	go broker.Run()

	// **** WHEN ****

	// prepare the test message
	expectedMessage := PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: ""}

	// encode message into the URL form values
	form := url.Values{}
	form.Set("token", "<dummy token>")
	form.Set("user", "<dummy user>")
	form.Set("message", "")
	formStr := form.Encode()

	// Prepare the POST request with form data
	urlStr := "https://localhost:" + strconv.Itoa(port) + "/1/messages.json"
	req, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(formStr))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formStr)))

	// initialize the client that does not check the certificates (for testing purposes only)
	tlsConfig := tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{TLSClientConfig: &tlsConfig}
	client := &http.Client{Transport: transport}

	// post the request
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
	pcm.AssertMessageAcceptedOnce(t, expectedMessage)
}
