package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/loader"
	"github.com/PeterCxy/gotgbot/support/types"
)

var Telegram *telegram.Telegram
var Commands types.CommandMap
var Default types.Command
var BotName string
var Debug bool

func main() {
	log.SetPrefix("Bot")
	// Load config
	if len(os.Args) <= 1 {
		log.Fatalln("Please provide path to the config file")
	}

	b, err := ioutil.ReadFile(os.Args[1])

	if (err != nil) || (b == nil) {
		log.Fatalf("Cannot read config file %s", os.Args[1])
	}

	var c interface{}
	err = json.Unmarshal(b, &c)

	if (err != nil) || (c == nil) {
		log.Fatalln("Failed to decode config file")
	}

	config := c.(map[string]interface{})

	Debug = false

	if config["debug"] != nil {
		Debug = config["debug"].(bool)
	}

	if Debug {
		log.Println(config)
	}

	if config["name"] == nil {
		log.Fatalln("Please provide the name of the bot")
	}

	BotName = config["name"].(string)

	Telegram = telegram.New(config["key"].(string), Debug)

	// Setup modules
	Commands, Default = loader.LoadModules(Telegram, config)

	// Boot up the server
	MainLoop()
}
