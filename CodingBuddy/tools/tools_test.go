package tools

import (
	"os"
	"strings"
	"testing"

	"net/http"

	"github.com/jwhenry28/LLMAgents/coding-buddy/utils"
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
		failMatch  bool
		contains []string
	}{
		{
			name:     "example.com",
			args:     []string{"http://localhost:8080/example.html"},
			failMatch:  false,
			contains: []string{"Example Domain", `<a href="https://www.iana.org/domains/example">More information...</a>`},
		},
		{
			name:     "hackernews",
			args:     []string{"http://localhost:8080/hackernews.html"},
			failMatch:  false,
			contains: []string{"Hacker News", `<a href="https://www.asimov.press/p/mitochondria">Mitochondria Are Alive</a>`, `<a href="https://igorstechnoclub.com/most-common-sqlite-misconception/" rel="nofollow">SQLite is not a single connection database</a>`, `<a href="https://mtlynch.io/why-i-quit-google/">I quit Google to work for myself (2018)</a>`},
		},
		{
			name:     "invalid url",
			args:     []string{"not-a-url"},
			failMatch:  true,
			contains: []string{},
		},
		{
			name:     "missing args",
			args:     []string{},
			failMatch:  true,
			contains: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetch := NewFetch(model.TextToolInput{
				Name: "fetch",
				Args: tt.args,
			})

			if !tt.failMatch && !fetch.Match() {
				t.Errorf("Fetch.Match() = false, want true")
				return
			}

			if tt.failMatch && fetch.Match() {
				t.Errorf("Fetch.Match() = true, want false")
				return
			}

			if !tt.failMatch {
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

func TestWrite(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		failMatch bool
		failExec bool
	}{
		{
			name:    "valid write",
			args:    []string{"test.txt", "hello world"},
			failMatch: false,
			failExec: false,
		},
		{
			name: "valid write multiline",
			args: []string{"test.txt", `hello world
this is a test`},
			failMatch: false,
			failExec: false,
		},
		{
			name:    "valid write nested",
			args:    []string{"dir/test.txt", "hello world"},
			failMatch: false,
			failExec: false,
		},
		{
			name:    "valid write nested deeper",
			args:    []string{"dir/subdir/test.txt", "hello world"},
			failMatch: false,
			failExec: false,
		},
		{
			name:    "absolute path",
			args:    []string{"/test.txt", "hello world"},
			failMatch: false,
			failExec: true,
		},
		{
			name:    "path traversal",
			args:    []string{"../test.txt", "hello world"},
			failMatch: false,
			failExec: true,
		},
		{
			name:    "missing args",
			args:    []string{"test.txt"},
			failMatch: true,
		},
		{
			name:    "no args",
			args:    []string{},
			failMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			write := NewWrite(model.TextToolInput{
				Name: "write",
				Args: tt.args,
			})

			if (tt.failMatch && write.Match()) || (!tt.failMatch && !write.Match() ){
				t.Errorf("Write.Match() = %v, want %v", write.Match(), tt.failMatch)
				return
			}

			if tt.failMatch {
				return
			}

			defer os.RemoveAll(utils.SANDBOX_DIR)

			result := write.Invoke()
			if tt.failExec != strings.HasPrefix(result, "error: ") {
				t.Errorf("Write.Invoke() = %v, want %v", result, tt.failExec)
				return
			}

			content, err := os.ReadFile(utils.SANDBOX_DIR + "/" + tt.args[0])
			if !tt.failExec && err != nil {
				t.Errorf("Failed to read file: %v", err)
				return
			}

			if !tt.failExec && string(content) != tt.args[1] {
				t.Errorf("File content = %v, want %v", string(content), tt.args[1])
				return
			}
		})
	}
}
