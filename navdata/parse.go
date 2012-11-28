package navdata

import (
	"bytes"
	"encoding/binary"
	"github.com/felixge/ardrone/navdata/options"
)

type Data struct {
	header         uint32
	sequenceNumber uint32
	state          uint32
	visionFlag     uint32
	Demo           *options.Demo
	Checksum       *options.Checksum
}

func Parse(raw []byte) (data Data, err error) {
	buf := bytes.NewBuffer(raw)
	binary.Read(buf, binary.LittleEndian, &data.header)

	if data.header != 0x55667788 {
		// @TODO: set err
		return
	}

	// @TODO: Needs error handling. Possibly using panic + defer + recover
	binary.Read(buf, binary.LittleEndian, &data.state)
	binary.Read(buf, binary.LittleEndian, &data.sequenceNumber)
	binary.Read(buf, binary.LittleEndian, &data.visionFlag)

	err = parseOptions(&data, buf)

	return
}

func parseOptions(data *Data, buf *bytes.Buffer) error {
	for buf.Len() > 0 {
		option := new(options.OptionBlock)
		binary.Read(buf, binary.LittleEndian, &option.Id)
		binary.Read(buf, binary.LittleEndian, &option.Length)

		optionData := make([]byte, option.Length-4)
		binary.Read(buf, binary.LittleEndian, optionData)
		optionBuf := bytes.NewBuffer(optionData)

		switch option.Id {
		case options.DEMO:
			data.Demo = options.ParseDemo(optionBuf)
		case options.CHECKSUM:
			data.Checksum = options.ParseChecksum(optionBuf)
		}
	}

	return nil
}
