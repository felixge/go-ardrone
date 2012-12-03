package navdata

import (
	"bytes"
	"testing"
)

// This is ~100x slower than my previous version using a Decoder type instance
// was, which is odd. Will need to figure this out.
func BenchmarkDecode(b *testing.B) {
	b.StopTimer()
	buf := fixtureBytes()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := Decode2(bytes.NewReader(buf))
		if err != nil {
			panic(err)
		}
	}
}
