package assets

import (
	"embed"
	"io/fs"
)

//go:embed static/*
var embeddedStatic embed.FS

var StaticFS fs.FS

func init() {
	sub, err := fs.Sub(embeddedStatic, "static")
	if err != nil {
		panic("static assets not found: " + err.Error())
	}
	StaticFS = sub
}
