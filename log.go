package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

var logging = false

// log out
func saveLog(platform string, channel string, info string) {
	// Current time
	t := time.Now().Local()
	// Build a timedate I like
	td := fmt.Sprintf("%04d/", t.Year())
	td += fmt.Sprintf("%02d/", t.Month())
	td += fmt.Sprintf("%02d ", t.Day())
	td += fmt.Sprintf("%02d:", t.Hour())
	td += fmt.Sprintf("%02d ", t.Minute())

	home, err := os.UserHomeDir()
	if err != nil {
		log.Println("Could not get home dir", err)
		return
	}
	// Create the directories
	mode := int(0755)
	home += "/logs"
	os.Mkdir(home, os.FileMode(mode))
	home += "/" + platform
	os.Mkdir(home, os.FileMode(mode))
	// Write to logfile
	file, err := os.OpenFile(home+"/"+channel+".log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Println("Failed opening file", err)
		return
	}
	defer file.Close()
	log.SetFlags(0)
	log.Println(td+platform, info)
	log.SetFlags(log.LstdFlags)
	w := bufio.NewWriter(file)
	w.WriteString(fmt.Sprintln(td+platform, info))
	w.Flush()
}
