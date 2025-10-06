package utils

import (
	"log"
	"sync/atomic"
	"time"
)

var (
	// FramesProcessed counts total frames processed
	FramesProcessed uint64
	// FramesBroadcast counts total frames broadcast to clients
	FramesBroadcast uint64
	// ClientConnections tracks current client count
	ClientConnections uint64
)

// Initialize stats reporting
func init() {
	go startStatsReporter()
}

// startStatsReporter runs periodic stats logging
func startStatsReporter() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		processed := atomic.LoadUint64(&FramesProcessed)
		broadcast := atomic.LoadUint64(&FramesBroadcast)
		clients := atomic.LoadUint64(&ClientConnections)
		
		if processed > 0 || broadcast > 0 {
			log.Printf("Stats - Processed: %d frames, Broadcast: %d frames, Clients: %d", 
				processed, broadcast, clients)
		}
	}
}

// IncrementProcessed increments processed frame counter
func IncrementProcessed() {
	atomic.AddUint64(&FramesProcessed, 1)
}

// IncrementBroadcast increments broadcast frame counter
func IncrementBroadcast() {
	atomic.AddUint64(&FramesBroadcast, 1)
}

// SetClientCount sets the current client count
func SetClientCount(count uint64) {
	atomic.StoreUint64(&ClientConnections, count)
}
