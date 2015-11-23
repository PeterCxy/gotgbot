package gank

import (
	"encoding/json"
	"fmt"
	"strings"
	"os"
	"log"
	"time"
	"io/ioutil"

	"github.com/ddliu/go-httpclient"

	"github.com/PeterCxy/gotelegram"
)

const api = "http://gank.avosapps.com/api/day/%d/%d/%d"

func Init(telegram *telegram.Telegram, modules map[string]bool, config map[string]interface{}) {
	if val, ok := modules["chan_gank"]; !ok || val {
		channel := config["chan_gank"].(string)

		log.Printf("gank.io: setting up channel @%s", channel)

		go func() {
			for {
				update(telegram, channel)
				time.Sleep(24 * time.Hour)
			}
		}()
	}
}

func update(tg *telegram.Telegram, channel string) {
	// Get current time
	t := time.Now()

	url := fmt.Sprintf(api, t.Year(), t.Month(), t.Day())

	log.Printf("gank.io update url: %s" , url)

	res, err := httpclient.Get(url, nil)

	if err != nil {
		return
	}

	b, _ := res.ReadAll()

	var data interface{}

	err = json.Unmarshal(b, &data)

	if err != nil {
		return
	}

	r := data.(map[string]interface{})

	if val, ok := r["error"]; !ok || val.(bool) {
		return
	}

	results := r["results"].(map[string]interface{})

	text := ""
	for _, c := range r["category"].([]interface{}) {
		category := c.(string)
		text += fmt.Sprintf("*%s*\n", category)

		articles := results[category].([]interface{})

		for _, article := range articles {
			ar := article.(map[string]interface{})

			text += fmt.Sprintf("[%s](%s) by %s\n",
				filter(ar["desc"].(string)),
				ar["url"].(string),
				telegram.Escape(ar["who"].(string)))
		}

		text += "\n\n"
	}

	if text != "" {
		tg.SendMessageChan(text, channel)
	} else {
		log.Println("gank.io: Oooooops, no update today!")
	}

	// Girls! Girls! Girls!
	// Exciting!!!!
	if _, ok := results["福利"]; !ok {
		return
	}

	pictures := results["福利"].([]interface{})

	for _, pic := range pictures {
		picture := pic.(map[string]interface{})
		picUrl := picture["url"].(string)

		res, err := httpclient.Get(picUrl, nil)

		if err != nil {
			continue
		}

		name := fmt.Sprintf("/tmp/fuli_%d.jpg", time.Now().Unix())

		defer func() {
			os.Remove(name)
		}()

		b, _ := res.ReadAll()

		err = ioutil.WriteFile(name, b, os.ModePerm)

		if err != nil {
			return
		}

		tg.SendPhotoChan(name, channel)
	}
}

func filter(str string) (ret string) {
	ret = strings.Replace(str, "[", "［", -1)
	ret = strings.Replace(str, "]", "］", -1)
	return
}
