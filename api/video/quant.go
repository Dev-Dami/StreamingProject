package video

import (
	"bytes"
	"image/jpeg"
	"log"
	"api/streaming"
)

func init() {
	go processFrame()
}

// Processes frames from the FrameChan, quantizes them, and broadcasts them.
func processFrame() {
	for frame := range FrameChan {
		quant, err := quantFrame(frame)
		if err != nil {
			log.Println("Quant error:", err)
			continue
		}
		streaming.Broadcast(quant)
	}
}

// Quantizes a JPEG image by reducing its quality.
func quantFrame(jpegData []byte) ([]byte, error) {
	img, err := jpeg.Decode(bytes.NewReader(jpegData))
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 30})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}