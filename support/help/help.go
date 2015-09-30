package help

import (
	"fmt"

	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"
)

type Help struct {
	tg *telegram.Telegram
	cmds *types.CommandMap
}

func Setup(t *telegram.Telegram, config map[string]interface{}, modules map[string]bool, cmds *types.CommandMap) types.Command {
	if val, ok := modules["help"]; !ok || val {
		(*cmds)["help"] = types.Command {
			Name: "help",
			Desc: "Show help information for this bot",
			ArgNum: 0,
			Processor: &Help {
				tg: t,
				cmds: cmds,
			},
		}
	}

	return types.Command{}
}

func (this *Help) Command(name string, msg telegram.TObject, args []string) {
	if name == "help" {
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
	}
}

func (this *Help) Default(name string, msg telegram.TObject) {

}
