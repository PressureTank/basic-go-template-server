// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package server

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"storj.io/qa-storj/static"
)

// ServeStatic will serve the requested static asset if possible. If served,
// will return true.
func (server *Server) ServeStatic(w http.ResponseWriter, r *http.Request) bool {
	data, err := static.FS.ReadFile(strings.TrimPrefix(r.URL.Path, "/"))
	if err != nil {
		return false
	}

	contentType := mime.TypeByExtension(filepath.Ext(r.URL.Path))
	if contentType == "" {
		contentType = http.DetectContentType(data)
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	_, err = w.Write(data)
	return err == nil
}
