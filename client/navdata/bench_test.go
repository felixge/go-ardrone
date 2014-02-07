package navdata

import (
	"testing"
)

func BenchmarkDecode(b *testing.B) {
	b.StopTimer()
	buf := fixtureBytes()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_, err := Decode(buf)
		if err != nil {
			panic(err)
		}
	}
}
