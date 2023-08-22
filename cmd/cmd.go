package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"

	"github.com/mjevans93308/brightwheel_takehome/internal/handlers"
)

var (
	port int
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "webserver",
		Short: "Start the Go web server",
		Run: func(cmd *cobra.Command, args []string) {
			startServer(port)
		},
	}

	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "web server port")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}

func startServer(port int) {
	r := chi.NewRouter()
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Use(middleware.BasicAuth("brightwheelrealm", map[string]string{
			"admin": "password1!",
		}))
		r.Get("/latest_timestamp", handlers.LatestTimestampHandler)
		r.Get("/cumulative_count", handlers.CumulativeHandler)
		r.Post("/device", handlers.DeviceHandler)
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Chi!"))
	})

	addr := fmt.Sprintf(":%d", port)
	server := http.Server{
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 3 * time.Second,
		Addr:         addr,
		Handler:      r,
	}

	fmt.Printf("Starting Chi web server on port %d...\n", port)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Shutting down the server...")

	if err := server.Shutdown(nil); err != nil {
		fmt.Printf("Error shutting down: %v\n", err)
	}
}
