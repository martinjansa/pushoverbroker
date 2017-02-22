package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
	"testing"
)

// TestAPI1MessageJSONShouldForwardToPushSenderAndReturnPostFailure is a test function for the REST API call
func TestAPI1MessageJSONShouldForwardToPushSender(t *testing.T) {

	var testcases = []struct {
		id                 string
		urlValues          map[string]string
		responseErr        error
		responseStatusCode int
		expectedMessage    PushNotification
		expectedStatusCode int
	}{
		{
			// checks that the success result is propagated if a call to push notification sender succeeds
			"ShouldReturnSuccessFromExternalAPI",
			map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": ""}, nil, 200,
			PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: ""}, 200,
		},
		{
			// checks that the success result 202 Accepted is returned if an attempt to push notification sender fails
			"ShouldReturnAcceptedOnPostError",
			map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": ""}, errors.New("posting failed, no internet"), 0,
			PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: ""}, 202,
		},
		{
			// checks that the success result 400 (Bad Request) is returned if the push notification sender API calls returns this status code
			"ShouldReturnBadRequestFromExternalAPI",
			map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": ""}, nil, 400,
			PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: ""}, 400,
		},
	}

	// **** GIVEN ****

	// get the certificate files path
	wd, _ := os.Getwd()
	certFilePath := path.Join(wd, "private", "server.cert.pem")
	keyFilePath := path.Join(wd, "private", "server.key.pem")

	// The REST API server is initialized and connected to the message handler mock
	pcm := NewPushNotificationsSenderMock()
	port := 8501
	broker := NewPushoverBroker(port, certFilePath, keyFilePath, pcm)

	// start the broker
	go broker.Run()

	// give the HTTP server enough time to start listening for the new connections
	time.Sleep(100*time.Millisecond)

	for _, tc := range testcases {

		t.Run(tc.id, func(t *testing.T) {

			// **** WHEN ****

			// encode message into the URL form values
			form := url.Values{}
			for name, value := range tc.urlValues {
				form.Set(name, value)
			}
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

			// force the moct to fail the posting of the push notification message
			pcm.ForceResponse(tc.responseErr, tc.responseStatusCode)

			// post the request
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("POST request failed with error %s, but was expected to succeed.", err)
				return
			}
			defer resp.Body.Close()

			// **** THEN ****

			// check the expected response code
			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("POST request returned status code %d and status message %s. Expected code %d.", resp.StatusCode, resp.Status, tc.expectedStatusCode)
			}
			body, _ := ioutil.ReadAll(resp.Body)
			t.Logf("POST request response body '%s'.", string(body))

			// the right message shoud be delivered to the mock
			pcm.AssertMessageAcceptedOnce(t, tc.expectedMessage)
		})
	}
}
