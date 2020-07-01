package http_request

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

func SendJson(url string, jsonByte []byte) (respBody []byte, err error) {
	body := bytes.NewBuffer(jsonByte)
	client := &http.Client{}

	defer client.CloseIdleConnections()

	client.Timeout = 1 * time.Second
	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	respBody, err = ioutil.ReadAll(resp.Body)
	return
}
