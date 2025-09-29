package main


import (
	"fmt"
	"log"
	"net/http"
	"video-streamer/streaming"
	"video-streamer/video"
)

func main() {
	fmt.Println("Starting Backend")

	video.StartPipeline("video/sample/FireForce-S1E3-360P.mp4")

	http.HandleFunc("/ws", streaming.ServeWS)

	log.Println("server is running")
	log.Fatal(http.ListenAndServe(":8080", nil))
}