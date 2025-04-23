package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
)

// Crea un archivo en el directorio actual.
// - Solo puede crear archivos en el directorio actual.
// - Crea el directorio y subdirectorios si no existen.
// - No puede crear un archivo si ya existe.
// - Crea el archivo con el contenido repetido el número de veces especificado.
// - Crea el archivo con el nombre especificado.
func CreateFile(filename string, content string, repeat int) error {
	// Obtener el directorio actual
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Construir la ruta completa del archivo
	path := filepath.Join(wd, filename)
	// Verificar si el archivo ya existe
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("file already exists: %s", path)
	}

	// Crear todos los directorios necesarios en la ruta
	err = os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	// Crear el archivo
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Escribir el contenido repetido en el archivo usando un buffer
	writer := bufio.NewWriter(file)
	for range repeat {
		_, err := writer.WriteString(content)
		if err != nil {
			return err
		}
	}
	writer.Flush()

	return nil
}

// Elimina un archivo en el directorio actual.
// - Solo puede eliminar archivos en el directorio actual.
// - No puede eliminar un archivo si no existe.
// - No puede eliminar directorios no vacíos.
// - Elimina el archivo con el nombre especificado.
func DeleteFile(filename string) error {
	// Obtener el directorio de trabajo actual
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Construir la ruta completa del archivo
	path := filepath.Join(wd, filename)
	// Verificar si el archivo existe
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", path)
	}

	// Eliminar el archivo
	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

// Maneja las solicitudes HTTP para crear archivos.
// Extrae los parámetros 'name', 'content' y 'repeat' de la consulta y valida los parámetros.
// Devuelve una respuesta HTTP indicando éxito o error.
func CreateFileHandler(request *HttpRequest) (*HttpResponse, error) {
	// Obtener parámetros de la consulta
	name := request.Target.Query().Get("name")
	if name == "" {
		return BadRequest().Text("name is required"), nil
	}

	content := request.Target.Query().Get("content")
	if content == "" {
		return BadRequest().Text("content is required"), nil
	}

	repeatStr := request.Target.Query().Get("repeat")
	if repeatStr == "" {
		return BadRequest().Text("repeat is required"), nil
	}

	// Convertir 'repeat' a entero
	repeat, err := strconv.Atoi(repeatStr)
	if err != nil {
		return BadRequest().Text("repeat must be a number"), nil
	}

	// Validar que 'repeat' sea positivo
	if repeat < 1 {
		return BadRequest().Text("repeat must be greater than 0"), nil
	}

	// Llamar a la función para crear el archivo
	err = CreateFile(name, content, repeat)
	if err != nil {
		// Registrar el error y devolver una respuesta de error interno del servidor
		slog.Error("Error creating file", "error", err)
		return NewHttpResponse(500, "Internal Server Error", "Error creating file"), nil
	}

	// Devolver una respuesta de éxito
	return Ok().Text("File created successfully"), nil
}

// Maneja las solicitudes HTTP para eliminar archivos.
// Extrae el parámetro 'name' de la consulta y valida el parámetro.
// Devuelve una respuesta HTTP indicando éxito o error.
func DeleteFileHandler(request *HttpRequest) (*HttpResponse, error) {
	// Obtener el parámetro 'name' de la consulta
	name := request.Target.Query().Get("name")
	if name == "" {
		return BadRequest().Text("name is required"), nil
	}

	// Llamar a la función para eliminar el archivo
	err := DeleteFile(name)
	if err != nil {
		// Registrar el error y devolver una respuesta de error interno del servidor
		slog.Error("Error deleting file", "error", err)
		return NewHttpResponse(500, "Internal Server Error", "Error deleting file"), nil
	}

	// Devolver una respuesta de éxito
	return Ok().Text("File deleted successfully"), nil
}
