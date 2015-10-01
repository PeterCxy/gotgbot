// Misc features
package misc

import (
	"strings"

	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"
	"github.com/PeterCxy/gotgbot/support/utils"
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

		// Remind
		(*cmds)["remind"] = types.Command {
			Name: "remind",
			ArgNum: 0,
			Desc: "Remind you of something after a period of time",
			Processor: misc,
		}

		// Cancel
		(*cmds)["cancel"] = types.Command {
			Name: "cancel",
			ArgNum: 0,
			Desc: "Cancel the current session with this bot",
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
		case "cancel":
			if utils.HasGrabber(msg.FromId(), msg.ChatId()) {
				utils.ReleaseGrabber(msg.FromId(), msg.ChatId())
				this.tg.SendMessage("Current session cancelled", msg.ChatId())
			}
		case "remind":
			this.tg.ReplyToMessage(msg.MessageId(), "What do you want me to remind you of?", msg.ChatId())
			utils.SetGrabber(types.Grabber {
				Name: "remind",
				Uid: msg.FromId(),
				Chat: msg.ChatId(),
				Processor: this,
			})
	}
}

func (this *Misc) Default(name string, msg telegram.TObject, state *map[string]interface{}) {
	if name == "remind" {
		if (*state)["remind"] == nil {
			(*state)["remind"] = msg["text"].(string)
		} else {
			// TODO finish implementing this stuff
			this.tg.SendMessage((*state)["remind"].(string), msg.ChatId())
		}
	}
}
