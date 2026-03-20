package main

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//go:embed VERSION
var versionFile string

var version = strings.TrimSpace(versionFile)

func main() {
	if os.Args[1] == "--version" {
		fmt.Println(version)
		return
	}

	if len(os.Args) < 3 || len(os.Args) > 4 {
		log.Fatalf("Usage: %s <directory> <port> [repo_name]\n", os.Args[0])
	}

	dir := os.Args[1]
	port := os.Args[2]
	var repoName string
	if len(os.Args) == 4 {
		repoName = os.Args[3]
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s\n", dir)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	fs := http.FileServer(http.Dir(dir))
	if repoName != "" {
		prefix := "/" + repoName
		r.Handle(prefix, http.RedirectHandler(prefix+"/", http.StatusMovedPermanently))
		r.Handle(prefix+"/*", http.StripPrefix(prefix, fs))
	} else {
		r.Handle("/*", fs)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		base := fmt.Sprintf("http://localhost:%s", port)
		if repoName != "" {
			base += "/" + repoName + "/"
		}
		log.Printf("Serving %s on %s\n", dir, base)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("forced shutdown: %v", err)
	}
}
