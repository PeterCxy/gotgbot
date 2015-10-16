// Barcode decoder
package barcode

import (
	"image"
	"image/jpeg"
	"strings"

	"github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gozbar"
	"github.com/ddliu/go-httpclient"
)

func (this *Barcode) Decode(msg telegram.TObject) {
	photos := msg.Photo()

	// A Photo must be sent
	if (photos == nil) || (len(photos) == 0) {
		this.tg.ReplyToMessage(msg.MessageId(), "I did not get anything to decode. o_O", msg.ChatId())
		return
	}

	this.tg.SendChatAction("typing", msg.ChatId())

	p := this.tg.GetFile(photos[len(photos)-1].FileId())

	if p == nil {
		return
	}

	url := this.tg.PathToUrl(p.FilePath())

	res, err := httpclient.Get(url, nil)

	if err != nil {
		return
	}

	defer res.Body.Close()

	var img image.Image

	img, err = jpeg.Decode(res.Body)

	if err != nil {
		return
	}

	// Create ZBar Image object
	zimg := zbar.FromImage(img)

	// Create Scanner
	scanner := zbar.NewScanner()
	scanner.SetConfig(0, zbar.CFG_ENABLE, 1)

	// Do scan
	if scanner.Scan(zimg) <= 0 {
		this.tg.ReplyToMessage(msg.MessageId(),
			"Failed to decode the code. Please make sure that there is one or more valid codes inside that picture.", msg.ChatId())
		return
	}

	// Iterate over the decoded symbols
	values := make([]string, 0)
	zimg.First().Each(func(s string) {
		values = append(values, s)
	})

	this.tg.ReplyToMessage(msg.MessageId(), strings.Join(values, "\n"), msg.ChatId())
}
