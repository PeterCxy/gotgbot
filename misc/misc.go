// Misc features
package misc

import (
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
	}

	return types.Command{}
}

func (this *Misc) Command(name string, msg telegram.TObject, args []string) {
	if name == "echo" {
		this.tg.SendMessage(args[0], msg.ChatId())
	}
}

func (this *Misc) Default(name string, msg telegram.TObject) {
}
