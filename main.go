package main

import (
	"os"
	"path"
)

//import (
//	"github.com/martinjansa/pushoverbroker"
//)

func main() {
	// get the certificate files path
	certFilePath := path.Join(path.Dir(os.Args[0]), "private", "server.cert.pem")
	keyFilePath := path.Join(path.Dir(os.Args[0]), "private", "server.key.pem")

	// initialize the server
	pushoverConnector := NewPushoverConnector()
	broker := NewPushoverBroker(8499, certFilePath, keyFilePath, pushoverConnector)
	broker.Run()
}
