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
	FLY_OK FlyState = iota
	FLY_LOST_ALTITUDE
	FLY_LOST_ALTITUDE_GO_DOWN
	FLY_ALTITUDE_OUT_ZONE
	FLY_COMBINED_YAW
	FLY_BRAKE
	FLY_NO_VISION
)

var flyStateStrings = map[FlyState]string{
	FLY_OK:                    "FLY_OK",
	FLY_LOST_ALTITUDE:         "FLY_LOST_ALTITUDE",
	FLY_LOST_ALTITUDE_GO_DOWN: "FLY_LOST_ALTITUDE_GO_DOWN",
	FLY_ALTITUDE_OUT_ZONE:     "FLY_ALTITUDE_OUT_ZONE",
	FLY_COMBINED_YAW:          "FLY_COMBINED_YAW",
	FLY_BRAKE:                 "FLY_BRAKE",
	FLY_NO_VISION:             "FLY_NO_VISION",
}

// String() returns the name of the given FlyState.
func (this FlyState) String() string {
	return flyStateStrings[this]
}

const (
	CONTROL_DEFAULT ControlState = iota
	CONTROL_INIT
	CONTROL_LANDED
	CONTROL_FLYING
	CONTROL_HOVERING
	CONTROL_TEST
	CONTROL_TAKOFF
	CONTROL_GOTOFIX
	CONTROL_LANDING
	CONTROL_LOOPING
)

var controlStrings = map[ControlState]string{
	CONTROL_DEFAULT:  "CONTROL_DEFAULT",
	CONTROL_INIT:     "CONTROL_INIT",
	CONTROL_LANDED:   "CONTROL_LANDED",
	CONTROL_FLYING:   "CONTROL_FLYING",
	CONTROL_HOVERING: "CONTROL_HOVERING",
	CONTROL_TEST:     "CONTROL_TEST",
	CONTROL_TAKOFF:   "CONTROL_TAKOFF",
	CONTROL_GOTOFIX:  "CONTROL_GOTOFIX",
	CONTROL_LANDING:  "CONTROL_LANDING",
	CONTROL_LOOPING:  "CONTROL_LOOPING",
}

// String() returns the name of the given ControlState.
func (this ControlState) String() string {
	return controlStrings[this]
}
