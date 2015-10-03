// Scripting support
package script

import (
	"strings"
	"time"

	"github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gobf"
	"github.com/PeterCxy/gotgbot/support/types"
)

type Script struct {
	tg *telegram.Telegram
	timeout int
}

func Setup(t *telegram.Telegram, config map[string]interface{}, modules map[string]bool, cmds *types.CommandMap) types.Command {
	if val, ok := modules["script"]; !ok || val {
		script := &Script{tg: t}
		if val, ok := config["script_timeout"]; ok {
			script.timeout = int(val.(float64))
		}

		// Brainfuck
		(*cmds)["brainfuck"] = types.Command {
			Name: "brainfuck",
			Args: "<code>",
			ArgNum: 1,
			Desc: "Execute branfuck code",
			Processor: script,
		}
	}

	return types.Command{}
}

func (this *Script) Command(name string, msg telegram.TObject, args []string) {
	if name == "brainfuck" {
		code := strings.Trim(args[0], " \n")

		if code == "" {
			this.tg.ReplyToMessage(msg.MessageId(), "Code is empty.", msg.ChatId())
		} else {
			end := time.Now().Add(time.Duration(int64(this.timeout)) * time.Second)
			res, err := brainfuck.New().SetInterrupter(func() bool {
				return time.Now().After(end)
			}).Exec(code)

			if err != nil {
				this.tg.ReplyToMessage(msg.MessageId(), err.Error(), msg.ChatId())
			} else {
				res = strings.Trim(res, " \n")

				if res == "" {
					res = "Empty"
				}

				this.tg.ReplyToMessage(msg.MessageId(), res, msg.ChatId())
			}
		}
	}
}

func (this *Script) Default(name string, msg telegram.TObject, state *map[string]interface{}) {
}
