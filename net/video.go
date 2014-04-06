package net

import (
	"encoding/binary"
	"io"
)

const VideoPort = "5555"

func NewPaVEDecoder(r io.Reader) *PaVEDecoder {
	return &PaVEDecoder{r}
}

type PaVEDecoder struct {
	r io.Reader
}

var paVEHeaderSize = binary.Size(PaVEHeader{})

func (d *PaVEDecoder) Decode(frame *PaVEFrame) error {
	err := binary.Read(d.r, binary.LittleEndian, &frame.PaVEHeader)
	if err != nil {
		return err
	}

	headerExtra := int(frame.HeaderSize) - paVEHeaderSize
	frame.PaVEHeaderExtra = grow(frame.PaVEHeaderExtra, headerExtra)
	frame.Payload = grow(frame.Payload, int(frame.PayloadSize))

	_, err = io.ReadFull(d.r, frame.PaVEHeaderExtra)
	if err != nil {
		return err
	}
	_, err = io.ReadFull(d.r, frame.Payload)
	if err != nil {
		return err
	}
	return nil
}

type PaVEHeader struct {
	Signature            [4]byte
	Version              uint8
	VideoCodec           uint8
	HeaderSize           uint16
	PayloadSize          uint32
	EncodedStreamWidth   uint16
	EncodedStreamHeight  uint16
	DisplayWidth         uint16
	DisplayHeight        uint16
	FrameNumber          uint32
	Timestamp            uint32
	TotalChunks          uint8
	ChunkIndex           uint8
	FrameType            uint8
	Control              uint8
	StreamBytePositionLw uint32
	StreamBytePositionUw uint32
	StreamId             uint16
	TotalSlices          uint8
	SliceIndex           uint8
	Header1Size          uint8
	Header2Size          uint8
	Reserved2            [2]byte
	AdvertisedSize       uint32
}

type PaVEFrame struct {
	PaVEHeader
	// PaVEHeaderExtra contains reserved / undocumented payload information
	// found at the end of the PaVEHeader in recent firmware versions.
	// see https://projects.ardrone.org/issues/show/159
	PaVEHeaderExtra []byte
	Payload         []byte
}

func grow(b []byte, length int) []byte {
	if cap(b) < length {
		b = make([]byte, 0, length*2)
	}
	return b[0:length]
}
