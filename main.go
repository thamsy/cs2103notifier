package main

import (
	"bytes"
	"cs2103notifier/constants"
	"cs2103notifier/secret"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
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
	setUpLog()
	var counter = 0
	for true {
		counter++
		fmt.Println("Iteration: " + strconv.Itoa(counter))
		getCurrentWebsite()
		compareForUpdates()
		time.Sleep(30 * time.Minute)
		if r := recover(); r != nil {
			log.Println("Recover")
		}
	}
}

func setUpLog() {
	f, err := os.OpenFile(constants.GetLogDir(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(f)
	log.Println("Log starting...")
}

func postSlack(msg string) error {
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

func getCurrentWebsite() error {

	for _, filename := range secret.Cs2103Subfiles {
		// Get Current Website State
		resp, err := http.Get(secret.Cs2103Website + filename)
		if err != nil {
			postSlack(err.Error())
		}
		defer resp.Body.Close()

		// Write into file
		out, err := os.OpenFile(constants.GetCurrentDir(filename), os.O_WRONLY, 0666)
		if err != nil {
			postSlack(err.Error())
			//panic(err)
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			postSlack(err.Error())
			//panic(err)
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
			log.Println(filename + " changed at" + time.Now().String() + ":\n" + string(output))
			if len(output) < 200 {
				postSlack("*File Changed: " + filename + "*\n" + string(output))
			} else {
				postSlack("*File Changed: " + filename + "*")
			}
			curr, _ := os.Open(constants.GetCurrentDir(filename))
			prev, _ := os.Create(constants.GetPrevDir(filename))
			defer curr.Close()
			defer prev.Close()
			_, err := io.Copy(prev, curr)
			if err != nil {
				postSlack(err.Error())
				//panic(err)
			}
		}
	}
}
