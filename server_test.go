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
	"testing"
	"time"
)

// implements the IncommingPushNotificationMessageHandler interface
type MessageHandlerMock struct {
	responseErr         error
	responseCode        int
	limits              *Limits
	handleMessageCalled int
	notification        PushNotification
}

// NewMessageHandlerMock initializes the mock
func NewMessageHandlerMock() *MessageHandlerMock {
	mh := new(MessageHandlerMock)
	mh.responseErr = nil
	mh.responseCode = 200
	mh.handleMessageCalled = 0
	return mh
}

func (mh *MessageHandlerMock) HandleMessage(message PushNotification) (error, int, *Limits) {
	mh.handleMessageCalled++
	mh.notification = message
	return mh.responseErr, mh.responseCode, mh.limits
}

func (mh *MessageHandlerMock) ForceResponse(responseErr error, reseponseCode int, limits *Limits) {
	mh.handleMessageCalled = 0
	mh.responseErr = responseErr
	mh.responseCode = reseponseCode
	mh.limits = limits
}

func (mh *MessageHandlerMock) AssertMessageAcceptedOnce(t *testing.T, message PushNotification) {
	if mh.handleMessageCalled != 1 {
		t.Errorf("1 message expected, %d received.", mh.handleMessageCalled)
	}
	if mh.notification != message {
		t.Error("The received push notification does not match the expected value.")
	}
}

// TestServer is a test function for the REST API call
func TestServer(t *testing.T) {

	// **** GIVEN ****

	// get the certificate files path
	wd, _ := os.Getwd()
	certFilePath := path.Join(wd, "private", "server.cert.pem")
	keyFilePath := path.Join(wd, "private", "server.key.pem")

	// The REST API server is initialized and connected to the message handler mock
	messageHandlerMock := NewMessageHandlerMock()
	port := 8502
	brokerServer := NewServer(port, certFilePath, keyFilePath, messageHandlerMock)

	// start the server
	go brokerServer.Run()

	// give the HTTP server enough time to start listening for the new connections
	time.Sleep(100 * time.Millisecond)

	t.Run("API1MessageJSONShouldForwardToMessageHandler", func(t *testing.T) {

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
				"ShouldReturnSuccessForEmptyMessage",
				map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": "<dummy message>"}, nil, 200,
				PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: "<dummy message>"}, 200,
			},
			{
				// checks that the success result is propagated if a call to push notification sender succeeds
				"ShouldReturnSuccessForNonEmptyMessage",
				map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": "<dummy message>"}, nil, 200,
				PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: "<dummy message>"}, 200,
			},
			{
				// checks that the success result 202 Accepted is returned if this code is passed from the message handler
				"ShouldReturn202FromNotificationHandler",
				map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": "<dummy message>"}, nil, 202,
				PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: "<dummy message>"}, 202,
			},
			{
				// checks that the error code 403 is returned if handling in the notification handler fails returns this status code
				"ShouldReturn403FromNotificationHandler",
				map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": "<dummy message>"}, nil, 403,
				PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: "<dummy message>"}, 403,
			},
			{
				// checks that the error code 429 is returned if handling in the notification handler fails returns this status code
				"ShouldReturn429FromNotificationHandler",
				map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": "<dummy message>"}, nil, 429,
				PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: "<dummy message>"}, 429,
			},
			{
				// checks that the error code 500 is returned if handling in the notification handler fails returns this status code
				"ShouldReturn500FromNotificationHandler",
				map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": "<dummy message>"}, nil, 500,
				PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: "<dummy message>"}, 500,
			},
			{
				// checks that the internal server error 500 is returned if handling in the notification handler fails
				"ShouldReturn500OnInternalError",
				map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": "<dummy message>"}, errors.New("any internal error"), 0,
				PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: "<dummy message>"}, 500,
			},
		}

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
				messageHandlerMock.ForceResponse(tc.responseErr, tc.responseStatusCode, nil)

				// post the request
				resp, err := client.Do(req)
				if err != nil {
					t.Errorf("POST request failed with error '%s', but was expected to succeed.", err)
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
				messageHandlerMock.AssertMessageAcceptedOnce(t, tc.expectedMessage)
			})
		}
	})

	t.Run("API1MessageJSONShouldReturnLimits", func(t *testing.T) {

		var testcases = []struct {
			id                 string
			responseErr        error
			responseStatusCode int
			limits             *Limits
			limitsExpected     bool
			expectedLimit      string
			expectedRemaining  string
			expectedReset      string
		}{
			{"ShouldReturnLimitsOnSuccess200", nil, 200, &Limits{limit: 1000, remaining: 500, reset: 123456789}, true, "1000", "500", "123456789"},
			{"ShouldReturnLimitsOnFailure400", nil, 400, nil, false, "", "", ""},
		}

		for _, tc := range testcases {

			t.Run(tc.id, func(t *testing.T) {

				// **** WHEN ****

				// encode message into the URL form values
				form := url.Values{}
				form.Set("token", "<dummy token>")
				form.Set("user", "<dummy user>")
				form.Set("message", "<dummy message>")
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
				messageHandlerMock.ForceResponse(tc.responseErr, tc.responseStatusCode, tc.limits)

				// post the request
				resp, err := client.Do(req)
				if err != nil {
					t.Errorf("POST request failed with error '%s', but was expected to succeed.", err)
					return
				}
				defer resp.Body.Close()

				// **** THEN ****

				// get the response header values
				limitValue := resp.Header.Get("X-Limit-App-Limit")
				remainingValue := resp.Header.Get("X-Limit-App-Remaining")
				resetValue := resp.Header.Get("X-Limit-App-Reset")

				// check the expected limits
				if tc.limitsExpected {

					// if any of the obtained values does not match the expected value
					if limitValue != tc.expectedLimit {
						t.Errorf("The received value of the X-Limi-App-Limit response header \"%s\" does not match the expected value \"%s\".", limitValue, tc.expectedLimit)
					}

					if remainingValue != tc.expectedRemaining {
						t.Errorf("The received value of the X-Limi-App-Remaining response header \"%s\" does not match the expected value \"%s\".", remainingValue, tc.expectedRemaining)
					}

					if resetValue != tc.expectedReset {
						t.Errorf("The received value of the X-Limi-App-Reset response header \"%s\" does not match the expected value \"%s\".", resetValue, tc.expectedReset)
					}

				} else { // the limits were not expected

					// if any of the values was obtained
					if limitValue != "" || remainingValue != "" || resetValue != "" {

						t.Errorf("The X-Limi-App-XXX response headers were received but not expected Limit=\"%s\", Remaining=\"%s\", Reset=\"%s\".", limitValue, remainingValue, resetValue)
					}
				}
			})
		}
	})
}
