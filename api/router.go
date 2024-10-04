package api

import (
	"log"
	"net/http"
	"time"

	"github.com/AsaHero/just-ask-bot/pkg/config"
	telegram_bot "github.com/AsaHero/just-ask-bot/pkg/telegram-bot"
)

type router struct {
	Config      *config.Config
	TelegramBot *telegram_bot.TelegramBot
}

// NewRouter initializes and returns a new HTTP handler with configured routes.
func NewRouter(config *config.Config, telegramBot *telegram_bot.TelegramBot) http.Handler {
	mux := http.NewServeMux()

	// Initialize the router with configuration and Telegram bot
	r := &router{
		Config:      config,
		TelegramBot: telegramBot,
	}

	// Set up the route for Telegram webhook
	mux.HandleFunc("/webhooks/telegram", corsMiddleware(r.TelegramWebhook))

	return loggerMiddleware(mux)
}

// corsMiddleware applies CORS headers to responses
func corsMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust this to allow specific domains
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		handler(w, r)
	}
}

// loggerMiddleware logs every HTTP request
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log the incoming request
		log.Printf("Received [%s] request from %s for %s", r.Method, r.RemoteAddr, r.URL.Path)

		// Use a response writer that allows us to capture the status code
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		// Log the completed request with the status code and response time
		log.Printf("Completed [%s] request from %s for %s with status %d in %s", r.Method, r.RemoteAddr, r.URL.Path, lrw.statusCode, time.Since(start))
	})
}

// loggingResponseWriter is a wrapper around http.ResponseWriter that captures the status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// WriteHeader captures the status code and delegates to the original ResponseWriter
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// newLoggingResponseWriter initializes a new loggingResponseWriter
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, 200}
}

// TelegramWebhook handles incoming updates from Telegram
func (h *router) TelegramWebhook(w http.ResponseWriter, r *http.Request) {
	h.TelegramBot.HandleUpdate(w, r)
}
