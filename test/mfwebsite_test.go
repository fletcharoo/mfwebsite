package test

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/gkampitakis/go-snaps/snaps"
)

func TestMain(m *testing.M) {
	r := m.Run()
	snaps.Clean(m, snaps.CleanOpts{Sort: true})
	os.Exit(r)
}

func Test_styleCSS(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/style.css")

	if err != nil {
		t.Fatalf("Failed to GET: %s", err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatalf("Failed to read body: %s", err)
	}

	snaps.MatchSnapshot(t, string(body))
}

func Test_index(t *testing.T) {
	resp, err := http.Get("http://localhost:8080")

	if err != nil {
		t.Fatalf("Failed to GET: %s", err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatalf("Failed to read body: %s", err)
	}

	snaps.MatchSnapshot(t, string(body))
}

func Test_other(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/other")

	if err != nil {
		t.Fatalf("Failed to GET: %s", err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatalf("Failed to read body: %s", err)
	}

	snaps.MatchSnapshot(t, string(body))
}

func Test_subdir(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/sub/subdir")

	if err != nil {
		t.Fatalf("Failed to GET: %s", err)
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		t.Fatalf("Failed to read body: %s", err)
	}

	snaps.MatchSnapshot(t, string(body))
}
