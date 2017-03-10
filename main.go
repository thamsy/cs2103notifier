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

func main() {
	slack.SendMsg(constants.START_MESSAGE)
	defer slack.SendMsg(constants.END_MESSAGE)
	initialize()
	var counter = 0
	for true {
		counter++
		fmt.Println("Iteration: " + strconv.Itoa(counter))
		website.GetCurrentWebsite()
		website.CompareForUpdates()
		time.Sleep(30 * time.Minute)
		if r := recover(); r != nil {
			log.Println("Recover")
		}
	}
}
