package main

import (
	"httpserver/core"
	"httpserver/handlers"
	"httpserver/service"
	"log/slog"
)

func main() {
	// Crea una nueva instancia del servidor HTTP.
	server := core.NewHttpServer()

	// Registra un manejador para la ruta GET "/fibonacci".
	server.Get("/fibonacci", service.FibonacciHandler)

	// Registra un manejador para la ruta POST "/createfile".
	server.Post("/createfile", service.CreateFileHandler)

	// Registra un manejador para la ruta DELETE "/deletefile".
	server.Delete("/deletefile", service.DeleteFileHandler)
	// Endpoints de cadenas
	server.Get("/reverse", handlers.ReverseHandler)
	server.Get("/toupper", handlers.ToUpperHandler)
	server.Get("/hash", handlers.HashHandler)

	// Inicia el servidor en el puerto 8080.
	err := server.Start(8080)

	if err != nil {
		slog.Error("Error starting or running server", "error", err)
	}
}
