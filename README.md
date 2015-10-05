gotgbot
---
Telegram bot example implementation in Golang

Build
---
```
go get github.com/PeterCxy/gotgbot/tgbot
```

Usage
---
```
tgbot /path/to/config.json
```

Config
---
Basic config

```json
{
	"key": "api_key",
	"name": "bot_name",
	"debug": true/false,
	"modules": [
		"module1": true/false,
		"module2": true/false,
		...
	],
}
```

To make the bot work, you will need extra configurations for modules.

Modules
---
Modules are defined in every sub-package in this repo. To add new modules, make a new package and register it in `support/loader/loader.go`

Some modules may need specific configuration in the config file to work. See the module sources for details.

License
---
See the file `LICENSE`
