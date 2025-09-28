package main

import (
	"fmt"
	"log"
	"net/http"
	"api/streaming"
	"api/video"
)

func main() {
	fmt.Println("Starting Backend")

	// Start the video processing pipeline.
	video.StartPipeline("video/sample/FireForce-S1E3-360P.mp4")

	// Handle WebSocket connections.
	http.HandleFunc("/ws", streaming.ServeWS)

	log.Println("server is running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}