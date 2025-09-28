<<<<<<< HEAD
package utils

import (
	"log"
	"sync/atomic"
	"time"
)

var ( 
    // Total frames processed since the application started.
    TotalFrames uint64
)

=======
package video

import (
	"log"
	"time"
)

>>>>>>> 6b02c2f6467c8679a14d337174326567ec50c5ae
func init() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
<<<<<<< HEAD
		for range ticker.C {
			log.Printf("Stats: %d frames processed", atomic.LoadUint64(&TotalFrames))
		}
	}()
}

// IncrementFrameCount increments the total frame count by one.
func IncrementFrameCount() {
	atomic.AddUint64(&TotalFrames, 1)
}
=======
		count := 0
		for range ticker.C {
			// WIP, track frame count, bandwidth, etc.
			log.Printf("Stats: %d frames processed (placeholder)", count)
			count++
		}
	}()
}
>>>>>>> 6b02c2f6467c8679a14d337174326567ec50c5ae
