package core

import (
	"bytes"
	"github.com/google/uuid"
)

type ID uuid.UUID

func (id ID) Equals(other ID) bool {
	return bytes.Equal(id[:], other[:])
}

func (id ID) GreaterThan(other ID) bool {
	return bytes.Compare(id[:], other[:]) == 1
}

func (id ID) LessThan(other ID) bool {
	return bytes.Compare(id[:], other[:]) == -1
}

func (id ID) String() string {
	return uuid.UUID(id).String()
}

func NewID() ID {
	return ID(uuid.New())
}