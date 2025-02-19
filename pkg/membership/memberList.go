package membership

import (
	"encoding/json"
	"time"
	"github.com/jpriverar/distributed-toolbox/pkg/core"

)

var PRUNE_TIMEOUT float32 = 20
var FAILED_TIMEOUT float32 = 10

type MemberList struct {
	entries []MemberListEntry
}

func NewMemberList(entries []MemberListEntry) *MemberList {
	list := MemberList{
		entries: make([]MemberListEntry, 0),
	}

	for _, entry := range entries {
		list.Add(entry)
	}

	return &list
}

func EqualEntries(entries1, entries2 []MemberListEntry) bool {
	if len(entries1) != len(entries2) {
		return false
	}

	var found bool
	for i := 0; i < len(entries1); i++ {
		found = false
		for j := 0; j < len(entries2); j++ {
			if entries1[i].Equals(entries2[j]) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (l MemberList) Empty() bool {
	return len(l.entries) == 0
}

func (l MemberList) Size() int {
	return len(l.entries)
}

func (l MemberList) Equals(other MemberList) bool {
	if l.Size() != other.Size() {
		return false
	}

	for i := 0; i < l.Size(); i++ {
		if !l.GetEntry(i).Equals(*other.GetEntry(i)) {
			return false
		}
	}
	return true
}

func (l MemberList) String() string {
	res := "[\n"
	for _, entry := range l.entries {
		res += "  " + entry.String() + ",\n"
	}
	res += "\n]"
	return res
}

func (l MemberList) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.entries)
}

func (l *MemberList) UnmarshalJSON(bytes []byte) error {
	err := json.Unmarshal(bytes, &l.entries)
	if err != nil {
		return err
	}
	return nil
}

func (l MemberList) GetMembers() []MemberListEntry {
	currTime := time.Now()
	goodEntries := make([]MemberListEntry, 0)
	for _, entry := range l.entries {
		if !entry.LastHeartbeat.Before(currTime.Add(time.Duration(-FAILED_TIMEOUT) * time.Second)) {
			goodEntries = append(goodEntries, entry)
		}
	}
	return goodEntries
}

func (l *MemberList) GetMember(id core.ID) *MemberListEntry {
	left, right := 0, l.Size()-1
	for left <= right {
		middle := int((left + right) / 2)
		currEntry := l.entries[middle]
		if currEntry.Id.Equals(id) {
			return &currEntry
		} else if currEntry.Id.GreaterThan(id) {
			right = middle - 1
		} else {
			left = middle + 1
		}
	}
	return nil
}

func (l MemberList) GetAllEntries() []MemberListEntry {
	return l.entries
}

func (l *MemberList) GetEntry(index int) *MemberListEntry {
	return &l.entries[index]
}

func (l *MemberList) Add(newEntry MemberListEntry) {
	if l.Empty() {
		l.entries = append(l.entries, newEntry)
		return
	}

	left, right := 0, l.Size()-1
	for left <= right {
		middle := int((left + right) / 2)
		currEntry := l.entries[middle]
		if currEntry.Id.Equals(newEntry.Id) {
			return
		} else if currEntry.Id.GreaterThan(newEntry.Id) {
			right = middle - 1
		} else {
			left = middle + 1
		}
	}

	if left < 0 {
		l.entries = append([]MemberListEntry{newEntry}, l.entries...)
	} else if left >= l.Size() {
		l.entries = append(l.entries, newEntry)
	} else {	
		l.entries = append(l.entries[:left+1], l.entries[left:]...)
		l.entries[left] = newEntry
	}
}

func (l *MemberList) Merge(other *MemberList) { 
	merged := new(MemberList)
	currTime := time.Now()
	i, j := 0, 0
	for i < l.Size() && j < other.Size() {
		if l.entries[i].Id.Equals(other.entries[j].Id) {
			merged.entries = append(merged.entries, l.entries[i])
			merged.entries[merged.Size()-1].LastHeartbeat = currTime

			if l.entries[i].HeartbeatCount >= other.entries[j].HeartbeatCount {
				merged.entries[merged.Size()-1].HeartbeatCount = l.entries[i].HeartbeatCount
			} else {
				merged.entries[merged.Size()-1].HeartbeatCount = other.entries[j].HeartbeatCount
			}
			i++; j++
		} else if l.entries[i].Id.LessThan(other.entries[j].Id) {
			merged.entries = append(merged.entries, l.entries[i])
			i++
		} else {
			merged.entries = append(merged.entries, other.entries[j])
			merged.entries[merged.Size()-1].LastHeartbeat = currTime
			j++
		}
	}

	for i < l.Size() {
		merged.entries = append(merged.entries, l.entries[i])
		i++
	}

	for j < other.Size() {
		merged.entries = append(merged.entries, other.entries[j])
		merged.entries[merged.Size()-1].LastHeartbeat = currTime
		j++
	}

	l.entries = make([]MemberListEntry, merged.Size())
	copy(l.entries, merged.entries)
}

func (l *MemberList) Prune() {
	currTime := time.Now()
	goodEntries := make([]MemberListEntry, 0)
	for _, entry := range l.entries {
		if !entry.LastHeartbeat.Before(currTime.Add(time.Duration(-PRUNE_TIMEOUT) * time.Second)) {
			goodEntries = append(goodEntries, entry) 
		}
	}
	l.entries = goodEntries
}