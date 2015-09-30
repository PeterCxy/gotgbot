package loader

import (
	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/types"
	"github.com/PeterCxy/gotgbot/support/help"
)

func LoadModules(tg *telegram.Telegram, config map[string]interface{}) (types.CommandMap, types.Command) {
	modules := parseModules(config["modules"].(map[string]interface{}))
	ret := make(types.CommandMap)
	def := types.Command{}

	// Help
	if d := help.Setup(tg, config, modules, &ret); d.Processor != nil {
		def = d
	}

	return ret, def
}

func parseModules(m map[string]interface{}) map[string]bool {
	ret := make(map[string]bool)
	for k, v := range m {
		ret[k] = v.(bool)
	}

	return ret
}
