package vkapi

import "github.com/tidwall/gjson"

var (
	commandHandlers = make(map[string]func(gjson.Result))
	messageHandlers = make(map[string]func(gjson.Result))
	eventHandlers   = make(map[string]func(gjson.Result))
	commandToken    = "/"
	token           string
	apiVersion      string
	groupID         string
)
