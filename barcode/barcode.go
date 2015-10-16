// Bar/QRcode reader / gendrator
package barcode

import (
	"github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"
	"github.com/PeterCxy/gotgbot/support/utils"
)

type Barcode struct {
	tg *telegram.Telegram
}

func Setup(t *telegram.Telegram, config map[string]interface{}, modules map[string]bool, cmds *types.CommandMap) types.Command {
	if val, ok := modules["barcode"]; !ok || val {
		barcode := &Barcode{tg: t}

		(*cmds)["barcode"] = types.Command{
			Name:      "barcode",
			Desc:      "Decode a barcode / qrcode. Reply to a message containing the picture of the code or call this command directly, I'll ask you for the picture.",
			ArgNum:    0,
			Processor: barcode,
		}
	}

	return types.Command{}
}

func (this *Barcode) Command(name string, msg telegram.TObject, args []string) {
	if name == "barcode" {
		if msg["reply_to_message"] != nil {
			// Decode the message replied to
			this.Decode(msg.ReplyToMessage())
		} else {
			// Decode from grabbed input
			this.tg.ReplyToMessage(msg.MessageId(), "Now send me the picture to decode.", msg.ChatId())
			utils.SetGrabber(types.Grabber{
				Name:      "barcode",
				Uid:       msg.FromId(),
				Chat:      msg.ChatId(),
				Processor: this,
			})
		}
	}
}

func (this *Barcode) Default(name string, msg telegram.TObject, state *map[string]interface{}) {
	if name == "barcode" {
		utils.ReleaseGrabber(msg.FromId(), msg.ChatId())
		this.Decode(msg)
	}
}
