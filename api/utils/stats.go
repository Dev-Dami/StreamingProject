package utils

import (
	"log"
	"sync/atomic"
	"time"
)

var (
	// TotalFrames is the number of frames processed since the application started.
	TotalFrames uint64
)

func init() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			log.Printf("Stats: %d frames processed", atomic.LoadUint64(&TotalFrames))
		}
	}()
}

// IncrementFrameCount increments the total frame count by one.
func IncrementFrameCount() {
	atomic.AddUint64(&TotalFrames, 1)
}
