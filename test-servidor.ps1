# test-servidor.ps1

$baseUrl = "http://localhost:8080"

function Test-Endpoint($description, $url) {
    Write-Host "==== $description ====" -ForegroundColor Cyan
    curl.exe --http1.0 $url
    Write-Host "`n"
}

Test-Endpoint "1. Reverse (texto invertido)" "$baseUrl/reverse?text=PriscilaJS"
Test-Endpoint "2. ToUpper (uppercase text)" "$baseUrl/toupper?text=HolaMundo"
Test-Endpoint "3. Hash (SHA-256)" "$baseUrl/hash?text=abc"
Test-Endpoint "4. Timestamp (hora actual)" "$baseUrl/timestamp"
Test-Endpoint "5. Random (números aleatorios)" "$baseUrl/random?count=5&min=10&max=100"
Test-Endpoint "6. Simulate (tarea de 2 segundos)" "$baseUrl/simulate?seconds=2&task=test"
Test-Endpoint "7. Sleep (espera 1 segundo)" "$baseUrl/sleep?seconds=1"
Test-Endpoint "8. LoadTest (3 tareas de 1 segundo)" "$baseUrl/loadtest?tasks=3&sleep=1"
Test-Endpoint "9. Status del servidor" "$baseUrl/status"
Test-Endpoint "10. Ayuda (help)" "$baseUrl/help"
Test-Endpoint "11. Ruta raíz /" "$baseUrl/"
