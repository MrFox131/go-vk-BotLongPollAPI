package vkapi

import "github.com/tidwall/gjson"

//AddCommandHandler adding command handler
func AddCommandHandler(handler func(gjson.Result), command string) {
	commandHandlers[command] = handler
}

//AddEventHandler adding event handler
func AddEventHandler(handler func(gjson.Result), event string) {
	eventHandlers[event] = handler
}

//AddMessageHandler adding message handler
func AddMessageHandler(handler func(gjson.Result), message string) {
	messageHandlers[message] = handler
}
