package commands

type Sequence struct {
	number   int
	commands []command
}

type command interface {
	String(number int) string
}

func (this *Sequence) Add(command command) {
	this.commands = append(this.commands, command)
}

func (this *Sequence) ReadMessage() string {
	var message string
	for _, command := range this.commands {
		this.number++
		message += command.String(this.number)
	}

	this.commands = nil

	return message
}
