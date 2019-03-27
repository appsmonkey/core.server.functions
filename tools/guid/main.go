package guid

import (
	"math/rand"
	"sync"
	"time"
)

const (
	pushChars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	numChars  = 24
)

var (
	// Timestamp of last push, used to prevent local collisions if you push twice in one ms.
	lastPushTimeMs int64
	// We generate 72-bits of randomness which get turned into `numChars` characters and appended to the
	// timestamp to prevent collisions with other clients.  We store the last characters we
	// generated because in the event of a collision, we'll use those same characters except
	// "incremented" by one.
	lastRandChars [numChars]int
	mu            sync.Mutex
	rnd           *rand.Rand
)

func init() {
	// seed to get randomness
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func genRandPart() {
	for i := 0; i < len(lastRandChars); i++ {
		lastRandChars[i] = rnd.Intn(62)
	}
}

// New creates a new random, unique id
func New() string {
	var id [8 + numChars]byte
	mu.Lock()
	timeMs := time.Now().UTC().UnixNano() / 1e6
	if timeMs == lastPushTimeMs {
		// increment lastRandChars
		for i := 0; i < numChars; i++ {
			lastRandChars[i]++
			if lastRandChars[i] < 62 {
				break
			}
			// increment the next byte
			lastRandChars[i] = 0
		}
	} else {
		genRandPart()
	}
	lastPushTimeMs = timeMs
	// put random as the second part
	for i := 0; i < numChars; i++ {
		id[(numChars+7)-i] = pushChars[lastRandChars[i]]
	}
	mu.Unlock()

	// put current time at the beginning
	for i := 7; i >= 0; i-- {
		n := int(timeMs % 62)
		id[i] = pushChars[n]
		timeMs = timeMs / 62
	}
	return string(id[:])
}
