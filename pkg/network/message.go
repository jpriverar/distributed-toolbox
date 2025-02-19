package network

import "fmt"

type Message struct {
	Type string
	Meta map[string]string
	From string
	Payload []byte
}

func (msg Message) String() string {
	return fmt.Sprintf("Type: %s\nFrom: %s\nPayload: %v", msg.Type, msg.From, msg.Payload)
}