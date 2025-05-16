# Proyecto HTTP Server (Go)

Autores:
- Joselyn Jiménez
- Katerine Guzmán
- Esteban Solano
---

Este proyecto implementa un servidor HTTP/1.0 concurrente en Go desde cero, sin usar net/http, y ofrece múltiples endpoints para cálculo, manipulación de strings, creación/elimnación de archivos y pruebas de carga.

## Características
- Servidor HTTP concurrente simple.
- Enrutamiento basado en método y ruta.
- Manejo de query parameters.
- Graceful shutdown al recibir señales SIGINT o SIGTERM.

### Estructura del código
```
/ (raíz del repositorio)
├─ main.go                # Inicialización del servidor y registro de rutas
├─ go.mod/go.sum          # Módulo Go y dependencias
├─ core/                  # Núcleo del servidor: parsing, routing, servidor TCP
│  ├─ http_server.go      # Lógica de aceptación de conexiones y dispatch
│  ├─ http_request.go     # Parseo de solicitudes HTTP
│  ├─ http_response.go    # Construcción y envío de respuestas HTTP
│  └─ router.go           # Emparejamiento de rutas y métodos
├─ handlers/              # Endpoints básicos (reverse, toupper, hash, root)
│  ├─ string.go
│  └─ string_test.go
├─ service/               # Lógica de negocio: createfile, deletefile y validaciones
│  ├─ file_service.go
│  └─ file_service_validation_test.go
├─ advanced/              # Endpoints avanzados (random, timestamp, simulate, sleep, loadtest, status, help)
│  ├─ advanced_integration_test.go
│  └─ advanced.go         # Implementación de handlers avanzados
├─ integration/           # Tests raw TCP de integración (código 200, 400, 404)
│  ├─ router_error_test.go
│  └─ advanced_integration_test.go
└─ core/                  # Tests unitarios de core (ReadRequest, WriteResponse, etc.)
   ├─ http_request_additional_test.go
   └─ http_response_additional_test.go
```

### Dependencias
- Go 1.20+ (compatible con módulos)
- No utiliza bibliotecas externas; sólo paquetes estándar (net, bufio, fmt, os, crypto/sha256, etc.)
- Para pruebas y cobertura: go test, go tool cover

## Instrucciones de Uso

### Compilar y ejecutar
1. Clonar el repositorio:
```bash
git clone https://github.com/KateGF/Http-Server-Project-SO.git
cd Http-Server-Project-SO
```

2. Compilar
```bash
go build -o server.exe main.go
```

3. Ejecutar
```bash
./server.exe
```

El servidor escuchará en el puerto 8080. Verás en consola:
```
INFO Server started address=[::]:8080
```

4. Detener
- Presiona 'Ctrl+C' en la ventana donde corre el servidor.

### Probar funcionalidades con curl
```
# 1. /help
curl -i http://localhost:8080/help

# 2. /status
curl -i http://localhost:8080/status

# 3. /fibonacci
curl -i "http://localhost:8080/fibonacci?num=10"

# 4. /createfile y /deletefile (GET)
curl -i "http://localhost:8080/createfile?name=test.txt&content=hola&repeat=3"
type test.txt
curl -i "http://localhost:8080/deletefile?name=test.txt"

# 5. /reverse
curl -i "http://localhost:8080/reverse?text=abcdef"

# 6. /toupper
curl -i "http://localhost:8080/toupper?text=holaGo"

# 7. /hash
curl -i "http://localhost:8080/hash?text=abc123"

# 8. /random
curl -i "http://localhost:8080/random?count=5&min=1&max=100"

# 9. /timestamp
curl -i http://localhost:8080/timestamp

# 10. /sleep
curl -i "http://localhost:8080/sleep?seconds=2"

# 11. /simulate
curl -i "http://localhost:8080/simulate?seconds=3&task=miTarea"

# 12. /loadtest
curl -i "http://localhost:8080/loadtest?tasks=10&sleep=1"
```

### Pruebas de error
```
# Parámetros faltantes -> Bad Request
curl -i http://localhost:8080/fibonacci

# Método no soportado -> Bad Request
curl -i -X POST http://localhost:8080/reverse?text=hola

# Ruta no existente -> Not Found
curl -i http://localhost:8080/no_such_route
```

### Ejecutar pruebas y medir cobertura

1. Pruebas unitarias y de integración:
```bash
go test ./... -timeout 30s -coverprofile=coverage
```

2. Ver cobertura total:
```bash
(go tool cover -func=coverage) | Select-String total:
```
Muestra un porcentaje > 90%.
