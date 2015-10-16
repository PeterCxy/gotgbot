// Picture source: gank.io
package pictures

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/ddliu/go-httpclient"

	"github.com/PeterCxy/gotelegram"
)

const api = "http://gank.avosapps.com/api/random/data/%E7%A6%8F%E5%88%A9/1"

func (this *Pictures) Meizhi(msg telegram.TObject) {
	this.tg.SendChatAction("upload_photo", msg.ChatId())

	res, err := httpclient.Get(api, nil)

	if err != nil {
		return
	}

	var data interface{}

	b, _ := res.ReadAll()

	err = json.Unmarshal(b, &data)

	if err != nil {
		return
	}

	r := data.(map[string]interface{})

	if r["error"].(bool) {
		return
	}

	url := r["results"].([]interface{})[0].(map[string]interface{})["url"].(string)

	res, err = httpclient.Get(url, nil)

	if err != nil {
		return
	}

	name := fmt.Sprintf("/tmp/meizhi_%d.jpg", time.Now().Unix())

	defer func() {
		os.Remove(name)
	}()

	b, _ = res.ReadAll()

	err = ioutil.WriteFile(name, b, os.ModePerm)

	if err != nil {
		return
	}

	this.tg.SendPhoto(name, msg.ChatId())
}
