package navdata

import (
	"bytes"
	"fmt"
	"io"
)

type Reader struct {
	r *binaryReader
}

type Navdata struct {
	Header   NavdataHeader
	Demo     *Demo
	Checksum Checksum
}

type NavdataHeader struct {
	Tag            HeaderTag
	State          uint32
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
		"navdata: Unknown header tag, expected: 0x%x, got: 0x%x",
		this.Expected,
		this.Got,
	)
}

func NewReader(r io.Reader) *Reader {
	return &Reader{r: newBinaryReader(r)}
}

// ReadNavdata blocks until the next navdata packets becomes available, which
// it then parses and returns.
func (this *Reader) ReadNavdata() (navdata *Navdata, err error) {
	// readOrPanic() panics, while not expected, should not propagate to the
	// caller, so we return them like regular errors instead.
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()

	navdata = &Navdata{}

	this.r.ReadOrPanic(&navdata.Header)

	if navdata.Header.Tag != DefaultHeaderTag {
		err = ErrUnknownHeaderTag{
			Expected: DefaultHeaderTag,
			Got:      navdata.Header.Tag,
		}
		return
	}

	this.readOptions(navdata)

	return navdata, nil
}

func (this *Reader) readOptions(navdata *Navdata) {
	for {
		header := &OptionHeader{}
		this.r.ReadOrPanic(header)

		// All navdata options can be extended (new values AT THE END) except
		// navdata_demo whose size must be constant across versions
		// -- ARDrone_SDK_2_0/ARDroneLib/Soft/Common/navdata_common.h
		//
		// For this reason we create a new bytes.Buffer for each option, which may
		// or may not end up reading all bytes of it.
		optionData := make([]byte, header.Length-4)
		this.r.ReadOrPanic(optionData)
		optionReader := newBinaryReader(bytes.NewReader(optionData))

		switch header.Tag {
		case DEMO:
			navdata.Demo = new(Demo)
			optionReader.ReadOrPanic(navdata.Demo)
		case CHECKSUM:
			optionReader.ReadOrPanic(&navdata.Checksum)
			return
		default:
			//fmt.Printf("Unknown Header: %+v\n", header)
		}
	}
}
