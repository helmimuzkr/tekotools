package main

import (
	"fmt"
	"time"
)

const (
	SYSTEM      = "SYSTEM"
	FORMAT_TIME = "2006-01-02 15:04:05"
)

func PrintStringErr(name string, pid int, err string) string {
	return fmt.Sprintf("%s - [%s | %d] : error: %s\n", time.Now().Format(FORMAT_TIME), name, pid, err)
}

func PrintErr(name string, pid int, err string) string {
	return fmt.Sprintf("%s - [%s | %d] : error: %s\n", time.Now().Format(FORMAT_TIME), name, pid, err)
}

func PrintLog(name string, pid int, messages ...string) {
	msg := messages[0]

	if len(messages) > 1 {
		for i := 1; i < len(messages); i++ {
			msg = msg + " | " + messages[i]
		}
	}

	fmt.Printf("%s - [%s | %d] :  %s\n", time.Now().Format(FORMAT_TIME), name, pid, msg)
}
