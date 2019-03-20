package vkapi

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

//Обработка простых сообщений
func messageHandling(event gjson.Result) {
	text := event.Get("object.text").String()
	for key, value := range messageHandlers {
		if key == text {
			value(event)
			break
		}
	}
	messageHandlers["default"](event)

}

//Обработка сообщений-команд
func commandHandling(event gjson.Result) {
	command := strings.Split(event.Get("object.text").String(), " ")[0][1:]
	command = strings.ToLower(command)
	for key, value := range commandHandlers {
		fmt.Println(command)
		if key == command {
			value(event)
			break
		}
	}
	commandHandlers["default"](event)

}

//Обработка иных событий
func eventHandling(event gjson.Result) {
	eventType := event.Get("type").String()
	for key, value := range eventHandlers {
		if key == eventType {
			value(event)
		}
	}
	eventHandlers["default"](event)
}
