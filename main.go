package pushoverbroker

//import (
//	"github.com/martinjansa/pushoverbroker"
//)

func main() {
	broker := NewPushoverBroker(nil)
	broker.Run(8500)
}
