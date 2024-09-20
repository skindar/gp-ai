package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const host = "localhost"
const port = 8080

func main() {
	frontendDir := "../frontend"
	audioDir := "audio-idle"
	fs := http.FileServer(http.Dir(frontendDir))
	audioFs := http.FileServer(http.Dir(audioDir))
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

	// Serve static files from the audio directory
	http.Handle("/audio-idle/", http.StripPrefix("/audio-idle/", audioFs))

	// Handle requests to /getLatestAudio
	http.HandleFunc("/getLatestAudio", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		files, err := filepath.Glob(filepath.Join(audioDir, "*"))
		if err != nil || len(files) == 0 {
			http.Error(w, "No audio files found", http.StatusNotFound)
			duration := time.Since(start)
			log.Printf("%s - %s %s - %d (%v) - No audio files found", r.RemoteAddr, r.Method, r.URL.Path, http.StatusNotFound, duration)
			return
		}

		// Log the list of files
		log.Printf("Available audio files: %v", files)

		// Create a new random generator
		randGen := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomIndex := randGen.Intn(len(files))
		randomFile := files[randomIndex]

		// Log the selected random file
		log.Printf("Selected random audio file: %s", randomFile)

		// Create the audio URL
		audioURL := fmt.Sprintf("/audio-idle/%s", filepath.Base(randomFile))

		// Create response JSON
		response := map[string]string{
			"url": audioURL,
		}

		// Set response headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Encode and write JSON response
		json.NewEncoder(w).Encode(response)

		duration := time.Since(start)
		log.Printf("%s - %s %s - %d (%v) - Returned audio URL: %s", r.RemoteAddr, r.Method, r.URL.Path, http.StatusOK, duration, audioURL)
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
