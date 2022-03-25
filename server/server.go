// Copyright (C) 2021 Storj Labs, Inc.
// See LICENSE for copying information.

package server

import (
	"embed"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/spacemonkeygo/monkit/v3"
	"github.com/zeebo/errs"
	"go.uber.org/zap"
)

var mon = monkit.Package()

// Config defines the configuration for the server.
type Config struct {
	ListenAddr string `default:":8080" help:"port to listen on"`
	JSONPath   string `default:"./TrivyContainerCheckResults.json"`
}

// Server defines the necessary components to host the proxy server.
type Server struct {
	log      *zap.Logger
	cfg      Config
	staticFS embed.FS
	index    *template.Template
}

// New creates a new server based on the provided config and static directory.
func New(log *zap.Logger, cfg Config, staticFS embed.FS) (s *Server, err error) {
	s = &Server{
		log:      log,
		cfg:      cfg,
		staticFS: staticFS,
	}

	// TODO add config for index.html
	s.index, err = template.ParseFiles(filepath.Join("./static/index.html"))
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Close closes the server.
func (server *Server) Close() error {
	return nil
}

// Serve starts the server.
func (server *Server) Serve() error {

	// construct a standard go http handler
	handler := http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/index.html" && req.Method == http.MethodGet {
				server.log.Debug("index page")
				server.handleIndex(w, req)
				return
			}
			if served := server.ServeStatic(w, req); served {
				return
			}

			server.serveJSONError(w, http.StatusNotFound, errors.New("page not found"))
		})

	return http.ListenAndServe(server.cfg.ListenAddr, handler)
}

type ContainerCheckResults struct {
	SchemaVersion int
	ArtifactName  string
	Metadata      Metadata
}

type Metadata struct {
	OS OS
}

type OS struct {
	Family string
	Name   string
}

func (server *Server) handleIndex(w http.ResponseWriter, req *http.Request) {
	// open JSON file

	data := ContainerCheckResults{
		SchemaVersion: 2,
		ArtifactName:  "test artifact name",
		Metadata: Metadata{
			OS: OS{
				Family: "fmaily test name",
				Name:   "test name",
			},
		},
	}
	// send template
	err := server.index.Execute(w, data)
	if err != nil {
		server.log.Error("index template could not be executed", zap.Error(err))
		return
	}

}

func (server *Server) logDebug(req *http.Request) {
	server.log.Debug("request",
		zap.String("host", req.Host),
		zap.String("hostname", req.URL.Hostname()),
		zap.String("url", req.URL.String()),
		zap.String("referer", req.Referer()),
		zap.String("user agent", req.UserAgent()),
		zap.String("forwarded for", req.Header.Get("X-Forwarded-For")))
}

// serveJSONError writes JSON error to response output stream.
func (server *Server) serveJSONError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	var response struct {
		Error string `json:"error"`
	}

	response.Error = err.Error()

	server.log.Debug("sending json error to client", zap.Error(err))

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		server.log.Error("failed to write json error response", zap.Error(errs.Wrap(err)))
	}
}
