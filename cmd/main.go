package main

import "net"
import "fmt"
import "time"

func printPackets() {
	recv, err := net.ListenPacket("udp", "1.0.0.1:8080")
	if err != nil {
		panic("Could not start listening to packets: " + err.Error())
	}
	defer recv.Close()

	buf := make([]byte, 2048)
	for {
		n, addr, err := recv.ReadFrom(buf)
		if err != nil {
			panic("Could not read data")
		}

		fmt.Println("Got", n, "bytes from", addr, "->", string(buf))
	}
}

func main() {
	go printPackets()
	time.Sleep(3 * time.Second)
	send, err := net.Dial("udp", "127.0.0.1:8080")
	if err != nil {
		panic("Could not dial localhost:8080")
	}

	for i := 0; i < 10; i++ {
		send.Write([]byte(fmt.Sprintf("Message %d", i)))
		time.Sleep(1 * time.Second)
	}
}