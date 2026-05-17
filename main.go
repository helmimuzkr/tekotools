package main

import (
	"flag"

	"tekotools/tekojar"
)

func main() {
	// go run . --settings /path/to/setting.json
	settingPath := flag.String("settings", "./mujar/settings.json", "path to settings file")
	flag.Parse()

	config, err := tekojar.LoadSetting(*settingPath)
	if err != nil {
		panic(err)
	}

	t := tekojar.New(config)

	go func() {
		for log := range t.WatchService("simple-loop.jar") {
			tekojar.PrintLog("UI", 0, log)
		}
	}()

	t.StartAll()
}
