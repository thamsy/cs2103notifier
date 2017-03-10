package main

import (
	"cs2103notifier/constants"
	"cs2103notifier/secret"
	"fmt"
	"log"
	"os"
)

func initialize() {
	initStorage()
	initLog()
}

func initStorage() {
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

func initLog() {
	f, err := os.OpenFile(constants.GetLogDir(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}
	log.SetOutput(f)
	log.Println("Log starting...")
}
