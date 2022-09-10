package main

import (
	"errors"
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
	m, err := paths(os.Args[1:])
	check(err)

	bind := os.Getenv("BIND")
	port := os.Getenv("PORT")
	if bind == "" || port == "" {
		check(errors.New("$BIND and $PORT must be set"))
	}
	log.Println("bind is", bind)
	log.Println("port is", port)

	s := http.Server{
		Addr:              fmt.Sprintf("%s:%s", bind, port),
		Handler:           router(m),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
	}
	check(s.ListenAndServe())
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
