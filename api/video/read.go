package video

import (
	"log"
	"os/exec"
	"time"
)

var FrameChan = make(chan []byte, 30)

func StartPipeline(videoPath string) {
	go readAndDecode(videoPath)
}

func readAndDecode(videoPath string) {
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-vf", "fps=10",
		"-f", "image2pipe",
		"-q:v", "1", "-",
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

