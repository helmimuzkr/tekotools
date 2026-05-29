package main

import (
	"embed"

	"tekotools/backend"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	tekojarApp := backend.NewTekojarApp()

	err := wails.Run(&options.App{
		Title:  "Tekotools",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  tekojarApp.Startup,
		OnShutdown: tekojarApp.Shutdown,
		Bind: []interface{}{
			tekojarApp,
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
