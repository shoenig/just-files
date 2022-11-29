package main

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/shoenig/test/must"
)

func TestRun(t *testing.T) {
	t.Setenv("BIND", "127.0.0.1")
	t.Setenv("PORT", "8787")
	dir, err := os.MkdirTemp("", "")
	must.NoError(t, err)
	filename := filepath.Join(dir, "hi.txt")
	err = os.WriteFile(filename, []byte("hello"), 0o644)
	must.NoError(t, err)

	t.Log("dir", dir)
	args := []string{"just-files", "/:" + dir}

	go func() {
		err := run(args)
		must.NoError(t, err)
	}()

	time.Sleep(500 * time.Millisecond)
	resp, err := http.Get("http://127.0.0.1:8787/hi.txt")
	must.NoError(t, err)
	must.Eq(t, 200, resp.StatusCode)
}
