package network

type NetworkLayer interface {
	Send(to string, message Message) error
	Receive(id string) (Message, error)
	Start()
	Stop()
}