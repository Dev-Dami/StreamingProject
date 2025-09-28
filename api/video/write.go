package video

import (
	"bufio"
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"os"
)

var QuantizedFrameChan = make(chan []byte, 30)

// QuantizeFrame takes a raw JPEG/PNG frame (or raw bytes), re-encodes it with
// lower quality, and sends it into the QuantizedFrameChan.
func QuantizeFrame(frame []byte, quality int) {
	// Decode frame
	img, _, err := image.Decode(bytes.NewReader(frame))
	if err != nil {
		log.Printf("QuantizeFrame: failed to decode frame: %v", err)
		return
	}

	// Re-encode with lower quality
	var buf bytes.Buffer
	options := &jpeg.Options{Quality: quality}
	if err := jpeg.Encode(&buf, img, options); err != nil {
		log.Printf("QuantizeFrame: failed to encode frame: %v", err)
		return
	}

	// Send quantized frame to channel (non-blocking fallback)
	select {
	case QuantizedFrameChan <- buf.Bytes():
	default:
		log.Println("QuantizeFrame: channel full, dropping frame")
	}
}

// SaveToDAT writes all quantized frames from the channel into a .dat file.
func SaveToDAT(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for frame := range QuantizedFrameChan {
		if _, err := writer.Write(frame); err != nil {
			return err
		}
	}
	return nil
}

// ReconstructMP4 will eventually split .dat into JPEGs and rebuild into MP4.
func ReconstructMP4(datFile, mp4File string) {
	// TODO: implement splitting .dat into JPEG frames then use ffmpeg/libx264
	log.Println("Reconstruction to be added")
}
