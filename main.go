package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	file := NewFile()
	// Create application with options
	err := wails.Run(&options.App{
		Title:  "slack-wails",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		OnStartup:        app.startup,
		MinWidth:         1280,
		MinHeight:        768,
		// Frameless:        true,
		Bind: []interface{}{
			app,
			file,
		},
		Windows: &windows.Options{
			WebviewBrowserPath: "", // 可以让windows使用默认浏览器打开链接
		},
		// Mac: &mac.Options{
		// 	TitleBar: mac.TitleBarHiddenInset(),
		// },
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
