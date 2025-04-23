package main

import "log/slog"

func main() {
	// Crea una nueva instancia del servidor HTTP.
	server := NewHttpServer()

	// Registra un manejador para la ruta GET "/fibonacci".
	server.Get("/fibonacci", FibonacciHandler)

	// Registra un manejador para la ruta POST "/createfile".
	server.Post("/createfile", CreateFileHandler)

	// Registra un manejador para la ruta DELETE "/deletefile".
	server.Delete("/deletefile", DeleteFileHandler)

	// Inicia el servidor en el puerto 8080.
	err := server.Start(8080)

	if err != nil {
		slog.Error("Error starting or running server", "error", err)
	}
}
