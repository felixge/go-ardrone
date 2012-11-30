package navdata

import (
	"testing"
)

func BenchmarkHello(b *testing.B) {
	b.StopTimer()
	reader := NewReader(fixture())
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		reader.ReadNavdata()
	}
}
