package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"api/streaming"
	"api/video"
	_ "api/utils" // Import for stats initialization
)

// CORS middleware
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"video-streamer"}`)) 
}

func findVideoFile() string {
	// Try to find a sample video file
	possiblePaths := []string{
		"video/sample/FireForce-S1E3-360P.mp4",
		"sample.mp4",
		"test.mp4",
		"video.mp4",
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// If no video found, create a placeholder message
	log.Println("No video file found. Please place a video file in one of these locations:")
	for _, path := range possiblePaths {
		abs, _ := filepath.Abs(path)
		log.Printf("  - %s", abs)
	}
	return ""
}

func main() {
	fmt.Println("Video Streamer Backend")
	fmt.Println("======================")

	// Find and validate video file
	videoPath := findVideoFile()
	if videoPath == "" {
		log.Fatal("No video file found! Place a video file (sample.mp4, test.mp4, or video.mp4) in the current directory.")
	}

	fmt.Printf("Using video file: %s\n", videoPath)

	// Import video package to trigger frame processing initialization
	_ = video.FrameChan

	// Start video processing pipeline
	video.StartPipeline(videoPath)

	// Setup HTTP routes
	http.HandleFunc("/ws", enableCORS(streaming.ServeWS))
	http.HandleFunc("/health", enableCORS(healthCheck))

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("WebSocket: ws://localhost:8080/ws")
	fmt.Println("Health: http://localhost:8080/health")
	fmt.Println("======================")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
