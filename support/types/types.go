package types

import (
	telegram "github.com/PeterCxy/gotelegram"
)

type Command struct {
	Name string
	Desc string
	Args string
	ArgNum int
	Debug bool
	Processor CommandProcessor
}

type CommandProcessor interface {
	Command(name string, msg telegram.TObject, args []string)
	Default(name string, msg telegram.TObject)
}

type CommandMap map[string]Command
