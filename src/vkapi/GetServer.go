package vkapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

/*GetBotLongPollServer returns secret key, server and ts for vkAPI services*/
func GetBotLongPollServer(userToken, groupID, userAPIVersion string, commandToken ...string) (secretKey, server, ts string, custErr error) {
	token = userToken
	apiVersion = userAPIVersion
	//Запрос к апи
	resp, err := http.Get(fmt.Sprintf("https://api.vk.com/method/groups.getLongPollServer?group_id=%s&access_token=%s&v=%s", groupID, token, apiVersion))
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
