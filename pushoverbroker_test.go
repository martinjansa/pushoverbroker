package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
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

// TestAPI1MessageJSON is a test function for the REST API call
func TestBrokerAPI1MessageJSON(t *testing.T) {

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
	time.Sleep(100 * time.Millisecond)

	t.Run("ShouldForwardToPushSender", func(t *testing.T) {

		var testcases = []struct {
			id                         string
			urlValues                  map[string]string
			responseErr                error
			responseStatusCode         int
			responseBody               string
			expectedMessage            PushNotification
			expectedStatusCode         int
			expectedResponseBodyStatus int
		}{
			{
				// checks that the success result is propagated if a call to push notification sender succeeds
				"ShouldReturnSuccessFromExternalAPI",
				map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": "<dummy message>"}, nil, 200, "{\"status\": 1}",
				PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: "<dummy message>"}, 200, 1,
			},
			{
				// checks that the success result 202 Accepted is returned if an attempt to push notification sender fails
				"ShouldReturnAcceptedOnPostError",
				map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": "<dummy message>"}, errors.New("posting failed, no internet"), 0, "{\"status\": 1}",
				PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: "<dummy message>"}, 202, 1,
			},
			{
				// checks that the success result 202 Accepted is returned if an attempt to push notification sender fails with temporary server error 500
				"ShouldReturnAcceptedOnPostError",
				map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": "<dummy message>"}, nil, 500, "{\"status\": 0}",
				PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: "<dummy message>"}, 202, 1,
			},
			{
				// checks that the success result 400 (Bad Request) is returned if the push notification sender API calls returns this status code
				"ShouldReturnBadRequestFromExternalAPI",
				map[string]string{"token": "<dummy token>", "user": "<dummy user>", "message": "<dummy message>"}, nil, 400, "{\"status\": 0}",
				PushNotification{Token: "<dummy token>", User: "<dummy user>", Message: "<dummy message>"}, 400, 0,
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
				pcm.ForceResponse(tc.responseErr, tc.responseStatusCode, &Limits{limit: 10000, remaining: 10, reset: 12345}, tc.responseBody)

				// post the request
				resp, err := client.Do(req)
				if err != nil {
					t.Errorf("POST request failed with error %s, but was expected to succeed.", err)
					return
				}
				defer resp.Body.Close()

				// **** THEN ****

				// the right message shoud be delivered to the mock
				pcm.AssertMessageAcceptedOnce(t, tc.expectedMessage)

				// check the expected response code
				if resp.StatusCode != tc.expectedStatusCode {
					t.Errorf("POST request returned status code %d and status message %s. Expected code %d.", resp.StatusCode, resp.Status, tc.expectedStatusCode)
					return
				}

				// check the response content type
				responseContentType := resp.Header.Get("Content-Type")
				if responseContentType != "application/json; charset=utf-8" {
					t.Errorf("POST request returned content type %s, expected application/json; charset=utf-8.", responseContentType)
					return
				}

				// get the content of the body
				//body, _ := ioutil.ReadAll(resp.Body)
				type ResponseJSONBodyContent struct {
					Status  int    `json:"status"`
					Request string `json:"request"`
				}
				var responseJSONBodyContent = ResponseJSONBodyContent{0, ""}
				err = json.NewDecoder(resp.Body).Decode(&responseJSONBodyContent)
				if err != nil {
					body, _ := ioutil.ReadAll(resp.Body)
					t.Errorf("POST request returned JSON \"%s\", which failed to decode with error %s.", body, err.Error())
					return
				}

				// check the status value
				if responseJSONBodyContent.Status != tc.expectedResponseBodyStatus {
					t.Errorf("POST request returned JSON with status %d, expected status was %d.", responseJSONBodyContent.Status, tc.expectedResponseBodyStatus)
					return
				}
			})
		}
	})

	t.Run("ShouldReturnLimits", func(t *testing.T) {

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
				pcm.ForceResponse(tc.responseErr, tc.responseStatusCode, tc.limits, "{\"status\": 1}")

				// post the request
				resp, err := client.Do(req)
				if err != nil {
					t.Errorf("POST request failed with error %s, but was expected to succeed.", err)
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
