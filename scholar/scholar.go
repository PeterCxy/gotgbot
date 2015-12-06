// scholar bot
package scholar

import (
	"fmt"
	"strings"

	"github.com/PeterCxy/gocalc"

	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"
	"github.com/PeterCxy/gotgbot/support/utils"
)

type Scholar struct {
	tg   *telegram.Telegram
	ipv6 bool
}

func Setup(t *telegram.Telegram, config map[string]interface{}, modules map[string]bool, cmds *types.CommandMap) types.Command {
	if val, ok := modules["scholar"]; !ok || val {
		scholar := &Scholar{tg: t}

		// Calc
		(*cmds)["calc"] = types.Command{
			Name:      "calc",
			Args:      "<expression>",
			ArgNum:    -1,
			Desc:      "Calculate <expression>, only math expression supported.",
			Processor: scholar,
		}

		// Google
		(*cmds)["google"] = types.Command{
			Name:      "google",
			Args:      "<query>",
			ArgNum:    -1,
			Desc:      "Search Google for <query>",
			Processor: scholar,
		}

		if val, ok := config["ipv6_google"]; !ok || val.(bool) {
			scholar.ipv6 = true
		} else {
			scholar.ipv6 = false
		}
	}

	return types.Command{}
}

func (this *Scholar) Command(name string, msg telegram.TObject, args []string) {
	if name == "calc" {
		res, err := calc.Calculate(strings.Join(args, " "))

		if err == nil {
			this.tg.ReplyToMessage(msg.MessageId(), fmt.Sprintf("%f", res), msg.ChatId())
		} else {
			this.tg.ReplyToMessage(msg.MessageId(), err.Error(), msg.ChatId())
		}
	} else if name == "google" {
		query := strings.Join(args, " ")

		if query == "" {
			this.tg.ReplyToMessage(msg.MessageId(), "Please provide something to search for.", msg.ChatId())
		} else {
			num := 5
			maxNum := 5
			irc := false

			if (msg.Chat()["title"] != nil) && strings.HasPrefix(msg.Chat()["title"].(string), "#") {
				num = 1 // Disable long output in IRC-connected groups
				irc = true
			}

			this.tg.SendChatAction("typing", msg.ChatId())
			res, hasNext := Google(query, 0, maxNum, this.ipv6)

			if len(res) > num {
				res = res[0:num]
			}

			if irc {
				hasNext = false
			}

			this.tg.SendMessageNoPreview(formatGoogle(res, hasNext), msg.ChatId())

			if hasNext {
				state := utils.SetGrabber(types.Grabber{
					Name:      "google",
					Uid:       msg.FromId(),
					Chat:      msg.ChatId(),
					Processor: this,
				})

				(*state)["start"] = len(res)
				(*state)["query"] = query
			}
		}
	}
}

func (this *Scholar) Default(name string, msg telegram.TObject, state *map[string]interface{}) {
	if name == "google" {
		if (msg["text"] == nil) || (strings.ToLower(msg["text"].(string)) != "next") {
			return
		}

		start := (*state)["start"].(int)
		query := (*state)["query"].(string)

		this.tg.SendChatAction("typing", msg.ChatId())
		res, hasNext := Google(query, start, 5, this.ipv6)

		this.tg.SendMessageNoPreview(formatGoogle(res, hasNext), msg.ChatId())

		if hasNext {
			(*state)["start"] = start + len(res)
		} else {
			utils.ReleaseGrabber(msg.FromId(), msg.ChatId())
		}
	}
}

func formatGoogle(ret []GoogleResult, hasNext bool) (str string) {
	for _, res := range ret {
		str += fmt.Sprintf("%s\n%s\n\n%s\n\n", res.url, res.title, res.summary)
	}

	if hasNext {
		str += "More results available. Send me 'Next' to see more."
	}

	return
}
