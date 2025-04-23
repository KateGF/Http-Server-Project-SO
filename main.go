package main

import "log/slog"

func main() {
	server := NewHttpServer()

	server.Get("/", func(request *HttpRequest) (*HttpResponse, error) {
		return Ok().Text("!"), nil
	})

	err := server.Start(8080)

	if err != nil {
		slog.Error("Error starting or running server", "error", err)
	}
}
