package main

import (
	"embed"
	"github.com/wyy-go/go-web-template/internal/routers"
)

//go:embed views
var Views embed.FS

//go:embed static
var Static embed.FS

//go:embed static/favicon.ico
var Favicon []byte

func init() {
	routers.Views = Views
	routers.Static = Static
	routers.Favicon = Favicon
}