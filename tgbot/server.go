package main

import (
	"log"
	"time"
	"strings"
	"fmt"

	telegram "github.com/PeterCxy/gotelegram"
	"github.com/PeterCxy/gotgbot/support/utils"
	//"github.com/PeterCxy/gotgbot/support/types"
)

func MainLoop() {
	// Set webhook to null first
	Telegram.SetWebhook("")

	// Start the grabber daemon
	go utils.GrabberDaemon()

	// Start the schedule daemon
	go utils.ScheduleDaemon()

	// Loop forever
	var offset int64 = 0
	for {
		updates := Telegram.GetUpdates(offset, 100, 20)

		if updates != nil {
			for i, v := range updates {
				update := telegram.TObject(v.(map[string]interface{}))
				if update["message"] == nil {
					continue
				}

				message := update.Message()

				// Process each message in seperate goroutines
				go handle(message)

				if i == (len(updates) - 1) {
					offset = update.UpdateId() + 1
				}
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func handle(msg telegram.TObject) {
	if Debug {
		log.Println(msg)
	}

	// Parse arguments (including the command)
	args := make([]string, 0)
	text := ""

	if msg["text"] != nil {
		text = strings.Trim(msg["text"].(string), " ")
	}

	if text != "" {
		args = telegram.ParseArgs(text)
	}

	if len(args) == 0 {
		return
	}

	// A command
	if strings.HasPrefix(args[0], "/") {
		cmd := args[0][1:]
		args = args[1:]

		if strings.Contains(cmd, "@") {
			if !strings.HasSuffix(cmd, "@" + BotName) {
				// This command is not our business
				return
			} else {
				cmd = cmd[0:strings.Index(cmd, "@")]
			}
		}

		command := Commands[cmd]
		if (command.ArgNum < 0) || (command.ArgNum == len(args)) {
			command.Processor.Command(cmd, msg, args)
		} else {
			str := fmt.Sprintf("Usage: /%s %s\nDescription: %s", command.Name, command.Args, command.Desc)
			Telegram.ReplyToMessage(
				msg.MessageId(),
				str,
				msg.ChatId())
		}
	} else if utils.HasGrabber(msg.FromId(), msg.ChatId()) {
		// Distribute to grabbers
		name, processor := utils.Grabber(msg.FromId(), msg.ChatId())
		processor.Default(name, msg, utils.GrabberState(msg.FromId(), msg.ChatId()))
	}

	// TODO Distribute to default processor
}
