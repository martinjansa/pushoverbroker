package main

import (
	"log"
	"net/http"
	"strconv"

	"fmt"

	"github.com/gorilla/schema"
)

// Limits represents the values of the message counts limits of the Pushover account
type Limits struct {
	limit     int
	remaining int
	reset     int
}

// PushNotificationHandlingResponse contains the response parameters from handling of the incomming push notification message
type PushNotificationHandlingResponse struct {
	responseCode     int     // HTTP response code
	limits           *Limits // information about the current accounts limits
	jsonResponseBody string  // JSON response body (if propagating from external service)
}

// IncommingPushNotificationMessageHandler handles message accepted by the REST API
type IncommingPushNotificationMessageHandler interface {
	HandleMessage(response *PushNotificationHandlingResponse, message PushNotification) error
}

// Server is the REST API server that handles the clients connections
type Server struct {
	mux          *http.ServeMux
	server       *http.Server
	certFilePath string
	keyFilePath  string
}

// NewServer creates a new server. Accepts the messageHandler that will handle all the received messages
func NewServer(port int, certFilePath string, keyFilePath string, messageHandler IncommingPushNotificationMessageHandler) *Server {
	s := new(Server)

	// create and inititalize the multiplexer
	s.mux = http.NewServeMux()
	s.certFilePath = certFilePath
	s.keyFilePath = keyFilePath

	// handler of the POST messages to /1/messages.json
	h1 := new(Post1MessageJSONHTTPHandler)
	h1.messageHandler = messageHandler
	// create a schema decoder
	h1.decoder = schema.NewDecoder()
	h1.decoder.IgnoreUnknownKeys(true)

	s.mux.Handle("/1/messages.json", h1)

	// create and initialize the HTTP server
	s.server = new(http.Server)
	s.server.Addr = ":" + strconv.Itoa(port)
	s.server.Handler = s.mux

	return s
}

// Run starts the HTTP server and listens and serves the incoming requests
func (s *Server) Run() {
	log.Fatal(s.server.ListenAndServeTLS(s.certFilePath, s.keyFilePath))
}

// Post1MessageJSONHTTPHandler handles the POST request at /1/messages.json
type Post1MessageJSONHTTPHandler struct {
	messageHandler IncommingPushNotificationMessageHandler
	decoder        *schema.Decoder
}

// WriteSuccessJSONResponse writes the response header and JSON body
func WriteSuccessJSONResponse(w http.ResponseWriter, responseCode int, request string) {
	responseBody := fmt.Sprintf("{\"status\": 1, \"request\": \"%s\"}", request)
	log.Printf("Writing response with status code %d and body %s.", responseCode, responseBody)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(responseCode)
	w.Write([]byte(responseBody))
}

// WriteErrorJSONResponse writes the response header and JSON body with error string
func WriteErrorJSONResponse(w http.ResponseWriter, responseCode int, request string, errorStr string) {
	responseBody := fmt.Sprintf("{\"status\": 0, \"request\": \"%s\", \"errors\":[\"%s\"]}", request, errorStr)
	log.Printf("Writing response with status code %d and body %s.", responseCode, responseBody)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(responseCode)
	w.Write([]byte(responseBody))
}

// handles the incomming request and forwards it to the message handler
func (h *Post1MessageJSONHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	request := "TODO: generate random"

	// if the request type is not POST
	if r.Method != "POST" {
		WriteErrorJSONResponse(w, 400, request, fmt.Sprintf("Received request of method '%s', expected 'POST'", r.Method))
		return
	}

	// does the request does not contain the requested content type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/x-www-form-urlencoded" {
		WriteErrorJSONResponse(w, 400, request, fmt.Sprintf("Received request with unsupported Content-Type %s, expected application/x-www-form-urlencoded", contentType))
		return
	}

	// parse the form
	err := r.ParseForm()
	if err != nil {
		WriteErrorJSONResponse(w, 400, request, fmt.Sprintf("The POST form parsing failed with error %s", err.Error()))
		return
	}

	// decode the POST form
	var pn PushNotification
	err = h.decoder.Decode(&pn, r.PostForm)
	if err != nil {
		WriteErrorJSONResponse(w, 400, request, fmt.Sprintf("The POST form decoding failed with error %s", err.Error()))
		return
	}
	//defer r.Body.Close()

	// if the message has all the mandatory fields token, user and message non empty
	err = pn.Validate()
	if err != nil {
		WriteErrorJSONResponse(w, 400, request, fmt.Sprintf("The POST form decoding failed with error %s. POST form content: '%s'", err.Error(), r.PostForm))
		return
	}
	// log the accepted message
	log.Printf("Received request with %s.", pn.DumpToString())

	// handle the message
	var response = PushNotificationHandlingResponse{}
	err = h.messageHandler.HandleMessage(&response, pn)

	// if the handling of the message failed
	if err != nil {

		// report the error
		WriteErrorJSONResponse(w, 500, request, fmt.Sprintf("Handling of the message %s failed with error %s, response code %d. Returning HTTP 500 (Internal Server Error)", pn.DumpToString(), err.Error(), response.responseCode))
		return
	}

	// if limits are provided
	if response.limits != nil {
		// construct the X-Limit-App-XXX headers
		w.Header().Set("X-Limit-App-Limit", strconv.Itoa(response.limits.limit))
		w.Header().Set("X-Limit-App-Remaining", strconv.Itoa(response.limits.remaining))
		w.Header().Set("X-Limit-App-Reset", strconv.Itoa(response.limits.reset))
	}

	// return the obtained response code
	WriteSuccessJSONResponse(w, response.responseCode, request)
}
