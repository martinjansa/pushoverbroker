package main

// PushoverBroker represents the main class constructing the Pushover broker. It initializes the REST API server, processing logc & database.
// Depends on the PushNotificationsSender interface
type PushoverBroker struct {
	server                  *Server
	processor               *Processor
	PushNotificationsSender PushNotificationsSender
}

// NewPushoverBroker creates an instance of the PushoverBroker
func NewPushoverBroker(port int, certFilePath string, keyFilePath string, PushNotificationsSender PushNotificationsSender) *PushoverBroker {
	pb := new(PushoverBroker)
	pb.PushNotificationsSender = PushNotificationsSender

	// create new message processor
	pb.processor = NewProcessor(PushNotificationsSender, NewLimitsCounterImpl())

	// create new HTTP server
	pb.server = NewServer(port, certFilePath, keyFilePath, pb.processor)
	return pb
}

// Run starts the server
func (pb *PushoverBroker) Run() error {
	pb.processor.Run()
	pb.server.Run()
	return nil
}
