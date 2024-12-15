package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Usage:
//   just-files [<url-path>:<file-path>, ...]
//
// Example:
//   $ BIND=127.0.0.1 PORT=8080 just-files /:.

func main() {
	if err := run(os.Args); err != nil {
		fmt.Println("failure:", err)
		os.Exit(1)
	}
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func run(args []string) error {
	m, err := paths(args[1:])
	if err != nil {
		return err
	}

	lockdown(m)

	bind := envOr("BIND", "127.0.0.1")
	port := envOr("PORT", "8000")

	log.Println("bind is", bind)
	log.Println("port is", port)

	s := http.Server{
		Addr:              fmt.Sprintf("%s:%s", bind, port),
		Handler:           router(m),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
	return s.ListenAndServe()
}

func paths(args []string) (map[string]string, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("must provide at least one path pair")
	}

	m := make(map[string]string, len(args))
	for _, arg := range args {
		tokens := strings.SplitN(arg, ":", 2)
		if len(tokens) != 2 {
			return nil, fmt.Errorf("arg format must be <url-path>:<file-path>")
		}
		m[tokens[0]] = tokens[1]
	}
	return m, nil
}

func router(m map[string]string) http.Handler {
	mux := http.NewServeMux()
	for urlPath, filePath := range m {
		log.Printf("mapping %s -> %s", urlPath, filePath)
		mux.Handle(urlPath, http.StripPrefix(urlPath, http.FileServer(http.Dir(filePath))))
	}
	return mux
}

func check(err error) {
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
}
