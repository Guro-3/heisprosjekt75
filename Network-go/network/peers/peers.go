package peers

import (
	"fmt"
	"net"
	"sort"
	"time"
)

type PeerUpdate struct {
	Peers []string
	New   string
	Lost  []string
}

const (
	interval  = 50 * time.Millisecond
	timeout   = 500 * time.Millisecond
	mcastIP   = "224.0.0.1" // multicast adresse
	mcastPort = 10334       // port alle noder bruker
)

// Transmitter sender ID-periodisk til multicast
func Transmitter(id string, transmitEnable <-chan bool) {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", mcastIP, mcastPort))
	if err != nil {
		panic(err)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			conn.Write([]byte(id))
		}
	}
}

// Receiver lytter på multicast og oppdaterer peers
func Receiver(myId string, peerUpdateCh chan<- PeerUpdate) {
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", mcastIP, mcastPort))
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	conn.SetReadBuffer(1024)

	var buf [1024]byte
	lastSeen := make(map[string]time.Time)
	var p PeerUpdate

	for {
		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, err := conn.ReadFromUDP(buf[:])
		if err != nil {
			// timeout ok
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				// fjerne døde peers
				updated := false
				p.Lost = nil
				for k, t := range lastSeen {
					if time.Since(t) > timeout {
						delete(lastSeen, k)
						p.Lost = append(p.Lost, k)
						updated = true
					}
				}
				if updated {
					p.Peers = make([]string, 0, len(lastSeen))
					for k := range lastSeen {
						p.Peers = append(p.Peers, k)
					}
					sort.Strings(p.Peers)
					sort.Strings(p.Lost)
					peerUpdateCh <- p
				}
				continue
			} else {
				panic(err)
			}
		}

		peerId := string(buf[:n])
		if peerId == "" || peerId == myId {
			continue
		}

		// Ny peer?
		updated := false
		p.New = ""
		if _, ok := lastSeen[peerId]; !ok {
			p.New = peerId
			updated = true
		}

		lastSeen[peerId] = time.Now()

		// fjern døde peers
		p.Lost = nil
		for k, t := range lastSeen {
			if time.Since(t) > timeout {
				delete(lastSeen, k)
				p.Lost = append(p.Lost, k)
				updated = true
			}
		}

		if updated {
			p.Peers = make([]string, 0, len(lastSeen))
			p.Peers = append(p.Peers, myId)
			for k := range lastSeen {
				p.Peers = append(p.Peers, k)
			}
			sort.Strings(p.Peers)
			sort.Strings(p.Lost)
			peerUpdateCh <- p
		}
	}
}
