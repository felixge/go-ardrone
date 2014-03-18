package hardware

import (
	"testing"
)

func TestNavboard(t *testing.T) {
	driver, err := OpenNavboard("/dev/ttyO1")
	if err != nil {
		t.Fatal(err)
	}
	var data Navdata
	for i := 0; i < 100; i++ {
		err = driver.Read(&data)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%+v\n", data)
	}
}
