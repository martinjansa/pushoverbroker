package pushoverbroker

// PushoverBroker represents the main class constructing the Pushover broker. It initializes the REST API server, processing logc & database.
// Depends on the PushoverConnector interface
type PushoverBroker struct {
	server            *Server
	processor         *Processor
	pushoverConnector PushoverConnector
}

// NewPushoverBroker creates an instance of the PushoverBroker
func NewPushoverBroker(port int, pushoverConnector PushoverConnector) *PushoverBroker {
	pb := new(PushoverBroker)
	pb.pushoverConnector = pushoverConnector

	// create new message processor
	pb.processor = NewProcessor(pushoverConnector)

	// create new HTTP server
	pb.server = NewServer(port, pb.processor)
	return pb
}

// Run starts the server
func (pb *PushoverBroker) Run() error {
	go pb.server.Run()
	go pb.processor.Run()
	return nil
}
