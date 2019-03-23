package vkapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

/*GetBotLongPollServer returns secret key, server and ts for vkAPI services*/
func GetBotLongPollServer(userToken, groupID, userAPIVersion string, commandToken ...string) (secretKey, server, ts string, custErr error) {
	token = userToken
	apiVersion = userAPIVersion
	//Запрос к апи
	log.Println("Requesting server from api.vk.com")
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/groups.getLongPollServer?group_id=%s&access_token=%s&v=%s", groupID, token, apiVersion))
	if err != nil { //проверка на ошибки
		log.Println("Server request error:", err.Error())
		return "", "", "", err
	}
	defer resp.Body.Close() //закрываем тело по завершении функции
	var buff bytes.Buffer   //буфер для копирования тела ответа
	io.Copy(&buff, resp.Body)

	var dat map[string]interface{}
	json.Unmarshal(buff.Bytes(), &dat) //разборка на запчасти

	response := dat["response"].(map[string]interface{}) // получаем сам ответ

	secretKey = response["key"].(string)
	server = response["server"].(string)
	ts = response["ts"].(string)

	return secretKey, server, ts, nil
}
