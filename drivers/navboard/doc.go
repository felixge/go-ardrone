// Package navboard implements the tty navboard interface.
package navboard

func NewDriver(ttyPath string) *Driver {
	return &Driver{}
}

type Driver struct {
}

func (d *Driver) Read(data *Data) (int, error) {
	return 0, nil
}

// Data holds the navboard data as read from the tty file. Based on
// https://github.com/RoboticaTUDelft/paparazzi/blob/minor1/sw/airborne/boards/ardrone/navdata.h
// but with some adjustements for values that seem to be signed rather than
// unsigned.
type Data struct {
	Seq uint16

	// Accelerometers
	Ax uint16
	Ay uint16
	Az uint16

	// Gyroscopes
	Gx int16
	Gy int16
	Gz int16

	TemperatureAcc  uint16
	TemperatureGyro uint16

	Ultrasound int16

	UsDebutEcho       uint16
	UsFinEcho         uint16
	UsAssociationEcho uint16
	UsDistanceEcho    uint16

	UsCurveTime  uint16
	UsCurveValue uint16
	UsCurveRef   uint16

	NbEcho uint16

	SumEcho  uint32
	Gradient int16

	FlagEchoIni uint16

	Pressure            int32
	TemperaturePressure int16

	Mx int16
	My int16
	Mz int16

	Checksum uint16
}
