package testtools

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// StartTestServer builds the server binary (using the Go build cache) and starts it
// as a subprocess pointed at config_test.json. It blocks until the health endpoint
// responds or the timeout elapses. The returned function stops the server and must
// be called before os.Exit in TestMain (defers do not run on os.Exit).
func StartTestServer() (func(), error) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("StartTestServer: %w", err)
	}

	binPath := filepath.Join(projectRoot, "tmp", "testserver-bin")

	if err := os.MkdirAll(filepath.Dir(binPath), 0o750); err != nil { //nolint:gomnd
		return nil, fmt.Errorf("StartTestServer: failed to create tmp dir: %w", err)
	}

	buildCmd := exec.Command("go", "build", "-o", binPath, "./cmd/server") //nolint:gosec
	buildCmd.Dir = projectRoot

	if out, err := buildCmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("StartTestServer: build failed: %w\n%s", err, out)
	}

	configPath := filepath.Join(projectRoot, "config_test.json")
	serverCmd := exec.Command(binPath, "-config", configPath) //nolint:gosec
	serverCmd.Dir = projectRoot

	// Capture output so startup errors are visible in test failure messages.
	var (
		mu  sync.Mutex
		buf bytes.Buffer
	)
	w := io.MultiWriter(os.Stderr, &buf)
	serverCmd.Stdout = w
	serverCmd.Stderr = w

	if err := serverCmd.Start(); err != nil {
		return nil, fmt.Errorf("StartTestServer: failed to start process: %w", err)
	}

	errCh := make(chan error, 1)

	go func() {
		if err := serverCmd.Wait(); err != nil {
			mu.Lock()
			output := buf.String()
			mu.Unlock()
			errCh <- fmt.Errorf("%w\n--- server output ---\n%s", err, output)
		}
	}()

	stop := func() {
		_ = serverCmd.Process.Kill()
	}

	// Allow up to 30s: db.New retries up to 5×5s before giving up.
	if err := waitForServer("http://localhost:8088/v1/health", 30*time.Second, errCh); err != nil {
		stop()
		return nil, err
	}

	return stop, nil
}

// waitForServer polls the given URL until it returns a 200 response, the timeout elapses,
// or the server goroutine reports a startup error.
func waitForServer(url string, timeout time.Duration, errCh chan error) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: time.Second}

	for time.Now().Before(deadline) {
		select {
		case err := <-errCh:
			return fmt.Errorf("test server failed to start: %w", err)
		default:
		}

		resp, err := client.Get(url) //nolint:noctx
		if err == nil {
			_ = resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				return nil
			}
		}

		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("test server did not become ready within %s", timeout)
}
