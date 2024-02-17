package assets

import (
	"embed"
	"net/http"
)

//go:embed assets/*
var assets embed.FS

func GetServer() http.Handler {
	return http.FileServer(http.FS(assets))
}
