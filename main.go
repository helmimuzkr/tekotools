package main

import (
	"flag"

	"mutill/mujar"
)

func main() {
	// go run . --settings /path/to/setting.json
	settingPath := flag.String("settings", "./mujar/settings.json", "path to settings file")
	flag.Parse()

	config, err := mujar.LoadSetting(*settingPath)
	if err != nil {
		panic(err)
	}

	m := mujar.New(config)

	go func() {
		for log := range m.WatchService("simple-loop.jar") {
			mujar.PrintLog("UI", 0, log)
		}
	}()

	m.StartAll()
}
