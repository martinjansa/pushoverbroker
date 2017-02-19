package main

//import (
//	"github.com/martinjansa/pushoverbroker"
//)

func main() {
	pushoverConnector := NewPushoverConnector()
	broker := NewPushoverBroker(8499, "", "", pushoverConnector)
	broker.Run()
}
