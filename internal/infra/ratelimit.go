package infra

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// ponytail: per-IP in-memory, global lock, map grows unbounded (negligible for blog traffic).
// Add Redis-backed + per-account buckets when >1 replica or under sustained attack.
var (
	rmu     sync.Mutex
	buckets = make(map[string]*tokenBucket)
)

type tokenBucket struct {
	tokens float64
	last   time.Time
}

func RateLimit(rps float64) gin.HandlerFunc {
	return func(c *gin.Context) {
		rmu.Lock()
		ip := c.ClientIP()
		b, ok := buckets[ip]
		now := time.Now()
		if !ok {
			b = &tokenBucket{tokens: rps, last: now}
			buckets[ip] = b
		}
		elapsed := now.Sub(b.last).Seconds()
		b.tokens += elapsed * rps
		b.last = now
		if b.tokens > rps {
			b.tokens = rps
		}
		if b.tokens < 1 {
			rmu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "muitas requisições. Tente novamente mais tarde."})
			return
		}
		b.tokens--
		rmu.Unlock()
		c.Next()
	}
}
