package jsonedit

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets/*
var embeddedFiles embed.FS

// GetEmbeddedFileSystem returns a http.FileSystem that serves files from the embedded filesystem
func GetEmbeddedFileSystem() http.FileSystem {
	// Create a sub-filesystem for the assets directory
	assetsFS, err := fs.Sub(embeddedFiles, "assets")
	if err != nil {
		panic(err)
	}
	return http.FS(assetsFS)
}
