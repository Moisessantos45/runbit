package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:            "Runbit",
		Width:            1200,
		Height:           720,
		MinWidth:         1024,
		MinHeight:        640,
		Frameless:        true,
		WindowStartState: options.Normal,
		CSSDragProperty:  "widows",
		CSSDragValue:     "1",
		BackgroundColour: &options.RGBA{R: 24, G: 26, B: 32, A: 255},

		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		OnStartup: app.startup,

		Windows: &windows.Options{
			Theme:        windows.SystemDefault,
			BackdropType: windows.Mica,
			DisableFramelessWindowDecorations: false,
		},
		Linux: &linux.Options{
			WindowIsTranslucent: true,
		},

		Bind: []any{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
