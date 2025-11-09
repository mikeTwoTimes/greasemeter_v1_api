package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Client struct {
    Bucket   *rate.Limiter
	LastSeen time.Time
}

type Limiter struct {
	sync.RWMutex
	clients     map[string]*Client
	maxRequests int
	frame       time.Duration
}

func Limit(l *Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !l.Allow(c.ClientIP()) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func NewLimiter(maxRequests int, frame time.Duration) *Limiter {
	l := Limiter{
		clients:     make(map[string]*Client),
		maxRequests: maxRequests,
		frame:       frame,
	}

	go l.cleanupClients()

	return &l
}

func (l *Limiter) Allow(ip string) bool {
	l.Lock()
	defer l.Unlock()

	client, exists := l.clients[ip]

	if !exists {
		limiter := rate.NewLimiter(
			rate.Every(l.frame / time.Duration(l.maxRequests)),
			l.maxRequests,
		)

		client = &Client{
			Bucket:   limiter,
			LastSeen: time.Now(),
		}

		l.clients[ip] = client
	}

	client.LastSeen = time.Now()

	return client.Bucket.Allow()
}

func (l *Limiter) cleanupClients() {
	for {
		time.Sleep(time.Minute)
		l.Lock()

		for ip, v := range l.clients {
			if time.Since(v.LastSeen) > 3 * time.Minute {
				delete(l.clients, ip)
			}
		}

		l.Unlock()
	}
}
