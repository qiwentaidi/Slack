package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	file := NewFile()
	db := NewDatabase()
	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Slack",
		Width:  1280,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			file.startup(ctx)
		},
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop: true,
		},
		OnDomReady: func(ctx context.Context) {
			runtime.OnFileDrop(ctx, func(x, y int, paths []string) {
				runtime.EventsEmit(ctx, "wails-drop", paths)
			})
		},
		MinWidth:  1280,
		MinHeight: 768,
		// Frameless:        true,
		Bind: []interface{}{
			app,
			file,
			db,
		},
		Windows: &windows.Options{
			WebviewBrowserPath: "", // 可以让windows使用默认浏览器打开链接
		},
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
