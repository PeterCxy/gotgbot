package blackgun

import (
	"time"

	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"

	// Blackguns
	"github.com/PeterCxy/gotgbot/blackgun/rich"
)

type BG struct {
	tg *telegram.Telegram
}

func Setup(t *telegram.Telegram, config map[string]interface{}, modules map[string]bool, cmds *types.CommandMap) types.Command {
	if val, ok := modules["blackgun"]; !ok || val {
		c := &BG{tg: t}

		initialize(config["neural_path"].(string))
		go save()

		return types.Command {
			Name: "bg",
			Processor: c,
		}
	}

	return types.Command{}
}

func (this *BG) Command(name string, msg telegram.TObject, args []string) {
}

func (this *BG) Default(name string, msg telegram.TObject, state *map[string]interface{}) {
	if msg["reply_to_message"] != nil {
		// Learn blackgun!
		if learn(msg) {
			return
		}
	}

	gun(this.tg, msg)
}

func initialize(path string) {
	rich.Init(path)
}

func learn(msg telegram.TObject) bool {
	if rich.Learn(msg) {
		return true
	}
	return false
}

func gun(tg *telegram.Telegram, msg telegram.TObject) {
	if rich.Gun(tg, msg) {
		return
	}
}

func save() {
	for {
		time.Sleep(10 * time.Minute)
		rich.Save()
	}
}
