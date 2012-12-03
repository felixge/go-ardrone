package navdata

import (
	"bytes"
	"fmt"
)

type Navdata struct {
	Header   NavdataHeader
	Demo     *Demo
	Checksum Checksum
}

type NavdataHeader struct {
	Tag            HeaderTag
	State          State
	SequenceNumber uint32
	VisionFlag     uint32
}

type HeaderTag uint32

const DefaultHeaderTag HeaderTag = 0x55667788

type ErrUnknownHeaderTag struct {
	Expected HeaderTag
	Got      HeaderTag
}

func (this ErrUnknownHeaderTag) Error() string {
	return fmt.Sprintf(
		"navdata: unknown header tag, expected: 0x%x, got: 0x%x",
		this.Expected,
		this.Got,
	)
}

type ErrBadChecksum struct {
	Expected Checksum
	Got      Checksum
}

func (this ErrBadChecksum) Error() string {
	return fmt.Sprintf(
		"navdata: bad checksum, expected: %d, got: %d",
		this.Expected,
		this.Got,
	)
}

// Decode blocks until the next navdata packets becomes available, which
// it then parses and returns.
func Decode(buf []byte) (navdata *Navdata, err error) {
	// readOrPanic() panics, while not expected, should not propagate to the
	// caller, so we return them like regular errors instead.
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	reader := newBinaryReader(bytes.NewReader(buf))
	navdata = &Navdata{}

	reader.readOrPanic(&navdata.Header)

	if navdata.Header.Tag != DefaultHeaderTag {
		err = ErrUnknownHeaderTag{
			Expected: DefaultHeaderTag,
			Got:      navdata.Header.Tag,
		}
		return
	}

	err = readOptions(reader, navdata)

	return
}

func readOptions(reader *binaryReader, navdata *Navdata) error {
	for {
		currentChecksum := reader.Checksum

		header := &OptionHeader{}
		reader.readOrPanic(header)

		// All navdata options can be extended (new values AT THE END) except
		// navdata_demo whose size must be constant across versions
		// -- ARDrone_SDK_2_0/ARDroneLib/Soft/Common/navdata_common.h
		//
		// For this reason we create a new bytes.Buffer for each option, which may
		// or may not end up reading all bytes of it.
		optionData := make([]byte, header.Length-4)
		reader.readOrPanic(optionData)
		optionReader := newBinaryReader(bytes.NewReader(optionData))

		switch header.Tag {
		case DEMO:
			navdata.Demo = new(Demo)
			optionReader.readOrPanic(navdata.Demo)
		case CHECKSUM:
			optionReader.readOrPanic(&navdata.Checksum)
			if (navdata.Checksum != currentChecksum) {
				return ErrBadChecksum{
					Expected: currentChecksum,
					Got: navdata.Checksum,
				}
			}
			return nil
		default:
			//fmt.Printf("Unknown Header: %+v\n", header)
		}
	}
	return nil
}
