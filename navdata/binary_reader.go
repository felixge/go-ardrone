package navdata

import (
	"encoding/binary"
	"io"
)

type binaryReader struct {
	r io.Reader
}

func newBinaryReader(r io.Reader) *binaryReader {
	return &binaryReader{r: r}
}

// readOrPanic is a helper function that triggers a panic when binary.Read()
// returns an error (EOF, ErrUnexpectedEOF, etc.). This allows us to unwind the
// stack in these cases without using `if err != nil` checks everywhere.
//
// see: ReadNavdata(), which stops the panic() from propagating to the user of
// this library.
func readOrPanic(r io.Reader, value interface{}) {
	if err := binary.Read(r, binary.LittleEndian, value); err != nil {
		panic(err)
	}
}

func (this *binaryReader) ReadOrPanic(value interface{}) {
	readOrPanic(this.r, value)
}
