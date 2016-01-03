// Are you rich?
package rich

import (
	"log"
	"strings"
	"os"

	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/blackgun/common"

	"github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/engine"
	"github.com/NOX73/go-neural/persist"
)

var network *neural.Network
var eng engine.Engine
var file string

func Init(path string) {
	file = path + "/rich.json"
	if _, err := os.Stat(file); err != nil {
		network = neural.NewNetwork(common.SampleLen, []int{1000, 1000, 2})
	} else {
		network = persist.FromFile(file)
	}
	//network.RandomizeSynapses()
	eng = engine.New(network)
	eng.Start()
}

func Save() {
	persist.ToFile(file, network)
}

func Learn(msg telegram.TObject) bool {
	if msg["text"] != nil {
		if strings.Contains(msg["text"].(string), "#RICH") {
			eng.Learn(common.TextToSample(msg["text"].(string)), []float64{1, 0}, 0.1)
			return true
		} else if strings.Contains(msg["text"].(string), "#POOR") {
			eng.Learn(common.TextToSample(msg["text"].(string)), []float64{0, 1}, 0.1)
			return true
		} else {
			eng.Learn(common.TextToSample(msg["text"].(string)), []float64{0, 0}, 0.1)
		}
	}

	return false
}

func Gun(tg *telegram.Telegram, msg telegram.TObject) bool {
	if msg["text"] != nil {
		out := eng.Calculate(common.TextToSample(msg["text"].(string)))
		log.Println(out)

		if (out[0] > 0.8) && (out[0] - out[1] > 0.5) {
			tg.ReplyToMessage(msg.MessageId(), "#RICH", msg.ChatId())
			return true
		}
	}

	return false
}
