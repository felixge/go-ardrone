package ardrone

import (
	"fmt"
)

type controlCommand interface{
	String() string
}

type ControlMessage struct{
	commands []controlCommand
}

func (this *ControlMessage) Add(command controlCommand) {
	this.commands = append(this.commands, command)
}

func (this *ControlMessage) String() string {
	var result string
	for _, command := range this.commands {
		result = fmt.Sprintf("%s%s", result, command.String())
	}

	return result
}
