package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

var (
	commandHandlers = make(map[string]func(gjson.Result))
	messageHandlers = make(map[string]func(gjson.Result))
	eventHandlers   = make(map[string]func(gjson.Result))
	commandToken    = "/"
	token           *string
	apiVersion      *string
	groupID         *string
)

func getBotLongPollServer(commandToken ...string) (secretKey, server, ts string, custErr error) {

	//Запрос к апи
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/groups.getLongPollServer?group_id=%s&access_token=%s&v=%s", *groupID, *token, *apiVersion))
	if err != nil { //проверка на ошибки
		println("Smth wrong")
		return "", "", "", errors.New("Ошибка при запросе")
	}
	defer resp.Body.Close() //закрываем тело по завершении функции
	var buff bytes.Buffer   //буфер для копирования тела ответа
	io.Copy(&buff, resp.Body)
	if err != nil {
		fmt.Println(fmt.Errorf("Smth is wrong: %g", err))
		return "", "", "", errors.New("Ошибка при получении тела ответа")
	}
	var dat map[string]interface{}
	json.Unmarshal(buff.Bytes(), &dat) //разборка на запчасти

	response := dat["response"].(map[string]interface{}) // получаем сам ответ

	secretKey = response["key"].(string)
	server = response["server"].(string)
	ts = response["ts"].(string)

	return secretKey, server, ts, nil
}

/*startPolling слушает ответы от LongPoll сервера вк и отправляет
все ивенты обработчик в виде []gjson.Result*/
func startPolling(secretKey, server, ts string, events chan []gjson.Result) {
	for true {
		resp, err := http.Get(fmt.Sprintf("%s?act=a_check&key=%s&ts=%s&wait=25", server, secretKey, ts))
		if err != nil { //проверка на ошибки
			println("Smth wrong: ")
		}
		var buffer bytes.Buffer
		io.Copy(&buffer, resp.Body)
		defer resp.Body.Close()
		updates := gjson.Get(buffer.String(), "updates").Array()
		if len(updates) != 0 {
			events <- updates
		}
		ts = gjson.Get(buffer.String(), "ts").String()
	}
}

//<Добавляем хендлеры>
func addCommandHandler(handler func(gjson.Result), command string) {
	commandHandlers[command] = handler
}

func addEventHandler(handler func(gjson.Result), event string) {
	eventHandlers[event] = handler
}

func addMessageHandler(handler func(gjson.Result), message string) {
	messageHandlers[message] = handler
}

//</Добавляем хендлеры>

//Распределение событий на соответствующие потоки(command/message/event)
func updatesHandling(events chan []gjson.Result) {
	//Установка заглушки на дефолтную прослушку.
	commandHandlers["default"] = func(gjson.Result) {}
	messageHandlers["default"] = func(gjson.Result) {}
	eventHandlers["default"] = func(gjson.Result) {}
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

//Функця обращение к API VK -> ответ сервера либо ошибку при подключении
func vkAPI(method string, args map[string]string) string {
	address := fmt.Sprintf("https://api.vk.com/method/%s?access_token=%s&v=%s", method, *token, *apiVersion)
	for key, value := range args {
		address = fmt.Sprintf("%s&%s=%s", address, key, url.QueryEscape(value))
	}
	response, err := http.Get(address)
	if err != nil {
		fmt.Printf("Smth wrong: %s", err.Error())
		return err.Error()
	}
	var buffer bytes.Buffer
	io.Copy(&buffer, response.Body)
	defer response.Body.Close()
	if len(gjson.Get(buffer.String(), "error.error_msg").String()) != 0 {
		return gjson.Get(buffer.String(), "error.error_msg").String()
	}
	return buffer.String()
}

//Полуение рандомного ID для некотрых методов API VK -> string
func getRandomID() string {
	rand.Seed(time.Now().UnixNano())
	if time.Now().UnixNano()%2 == 0 {
		return strconv.Itoa(int(rand.Int31()))
	}
	return strconv.Itoa(int(rand.Int31() * -1))
}

//пример хендлера
func handler(event gjson.Result) {
	if time.Now().UnixNano()%2 == 0 {
		vkAPI("messages.send", map[string]string{"random_id": getRandomID(),
			"message": "Смысол есть",
			"peer_id": event.Get("object.peer_id").String()})
		return
	}
	vkAPI("messages.send", map[string]string{"random_id": getRandomID(),
		"message": "Смысола нет",
		"peer_id": event.Get("object.peer_id").String()})
}

func main() {
	token = flag.String("token", "***", "BotLongPollToken") //Заменить на свой токен звездочки
	apiVersion = flag.String("apiVersion", "5.92", "VkAPIVersion")
	groupID = flag.String("groupID", "***", "Vk Group ID") //Заменить на group id звездочки
	flag.Parse()
	if *token == "" || *groupID == "" {
		fmt.Println("Lack of args")
		return
	}
	events := make(chan []gjson.Result)
	//запрос LongPoll сервера
	secretKey, server, ts, err := getBotLongPollServer()
	if err != nil {
		fmt.Printf("Can't get LongPoll Server, exiting.\nError: %s\n", err)
		return
	}
	fmt.Printf("%s, %s, %s\n", secretKey, server, ts)
	go updatesHandling(events)
	addCommandHandler(handler, "смысол") //пример назанчения хендлера
	startPolling(secretKey, server, ts, events)
	return
}
