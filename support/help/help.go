package help

import (
	"fmt"
	"strings"

	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"
)

type Help struct {
	tg *telegram.Telegram
	cmds *types.CommandMap
}

func Setup(t *telegram.Telegram, config map[string]interface{}, modules map[string]bool, cmds *types.CommandMap) types.Command {
	if val, ok := modules["help"]; !ok || val {
		help := &Help {
			tg: t,
			cmds: cmds,
		}

		(*cmds)["help"] = types.Command {
			Name: "help",
			Desc: "Show help information for this bot",
			ArgNum: 0,
			Processor: help,
		}

		(*cmds)["father"] = types.Command {
			Name: "father",
			Desc: "For @BotFather",
			ArgNum: 0,
			Debug: true,
			Processor: help,
		}
	}

	return types.Command{}
}

func (this *Help) Command(name string, msg telegram.TObject, args []string) {
	if name == "help" {
		if !msg.Chat().IsGroup() {
			str := ""
			for _, v := range (*this.cmds) {
				// Skip debug functions
				if v.Debug {
					continue
				}

				str += fmt.Sprintf(
					"/%s %s\n%s\n\n",
					v.Name, v.Args, v.Desc)
			}
			this.tg.ReplyToMessage(msg.MessageId(), str, msg.ChatId())
		} else {
			this.tg.ReplyToMessage(msg.MessageId(), "Help only available in private chats.", msg.ChatId())
		}
	} else if name == "father" {
		if !msg.Chat().IsGroup() {
			str := ""
			for _, v := range (*this.cmds) {
				if v.Debug { continue }

				str += fmt.Sprintf(
					"%s - %s %s\n",
					v.Name, v.Args, strings.Split(v.Desc, "\n")[0])
			}
			this.tg.ReplyToMessage(msg.MessageId(), str, msg.ChatId())
		}
	}
}

func (this *Help) Default(name string, msg telegram.TObject, state *map[string]interface{}) {

}
