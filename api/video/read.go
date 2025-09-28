package video

import (
	"log"
	"os"
	"os/exec"
	"time"
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
	buf := make([]byte, 1024*1024)
	for {
		n, err := stdout.Read(buf)
		if err != nil {
			break
		}
		frame := make([]byte, n)
		copy(frame, buf[:n])

		select {
		case FrameChan <- frame:
		default:
		}
		time.Sleep(10 * time.Millisecond)
	}
	cmd.Wait()
	close(FrameChan)
}

