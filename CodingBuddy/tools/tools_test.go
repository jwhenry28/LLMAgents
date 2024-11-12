package tools

import (
	"strings"
	"testing"

	"net/http"

	"github.com/jwhenry28/LLMAgents/shared/model"
)

type testServer struct {
	server *http.Server
}

func newTestServer() *testServer {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("../testdata"))
	mux.Handle("/", fs)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	ts := &testServer{
		server: server,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return ts
}

func (ts *testServer) shutdown() error {
	return ts.server.Close()
}

func TestFetch(t *testing.T) {
	ts := newTestServer() // comment this out and use python http.server if debugging
	defer ts.shutdown()

	tests := []struct {
		name     string
		args     []string
		wantErr  bool
		contains []string
	}{
		{
			name:     "example.com",
			args:     []string{"http://localhost:8080/example.html"},
			wantErr:  false,
			contains: []string{"Example Domain", `<a href="https://www.iana.org/domains/example">More information...</a>`},
		},
		{
			name:     "hackernews",
			args:     []string{"http://localhost:8080/hackernews.html"},
			wantErr:  false,
			contains: []string{"Hacker News", `<a href="https://www.asimov.press/p/mitochondria">Mitochondria Are Alive</a>`, `<a href="https://igorstechnoclub.com/most-common-sqlite-misconception/" rel="nofollow">SQLite is not a single connection database</a>`, `<a href="https://mtlynch.io/why-i-quit-google/">I quit Google to work for myself (2018)</a>`},
		},
		{
			name:     "invalid url",
			args:     []string{"not-a-url"},
			wantErr:  true,
			contains: []string{},
		},
		{
			name:     "missing args",
			args:     []string{},
			wantErr:  true,
			contains: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetch := NewFetch(model.TextToolInput{
				Name: "fetch",
				Args: tt.args,
			})

			if !tt.wantErr && !fetch.Match() {
				t.Errorf("Fetch.Match() = false, want true")
				return
			}

			if tt.wantErr && fetch.Match() {
				t.Errorf("Fetch.Match() = true, want false")
				return
			}

			if !tt.wantErr {
				result := fetch.Invoke()
				for _, expected := range tt.contains {
					if !strings.Contains(result, expected) {
						t.Errorf("missing expected string: %s", expected)
					}
				}
			}
		})
	}
}
