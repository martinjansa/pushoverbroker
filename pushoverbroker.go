package pushoverbroker

import "errors"

// PushoverBroker represents the main class constructing the Pushover broker. It initializes the REST API server, processing logc & database.
// Depends on the PushoverConnector interface
type PushoverBroker struct {
	pushoverConnector PushoverConnector
}

// NewPushoverBroker creates an instance of the PushoverBroker
func NewPushoverBroker(pushoverConnector PushoverConnector) *PushoverBroker {
	pb := new(PushoverBroker)
	pb.pushoverConnector = pushoverConnector
	return pb
}

// Run starts the server
func (pb *PushoverBroker) Run(port int) error {
	return errors.New("Not implemented")
}
