package video

import (
	"bytes"
	"image/jpeg"
	"log"
	"api/streaming"
)

// Start frame processing when package is initialized
func init() {
	go processFrames()
}

// processFrames continuously processes frames from FrameChan
func processFrames() {
	framesSent := 0
	log.Println("Frame processing started...")
	
	for frame := range FrameChan {
		// Try to compress the frame
		compressedFrame, err := compressFrame(frame)
		if err != nil {
			// If compression fails, send original frame
			log.Printf("Compression failed: %v, sending original", err)
			compressedFrame = frame
		}
		
		// Send to WebSocket clients
		streaming.Broadcast(compressedFrame)
		framesSent++
		
		if framesSent%50 == 0 {
			log.Printf("Sent %d frames to clients", framesSent)
		}
	}
	
	log.Println("Frame processing stopped")
}

// compressFrame attempts to compress a JPEG frame
func compressFrame(jpegData []byte) ([]byte, error) {
	// Validate that this looks like a JPEG
	if len(jpegData) < 4 {
		return jpegData, nil // Too small, return as-is
	}
	
	if jpegData[0] != 0xFF || jpegData[1] != 0xD8 {
		return jpegData, nil // Not a JPEG, return as-is
	}
	
	// Try to decode and re-encode with lower quality
	img, err := jpeg.Decode(bytes.NewReader(jpegData))
	if err != nil {
		return jpegData, err // Return original on decode error
	}
	
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 40}) // Higher quality than before
	if err != nil {
		return jpegData, err // Return original on encode error
	}
	
	return buf.Bytes(), nil
}
