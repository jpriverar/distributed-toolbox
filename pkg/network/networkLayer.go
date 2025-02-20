package network

type NetworkLayer interface {
	Register(id, endpoint string) error
	Unregister(ip string) error
	Send(to string, message Message) error
	Receive(id string) (Message, error)
}