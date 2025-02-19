package core

import "github.com/jpriverar/distributed-toolbox/pkg/network"

type Protocol interface {
	Start()
	Stop()
	HandleMessage(msg network.Message)
}