package tekojar

import (
	"flag"
)

func main() {
	// go run . --settings /path/to/setting.json
	settingPath := flag.String("settings", "./mujar/settings.json", "path to settings file")
	flag.Parse()

	config, err := LoadSetting(*settingPath)
	if err != nil {
		panic(err)
	}

	t := New(config)

	go func() {
		for log := range t.WatchService("simple-loop.jar") {
			PrintLog("UI", 0, log)
		}
	}()

	t.StartAll()
}
