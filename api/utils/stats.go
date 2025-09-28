package video

import (
	"log"
	"time"
)

func init() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		count := 0
		for range ticker.C {
			// WIP, track frame count, bandwidth, etc.
			log.Printf("Stats: %d frames processed (placeholder)", count)
			count++
		}
	}()
}