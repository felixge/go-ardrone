package navdata

import (
	"io"
)

type Checksum uint32

type checkSumReader struct {
	r io.Reader
}

func newChecksumReader(r io.Reader) *checkSumReader {
	return &checkSumReader{r: r}
}
