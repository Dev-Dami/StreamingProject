package video

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

// FrameChan is a channel that receives video frames.
var FrameChan = make(chan []byte, 30)

// StartPipeline starts the video processing pipeline.
func StartPipeline(videoPath string) {
	go readAndDecode(videoPath)
}

// readAndDecode reads a video file, decodes it into frames, and sends them to the FrameChan.
func readAndDecode(videoPath string) {
	if _, err := os.Stat(videoPath); os.IsNotExist(err) {
		log.Fatalf("video file not found: %s", videoPath)
	}

	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-f", "image2pipe",
		"-vcodec", "mjpeg",
		"-vf", "fps=10",
		"-q:v", "1",
		"-",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("ffmpeg stdout Pipe:", err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal("ffmpeg Start:", err)
	}

	buffer := make([]byte, 1024*1024)
	var frameEnd int
	var frame []byte

	for {
		n, err := stdout.Read(buffer[frameEnd:])
		if err != nil {
			break
		}
		frameEnd += n

		// Find start of image marker
		soi := bytes.Index(buffer[:frameEnd], []byte{0xff, 0xd8})
		if soi == -1 {
			continue
		}

		// Find end of image marker
		eoi := bytes.Index(buffer[soi:frameEnd], []byte{0xff, 0xd9})
		if eoi == -1 {
			continue		}


		frame = make([]byte, soi+eoi+2)
		copy(frame, buffer[soi:soi+eoi+2])

		select {
		case FrameChan <- frame:
		default:
		}

		// Move remaining data to the beginning of the buffer
		copy(buffer, buffer[soi+eoi+2:frameEnd])
		frameEnd = frameEnd - (soi + eoi + 2)
	}

	cmd.Wait()
	close(FrameChan)
}

