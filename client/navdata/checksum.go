package navdata

type Checksum uint32

func (c *Checksum) Add(buf []byte) {
	for i := 0; i < len(buf); i++ {
		*c += Checksum(buf[i])
	}
}

func (c *Checksum) Sub(buf []byte) {
	for i := 0; i < len(buf); i++ {
		*c -= Checksum(buf[i])
	}
}
