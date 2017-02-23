package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/schema"
)

// Limits represents the values of the message counts limits of the Pushover account
type Limits struct {
	limit     int
	remaining int
	reset     int
}

// IncommingPushNotificationMessageHandler handles message accepted by the REST API
type IncommingPushNotificationMessageHandler interface {
	HandleMessage(message PushNotification) (error, int, *Limits)
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

// handles the incomming request and forwards it to the message handler
func (h *Post1MessageJSONHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// parse the form
	err := r.ParseForm()
	if err != nil {
		log.Panicf("POST form parsing failed with error %s.", err.Error())
	}

	// decode the POST form
	var pn PushNotification
	err = h.decoder.Decode(&pn, r.PostForm)
	if err != nil {
		log.Panicf("POST form decoding failed with error %s.", err.Error())
	}
	//defer r.Body.Close()

	// log the accepted message
	log.Printf("Received request with %s.", pn.DumpToString())

	// handle the message
	err, responseCode, limits := h.messageHandler.HandleMessage(pn)

	// if the handling of the message failed
	if err != nil {

		// report the error
		log.Printf("Handling of the message %s failed with error %s, response code %d. Returning HTTP 500 (Internal Server Error).", pn.DumpToString(), err.Error(), responseCode)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	// if limits are provided
	if limits != nil {
		// construct the X-Limit-App-XXX headers
		w.Header().Set("X-Limit-App-Limit", strconv.Itoa(limits.limit))
		w.Header().Set("X-Limit-App-Remaining", strconv.Itoa(limits.remaining))
		w.Header().Set("X-Limit-App-Reset", strconv.Itoa(limits.reset))
	}

	// return the obtained response code
	w.WriteHeader(responseCode)
}
