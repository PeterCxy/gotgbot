// Chinese language model
package chinese

import (
	"strings"
	"errors"

	"gopkg.in/redis.v3"
	"github.com/huichen/sego"

	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"
)

type Chinese struct {
	tg *telegram.Telegram
	redis *redis.Client
	seg sego.Segmenter
	debug bool
}

func Setup(t *telegram.Telegram, config map[string]interface{}, modules map[string]bool, cmds *types.CommandMap) types.Command {
	if val, ok := modules["chinese"]; !ok || val {
		var s sego.Segmenter
		s.LoadDictionary(config["dict"].(string))
		c := &Chinese {
			tg: t,
			redis: redis.NewClient(&redis.Options {
				Addr: config["redis"].(string),
				Password: "",
				DB: int64(config["redis_db"].(float64)),
			}),
			seg: s,
		}

		if config["debug"] != nil {
			c.debug = config["debug"].(bool)
		}

		(*cmds)["learn"] = types.Command {
			Name: "learn",
			Args: "<expr>",
			ArgNum: -1,
			Desc: "Learn a Chinese expression",
			Processor: c,
		}

		(*cmds)["speak"] = types.Command {
			Name: "speak",
			ArgNum: 0,
			Desc: "Speak a Chinese sentence based on previously learned data",
			Processor: c,
		}

		pong, err := c.redis.Ping().Result()

		if (err != nil) || (pong != "PONG") {
			panic(errors.New("Cannot PING redis"))
		}

		return types.Command {
			Name: "chn",
			Processor: c,
		}
	}

	return types.Command{}
}

func (this *Chinese) Command(name string, msg telegram.TObject, args []string) {
	if name == "learn" {
		this.Learn(strings.Join(args, " "), msg.ChatId())
	} else if name == "speak" {
		this.tg.SendMessage(this.Speak(msg.ChatId()), msg.ChatId())
	}
}

func (this *Chinese) Default(name string, msg telegram.TObject, state *map[string]interface{}) {
	if msg["text"] != nil {
		text := msg["text"].(string)
		if !strings.HasPrefix(text, "/") {
			this.Learn(text, msg.ChatId())
		}
	}
}
