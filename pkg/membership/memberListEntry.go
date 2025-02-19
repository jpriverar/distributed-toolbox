package membership

import (
	"fmt"
	"time"
	"sort"
	"github.com/jpriverar/distributed-toolbox/pkg/core"

)

type MemberListEntry struct {
	Id core.ID
	Addr string
	HeartbeatCount uint64
	LastHeartbeat time.Time
}

func NewMemberlistEntry(id core.ID, addr string) *MemberListEntry{
	return &MemberListEntry{
		Id: id,
		Addr: addr,
		HeartbeatCount: 0,
		LastHeartbeat: time.Now(),
	}
}

func (entry *MemberListEntry) Heartbeat() {
	entry.HeartbeatCount++
	entry.LastHeartbeat = time.Now()
}

func (e MemberListEntry) Equals(other MemberListEntry) bool {
	return e.Id.Equals(other.Id) &&
		   e.Addr == other.Addr &&
		   e.HeartbeatCount == other.HeartbeatCount
}

func (e MemberListEntry) String() string {
	return fmt.Sprintf("[%s %s %d]", e.Id, e.Addr, e.HeartbeatCount)
}

func SortMemberListEntries(entries []MemberListEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Id.LessThan(entries[j].Id)
	})
}

func SortedMemberListEntries(entries []MemberListEntry) []MemberListEntry {
	sorted := make([]MemberListEntry, len(entries))
	copy(sorted, entries)

	SortMemberListEntries(sorted)
	return sorted
}