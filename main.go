package pushoverbroker

//import (
//	"github.com/martinjansa/pushoverbroker"
//)

func main() {
	broker := NewPushoverBroker(8500, nil)
	broker.Run()
}
