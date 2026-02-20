package network

import (
	"fmt"
	"heisprosjekt75/Network-go/network/peers"
	"os"
)

func NetworkInit() (id string, peerUpdateCh <-chan peers.PeerUpdate) {
	id = fmt.Sprintf("%s-%d", "localhost", os.Getpid())
	ch := make(chan peers.PeerUpdate, 10)

	peerTxEnable := make(chan bool, 1)
	peerTxEnable <- true

	go peers.Transmitter(id, peerTxEnable)
	go peers.Receiver(id, ch)

	fmt.Println("[NET] started peers with id:", id)
	return id, ch
}
