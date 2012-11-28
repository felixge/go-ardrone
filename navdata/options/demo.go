package options

import (
	"bytes"
	"encoding/binary"
)

type Demo struct {
	FlyState     FlyState
	ControlState ControlState
}

type FlyState uint16
type ControlState uint16

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

func (this ControlState) String() string {
	return controlStrings[this]
}

func ParseDemo(buf *bytes.Buffer) *Demo {
	demo := &Demo{}

	binary.Read(buf, binary.LittleEndian, &demo.FlyState)
	binary.Read(buf, binary.LittleEndian, &demo.ControlState)

	return demo
}
