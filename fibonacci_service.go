package main

import "strconv"

// Almacena los números de Fibonacci ya calculados para evitar recálculos.
var memo = map[int]int{}

// Calcula el n-ésimo número de Fibonacci utilizando recursión con memoización.
func Fibonacci(num int) int {
	// Casos base de la recursión.
	if num <= 0 {
		return 0
	}
	if num == 1 {
		return 1
	}

	// Verifica si el resultado ya está en la caché (memoización).
	if v, exists := memo[num]; exists {
		return v
	}

	// Si no está en caché, calcula recursivamente, almacena en caché y retorna.
	memo[num] = Fibonacci(num-1) + Fibonacci(num-2)

	return memo[num]
}

// Extrae el parámetro 'num' de la consulta, calcula el número de Fibonacci correspondiente y retorna la respuesta HTTP.
func FibonacciHandler(request *HttpRequest) (*HttpResponse, error) {
	// Obtiene el valor del parámetro 'num' de la URL query.
	numStr := request.Target.Query().Get("num")

	// Valida si el parámetro 'num' está presente.
	if numStr == "" {
		return BadRequest().Text("num is required"), nil
	}

	// Convierte el parámetro 'num' de string a entero.
	num, err := strconv.Atoi(numStr)

	// Valida si la conversión fue exitosa.
	if err != nil {
		return BadRequest().Text("num must be a number"), nil
	}

	// Valida si el número está dentro del rango permitido (0 a 92).
	// El límite 92 se debe a que Fibonacci(93) excede el máximo valor de int64.
	if num < 0 || num > 92 {
		return BadRequest().Text("num must be between 0 and 92"), nil
	}

	// Calcula el número de Fibonacci usando la función optimizada.
	v := Fibonacci(num)

	// Crea una respuesta HTTP 200 OK con el resultado como texto plano.
	return Ok().Text(strconv.Itoa(v)), nil
}
