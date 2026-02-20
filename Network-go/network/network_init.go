package network

import (
	"fmt"
	"heisprosjekt75/Network-go/network/localip"
	"heisprosjekt75/Network-go/network/peers"
	"os"
)

func NetworkInit() (id string, peerUpdateCh <-chan peers.PeerUpdate) {
	// // Our id can be anything. Here we pass it on the command line, using
	// //  `go run main.go -id=our_id`
	// var id string
	// flag.StringVar(&id, "id", "", "id of this peer")
	// flag.Parse()

	// // ... or alternatively, we can use the local IP address.
	// // (But since we can run multiple programs on the same PC, we also append the
	// //  process ID)

	// 1) Lag ID
	ip, err := localip.LocalIP()
	if err != nil {
		ip = "DISCONNECTED"
	}
	id = fmt.Sprintf("%s-%d", ip, os.Getpid())

	// 2) Lag kanal som main kan lese fra
	ch := make(chan peers.PeerUpdate, 10)

	// 3) Start peers
	peerTxEnable := make(chan bool, 1)
	peerTxEnable <- true

	go peers.Transmitter(10334, id, peerTxEnable)
	go peers.Receiver(10334, id, ch)

	fmt.Println("[NET] started peers with id:", id)

	return id, ch
}

// // We make a channel for receiving updates on the id's of the peers that are
// //  alive on the network

// // We can disable/enable the transmitter after it has been started.
// // This could be used to signal that we are somehow "unavailable".
// peerTxEnable := make(chan bool)
// go peers.Transmitter(15647, id, peerTxEnable)
// go peers.Receiver(15647, peerUpdateCh)

// // ... and start the transmitter/receiver pair on some port
// // These functions can take any number of channels! It is also possible to
// //  start multiple transmitters/receivers on the same port.
// go bcast.Transmitter(16569, UDPHeartbeatTx)
// go bcast.Receiver(16569, UDPHeartbeatRx)

// // The example message. We just send one of these every second.
// go sendHeartbeat() {
// 	/*helloMsg := HelloMsg{"Hello from " + id, 0}
// 	for {
// 		helloMsg.Iter++
// 		helloTx <- helloMsg
// 		time.Sleep(1 * time.Second)
// 	}*/
// }()
