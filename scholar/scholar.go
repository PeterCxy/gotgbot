// scholar bot
package scholar

import (
	"fmt"
	"strings"

	"github.com/PeterCxy/gocalc"

	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"
)

type Scholar struct {
	tg *telegram.Telegram
}

func Setup(t *telegram.Telegram, config map[string]interface{}, modules map[string]bool, cmds *types.CommandMap) types.Command {
	if val, ok := modules["scholar"]; !ok || val {
		scholar := &Scholar{tg: t}

		// Calc
		(*cmds)["calc"] = types.Command {
			Name: "calc",
			Args: "<expression>",
			ArgNum: -1,
			Desc: "Calculate <expression>, only math expression supported.",
			Processor: scholar,
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
	}
}

func (this *Scholar) Default(name string, msg telegram.TObject, state *map[string]interface{}) {
}
