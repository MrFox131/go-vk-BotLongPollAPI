package vkapi

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

//VkAPI using method with args and returns server answer
func VkAPI(method string, args map[string]string) string {
	address := fmt.Sprintf("https://api.vk.com/method/%s?access_token=%s&v=%s", method, token, apiVersion)

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
