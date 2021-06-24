package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// $ static-files [<url-path>:<file-path>, ...]

func main() {
	m, err := args(os.Args[1:])
	check(err)

	mux := router(m)

	_ = mux
	s := http.Server{
		Addr:              "0.0.0.0:8000",
		Handler:           mux,
		TLSConfig:         nil,
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	check(s.ListenAndServe())
}

func args(args []string) (map[string]string, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("must provide at least one path pair")
	}

	m := make(map[string]string, len(args))
	for _, arg := range args {
		tokens := strings.SplitN(arg, ":", 2)
		if len(tokens) != 2 {
			return nil, fmt.Errorf("arg format must be <url-path>:<file-path>")
		}
		m[tokens[0]] = m[tokens[1]]
	}
	return m, nil
}

func router(m map[string]string) http.Handler {
	mux := http.NewServeMux()
	for urlPath, filePath := range m {
		mux.Handle(urlPath, http.StripPrefix(filePath, http.FileServer(http.Dir(filePath))))
	}
	return mux
}

func check(err error) {
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
}
