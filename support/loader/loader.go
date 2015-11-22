package loader

import (
	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/help"
	"github.com/PeterCxy/gotgbot/support/types"

	"github.com/PeterCxy/gotgbot/barcode"
	"github.com/PeterCxy/gotgbot/chinese"
	"github.com/PeterCxy/gotgbot/misc"
	"github.com/PeterCxy/gotgbot/pictures"
	"github.com/PeterCxy/gotgbot/scholar"
	"github.com/PeterCxy/gotgbot/script"

	"github.com/PeterCxy/gotgbot/channels/gank"
)

func LoadModules(tg *telegram.Telegram, config map[string]interface{}) (types.CommandMap, types.Command) {
	modules := parseModules(config["modules"].(map[string]interface{}))
	ret := make(types.CommandMap)
	def := types.Command{}

	// Help
	if d := help.Setup(tg, config, modules, &ret); d.Processor != nil {
		def = d
	}

	// Misc
	if d := misc.Setup(tg, config, modules, &ret); d.Processor != nil {
		def = d
	}

	// Scholar
	if d := scholar.Setup(tg, config, modules, &ret); d.Processor != nil {
		def = d
	}

	// Chinese
	if d := chinese.Setup(tg, config, modules, &ret); d.Processor != nil {
		def = d
	}

	// Script
	if d := script.Setup(tg, config, modules, &ret); d.Processor != nil {
		def = d
	}

	// Pictures
	if d := pictures.Setup(tg, config, modules, &ret); d.Processor != nil {
		def = d
	}

	// Barcode
	if d := barcode.Setup(tg, config, modules, &ret); d.Processor != nil {
		def = d
	}

	// Load Channels
	gank.Init(tg, modules, config)

	return ret, def
}

func parseModules(m map[string]interface{}) map[string]bool {
	ret := make(map[string]bool)
	for k, v := range m {
		ret[k] = v.(bool)
	}

	return ret
}
