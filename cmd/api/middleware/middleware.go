package middleware

import (
	"html"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Middleware func(http.Handler) http.Handler

func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}

// Logger logs every request with method, path, and duration
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// CORS restricts origin to allowed list
func CORS(next http.Handler) http.Handler {
	allowedOrigins := map[string]bool{
		"http://localhost:3000":  true,
		"http://localhost:3001":  true,
		"http://localhost:5173":  true,
		"https://vocal-licorice-48fc8a.netlify.app": true,
		"https://alona-nonforested-remedios.ngrok-free.dev": true,
	}

	// Check environment for production origin
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else if strings.HasSuffix(origin, ".onrender.com") || strings.HasSuffix(origin, ".vercel.app") {
			// Allow render/vercel subdomains in production
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, ngrok-skip-browser-warning")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Recovery catches panics and returns 500
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: %v", err)
				http.Error(w, `{"error":"internal server error"}`, http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// SecurityHeaders adds common security headers
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		next.ServeHTTP(w, r)
	})
}

// RateLimiter implements a token bucket per IP
type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
}

type visitor struct {
	tokens    int
	lastSeen  time.Time
}

var limiter = &rateLimiter{
	visitors: make(map[string]*visitor),
}

func init() {
	// Cleanup stale visitors every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			limiter.mu.Lock()
			for ip, v := range limiter.visitors {
				if time.Since(v.lastSeen) > 10*time.Minute {
					delete(limiter.visitors, ip)
				}
			}
			limiter.mu.Unlock()
		}
	}()
}

func getVisitor(ip string, maxTokens int) *visitor {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	v, exists := limiter.visitors[ip]
	if !exists {
		v = &visitor{tokens: maxTokens, lastSeen: time.Now()}
		limiter.visitors[ip] = v
		return v
	}

	// Refill tokens: 1 token per second, max bucket size
	elapsed := time.Since(v.lastSeen)
	refill := int(elapsed.Seconds())
	if refill > 0 {
		v.tokens += refill
		if v.tokens > maxTokens {
			v.tokens = maxTokens
		}
		v.lastSeen = time.Now()
	}

	return v
}

// RateLimit limits requests per IP (general: 60/min)
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := extractIP(r)
		v := getVisitor(ip, 60)

		limiter.mu.Lock()
		if v.tokens <= 0 {
			limiter.mu.Unlock()
			w.Header().Set("Retry-After", "60")
			http.Error(w, `{"error":"too many requests","code":429}`, http.StatusTooManyRequests)
			return
		}
		v.tokens--
		limiter.mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

// AuthRateLimit stricter limit for auth endpoints (10/min)
func AuthRateLimit(next http.Handler) http.Handler {
	authLimiter := &rateLimiter{visitors: make(map[string]*visitor)}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only limit POST auth endpoints
		if r.Method == http.MethodPost && (strings.Contains(r.URL.Path, "/auth/login") || strings.Contains(r.URL.Path, "/auth/register")) {
			ip := extractIP(r)

			authLimiter.mu.Lock()
			v, exists := authLimiter.visitors[ip]
			if !exists {
				v = &visitor{tokens: 10, lastSeen: time.Now()}
				authLimiter.visitors[ip] = v
			} else {
				elapsed := time.Since(v.lastSeen)
				refill := int(elapsed.Seconds() / 6) // 1 token per 6 seconds = 10/min
				if refill > 0 {
					v.tokens += refill
					if v.tokens > 10 {
						v.tokens = 10
					}
					v.lastSeen = time.Now()
				}
			}

			if v.tokens <= 0 {
				authLimiter.mu.Unlock()
				w.Header().Set("Retry-After", "60")
				http.Error(w, `{"error":"too many login attempts, try again later","code":429}`, http.StatusTooManyRequests)
				return
			}
			v.tokens--
			authLimiter.mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}

// SanitizeInput escapes HTML entities in common string fields
func SanitizeInput(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Reject requests with suspicious content-type
		if r.Method != http.MethodGet && r.Method != http.MethodOptions && r.Method != http.MethodDelete {
			ct := r.Header.Get("Content-Type")
			if ct != "" && !strings.Contains(ct, "application/json") && !strings.Contains(ct, "multipart/form-data") {
				http.Error(w, `{"error":"unsupported content type","code":415}`, http.StatusUnsupportedMediaType)
				return
			}
		}

		// Limit request body size (5MB max)
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, 5*1024*1024)
		}

		next.ServeHTTP(w, r)
	})
}

// SanitizeString escapes HTML in user input (use in handlers)
func SanitizeString(s string) string {
	s = html.EscapeString(s)
	// Remove null bytes
	s = strings.ReplaceAll(s, "\x00", "")
	return s
}

func extractIP(r *http.Request) string {
	// Check X-Forwarded-For (behind proxy)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.SplitN(xff, ",", 2)
		return strings.TrimSpace(parts[0])
	}
	if xri := r.Header.Get("X-Real-Ip"); xri != "" {
		return xri
	}
	// Fallback to remote addr
	parts := strings.SplitN(r.RemoteAddr, ":", 2)
	return parts[0]
}
