package main

import (
	"flag"
	"fmt"
	"time"

	"./vkapi"

	"github.com/tidwall/gjson"
)

func handler(event gjson.Result) {
	if time.Now().UnixNano()%2 == 0 {
		vkapi.VkAPI("messages.send", map[string]string{"random_id": vkapi.GetRandomID(),
			"message": "Смысол есть",
			"peer_id": event.Get("object.peer_id").String()})
		return
	}
	vkapi.VkAPI("messages.send", map[string]string{"random_id": vkapi.GetRandomID(),
		"message": "Смысола нет",
		"peer_id": event.Get("object.peer_id").String()})
}

func main() {
	token := flag.String("token", "2335fc759f1e25c246b3188c793ce487d5e71513fa8a15980e92c8d126f60a7e254775269fc11561ebd76", "BotLongPollToken")
	apiVersion := flag.String("apiVersion", "5.92", "vkapi.VkAPIVersion")
	groupID := flag.String("groupID", "179027597", "Vk Group ID")

	flag.Parse()
	if *token == "" || *groupID == "" {
		fmt.Println("Lack of args")
		return
	}

	//запрос LongPoll сервера
	secretKey, server, ts, err := vkapi.GetBotLongPollServer(*token, *groupID, *apiVersion)
	if err != nil {
		fmt.Printf("Can't get LongPoll Server, exiting.\nError: %s\n", err)
		return
	}

	fmt.Printf("%s, %s, %s\n", secretKey, server, ts)
	vkapi.AddCommandHandler(handler, "смысол")
	vkapi.StartPolling(secretKey, server, ts)
	return
}
