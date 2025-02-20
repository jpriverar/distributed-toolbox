package core

import (
	"net"
	"github.com/jpriverar/distributed-toolbox/pkg/network"
)

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