package membership

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/jpriverar/distributed-toolbox/pkg/core"
	"github.com/jpriverar/distributed-toolbox/pkg/network"
)

var GOSSIP string = "MEMBER_GOSSIP"
var JOINREQ string = "MEMBER_JOIN_REQUEST"
var JOINRES string = "MEMBER_JOIN_RESPONSE"

type Member struct {
	Joined bool
	Failed bool
	MemberList *MemberList
	node *core.Node
}

func NewMember(node *core.Node) *Member {
	return &Member{
		Joined: false,
		Failed: false,
		MemberList: NewMemberList(nil),
		node: node,
	}
}

func (m *Member) Start() {
	m.join()
	for !m.Joined {
		time.Sleep(1 * time.Second)
	}
}

func (m *Member) Stop() {
	
}

func (m *Member) handleMessage(msg network.Message) {
	switch msg.Type {
	case GOSSIP:
		fmt.Println("Got a gossip message")
	case JOINREQ:
		encoded, _ := json.Marshal(m.MemberList)
		m.node.SendMessage(string(msg.From), network.Message{
			Type: JOINRES,
			Meta: map[string]string{},
			From: m.node.Id.String(),
			Payload: encoded,
		})
	case JOINRES:
		m.Joined = true
	}
}

func (m *Member) join() {
	m.node.SendMessage("13", network.Message{
		Type: JOINREQ,
		Meta: map[string]string{},
		From: m.node.Id.String(),
		Payload: []byte{},
	})
}

// func (m *Member) gossip() {
// 	m.gossip()
// 	m.node.SendMessage("13", network.Message{
// 		Type: GOSSIP,
// 		Meta: map[string]string{},
// 		From: m.node.Id,
// 		Payload: []byte{},
// 	})
// }

// func (m *Member) heartbeat() {
// 	m.memberList.GetById(m.node.Id).Heartbeat++
// 	m.memberList.GetById(m.node.Id).Timestamp = time.Now()
// }