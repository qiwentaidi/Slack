package main

import (
	"context"
	"embed"

	rt "runtime"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
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
	winDropOptions := &options.DragAndDrop{
		EnableFileDrop:     true,
		DisableWebViewDrop: true,
	}
	macDropOptions := &options.DragAndDrop{
		EnableFileDrop: true,
	}
	err := wails.Run(&options.App{
		Title:  "Slack",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			file.startup(ctx)
			db.startup(ctx)
		},
		DragAndDrop: func() *options.DragAndDrop {
			if rt.GOOS == "windows" {
				return winDropOptions
			} else {
				return macDropOptions
			}
		}(),
		OnDomReady: func(ctx context.Context) {
			runtime.OnFileDrop(ctx, func(x, y int, paths []string) {

			})
		},
		MinWidth:  1280,
		MinHeight: 768,
		Bind: []interface{}{
			app,
			file,
			db,
		},
		Mac: &mac.Options{
			TitleBar:            mac.TitleBarHiddenInset(),
			WindowIsTranslucent: true,
		},
		Windows: &windows.Options{
			WebviewBrowserPath: "", // 可以让windows使用默认浏览器打开链接
		},
		Frameless: rt.GOOS != "darwin",
	})
	if err != nil {
		println("Error:", err.Error())
	}
}
