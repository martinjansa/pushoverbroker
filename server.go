package pushoverbroker

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

// IncommingPushNotificationMessageHandler handles message accepted by the REST API
type IncommingPushNotificationMessageHandler interface {
	HandleMessage(message PushNotification) error
}

// Server is the REST API server that handles the clients connections
type Server struct {
	mux    *http.ServeMux
	server *http.Server
}

// NewServer creates a new server. Accepts the messageHandler that will handle all the received messages
func NewServer(port int, messageHandler IncommingPushNotificationMessageHandler) *Server {
	s := new(Server)

	// create and inititalize the multiplexer
	s.mux = http.NewServeMux()

	// handler of the POST messages to /1/messages.json
	h1 := new(Post1MessageJSONHTTPHandler)
	h1.messageHandler = messageHandler
	s.mux.Handle("/1/messages.json", h1)

	// create and initialize the HTTP server
	s.server = new(http.Server)
	s.server.Addr = ":" + strconv.Itoa(port)
	s.server.Handler = s.mux

	return s
}

// Run starts the HTTP server and listens and serves the incoming requests
func (s *Server) Run() {
	s.server.ListenAndServe()
}

// Post1MessageJSONHTTPHandler handles the POST request at /1/messages.json
type Post1MessageJSONHTTPHandler struct {
	messageHandler IncommingPushNotificationMessageHandler
}

// handles the incomming request and forwards it to the message handler
func (h *Post1MessageJSONHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// parse the JSON
	decoder := json.NewDecoder(r.Body)
	var pn PushNotification
	err := decoder.Decode(&pn)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	log.Println(pn)

	// handle the message
	err = h.messageHandler.HandleMessage(pn)
	if err != nil {
		// report the error
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}

	// succeeded response
	w.WriteHeader(200)
}
