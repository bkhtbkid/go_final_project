package server

import (
	"embed"
	"io/fs"
	"net/http"
)

func FsHandler(dir embed.FS) http.Handler {
	sub, err := fs.Sub(dir, "web")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(sub))
}
