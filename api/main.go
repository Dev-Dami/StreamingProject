package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"video-streamer/streaming"
	"video-streamer/video"
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
		log.Printf("  - %s", filepath.Abs(path))
	}
	return ""
}

func main() {
	fmt.Println("🚀 Starting Video Streamer Backend")
	fmt.Println("=================================")

	// Find video file
	videoPath := findVideoFile()
	if videoPath != "" {
		fmt.Printf("📹 Using video file: %s\n", videoPath)
		video.StartPipeline(videoPath)
	} else {
		fmt.Println("⚠️  No video file found - WebSocket server will start but no video will stream")
	}

	// Setup routes
	http.HandleFunc("/ws", enableCORS(streaming.ServeWS))
	http.HandleFunc("/health", enableCORS(healthCheck))

	fmt.Println("🌐 Server starting on http://localhost:8080")
	fmt.Println("   - WebSocket endpoint: ws://localhost:8080/ws")
	fmt.Println("   - Health check: http://localhost:8080/health")
	fmt.Println("=================================")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
