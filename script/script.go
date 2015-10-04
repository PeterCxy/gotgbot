// Scripting support
package script

import (
	"strings"
	"time"
	"fmt"

	"github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gobf"
	"github.com/PeterCxy/gotgbot/support/types"
	"github.com/PeterCxy/gotgbot/support/utils"
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
			}).SetInput(func(out string) string {
				this.tg.ReplyToMessage(msg.MessageId(),
					fmt.Sprintf("Output: %s\nInput needed. Now send me the input data in 30 seconds. If nothing received, the interpreter will be interrupted.\nSend /cancel to force interrupt.", out),
					msg.ChatId())

				status := utils.SetGrabber(types.Grabber {
					Name: "brainfuck",
					Uid: msg.FromId(),
					Chat: msg.ChatId(),
					Processor: this,
				})

				input := make(chan string)
				(*status)["chan"] = input

				now := time.Now()
				var dur time.Duration
				if now.Before(end) {
					dur = end.Sub(now)
				} else {
					// WTF??
					dur = time.Duration(0)
				}
				select {
					case result := <-input:
						end = time.Now().Add(dur)
						return result
					case <-time.After(30 * time.Second):
						utils.ReleaseGrabber(msg.FromId(), msg.ChatId())
						this.tg.ReplyToMessage(msg.MessageId(), "Input timed out. Interrupting.", msg.ChatId())
						return string(0)
				}

				return string(0)
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
	if name == "brainfuck" {
		if msg["text"] != nil {
			text := strings.Trim(msg["text"].(string), " ")

			if text != "" {
				input := (*state)["chan"].(chan string)
				input <- text
				utils.ReleaseGrabber(msg.FromId(), msg.ChatId())
			}
		}
	}
}
