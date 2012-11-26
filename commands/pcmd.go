package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Pcmd struct {
	LeftRight float32
	FrontBack float32
	UpDown    float32
	ClockSpin float32
}

func (this *Pcmd) String(number int) string {
	flags := 0
	args := []string{
		convenientFloatToString(this.LeftRight),
		convenientFloatToString(-this.FrontBack),
		convenientFloatToString(this.UpDown),
		convenientFloatToString(this.ClockSpin),
	}

	for _, val := range args {
		if val != "0" {
			flags = 1
		}
	}

	return fmt.Sprintf(
		"AT*CMD=%d,%d,%s,%s,%s,%s\r",
		number,
		flags,
		args[0],
		args[1],
		args[2],
		args[3],
	)
}

// convenientFloatToString returns the same result as floatToString, except
// for two cases.
//
// a) If there was an error, "0" is returned
// b) If -0 was passed in, "0" is returned
func convenientFloatToString(floatVal float32) string {
	if floatVal == -0 {
		floatVal = 0
	}

	result, err := floatToString(floatVal)

	if err != nil {
		return "0"
	}

	return result
}

func floatToString(floatVal float32) (result string, err error) {
	buf := new(bytes.Buffer)

	err = binary.Write(buf, binary.BigEndian, floatVal)
	if err != nil {
		return
	}

	var intVal int32
	err = binary.Read(buf, binary.BigEndian, &intVal)
	if err != nil {
		return
	}

	return fmt.Sprintf("%d", intVal), nil
}
