package commands

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Pcmd struct {
	Pitch    float64
	Roll     float64
	Yaw      float64
	Vertical float64
}

func (pcmd *Pcmd) String(number int) string {
	flags := 0
	args := []string{
		convenientFloatToString(pcmd.Roll),
		convenientFloatToString(-pcmd.Pitch),
		convenientFloatToString(pcmd.Vertical),
		convenientFloatToString(pcmd.Yaw),
	}

	for _, val := range args {
		if val != "0" {
			flags = 1
		}
	}

	return fmt.Sprintf(
		"AT*PCMD=%d,%d,%s,%s,%s,%s\r",
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
func convenientFloatToString(floatVal float64) string {
	if floatVal == -0 {
		floatVal = 0
	}

	result, err := floatToString(floatVal)

	if err != nil {
		return "0"
	}

	return result
}

func floatToString(floatVal float64) (result string, err error) {
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
