package network

import (
	"strconv"
	"testing"
)

func TestSingleMessage(t *testing.T) {
	network := NewEmulnet(&EmulnetConfig{
							avgLatency: 0,
							dropRate: 0,
						})

	network.Register("Node0")

	go network.Send("Node0", Message{
		Type: "test",
		From: "coordinator",
		Payload: []byte("test message"),
	})


	msg, err := network.Receive("Node0")
	if err != nil {
		t.Fatalf("error receiving message for Node0: %v", err)
	}

	payload := string(msg.Payload)
	if payload != "test message" {
		t.Errorf("expected payload to be 'test message', got %s", payload)
	}
}

func TestMultipleMessages(t *testing.T) {
	network := NewEmulnet(&EmulnetConfig{
							avgLatency: 0,
							dropRate: 0,
						})

	nodeCount := 10
	for i := 0; i < nodeCount; i++ {
		network.Register("Node" + strconv.Itoa(i))
	}

	messageCount := 100
	for i := 0; i < nodeCount; i++ {
		for j := 0; j < messageCount; j++ {
			go network.Send("Node" + strconv.Itoa(i), Message{
				Type: "test",
				From: "coordinator",
				Payload: []byte{byte(i), byte(j)},
			})
		}
	}

	for i := 0; i < nodeCount; i++ {
		receivedMsgs := make([]int, messageCount)
		for j := 0; j < messageCount; j++ {
			msg, err := network.Receive("Node" + strconv.Itoa(i))
			if err != nil {
				t.Fatalf("error receiving message for Node%d: %v", i, err)
			}

			dstNode := msg.Payload[0]
			msgNum := msg.Payload[1]
			if dstNode != byte(i) {
				t.Errorf("received message for Node%d, expected Node%d", dstNode, i)
			}

			receivedMsgs[msgNum]++
		}

		for j := 0; j < messageCount; j++ {
			if receivedMsgs[j] == 0 {
				t.Errorf("Missing expected message %d for Node%d", j, i)
			} else if receivedMsgs[j] > 1 {
				t.Errorf("Received duplicated message %d for Node%d", j, i)
			}
		}
	}
}