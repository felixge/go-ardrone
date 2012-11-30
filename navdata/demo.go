package navdata

type Demo struct {
	FlyState     FlyState
	ControlState ControlState
	Battery      uint32 // percentage value, between 0 - 100
	Orientation
	Altitude uint32
	Speed
}

type FlyState uint16
type ControlState uint16

// Orientation values are given in milli-degree.
type Orientation struct {
	FrontBack float32
	LeftRight float32
	ClockSpin float32
}

// Speed stores the estimated linear velocity of the drone. I am not sure yet
// what unit these values are in.
// see: ARDrone_SDK_2_0/ARDroneLib/Soft/Common/navdata_common.h
type Speed struct {
	LeftRight float32
	FrontBack float32
	UpDown    float32
}

const (
	OK FlyState = iota
	LOST_ALTITUDE
	LOST_ALTITUDE_GO_DOWN
	ALTITUDE_OUT_ZONE
	COMBINED_YAW
	BRAKE
	NO_VISION
)

var flyStateStrings = map[FlyState]string{
	OK:                    "OK",
	LOST_ALTITUDE:         "LOST_ALTITUDE",
	LOST_ALTITUDE_GO_DOWN: "LOST_ALTITUDE_GO_DOWN",
	ALTITUDE_OUT_ZONE:     "ALTITUDE_OUT_ZONE",
	COMBINED_YAW:          "COMBINED_YAW",
	BRAKE:                 "BRAKE",
	NO_VISION:             "NO_VISION",
}

// String() returns the name of the given FlyState.
func (this FlyState) String() string {
	return flyStateStrings[this]
}

const (
	DEFAULT ControlState = iota
	INIT
	LANDED
	FLYING
	HOVERING
	TEST
	TAKOFF
	GOTOFIX
	LANDING
	LOOPING
)

var controlStrings = map[ControlState]string{
	DEFAULT:  "DEFAULT",
	INIT:     "INIT",
	LANDED:   "LANDED",
	FLYING:   "FLYING",
	HOVERING: "HOVERING",
	TEST:     "TEST",
	TAKOFF:   "TAKOFF",
	GOTOFIX:  "GOTOFIX",
	LANDING:  "LANDING",
	LOOPING:  "LOOPING",
}

// String() returns the name of the given ControlState.
func (this ControlState) String() string {
	return controlStrings[this]
}
