package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"strings"
)

// Representa una solicitud HTTP recibida.
// Contiene el método, el objetivo (URL), las cabeceras y el cuerpo de la solicitud.
type HttpRequest struct {
	Method string            // Método HTTP (GET, POST, etc.)
	Target *url.URL          // URL objetivo de la solicitud
	Header map[string]string // Cabeceras HTTP como un mapa de clave-valor
	Body   string            // Cuerpo de la solicitud (si existe)
}

// Crea una nueva instancia de HttpRequest.
func NewHttpRequest(method string, target *url.URL, header map[string]string, body string) *HttpRequest {
	return &HttpRequest{
		Method: method,
		Target: target,
		Header: header,
		Body:   body,
	}
}

// lee una solicitud HTTP completa desde una conexión de red.
// Devuelve un puntero a HttpRequest o un error si ocurre algún problema.
func ReadRequest(conn net.Conn) (*HttpRequest, error) {
	lines := make([]string, 0)

	reader := bufio.NewReader(conn)

	// Lee las líneas de la cabecera hasta encontrar una línea vacía
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			// Si es fin de archivo (EOF), puede ser normal si la conexión se cierra
			if errors.Is(err, io.EOF) {
				break
			}

			// Otro error durante la lectura
			return nil, err
		}

		// Elimina los sufijos de retorno de carro y nueva línea
		line = strings.TrimSuffix(line, "\r\n")
		line = strings.TrimSuffix(line, "\n")

		// Una línea vacía indica el final de las cabeceras
		if line == "" {
			break
		}

		lines = append(lines, line)
	}

	// Si no se leyeron líneas, la solicitud está vacía
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty request")
	}

	// Une las líneas de la cabecera y añade el doble salto de línea final
	headersPart := strings.Join(lines, "\r\n") + "\r\n\r\n"

	// Parsea la parte de las cabeceras para obtener la estructura HttpRequest inicial
	request, err := ParseRequest(headersPart)
	if err != nil {
		return nil, fmt.Errorf("can't parse request: %w", err)
	}

	// Comprueba si existe la cabecera Content-Length para leer el cuerpo
	if contentLengthStr, ok := request.Header["Content-Length"]; ok {
		var contentLength int

		// Convierte el valor de Content-Length a entero
		_, err := fmt.Sscan(contentLengthStr, &contentLength)
		if err != nil {
			return nil, fmt.Errorf("bad content length format: %w", err)
		}

		// Si hay longitud de contenido, lee el cuerpo
		if contentLength > 0 {
			body := make([]byte, contentLength)

			// Lee exactamente contentLength bytes desde el reader
			_, err := io.ReadAtLeast(reader, body, contentLength)
			if err != nil {
				return nil, fmt.Errorf("can't read body: %w", err)
			}

			request.Body = string(body)
		}
	}

	return request, nil
}

// Analiza la parte de las cabeceras de una solicitud HTTP (como string).
// Devuelve un puntero a HttpRequest (sin el cuerpo) o un error.
func ParseRequest(headersPart string) (*HttpRequest, error) {
	// Encuentra el final de las cabeceras (doble salto de línea)
	headerEndIndex := strings.Index(headersPart, "\n\n")
	if headerEndIndex == -1 {
		return nil, fmt.Errorf("no header end")
	}

	// Extrae la parte de las cabeceras
	headerPart := headersPart[:headerEndIndex]

	// Divide las cabeceras en líneas individuales
	lines := strings.Split(headerPart, "\n")
	if len(lines) < 1 {
		return nil, fmt.Errorf("no start line")
	}

	// Parsea la línea de inicio (ej: "GET /path HTTP/1.1")
	start := strings.SplitN(lines[0], " ", 3)
	if len(start) < 2 {
		return nil, fmt.Errorf("no method or target")
	}

	// Extrae el método
	method := start[0]
	if method == "" {
		return nil, fmt.Errorf("no method")
	}

	// Extrae el target (URL)
	targetStr := start[1]
	if targetStr == "" {
		return nil, fmt.Errorf("no target")
	}

	// Parsea el target como una URL
	target, err := url.Parse(targetStr)
	if err != nil {
		return nil, fmt.Errorf("bad target format: %w", err)
	}

	// Crea un mapa para almacenar las cabeceras
	headers := make(map[string]string)

	// Procesa cada línea de cabecera (a partir de la segunda línea)
	for _, line := range lines[1:] {
		// Divide la línea en clave y valor por el primer ":"
		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {
			// Ignora líneas mal formadas
			continue
		}

		// Limpia espacios en blanco de clave y valor
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		if k == "" {
			// Ignora cabeceras con clave vacía
			continue
		}

		headers[k] = v
	}

	body := ""

	// Crea y devuelve el objeto HttpRequest con los datos parseados
	return NewHttpRequest(method, target, headers, body), nil
}
