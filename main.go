package main

import (
	"cs2103notifier/constants"
	"cs2103notifier/slack"
	"cs2103notifier/website"
	"fmt"
	"log"
	"strconv"
	"time"
)

func init() {
	initialize()
	slack.SendMsg(constants.START_MESSAGE)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recover")
			main()
		}
		slack.SendMsg(constants.END_MESSAGE)
	}()

	var counter = 0
	for true {
		counter++
		fmt.Println("Iteration: " + strconv.Itoa(counter))
		website.GetCurrentWebsite()
		website.CompareForUpdates()
		time.Sleep(30 * time.Minute)
	}
}
