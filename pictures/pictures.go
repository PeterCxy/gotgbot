// Pictures fetcher
package pictures

import (
	"github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"
)

type Pictures struct {
	tg *telegram.Telegram
}

func Setup(t *telegram.Telegram, config map[string]interface{}, modules map[string]bool, cmds *types.CommandMap) types.Command {
	if val, ok := modules["pictures"]; !ok || val {
		pictures := &Pictures{tg: t}

		// Meizhi (gank.io girls)
		(*cmds)["meizhi"] = types.Command {
			Name: "meizhi",
			ArgNum: 0,
			Desc: "Random picture of girls from gank.io",
			Processor: pictures,
		}
	}

	return types.Command{}
}

func (this *Pictures) Command(name string, msg telegram.TObject, args []string) {
	if name == "meizhi" {
		this.Meizhi(msg)
	}
}

func (this *Pictures) Default(name string, msg telegram.TObject, state *map[string]interface{}) {
}
