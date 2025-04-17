package main

import (
	"fmt"
	"net"
	"sort"
	"strings"
)

const DEFAULT_BUFFER_SIZE = 1024

const MAX_REQUEST_SIZE = 4096

type Handle func(request *HttpRequest) (*HttpResponse, error)

type Handler struct {
	Method string
	Path   string
	Handle Handle
}

type HttpServer struct {
	Handlers []Handler
}

func NewHttpServer() *HttpServer {
	return &HttpServer{
		Handlers: []Handler{},
	}
}

func (server *HttpServer) AddHandler(method, path string, handle Handle) {
	handler := Handler{
		Method: method,
		Path:   path,
		Handle: handle,
	}

	server.Handlers = append(server.Handlers, handler)
}

func (server *HttpServer) Get(path string, handle Handle) {
	server.AddHandler("GET", path, handle)
}

func (server *HttpServer) Post(path string, handle Handle) {
	server.AddHandler("POST", path, handle)
}

func (server *HttpServer) SortHandlers() {
	sort.Slice(server.Handlers, func(i, j int) bool {
		iCount := strings.Count(server.Handlers[i].Path, "/")
		jCount := strings.Count(server.Handlers[j].Path, "/")

		return iCount > jCount
	})
}

func MatchPath(requestPath, handlerPath string) bool {
	if requestPath == handlerPath {
		return true
	}

	handlerPath += "/"

	// TODO: Check if requestPath ends with /
	fmt.Printf("requestPath: %s, handlerPath: %s\n", requestPath, handlerPath)

	return strings.HasPrefix(requestPath, handlerPath)
}

func (server *HttpServer) Start(port int) error {
	server.SortHandlers()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))

	if err != nil {
		return err
	}

	return nil
}

func (server *HttpServer) Handle(conn net.Conn) error {
	defer conn.Close()

	return nil
}
