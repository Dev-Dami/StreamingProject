package main


import (
	"fmt"
	"log"
	"net/http"
	"api/streaming"
	"api/video"
)

func main () {
	fmt.println("Starting Backend")

	video.StartPipeline("sample/input.mp4")

	http.HandleFunc("/ws", Streaming.ServeWS)

	log.println("server is running")
	log.Fatal(http.ListenAndServe(":8080, nil"))
}