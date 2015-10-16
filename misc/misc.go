// Misc features
package misc

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"
	"github.com/PeterCxy/gotgbot/support/utils"
)

type Misc struct {
	tg *telegram.Telegram
}

func Setup(t *telegram.Telegram, config map[string]interface{}, modules map[string]bool, cmds *types.CommandMap) types.Command {
	if val, ok := modules["misc"]; !ok || val {
		misc := &Misc{tg: t}

		// Echo
		(*cmds)["echo"] = types.Command{
			Name:      "echo",
			Args:      "<text>",
			ArgNum:    1,
			Desc:      "Echo <text>",
			Processor: misc,
		}

		// Remind
		(*cmds)["remind"] = types.Command{
			Name:      "remind",
			ArgNum:    0,
			Desc:      "Remind you of something after a period of time",
			Processor: misc,
		}

		// Choose
		(*cmds)["choose"] = types.Command{
			Name:   "choose",
			Args:   "<choices> <format>...",
			ArgNum: -1,
			Desc: `Choose from <choices> and format the result with <format>
			<choices> format: item1choice1,item1hoice2;item2choice1,item2choice2;[start-end]. Wrap with quotes if containing spaces.
			<format> is a printf-like format string. Each item in <choice> is mapped to a parameter for printf.`,
			Processor: misc,
		}

		// Cancel
		(*cmds)["cancel"] = types.Command{
			Name:      "cancel",
			ArgNum:    0,
			Desc:      "Cancel the current session with this bot",
			Processor: misc,
		}

		// Debug parse
		(*cmds)["parse"] = types.Command{
			Name:      "parse",
			Args:      "arguments",
			ArgNum:    -1,
			Desc:      "Parse argument list [debug]",
			Debug:     true,
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
		utils.SetGrabber(types.Grabber{
			Name:      "remind",
			Uid:       msg.FromId(),
			Chat:      msg.ChatId(),
			Processor: this,
		})
	case "choose":
		this.choose(msg, args)
	}
}

func (this *Misc) Default(name string, msg telegram.TObject, state *map[string]interface{}) {
	if name == "remind" {
		if (*state)["remind"] == nil {
			(*state)["remind"] = msg["text"].(string)
			this.tg.ReplyToMessage(msg.MessageId(), "How long after now should I remind you?", msg.ChatId())
		} else {
			duration, err := time.ParseDuration(msg["text"].(string))
			if err != nil {
				this.tg.ReplyToMessage(msg.MessageId(), "Invalid time. Supported format: BhCmDsEmsFnsGus", msg.ChatId())
			} else {
				text := (*state)["remind"].(string)

				this.tg.ReplyToMessage(msg.MessageId(), "Yes, sir!", msg.ChatId())

				utils.PostDelayed(func() {
					this.tg.SendMessage("@"+msg.From()["username"].(string)+" "+text, msg.ChatId())
				}, duration)

				utils.ReleaseGrabber(msg.FromId(), msg.ChatId())
			}
		}
	}
}

func (this *Misc) choose(msg telegram.TObject, args []string) {
	if len(args) <= 1 {
		this.tg.ReplyToMessage(msg.MessageId(), "Please provide the output format.", msg.ChatId())
	} else {
		items := strings.Split(args[0], ";")
		results := make([]interface{}, len(items))
		for i, v := range items {
			// Random number in a range
			if strings.HasPrefix(v, "[") && strings.HasSuffix(v, "]") && strings.Contains(v, "-") {
				v = v[1 : len(v)-1]
				a := strings.Split(v, "-")
				if len(a) == 2 {
					start, err1 := strconv.ParseInt(a[0], 10, 64)
					end, err2 := strconv.ParseInt(a[1], 10, 64)

					if (err1 == nil) && (err2 == nil) {
						r := float64(start) + rand.Float64()*float64(end-start)
						results[i] = r
						continue
					}
				}

				this.tg.ReplyToMessage(msg.MessageId(), "Range format: [start-end]", msg.ChatId())
				return
			} else {
				a := strings.Split(v, ",")
				results[i] = a[rand.Intn(len(a))]
			}
		}

		format := strings.Join(args[1:], " ")
		tokens := ParseFormat(format)

		// Now let's do the heavy type conversion stuff
		i := 0
		for _, t := range tokens {
			if i >= len(results) {
				break
			}

			if len(t) < 1 {
				continue
			}

			switch t[len(t)-1] {
			case 'b', 'c', 'd', 'o', 'q', 'x', 'X', 'U':
				// Integer
				switch t := results[i].(type) {
				case string:
					results[i], _ = strconv.ParseInt(results[i].(string), 10, 64)
				case float64:
					results[i] = int64(results[i].(float64))
				default:
					_ = t
				}
			case 'e', 'E', 'f', 'F', 'g', 'G':
				// Float
				switch t := results[i].(type) {
				case string:
					results[i], _ = strconv.ParseFloat(results[i].(string), 64)
				case int64:
					results[i] = float64(results[i].(int64))
				default:
					_ = t
				}
			}

			if t != "%" {
				i += 1
			}
		}

		this.tg.ReplyToMessage(
			msg.MessageId(),
			fmt.Sprintf(format, results...),
			msg.ChatId())
	}
}
