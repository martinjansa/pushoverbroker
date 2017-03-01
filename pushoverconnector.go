package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/schema"
)

// PushoverConnector sends push notifications to Pushover service
type PushoverConnector struct {
	client  *http.Client
	encoder *schema.Encoder
}

// NewPushoverConnector creates a new pushover connector
func NewPushoverConnector() *PushoverConnector {
	pc := new(PushoverConnector)
	pc.client = &http.Client{}
	pc.encoder = schema.NewEncoder()
	return pc
}

// PostPushNotificationMessage post a message to the Pushover server and returns error if ocurred (or nil) and response code (or 0 on POST error)
func (pc *PushoverConnector) PostPushNotificationMessage(response *PushNotificationHandlingResponse, message PushNotification) error {

	// encode message into the URL form values
	form := url.Values{}
	err := pc.encoder.Encode(message, form)
	//if err != nil {
	//	return fmt.Errorf("encoding of the message \"%s\" failed with error %s", message, err)
	//}
	formStr := form.Encode()

	// Prepare the POST request with form data
	url := "https://api.pushover.net/1/messages.json"
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(formStr))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(formStr)))

	resp, err := pc.client.Do(req)
	if err != nil {
		response.responseCode = 0
		response.limits = nil
		return fmt.Errorf("sending the Pushover API POST request at %s with foorm \"%s\" failed with error %s", url, form, err)
	}
	defer resp.Body.Close()

	// get the body
	body, _ := ioutil.ReadAll(resp.Body)

	// the request should respond correctly
	if resp.StatusCode != 200 {
		response.responseCode = 0
		response.limits = nil
		return fmt.Errorf("processing of the Pushover API POST request at %s with form \"%s\" returned status code %d, status message %s, body \"%s\"", url, form, resp.StatusCode, resp.Status, string(body))
	}

	// get the response header values
	limitValue := resp.Header.Get("X-Limit-App-Limit")
	remainingValue := resp.Header.Get("X-Limit-App-Remaining")
	resetValue := resp.Header.Get("X-Limit-App-Reset")

	// convert the limits to numbers
	limitValueInt, err := strconv.Atoi(limitValue)
	if err != nil {
		fmt.Printf("Obtained X-Limit-App-Limit value \"%s\" failed to be converted to number with error \"%s\".", limitValue, err.Error())
	}
	remainingValueInt, err := strconv.Atoi(remainingValue)
	if err != nil {
		fmt.Printf("Obtained X-Limit-App-Remaining value \"%s\" failed to be converted to number with error \"%s\".", remainingValue, err.Error())
	}
	resetValueInt, err := strconv.Atoi(resetValue)
	if err != nil {
		fmt.Printf("Obtained X-Limit-App-Reset value \"%s\" failed to be converted to number with error \"%s\".", resetValue, err.Error())
	}
	response.limits = &Limits{limitValueInt, remainingValueInt, resetValueInt}
	response.responseCode = 0
	response.jsonResponseBody = string(body)

	//t.Logf("POST request response body '%s'.", string(body))

	return nil
}
