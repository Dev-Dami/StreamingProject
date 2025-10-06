package video

import (
	"bytes"
	"image/jpeg"
	"log"
	"api/streaming"
)

var QuantFrameChan = make(chan []byte, 30)

func init() {
	go processFrame()
}

func processFrame() {
	for frame := range FrameChan {
		quant, err := quantFrame(frame)
		if err != nil {
			log.Println("Quant error:", err)
			continue
		}
		QuantFrameChan <- quant
		streaming.Broadcast(quant)
	}
}

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