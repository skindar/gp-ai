package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const host = "localhost"
const port = 8080

func main() {
	frontendDir := "../frontend"
	fs := http.FileServer(http.Dir(frontendDir))
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds)

	// Handle requests to the root path
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var statusCode int
		if r.URL.Path == "/" {
			filePath := filepath.Join(frontendDir, "index.html")
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				http.NotFound(w, r)
				statusCode = http.StatusNotFound
			} else {
				http.ServeFile(w, r, filePath)
				statusCode = http.StatusOK
			}
		} else {
			// For other paths, use the file server
			originalWriter := w
			w = &responseWriter{ResponseWriter: originalWriter}
			fs.ServeHTTP(w, r)
			statusCode = w.(*responseWriter).statusCode
			if statusCode == 0 {
				statusCode = http.StatusOK
			}
		}

		// Log the completed request
		duration := time.Since(start)
		log.Printf("%s - %s %s - %d (%v)", r.RemoteAddr, r.Method, r.URL.Path, statusCode, duration)
	})

	// Start the server
	fmt.Printf("Server is running on http://%v:%d\n", host, port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
