package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/connect"
)

// $ connect-static-files [<url-path>:<file-path>, ...]

func main() {
	m, err := paths(os.Args[1:])
	check(err)

	service := getStringOr("SERVICE", "connect-static")
	bind := getStringOr("BIND", "0.0.0.0")
	port := mustGetInt("PORT")

	cs, err := connect.NewService(service, consul())
	check(err)

	s := http.Server{
		Addr:              fmt.Sprintf("%s:%d", bind, port),
		Handler:           router(m),
		TLSConfig:         cs.ServerTLSConfig(),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 1 * time.Second,
		WriteTimeout:      5 * time.Second,
	}

	check(s.ListenAndServeTLS("", ""))
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
		fmt.Printf("mapping %s -> %s\n", urlPath, filePath)
		mux.Handle(urlPath, http.StripPrefix(urlPath, http.FileServer(http.Dir(filePath))))
	}
	return mux
}

func check(err error) {
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
}

// todo: refactor these helpers

func consul() *api.Client {
	logEnvironment("CONSUL_HTTP_ADDR")
	logEnvironment("CONSUL_NAMESPACE")
	logEnvironment("CONSUL_CACERT")
	logEnvironment("CONSUL_CLIENT_CERT")
	logEnvironment("CONSUL_CLIENT_KEY")
	logEnvironment("CONSUL_HTTP_SSL")
	logEnvironment("CONSUL_HTTP_SSL_VERIFY")
	logEnvironment("CONSUL_TLS_SERVER_NAME")
	logEnvironment("CONSUL_GRPC_ADDR")
	logEnvironment("CONSUL_HTTP_TOKEN_FILE")
	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatal("failed to make consul client:", err)
	}
	return consulClient
}

func logEnvironment(name string) {
	value := os.Getenv(name)
	if value == "" {
		value = "<unset>"
	}
	log.Printf("environment %s = %s", name, value)
}

func mustGetInt(name string) int {
	if s := os.Getenv(name); s != "" {
		p, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(name + " must be a number")
		}
		return p
	}
	log.Fatal(name + " must be set")
	return -1
}

func getStringOr(name, alt string) string {
	if s := os.Getenv(name); s != "" {
		return s
	}
	return alt
}
