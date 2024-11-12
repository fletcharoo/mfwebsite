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
