package navdata

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"runtime"
)

var _, filename, _, _ = runtime.Caller(0)
var fixturePath = path.Dir(filename) + "/decoder_fixture.bin"

func ExampleDecoder_ReadNavdata_Ok() {
	reader := NewDecoder(fixture())
	navdata, _ := reader.ReadNavdata()
	json, _ := json.MarshalIndent(navdata, "", "\t")
	fmt.Print(string(json))

	// Output:
	// {
	// 	"Header": {
	// 		"Tag": 1432778632,
	// 		"State": 1333788880,
	// 		"SequenceNumber": 300711,
	// 		"VisionFlag": 1
	// 	},
	// 	"Demo": {
	// 		"FlyState": 0,
	// 		"ControlState": 2,
	// 		"Battery": 50,
	// 		"Altitude": 0
	// 	},
	//	"Checksum": 36358
	// }
}

func ExampleDecoder_ReadNavdata_ErrUnknownHeaderTag() {
	badHeader := make([]byte, 16)
	badHeader[0] = 0x01
	badHeader[1] = 0x02
	badHeader[2] = 0x03
	badHeader[3] = 0x04

	reader := NewDecoder(bufio.NewReader(bytes.NewReader(badHeader)))
	_, err := reader.ReadNavdata()
	fmt.Print(err.Error())

	// Output:
	// navdata: unknown header tag, expected: 0x55667788, got: 0x4030201
}

func ExampleDecoder_ReadNavdata_ErrUnexpectedEof() {
	reader := NewDecoder(bufio.NewReader(bytes.NewReader([]byte{0x00})))
	_, err := reader.ReadNavdata()
	fmt.Print(err.Error())

	// Output:
	// unexpected EOF
}

func ExampleDecoder_ReadNavdata_ErrBadChecksum() {
	data := fixtureBytes()
	// corrupt a byte
	data[20] = data[20] + 1;

	_, err := Parse(data)

	fmt.Print(err.Error())

	// Output:
	// navdata: bad checksum, expected: 36359, got: 36358
}

func fixture() io.Reader {
	return bytes.NewReader(fixtureBytes())
}

func fixtureBytes() []byte {
	data, err := ioutil.ReadFile(fixturePath)
	if err != nil {
		panic(err)
	}

	return data
}
