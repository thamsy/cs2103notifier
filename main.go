package main

import (
	"bytes"
	"cs2103notifier/constants"
	"cs2103notifier/secret"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func init() {
	os.Mkdir(constants.GetStorageDir(), 0777)
	for _, filename := range secret.Cs2103Subfiles {
		if _, err := os.Stat(constants.GetCurrentDir(filename)); os.IsNotExist(err) {
			os.Create(constants.GetCurrentDir(filename))
		}
		if _, err := os.Stat(constants.GetPrevDir(filename)); os.IsNotExist(err) {
			os.Create(constants.GetPrevDir(filename))
		}
	}
}

func main() {
	postSlack(constants.START_MESSAGE)
	defer postSlack(constants.END_MESSAGE)
	for true {
		getCurrentWebsite()
		compareForUpdates()
		time.Sleep(5 * time.Second)
	}
}

func postSlack(msg string) error {
	var json = []byte(`{"username":"cs2103notifier",
    "icon_emoji": ":innocent:",
    "text":"` + msg + `"}`)
	req, _ := http.NewRequest("POST", secret.SlackBot, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	return nil
}

func getCurrentWebsite() error {
	for _, filename := range secret.Cs2103Subfiles {
		// Get Current Website State
		resp, err := http.Get(secret.Cs2103Website + filename)
		defer resp.Body.Close()
		if err != nil {
			postSlack(err.Error())
			panic(err)
		}

		// Write into file
		out, err := os.OpenFile(constants.GetCurrentDir(filename), os.O_WRONLY, 0666)
		defer out.Close()
		if err != nil {
			postSlack(err.Error())
			panic(err)
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			postSlack(err.Error())
			panic(err)
		}
	}
	return nil
}

func compareForUpdates() {
	for _, filename := range secret.Cs2103Subfiles {
		// Execute diff command
		cmd := exec.Command("diff", constants.GetCurrentDir(filename), constants.GetPrevDir(filename))
		output, _ := cmd.Output()

		// Post to Slack if there is diff, and copy current.txt to prev.txt
		if len(output) != 0 {
			if len(output) < 200 {
				postSlack("*File Changed: " + filename + "*\n" + string(output))
			} else {
				postSlack("File Changed: " + filename)
			}
			curr, _ := os.Open(constants.GetCurrentDir(filename))
			prev, _ := os.Create(constants.GetPrevDir(filename))
			defer curr.Close()
			defer prev.Close()
			_, err := io.Copy(prev, curr)
			if err != nil {
				postSlack(err.Error())
				panic(err)
			}
		}
	}
}
