package vkapi

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

func updatesHandling(events chan []gjson.Result) {
	//Установка заглушки на дефолтную прослушку.
	if commandHandlers["default"] == nil {
		commandHandlers["default"] = func(gjson.Result) {}
	}
	if messageHandlers["default"] == nil {
		messageHandlers["default"] = func(gjson.Result) {}
	}
	if eventHandlers["default"] == nil {
		eventHandlers["default"] = func(gjson.Result) {}
	}
	for true {
		for _, item := range <-events {
			switch item.Get("type").String() {
			case "message_new":
				fmt.Println(strings.HasPrefix(item.Get("object.text").String(), commandToken))
				if strings.HasPrefix(item.Get("object.text").String(), commandToken) {
					go commandHandling(item)
				} else {
					go messageHandling(item)
				}

			default:
				go eventHandling(item)
			}

		}
	}
}
