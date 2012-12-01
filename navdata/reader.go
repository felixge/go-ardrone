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

// Parse provides a simple interface to Reader.ReadNavdata(). However, it is
// ~100x slower for consecutive parsing then using the Reader interface
// directly (see bench_test.go).
func Parse(buf []byte) (navdata *Navdata, err error) {
	reader := NewReader(bytes.NewReader(buf))
	navdata, err = reader.ReadNavdata()
	return
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

	this.r.readOrPanic(&navdata.Header)

	if navdata.Header.Tag != DefaultHeaderTag {
		err = ErrUnknownHeaderTag{
			Expected: DefaultHeaderTag,
			Got:      navdata.Header.Tag,
		}
		return
	}

	err = this.readOptions(navdata)

	return
}

func (this *Reader) readOptions(navdata *Navdata) error {
	for {
		currentChecksum := this.r.Checksum

		header := &OptionHeader{}
		this.r.readOrPanic(header)

		// All navdata options can be extended (new values AT THE END) except
		// navdata_demo whose size must be constant across versions
		// -- ARDrone_SDK_2_0/ARDroneLib/Soft/Common/navdata_common.h
		//
		// For this reason we create a new bytes.Buffer for each option, which may
		// or may not end up reading all bytes of it.
		optionData := make([]byte, header.Length-4)
		this.r.readOrPanic(optionData)
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
