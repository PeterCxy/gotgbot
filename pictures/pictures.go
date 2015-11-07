// Pictures fetcher
package pictures

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"

	"github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"
)

type Pictures struct {
	tg    *telegram.Telegram
	pic   string
	debug bool
}

func Setup(t *telegram.Telegram, config map[string]interface{}, modules map[string]bool, cmds *types.CommandMap) types.Command {
	if val, ok := modules["pictures"]; !ok || val {
		pictures := &Pictures{tg: t, pic: config["gallery"].(string)}

		if debug, ok := config["debug"]; ok {
			pictures.debug = debug.(bool)
		}

		// Meizhi (gank.io girls)
		(*cmds)["meizhi"] = types.Command{
			Name:      "meizhi",
			ArgNum:    0,
			Desc:      "Random picture of girls from gank.io",
			Processor: pictures,
		}

		// Local collection
		(*cmds)["pic"] = types.Command{
			Name:      "pic",
			ArgNum:    0,
			Desc:      "Fetch a picture randomly from Peter's personal collection (WARNING: Might be NSFW)",
			Processor: pictures,
		}
	}

	return types.Command{}
}

func (this *Pictures) Command(name string, msg telegram.TObject, args []string) {
	switch name {
	case "meizhi":
		this.Meizhi(msg)
	case "pic":
		this.Pic(msg)
	}
}

func (this *Pictures) Default(name string, msg telegram.TObject, state *map[string]interface{}) {
}

func (this *Pictures) Pic(msg telegram.TObject) {
	this.tg.SendChatAction("upload_photo", msg.ChatId())
	files, _ := ioutil.ReadDir(this.pic)
	f := files[rand.Intn(len(files))]
	name := fmt.Sprintf("%s/%s", this.pic, f.Name())

	if this.debug {
		log.Println(name)
	}

	this.tg.SendPhoto(name, msg.ChatId())
}
