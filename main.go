package main

//import (
//	"github.com/martinjansa/pushoverbroker"
//)

func main() {
	broker := NewPushoverBroker(8499, nil)
	broker.Run()
}
