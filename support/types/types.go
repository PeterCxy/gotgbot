package types

import (
	telegram "github.com/PeterCxy/gotelegram"
)

type Command struct {
	Name      string
	Desc      string
	Args      string
	ArgNum    int
	Debug     bool
	Processor CommandProcessor
}

type Grabber struct {
	Name      string
	Uid       int64
	Chat      int64
	Processor CommandProcessor
}

type CommandProcessor interface {
	Command(name string, msg telegram.TObject, args []string)
	Default(name string, msg telegram.TObject, state *map[string]interface{})
}

type CommandMap map[string]Command
