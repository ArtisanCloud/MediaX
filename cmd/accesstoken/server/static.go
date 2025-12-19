package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/*
var debugStaticFS embed.FS

func (s *accessTokenServer) staticFileHandler() http.Handler {
	sub, err := fs.Sub(debugStaticFS, "static")
	if err != nil {
		return http.NotFoundHandler()
	}
	return http.FileServer(http.FS(sub))
}
