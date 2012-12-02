package navdata

import (
	"testing"
)

func BenchmarkDecoder_Decode(b *testing.B) {
	b.StopTimer()
	decoder := NewDecoder(fixture())
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		decoder.Decode()
	}
}

func BenchmarkDecode(b *testing.B) {
	b.StopTimer()
	buf := fixtureBytes()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Decode(buf)
	}
}
