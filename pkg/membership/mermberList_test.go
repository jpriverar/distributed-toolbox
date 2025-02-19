package membership

import (
	"testing"
	"time"
	"encoding/json"
	"github.com/jpriverar/distributed-toolbox/pkg/core"
)

func TestGetEntries(t *testing.T) {
	ids := []core.ID{core.NewID(),
					 core.NewID(),
					 core.NewID(),
					 core.NewID()}
	dummyTime := time.Now()
	tests := []struct {
		list MemberList
		expectedEntries []MemberListEntry
	}{
		{
			*NewMemberList(
				[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
				},
			),
			[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
			},
		},
		{
			*NewMemberList(
				[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					{ids[1], "192.168.0.24", 2, dummyTime.Add(time.Duration(-FAILED_TIMEOUT*1.1) * time.Second)},
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
					{ids[3], "192.168.0.26", 2, dummyTime.Add(time.Duration(-FAILED_TIMEOUT*1.1) * time.Second)},
				},
			),
			[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
			},
		},
		{
			MemberList{
				entries: []MemberListEntry{
					{ids[0], "192.168.0.23", 1, dummyTime.Add(time.Duration(-(PRUNE_TIMEOUT*1.1)) * time.Second)},
					{ids[1], "192.168.0.24", 2, dummyTime.Add(time.Duration(-(PRUNE_TIMEOUT*1.1)) * time.Second)},
					{ids[2], "192.168.0.25", 2, dummyTime.Add(time.Duration(-(PRUNE_TIMEOUT*1.1)) * time.Second)},
				},
			},
			[]MemberListEntry{},
		},
	}

	for _, tt := range tests {
		if !EqualEntries(tt.list.GetMembers(), tt.expectedEntries) {
			t.Errorf("error getting entries:\nexpected: %v\ngot: %v", tt.expectedEntries, tt.list.GetMembers())
		}
	}
}

func TestAdd(t *testing.T) {
	ids := []core.ID{core.NewID(),
					 core.NewID(),
					 core.NewID(),
					 core.NewID(),
					 core.NewID()}
	tests := []struct {
		list MemberList
		newEntry MemberListEntry
		expected MemberList
	}{
		{
			MemberList{},
			*NewMemberlistEntry(ids[0], "192.168.0.23"),
			MemberList{
				entries: SortedMemberListEntries(
					[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
				}),
			},
		},
		{
			MemberList{
				entries: SortedMemberListEntries([]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[3], "192.168.0.26"),
					*NewMemberlistEntry(ids[4], "192.168.0.27"),
				}),
			},
			*NewMemberlistEntry(ids[2], "192.168.0.25"),
			MemberList{
				entries: SortedMemberListEntries([]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
					*NewMemberlistEntry(ids[3], "192.168.0.26"),
					*NewMemberlistEntry(ids[4], "192.168.0.27"),
				}),
			},
		},
		{
			MemberList{
				entries: SortedMemberListEntries([]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
					*NewMemberlistEntry(ids[3], "192.168.0.26"),
				}),
			},
			*NewMemberlistEntry(ids[4], "192.168.0.27"),
			MemberList{
				entries: SortedMemberListEntries([]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
					*NewMemberlistEntry(ids[3], "192.168.0.26"),
					*NewMemberlistEntry(ids[4], "192.168.0.27"),
				}),
			},
		},
		{
			MemberList{
				entries: SortedMemberListEntries([]MemberListEntry{
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
					*NewMemberlistEntry(ids[3], "192.168.0.26"),
					*NewMemberlistEntry(ids[4], "192.168.0.27"),
				}),
			},
			*NewMemberlistEntry(ids[0], "192.168.0.23"),
			MemberList{
				entries: SortedMemberListEntries([]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
					*NewMemberlistEntry(ids[3], "192.168.0.26"),
					*NewMemberlistEntry(ids[4], "192.168.0.27"),
				}),
			},
		},
	}

	for _, tt := range tests {	
		tt.list.Add(tt.newEntry)
		if !tt.list.Equals(tt.expected) {
			t.Errorf("error adding entry %v to list:\nexpected: %v\ngot: %v", tt.newEntry, tt.expected, tt.list)
		}
	}
}

