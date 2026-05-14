package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("start application")

	// go run . --config /path/to/config.yml
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	config, err := LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	mutill := Mutill{}
	mutill.RegisterService(config)
	mutill.StartAll()

	go func() {
		for log := range mutill.ReadLog() {
			fmt.Println(log)
		}
	}()

	// listen for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	mutill.StopAll()
	fmt.Println("stop application")
}
