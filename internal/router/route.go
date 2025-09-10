package router

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jexroid/gopi/internal/crud"
	"github.com/jexroid/gopi/internal/database"
	"github.com/jexroid/gopi/internal/handler"
	httpSwagger "github.com/swaggo/http-swagger"
)

func Handler(r *chi.Mux) {
	var DB = database.Init()
	AuthDB := handler.InitDB(DB)
	CrudDB := crud.InitDB(DB)

	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // This should work
	))

	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.ServeFile(w, r, "docs/swagger.json")
	})

	r.Route("/auth", func(router chi.Router) {
		router.Group(func(r chi.Router) {
			r.Use(RateLimit(3, 10*time.Minute))
			r.Post("/otp", AuthDB.OTP)
		})

		router.Post("/signup", AuthDB.Signup)
		router.Post("/signin", AuthDB.Signin)
		router.Post("/validate", AuthDB.Validate)
		router.Get("/logout", handler.Logout)
		router.Post("/verify-otp", AuthDB.VerifyOTP)
	})

	// User CRUD routes
	r.Route("/user", func(router chi.Router) {
		router.Get("/{uuid}", CrudDB.Read)      // Read a user by UUID
		router.Get("/", CrudDB.All)             // Read all users
		router.Put("/{uuid}", CrudDB.Update)    // Update a user by UUID
		router.Delete("/{uuid}", CrudDB.Delete) // Delete a user by UUID
	})
}

func RateLimit(requests int, window time.Duration) func(http.Handler) http.Handler {
	type client struct {
		requests int
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, cl := range clients {
				if time.Since(cl.lastSeen) > window {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getClientIP(r)

			mu.Lock()
			defer mu.Unlock()

			cl, exists := clients[ip]
			if !exists {
				cl = &client{
					requests: 0,
					lastSeen: time.Now(),
				}
				clients[ip] = cl
			}

			if time.Since(cl.lastSeen) > window {
				cl.requests = 0
				cl.lastSeen = time.Now()
			}

			if cl.requests >= requests {
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", int(window.Seconds())))
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"error": "Rate limit exceeded", "retry_after": "10 minutes"}`))
				return
			}

			cl.requests++
			cl.lastSeen = time.Now()

			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", requests))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", requests-cl.requests))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", int(window.Seconds())))

			next.ServeHTTP(w, r)
		})
	}
}

func getClientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
