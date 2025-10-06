package video

import (
	"bufio"
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"os"
)

// SaveFramesToFile saves video frames to a file for debugging
func SaveFramesToFile(filename string, maxFrames int) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filename, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	frameCount := 0
	log.Printf("Starting to save frames to %s (max %d frames)", filename, maxFrames)

	// Read frames from the channel and save them
	for frame := range FrameChan {
		if frameCount >= maxFrames {
			log.Printf("Reached maximum frame limit (%d), stopping save", maxFrames)
			break
		}

		// Write frame length as header (for parsing later)
		lengthBytes := make([]byte, 4)
		lengthBytes[0] = byte(len(frame) >> 24)
		lengthBytes[1] = byte(len(frame) >> 16)
		lengthBytes[2] = byte(len(frame) >> 8)
		lengthBytes[3] = byte(len(frame))

		if _, err := writer.Write(lengthBytes); err != nil {
			return fmt.Errorf("failed to write frame length: %v", err)
		}

		// Write frame data
		if _, err := writer.Write(frame); err != nil {
			return fmt.Errorf("failed to write frame data: %v", err)
		}

		frameCount++
		if frameCount%10 == 0 {
			log.Printf("Saved %d frames to file", frameCount)
		}
	}

	log.Printf("Finished saving %d frames to %s", frameCount, filename)
	return nil
}

// ValidateJPEGFrame checks if a frame is a valid JPEG
func ValidateJPEGFrame(frame []byte) bool {
	if len(frame) < 4 {
		return false
	}

	// Check for JPEG magic bytes (FF D8 at start, FF D9 at end)
	if frame[0] != 0xFF || frame[1] != 0xD8 {
		return false
	}

	if frame[len(frame)-2] != 0xFF || frame[len(frame)-1] != 0xD9 {
		return false
	}

	// Try to decode it
	_, err := jpeg.Decode(bytes.NewReader(frame))
	return err == nil
}

// GetFrameInfo returns basic info about a JPEG frame
func GetFrameInfo(frame []byte) (width, height int, size int) {
	size = len(frame)
	if !ValidateJPEGFrame(frame) {
		return 0, 0, size
	}

	img, err := jpeg.Decode(bytes.NewReader(frame))
	if err != nil {
		return 0, 0, size
	}

	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), size
}
