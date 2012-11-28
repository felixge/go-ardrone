package options

import (
	"bytes"
	"encoding/binary"
)

type Checksum struct {
	Value uint32
}

func ParseChecksum(buf *bytes.Buffer) *Checksum {
	checksum := &Checksum{}

	binary.Read(buf, binary.LittleEndian, &checksum.Value)

	return checksum
}
