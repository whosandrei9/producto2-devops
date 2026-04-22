package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Respuesta struct {
	Mensaje string `json:"mensaje"`
	Alumno  string `json:"alumno"`
	Entorno string `json:"entorno"`
	Pod     string `json:"pod"`
}

func main() {
	http.HandleFunc("/", inicioHandler)
	http.HandleFunc("/saludo", saludoHandler)
	http.HandleFunc("/api/mensaje", apiMensajeHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/info", infoHandler)
	http.HandleFunc("/imagen", imagenHandler)

	fmt.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func obtenerEntornoYPod() (string, string) {
	entorno := os.Getenv("APP_ENV")
	if entorno == "" {
		entorno = "local"
	}

	pod, err := os.Hostname()
	if err != nil {
		pod = "desconocido"
	}

	return entorno, pod
}

func inicioHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	entorno, pod := obtenerEntornoYPod()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `<!DOCTYPE html>
<html lang="es">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Producto 2 - DevOps</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			background-color: #f3f4f6;
			margin: 0;
			padding: 40px;
			text-align: center;
		}
		.contenedor {
			max-width: 900px;
			margin: 0 auto;
			background: white;
			padding: 30px;
			border-radius: 14px;
			box-shadow: 0 4px 12px rgba(0,0,0,0.1);
		}
		h1 {
			color: #0b5cab;
		}
		img {
			max-width: 320px;
			margin-top: 20px;
			border-radius: 10px;
		}
		.datos {
			margin-top: 20px;
			padding: 15px;
			background: #eef4ff;
			border-radius: 10px;
		}
		.enlaces {
			margin-top: 25px;
		}
		.enlaces a {
			display: block;
			margin: 8px 0;
			color: #0b5cab;
			text-decoration: none;
			font-weight: bold;
		}
	</style>
</head>
<body>
	<div class="contenedor">
		<h1>Soy alumno de la UOC</h1>
		<p>Aplicación desplegada con Jenkins y Kubernetes para el Producto 2 - rama test.</p>

		<div class="datos">
			<p><strong>Entorno:</strong> %s</p>
			<p><strong>Pod:</strong> %s</p>
		</div>

		<img src="/imagen" alt="Imagen generada por la aplicación">

		<div class="enlaces">
			<a href="/saludo">/saludo</a>
			<a href="/api/mensaje">/api/mensaje</a>
			<a href="/health">/health</a>
			<a href="/info">/info</a>
		</div>
	</div>
</body>
</html>`, entorno, pod)
}

func saludoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "Soy alumno de la UOC")
}

func apiMensajeHandler(w http.ResponseWriter, r *http.Request) {
	entorno, pod := obtenerEntornoYPod()

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	respuesta := Respuesta{
		Mensaje: "Soy alumno de la UOC",
		Alumno:  "Andrei",
		Entorno: entorno,
		Pod:     pod,
	}

	json.NewEncoder(w).Encode(respuesta)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	entorno, pod := obtenerEntornoYPod()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "entorno=%s pod=%s", entorno, pod)
}

func imagenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")

	svg := `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" width="500" height="250" viewBox="0 0 500 250">
	<rect width="500" height="250" fill="#0b5cab"/>
	<circle cx="110" cy="125" r="65" fill="#ffffff"/>
	<text x="300" y="100" font-size="34" text-anchor="middle" fill="#ffffff" font-family="Arial">UOC</text>
	<text x="300" y="145" font-size="22" text-anchor="middle" fill="#ffffff" font-family="Arial">Producto 2</text>
	<text x="300" y="180" font-size="18" text-anchor="middle" fill="#ffffff" font-family="Arial">Jenkins + Kubernetes</text>
</svg>`

	fmt.Fprint(w, svg)
}