func TestMerge(t *testing.T) {
	ids := []core.ID{core.NewID(),
					 core.NewID(),
					 core.NewID(),
					 core.NewID()}
	dummyTime := time.Now()
	tests := []struct {
		list1 MemberList
		list2 MemberList
		expected MemberList
	}{
		{
			MemberList{},
			*NewMemberList(
				[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
					*NewMemberlistEntry(ids[3], "192.168.0.26"),
				},
			),
			*NewMemberList(
				[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
					*NewMemberlistEntry(ids[3], "192.168.0.26"),
				},
			),
		},
		{
			*NewMemberList(
				[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
				},
			),
			*NewMemberList(
				[]MemberListEntry{
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[3], "192.168.0.26"),
				},
			),
			*NewMemberList(
				[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
					*NewMemberlistEntry(ids[3], "192.168.0.26"),
				},
			),
		},
		{
			*NewMemberList(
				[]MemberListEntry{
					{ids[0], "192.168.0.23", 0, dummyTime},
					{ids[0], "192.168.0.24", 1, dummyTime},
					{ids[0], "192.168.0.25", 2, dummyTime},
				},
			),
			*NewMemberList(
				[]MemberListEntry{
					{ids[0], "192.168.0.23", 1, dummyTime},
					{ids[0], "192.168.0.24", 2, dummyTime},
					{ids[0], "192.168.0.25", 0, dummyTime},
				},
			),
			*NewMemberList(
				[]MemberListEntry{
					{ids[0], "192.168.0.23", 1, dummyTime},
					{ids[0], "192.168.0.24", 2, dummyTime},
					{ids[0], "192.168.0.25", 2, dummyTime},
				},
			),
		},
	}

	for _, tt := range tests {	
		tt.list1.Merge(&tt.list2)
		if !tt.list1.Equals(tt.expected) {
			t.Errorf("error merging lists\nexpected: %v\ngot: %v", tt.expected, tt.list1)
		}
	}
}

func TestPrune(t *testing.T) {
	ids := []core.ID{core.NewID(),
					 core.NewID(),
					 core.NewID(),
					 core.NewID()}
	dummyTime := time.Now()
	tests := []struct {
		list MemberList
		expected MemberList
	}{
		{
			*NewMemberList(
				[]MemberListEntry{
					{ids[0], "192.168.0.23", 1, dummyTime.Add(time.Duration(-(PRUNE_TIMEOUT*1.1)) * time.Second)},
					{ids[1], "192.168.0.24", 2, dummyTime.Add(time.Duration(-(PRUNE_TIMEOUT*1.1)) * time.Second)},
					{ids[2], "192.168.0.25", 2, dummyTime.Add(time.Duration(-(PRUNE_TIMEOUT*1.1)) * time.Second)},
				},
			),
			MemberList{},
		},
		{
			*NewMemberList(
				[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
				},
			),
			*NewMemberList(
				[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
				},
			),
		},
		{
			*NewMemberList(
				[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					{ids[1], "192.168.0.24", 0, dummyTime.Add(time.Duration(-(PRUNE_TIMEOUT*1.1)) * time.Second)},
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
					{ids[3], "192.168.0.26", 0, dummyTime.Add(time.Duration(-(PRUNE_TIMEOUT*1.1)) * time.Second)},
				},
			),
			*NewMemberList(
				[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
				},
			),
		},
	}

	for _, tt := range tests {
		tt.list.Prune()
		if !tt.list.Equals(tt.expected) {
			t.Errorf("error pruning list\nexpected: %v\ngot: %v", tt.expected, tt.list)
		}
	}
}

func TestEncoding(t *testing.T) {
	ids := []core.ID{core.NewID(),
					 core.NewID(),
					 core.NewID(),
					 core.NewID(),
					 core.NewID()}
	tests := []struct {
		list MemberList
	}{
		{
			*NewMemberList(
				[]MemberListEntry{
					*NewMemberlistEntry(ids[0], "192.168.0.23"),
					*NewMemberlistEntry(ids[1], "192.168.0.24"),
					*NewMemberlistEntry(ids[2], "192.168.0.25"),
					*NewMemberlistEntry(ids[3], "192.168.0.26"),
					*NewMemberlistEntry(ids[4], "192.168.0.27"),
				},
			),
		},
		{
			MemberList{},
		},
	}

	for _, tt := range tests {
		encoded, err := json.Marshal(tt.list)
		if err != nil {
			t.Errorf("error encoding list: %v", err)
			t.Log(string(encoded))
		}

		decodedList := NewMemberList(nil)
		err = json.Unmarshal(encoded, decodedList)
		if err != nil {
			t.Errorf("error encoding list: %v", err)
			t.Log(decodedList)
		}

		if !tt.list.Equals(*decodedList) {
			t.Errorf("error encoding list\nexpected: %v\ngot: %v", tt.list, decodedList)
		}
	}
}