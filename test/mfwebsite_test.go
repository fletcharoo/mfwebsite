package test

import (
	"io"
	"net/http"
	"testing"
)

func Test_HelloWorld(t *testing.T) {
	resp, err := http.Get("http://localhost:8080")

	if err != nil {
		t.Fatalf("failed to GET: %s", err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatalf("failed to read body: %s", err)
	}

	if string(body) != "Hello, world!" {
		t.Fatalf("unexpected response: %s", err)
	}
}
