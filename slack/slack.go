package slack

import (
	"bytes"
	"cs2103notifier/secret"
	"net/http"
)

func SendMsg(msg string) error {
	var urls = []string{secret.SlackBot, secret.SlackBot2}
	for _, url := range urls {
		err := postSlackToURL(msg, url)
		if err != nil {
			return err
		}
	}
	return nil
}

func postSlackToURL(msg string, url string) error {
	var json = []byte(`{"username":"cs2103notifier",
    "icon_emoji": ":innocent:",
    "text":"` + msg + `"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
