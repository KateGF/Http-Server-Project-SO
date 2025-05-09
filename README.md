# Http-Server-Project-SO

Este proyecto implementa un servidor HTTP básico en Go. El servidor es capaz de manejar solicitudes.

## Características

- Servidor HTTP concurrente simple.
- Enrutamiento basado en método y ruta.
- Manejo de query parameters.
- Graceful shutdown al recibir señales SIGINT o SIGTERM.

## Instalación

Clona el repositorio:

```bash
git clone https://github.com/KateGF/Http-Server-Project-SO.git
cd Http-Server-Project-SO
```

## Uso

Para crear una instancia del servidor, definir una función de manejador, registrar el manejador e iniciar el servidor, puedes usar el siguiente código de ejemplo en Go:

```go
package main

import "log/slog"

func main() {
	server := NewHttpServer()

	server.Get("/", func(request *HttpRequest) (*HttpResponse, error) {
		name := request.Target.Query().Get("name")

		if name == "" {
			name = "World"
		}

		return Ok().Text("Hello, " + name + "!"), nil
	})

	err := server.Start(8080)

	if err != nil {
		slog.Error("Error starting or running server", "error", err)
	}
}
```

Para iniciar el servidor, ejecuta el siguiente comando en la raíz del proyecto:

```bash
go run .
```

Por defecto, el servidor se iniciará en el puerto `8080`. Verás un mensaje indicando que el servidor ha comenzado:

```
INFO Server started address=[::]:8080
```

Para detener el servidor, presiona `Ctrl+C`.

## Pruebas

Para ejecutar las pruebas unitarias:

```bash
go test
```

Para generar un informe de cobertura de pruebas:

```bash
make coverage
```

Para correr las pruebas de integración se usa el siguiente comando:

```bash
go test -timeout 30s -run Integration 
```

Esto ejecutará las pruebas, generará un archivo `coverage.out` y creará un informe HTML `coverage.html` que puedes abrir en tu navegador para ver la cobertura detallada.
