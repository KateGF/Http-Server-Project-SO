package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
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

	return strings.HasPrefix(requestPath, handlerPath)
}

func (server *HttpServer) Start(port int) error {
	server.SortHandlers()

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println()
		ln.Close()
	}()

	slog.Info("Server started", "address", ln.Addr().String())

	for {
		conn, err := ln.Accept()

		if errors.Is(err, net.ErrClosed) {
			slog.Info("Server stopped")
			return nil
		}

		if err != nil {
			return err
		}

		go server.HandleWithError(conn)
	}
}

func (server *HttpServer) HandleWithError(conn net.Conn) {
	err := server.Handle(conn)
	if err != nil {
		slog.Error("Error", "error", err)
	}
}

func (server *HttpServer) Handle(conn net.Conn) error {
	defer conn.Close()

	request, err := ReadRequest(conn)
	if err != nil {
		return err
	}

	slog.Info("Request", "address", conn.RemoteAddr().String(), "method", request.Method, "path", request.Target.Path)

	handled := false

	for _, handler := range server.Handlers {
		if request.Method != handler.Method || !MatchPath(request.Target.Path, handler.Path) {
			continue
		}

		response, err := handler.Handle(request)
		if err != nil {
			return err
		}

		err = response.WriteResponse(conn)
		if err != nil {
			return err
		}

		handled = true
		break
	}

	if handled {
		return nil
	}

	response := &HttpResponse{
		StatusCode: 404,
		StatusText: "Not Found",
		Headers:    map[string]string{},
		Body:       "404 Not Found",
	}

	err = response.WriteResponse(conn)
	if err != nil {
		return err
	}

	return nil
}
