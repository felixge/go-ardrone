package commands

import (
	"fmt"
)

type Ref struct{
	Fly bool
	Emergency bool
}

func (this *Ref) String(number int) string {
	ref := 0

	if this.Emergency {
		ref = ref | (1 << 8)
	}

	if this.Fly {
		ref = ref | (1 << 9)
	}

	return fmt.Sprintf("AT*REF=%d,%d\r", number, ref)
}
