# API

This project is a real-time video streamer that uses ffmpeg to extract frames from a video file, quantizes them, and broadcasts them to connected clients using WebSockets.

## Architecture

The application is divided into three main packages:

- **`main`**: The entry point of the application. It starts the video processing pipeline and the WebSocket server.
- **`video`**: This package is responsible for reading the video file, decoding it into frames, and quantizing them.
- **`streaming`**: This package handles the WebSocket connections and broadcasts the frames to the clients.

## How to Run

1.  **Install dependencies**: Make sure you have Go and ffmpeg installed on your system.
2.  **Run the application**: Navigate to the `api` directory and run the following command:
    ```sh
    go run .
    ```
3.  **Connect a client**: Open a WebSocket client and connect to `ws://localhost:8080/ws` to receive the video stream.