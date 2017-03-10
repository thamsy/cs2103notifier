package website

import (
	"cs2103notifier/constants"
	"cs2103notifier/secret"
	"cs2103notifier/slack"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func GetCurrentWebsite() error {
	for _, filename := range secret.Cs2103Subfiles {
		// Get Current Website State
		resp, err := http.Get(secret.Cs2103Website + filename)
		if err != nil {
			slack.SendMsg(err.Error())
		}
		defer resp.Body.Close()

		// Write into file
		out, err := os.OpenFile(constants.GetCurrentDir(filename), os.O_WRONLY, 0666)
		if err != nil {
			slack.SendMsg(err.Error())
			//panic(err)
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			slack.SendMsg(err.Error())
			//panic(err)
		}
	}
	return nil
}

func CompareForUpdates() {
	for _, filename := range secret.Cs2103Subfiles {
		// Execute diff command
		cmd := exec.Command("diff", constants.GetCurrentDir(filename), constants.GetPrevDir(filename))
		output, _ := cmd.Output()

		// Post to Slack if there is diff, and copy current.txt to prev.txt
		if len(output) != 0 {
			log.Println(filename + " changed at" + time.Now().String() + ":\n" + string(output))
			if len(output) < 200 {
				slack.SendMsg("*File Changed: " + filename + "*\n" + string(output))
			} else {
				slack.SendMsg("*File Changed: " + filename + "*")
			}
			curr, _ := os.Open(constants.GetCurrentDir(filename))
			prev, _ := os.Create(constants.GetPrevDir(filename))
			defer curr.Close()
			defer prev.Close()
			_, err := io.Copy(prev, curr)
			if err != nil {
				slack.SendMsg(err.Error())
				//panic(err)
			}
		}
	}
}
