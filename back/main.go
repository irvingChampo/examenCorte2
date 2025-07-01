package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"A_Lexico/back/analyzer"
)

// Estructura para deserializar el JSON del cliente
type Request struct {
	Code string `json:"code"`
}

// Middleware para habilitar CORS
func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*") // O usa "http://localhost:3000" si quieres restringir
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
}

// Manejador principal
func analyzeHandler(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	// Si es una petición preflight (OPTIONS), solo responde OK
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	log.Println("📨 Solicitud recibida")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("❌ Error leyendo body:", err)
		http.Error(w, "No se pudo leer el cuerpo", http.StatusBadRequest)
		return
	}

	if len(body) == 0 {
		log.Println("❌ Body vacío")
		http.Error(w, "El cuerpo de la solicitud está vacío", http.StatusBadRequest)
		return
	}

	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		log.Println("❌ Error unmarshal:", err)
		log.Println("🧾 Body recibido:", string(body))
		http.Error(w, "JSON inválido o malformado", http.StatusBadRequest)
		return
	}

	log.Println("📦 Código recibido:", req.Code)

	result := analyzer.AnalyzeCode(req.Code)
	fmt.Fprint(w, result)
}

// Función principal del servidor
func main() {
	http.HandleFunc("/analyze", analyzeHandler)
	log.Println("🚀 Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
