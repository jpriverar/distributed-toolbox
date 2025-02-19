package network

import (
	"fmt"
	"sync"
)

type EmulnetConfig struct {
	avgLatency float64
	dropRate   float64
}

type Emulnet struct {
	config *EmulnetConfig
	nodes  map[string]chan Message
	lock   sync.RWMutex
}

func NewEmulnet(config *EmulnetConfig) *Emulnet {
	return &Emulnet{
		config: config,
		nodes:  make(map[string]chan Message),
	}
}

func (net *Emulnet) RegisterNode(id string) error {
	net.lock.Lock()
	defer net.lock.Unlock()

	if _, exists := net.nodes[id]; exists {
		return fmt.Errorf("node with ID %s is already registered in the network", id)
	}

	incoming := make(chan Message, 100)
	net.nodes[id] = incoming
	return nil
}

func (net *Emulnet) UnregisterNode(id string) error {
	net.lock.Lock()
	defer net.lock.Unlock()

	if incoming, exists := net.nodes[id]; exists {
		close(incoming)
		delete(net.nodes, id)
		return nil
	} else {
		return fmt.Errorf("node with ID %s is not registered in the network", id)
	}
}

func (net *Emulnet) Send(to string, msg Message) error {
	net.lock.RLock()
	defer net.lock.RUnlock()

	if dstChan, exists := net.nodes[to]; exists {
		dstChan <- msg
		return nil
	} else {
		return fmt.Errorf("node with ID %s is not registered in the network", to)
	}

}

func (net *Emulnet) Receive(id string) (Message, error) {
	net.lock.RLock()
	defer net.lock.RUnlock()

	if dstChan, exists := net.nodes[id]; exists {
		return <- dstChan, nil
	} else {
		return Message{}, fmt.Errorf("node with ID %s is not registered in the network", id)
	}
}

func (net *Emulnet) Start() {

}

func (net *Emulnet) Stop() {

}