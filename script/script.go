// Scripting support
package script

import (
	"fmt"
	"strings"
	"time"

	"github.com/PeterCxy/gobf"
	"github.com/PeterCxy/gogmh"
	"github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"
	"github.com/PeterCxy/gotgbot/support/utils"
)

type Script struct {
	tg      *telegram.Telegram
	timeout int
}

func Setup(t *telegram.Telegram, config map[string]interface{}, modules map[string]bool, cmds *types.CommandMap) types.Command {
	if val, ok := modules["script"]; !ok || val {
		script := &Script{tg: t}
		if val, ok := config["script_timeout"]; ok {
			script.timeout = int(val.(float64))
		}

		// Code
		(*cmds)["code"] = types.Command{
			Name:      "code",
			ArgNum:    0,
			Desc:      "Store some code in memory and execute it as a script.",
			Processor: script,
		}

		// Brainfuck
		(*cmds)["brainfuck"] = types.Command{
			Name:      "brainfuck",
			Args:      "<code>",
			ArgNum:    1,
			Desc:      "Execute brainfuck code",
			Processor: script,
		}

		// Grass-Mud-Horse
		(*cmds)["cnm"] = types.Command{
			Name:      "cnm",
			Args:      "<code>...",
			ArgNum:    -1,
			Desc:      "Execute Grass-Mud-Horse commands splitted with spaces. https://code.google.com/p/grass-mud-horse/wiki/A_Brife_To_GrassMudHorse_Language",
			Processor: script,
		}
	}

	return types.Command{}
}

func (this *Script) Command(name string, msg telegram.TObject, args []string) {
	if name == "code" {
		this.tg.ReplyToMessage(msg.MessageId(), "Now send me some code. You can split it into multiple messages.\nAfter finishing, send 'exec script_type' to execute it as a script.", msg.ChatId())
		status := utils.SetGrabber(types.Grabber{
			Name:      "code",
			Uid:       msg.FromId(),
			Chat:      msg.ChatId(),
			Processor: this,
		})

		(*status)["code"] = ""
	} else if (name == "brainfuck") || (name == "cnm") {
		code := ""
		if name == "brainfuck" {
			code = strings.Trim(args[0], " \n")
		} else if name == "cnm" {
			code = strings.Trim(strings.Join(args, " "), " \n")
		}

		if code == "" {
			this.tg.ReplyToMessage(msg.MessageId(), "Code is empty.", msg.ChatId())
		} else {
			end := time.Now().Add(time.Duration(int64(this.timeout)) * time.Second)
			interrupterFunc := func() bool {
				return time.Now().After(end)
			}
			inputFunc := func(out string) string {
				this.tg.ReplyToMessage(msg.MessageId(),
					fmt.Sprintf("Output: %s\nInput needed. Now send me the input data in 60 seconds. If nothing received, the interpreter will be interrupted.\nSend /cancel to force interrupt.", out),
					msg.ChatId())

				status := utils.SetGrabber(types.Grabber{
					Name:      "input",
					Uid:       msg.FromId(),
					Chat:      msg.ChatId(),
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
				case <-time.After(60 * time.Second):
					utils.ReleaseGrabber(msg.FromId(), msg.ChatId())
					this.tg.ReplyToMessage(msg.MessageId(), "Input timed out. Interrupting.", msg.ChatId())
					return string(0)
				}

				return string(0)
			}

			var res string
			var err error

			if name == "brainfuck" {
				res, err = brainfuck.New().
					SetInterrupter(interrupterFunc).
					SetInput(inputFunc).
					Exec(code)
			} else if name == "cnm" {
				res, err = gmh.New().
					SetInterrupter(interrupterFunc).
					SetInput(inputFunc).
					Exec(strings.FieldsFunc(strings.Join(args, " "), func(r rune) bool {
					return (r == ' ') || (r == '\n')
				}))
			}

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
	if name == "code" {
		if msg["text"] == nil {
			return
		}

		str := msg["text"].(string)

		if !strings.HasPrefix(strings.ToLower(str), "exec ") {
			(*state)["code"] = (*state)["code"].(string) + "\n" + str

			if len((*state)["code"].(string)) >= 131070 {
				this.tg.SendMessage("Code too long. Maximum: 131KiB", msg.ChatId())
				utils.ReleaseGrabber(msg.FromId(), msg.ChatId())
			} else {
				this.tg.SendMessage("Code received. If you have finished, send 'exec script_type' to execute. Send /cancel to cancel.", msg.ChatId())
			}
		} else {
			code := (*state)["code"].(string)
			utils.ReleaseGrabber(msg.FromId(), msg.ChatId())
			this.Command(strings.ToLower(str[5:]), msg, []string{code})
		}

	} else if name == "input" {
		if msg["text"] != nil {
			text := strings.Trim(msg["text"].(string), " ")

			if text != "" {
				input := (*state)["chan"].(chan string)
				input <- strings.Replace(text, "\\n", "\n", -1)
				utils.ReleaseGrabber(msg.FromId(), msg.ChatId())
			}
		}
	}
}
