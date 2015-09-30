// Misc features
package misc

import (
	"strings"

	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"
)

type Misc struct {
	tg *telegram.Telegram
}

func Setup(t *telegram.Telegram, config map[string]interface{}, modules map[string]bool, cmds *types.CommandMap) types.Command {
	if val, ok := modules["help"]; !ok || val {
		misc := &Misc{tg: t}

		// Echo
		(*cmds)["echo"] = types.Command {
			Name: "echo",
			Args: "<text>",
			ArgNum: 1,
			Desc: "Echo <text>",
			Processor: misc,
		}

		// Debug parse
		(*cmds)["parse"] = types.Command {
			Name: "parse",
			Args: "arguments",
			ArgNum: -1,
			Desc: "Parse argument list [debug]",
			Debug: true,
			Processor: misc,
		}
	}

	return types.Command{}
}

func (this *Misc) Command(name string, msg telegram.TObject, args []string) {
	switch name {
		case "echo":
			this.tg.SendMessage(args[0], msg.ChatId())
		case "parse":
			this.tg.ReplyToMessage(msg.MessageId(), strings.Join(args, "\n"), msg.ChatId())
	}
}

func (this *Misc) Default(name string, msg telegram.TObject) {
}
