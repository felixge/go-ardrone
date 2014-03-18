package hardware

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// NavboardTTY is the default location of the navboard tty file.
const NavboardTTY = "/dev/ttyO1"

const packetSize = 0x3a

var packetHeader = []byte{packetSize, 0x00}

// OpenNavboard opens the navboard tty file at the given location and returns a
// Navboard struct on success.
func OpenNavboard(ttyPath string) (*Navboard, error) {
	n := &Navboard{ttyPath: ttyPath, buf: &bytes.Buffer{}}
	return n, n.open()
}

// Navboard provides access to the navboard. Must be used from a single
// goroutine.
type Navboard struct {
	ttyPath string
	file    *os.File
	reader  *bufio.Reader
	buf     *bytes.Buffer
}

// open opens the navboard tty file.
func (b *Navboard) open() error {
	file, err := os.OpenFile(b.ttyPath, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	b.file = file
	b.reader = bufio.NewReader(file)
	return nil
}

// Read reads the next packet of navdata into the given data struct.
func (b *Navboard) Read(data *Navdata) error {
	// The loop below is used to sync with the packet stream received from the
	// navboard. This has to be done as the first bytes we read will often be in
	// the middle of a packet. A better approach would be to flush any buffered
	// data from the tty before the first read, but previous attempts of doing
	// this with ioctl and TCFLSH have not been successful. So for now we use the
	// known packetHeader as a marker to find in the stream. This can fail if
	// this value occurs in the middle of a packet and/or if Parrot adds
	// additional data to the packet in the future. Improvements are welcome!
	// Note: I may have had the parrot firmware running while experimenting with
	// TCFLSH, so this approach should be retried to make sure.
	i := 0
	for {
		c, err := b.reader.ReadByte()
		if err != nil {
			return err
		}
		if c == packetHeader[i] {
			i++
			if i == len(packetHeader) {
				break
			}
		} else {
			i = 0
		}
	}

	if _, err := io.CopyN(b.buf, b.reader, packetSize); err != nil {
		return err
	}
	sum := uint16(0)
	buf := b.buf.Bytes()
	for i := 0; i < len(buf)-2; i += 2 {
		sum += uint16(buf[i]) + (uint16(buf[i+1]) << 8)
	}
	if err := binary.Read(b.buf, binary.LittleEndian, data); err != nil {
		return err
	}
	if sum != data.Checksum {
		return fmt.Errorf("Bad checksum. expected=%d got=%d", data.Checksum, sum)
	}
	return nil
}

// Close closes the underlaying tty file.
func (d *Navboard) Close() error {
	err := d.file.Close()
	d.file = nil
	d.reader = nil
	return err
}

// Navdata holds the navboard data as read from the tty file. Based on
// https://github.com/paparazzi/paparazzi/blob/master/sw/airborne/boards/ardrone/navdata.h
// but with some adjustements for values that seem to be signed rather than
// unsigned.
type Navdata struct {
	Seq uint16

	// Accelerometers
	Ax uint16
	Ay uint16
	Az uint16

	// Gyroscopes
	Gx int16
	Gy int16
	Gz int16

	TemperatureAcc  uint16
	TemperatureGyro uint16

	Ultrasound uint16

	UsDebutEcho       uint16
	UsFinEcho         uint16
	UsAssociationEcho uint16
	UsDistanceEcho    uint16

	UsCurveTime  uint16
	UsCurveValue uint16
	UsCurveRef   uint16

	NbEcho uint16

	SumEcho  uint32
	Gradient int16

	FlagEchoIni uint16

	Pressure            int32
	TemperaturePressure int16

	Mx int16
	My int16
	Mz int16

	Checksum uint16
}
