package video

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

// FrameChan is a channel that receives video frames.
var FrameChan = make(chan []byte, 100)

// StartPipeline starts the video processing pipeline.
func StartPipeline(videoPath string) {
	go readAndDecode(videoPath)
}

// readAndDecode reads a video file, decodes it into frames, and sends them to the FrameChan.
func readAndDecode(videoPath string) {
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		log.Fatalf("video file not found: %s", videoPath)
	}

	log.Printf("Starting video processing for: %s", videoPath)

	// Use simpler FFmpeg command that works reliably
	cmd := exec.Command("ffmpeg",
		"-re", // Read at native frame rate
		"-i", videoPath,
		"-vf", "fps=10,scale=640:360", // 10 FPS and scale down
		"-c:v", "mjpeg", // Use MJPEG codec
		"-q:v", "8", // Medium quality for smaller size
		"-f", "mjpeg", // Output MJPEG format
		"pipe:1", // Output to stdout
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("ffmpeg stdout Pipe:", err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal("ffmpeg Start:", err)
	}

	log.Println("FFmpeg started successfully - processing frames...")

	// Process frames with a much simpler approach
	frameCount := 0
	buffer := make([]byte, 64*1024) // 64KB buffer
	var frameData []byte

	for {
		n, err := stdout.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Printf("Video processing completed. Total frames: %d", frameCount)
			} else {
				log.Printf("Read error: %v. Frames processed: %d", err, frameCount)
			}
			break
		}

		// Append new data to our buffer
		frameData = append(frameData, buffer[:n]...)

		// Extract complete JPEG frames
		for len(frameData) > 4 {
			// Look for JPEG start marker (FF D8)
			startPos := bytes.Index(frameData, []byte{0xFF, 0xD8})
			if startPos == -1 {
				// No start marker found, clear some old data
				if len(frameData) > 1024*1024 { // 1MB
					frameData = frameData[len(frameData)/2:]
				}
				break
			}

			// Look for JPEG end marker (FF D9) after start
			endPos := bytes.Index(frameData[startPos+2:], []byte{0xFF, 0xD9})
			if endPos == -1 {
				// No end marker yet, wait for more data
				break
			}

			// Calculate actual end position
			actualEndPos := startPos + 2 + endPos + 2 // include end marker

			// Extract the complete frame
			frame := make([]byte, actualEndPos-startPos)
			copy(frame, frameData[startPos:actualEndPos])

			// Try to send frame (non-blocking)
			select {
			case FrameChan <- frame:
				frameCount++
				if frameCount%30 == 0 {
					log.Printf("Processed %d frames (latest: %d bytes)", frameCount, len(frame))
				}
			default:
				// Channel is full, skip this frame
			}

			// Remove processed frame from buffer
			frameData = frameData[actualEndPos:]
		}

		// Prevent buffer from growing too large
		if len(frameData) > 5*1024*1024 { // 5MB max
			log.Println("Buffer too large, resetting...")
			frameData = nil
		}

		// Small delay to prevent CPU spinning
		time.Sleep(1 * time.Millisecond)
	}

	cmd.Wait()
	close(FrameChan)
	log.Println("Video processing pipeline closed")
}

