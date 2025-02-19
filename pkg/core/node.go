package core

import (
	"net"
	"bytes"
	"github.com/google/uuid"
	"github.com/jpriverar/distributed-toolbox/pkg/network"
)

type ID uuid.UUID

func (id ID) Equals(other ID) bool {
	return bytes.Equal(id[:], other[:])
}

func (id ID) GreaterThan(other ID) bool {
	return bytes.Compare(id[:], other[:]) == 1
}

func (id ID) LessThan(other ID) bool {
	return bytes.Compare(id[:], other[:]) == -1
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

func NewID() ID {
	return ID(uuid.New())
}

type Node struct {
	Id ID
	Addr net.IP
	protocols map[string]Protocol
	network network.NetworkLayer
}

func NewNode(id ID, addr net.IP, network network.NetworkLayer) *Node {
	return &Node{
		Id: id,
		Addr: addr,
		protocols: make(map[string]Protocol),
		network: network,
	}
}

func (n *Node) SendMessage(to string, msg network.Message) error {
	return n.network.Send(to, msg)
}