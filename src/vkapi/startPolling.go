package vkapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/tidwall/gjson"
)

/*StartPolling start getting events from server*/
func StartPolling(secretKey, server, ts string) {
	events := make(chan []gjson.Result)
	go updatesHandling(events)
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
