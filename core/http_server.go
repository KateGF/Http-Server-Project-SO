package core

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

// Define el tipo para las funciones que manejan las solicitudes HTTP.
// Recibe un puntero a HttpRequest y devuelve un puntero a HttpResponse y un error.
type Handle func(request *HttpRequest) (*HttpResponse, error)

// Representa un manejador para una ruta y método HTTP específicos.
type Handler struct {
	Method string // Método HTTP (ej. "GET", "POST")
	Path   string // Ruta de la URL (ej. "/users")
	Handle Handle // Función que manejará la solicitud
}

// Representa el servidor HTTP.
type HttpServer struct {
	Handlers []Handler    // Lista de manejadores registrados
	Listener net.Listener // Listener para aceptar conexiones
}

// Crea una nueva instancia de HttpServer.
func NewHttpServer() *HttpServer {
	return &HttpServer{
		Handlers: []Handler{},
	}
}

// Agrega un nuevo manejador al servidor.
func (server *HttpServer) AddHandler(method, path string, handle Handle) {
	handler := Handler{
		Method: method,
		Path:   path,
		Handle: handle,
	}

	server.Handlers = append(server.Handlers, handler)
}

// Un atajo para agregar un manejador para el método GET.
func (server *HttpServer) Get(path string, handle Handle) {
	server.AddHandler("GET", path, handle)
}

// Un atajo para agregar un manejador para el método POST.
func (server *HttpServer) Post(path string, handle Handle) {
	server.AddHandler("POST", path, handle)
}

// Un atajo para agregar un manejador para el método DELETE.
func (server *HttpServer) Delete(path string, handle Handle) {
	server.AddHandler("DELETE", path, handle)
}

// Ordena los manejadores por la especificidad de la ruta (más segmentos primero).
func (server *HttpServer) SortHandlers() {
	sort.Slice(server.Handlers, func(i, j int) bool {
		// Cuenta el número de '/' en cada ruta.
		iCount := strings.Count(server.Handlers[i].Path, "/")
		jCount := strings.Count(server.Handlers[j].Path, "/")

		if iCount != jCount {
			// Ordena de forma descendente por el número de segmentos.
			return iCount > jCount
		}

		iLen := len(server.Handlers[i].Path)
		jLen := len(server.Handlers[j].Path)

		return iLen > jLen
	})
}

// Verifica si la ruta de la solicitud coincide con la ruta del manejador.
// Permite coincidencias exactas o coincidencias de prefijo si la ruta del manejador termina en '/'.
func MatchPath(requestPath, handlerPath string) bool {
	// Coincidencia exacta.
	if requestPath == handlerPath {
		return true
	}

	// Añade '/' al final de la ruta del manejador para comprobar prefijos.
	handlerPath += "/"

	// Comprueba si la ruta de la solicitud comienza con la ruta del manejador (seguida de '/').
	return strings.HasPrefix(requestPath, handlerPath)
}

// Inicia el servidor HTTP en el puerto especificado.
func (server *HttpServer) Start(port int) error {
	// Ordena los manejadores antes de empezar a escuchar.
	server.SortHandlers()

	// Empieza a escuchar conexiones TCP en el puerto dado.
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	// Asigna el listener al servidor.
	server.Listener = ln

	// Canal para recibir señales del sistema operativo (SIGINT, SIGTERM).
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine para manejar el cierre ordenado del servidor.
	go func() {
		// Espera una señal de interrupción o terminación.
		<-sigCh
		fmt.Println()
		// Cierra el listener para detener la aceptación de nuevas conexiones.
		server.Stop()
	}()

	slog.Info("Server started", "address", ln.Addr().String())

	// Bucle principal para aceptar conexiones entrantes.
	for {
		// Acepta una nueva conexión.
		conn, err := ln.Accept()

		// Si el error es porque el listener fue cerrado, termina limpiamente.
		if errors.Is(err, net.ErrClosed) {
			slog.Info("Server stopped")
			return nil
		}

		// Si hay otro error al aceptar la conexión, lo devuelve.
		if err != nil {
			return err
		}

		// Maneja cada conexión en una goroutine separada.
		// HandleWithError se asegura de que los errores se registren.
		go server.HandleWithError(conn)
	}
}

// Detiene el servidor HTTP.
func (server *HttpServer) Stop() {
	if server.Listener != nil {
		server.Listener.Close()
	}
}

// Un envoltorio para Handle que registra cualquier error ocurrido durante el manejo de la conexión.
func (server *HttpServer) HandleWithError(conn net.Conn) {
	err := server.Handle(conn)
	if err != nil {
		slog.Error("Error", "error", err)
	}
}

// Maneja una conexión individual.
func (server *HttpServer) Handle(conn net.Conn) error {
	// Asegura que la conexión se cierre al final de la función.
	defer conn.Close()

	// Lee y parsea la solicitud HTTP de la conexión.
	request, err := ReadRequest(conn)
	if err != nil {
		// En lugar de cerrar sin responder, devolvemos 400 Bad Request con el mensaje de error
		resp := BadRequest().Text(err.Error())
		resp.WriteResponse(conn)
		return nil
	}

	slog.Info("Request", "address", conn.RemoteAddr().String(), "method", request.Method, "path", request.Target.Path)

	// Bandera para saber si la solicitud fue manejada.
	handled := false

	// Itera sobre los manejadores registrados (ya ordenados).
	for _, handler := range server.Handlers {
		// Comprueba si el método y la ruta coinciden.
		if request.Method != handler.Method || !MatchPath(request.Target.Path, handler.Path) {
			// Pasa al siguiente manejador si no coincide.
			continue
		}

		// Llama a la función Handle asociada al manejador.
		response, err := handler.Handle(request)
		if err != nil {
			// Si hay un error, crea una respuesta 500 Internal Server Error.
			response = &HttpResponse{
				StatusCode: 500,
				StatusText: "Internal Server Error",
				Headers:    map[string]string{},
				Body:       "500 Internal Server Error",
			}
		}

		// Escribe la respuesta generada por el manejador en la conexión.
		err = response.WriteResponse(conn)
		if err != nil {
			return err
		}

		// Marca la solicitud como manejada.
		handled = true
		// Deja de buscar manejadores una vez que uno coincide.
		break
	}

	// La conexión se manejó correctamente.
	if handled {
		return nil
	}

	// Si ningún manejador coincidió con la solicitud.
	// Crea una respuesta 404 Not Found.
	response := &HttpResponse{
		StatusCode: 404,
		StatusText: "Not Found",
		Headers:    map[string]string{},
		Body:       "404 Not Found",
	}

	// Escribe la respuesta 404 en la conexión.
	err = response.WriteResponse(conn)
	if err != nil {
		return err
	}

	// Se envió un 404.
	return nil
}
