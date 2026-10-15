package testtools

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/doublehops/dh-go-framework/internal/config"
)

// startupTimeout bounds both how long the subprocess has to report its listening port and how
// long its health endpoint then has to start responding. db.New retries up to 5×5s before giving
// up, so this must comfortably exceed that.
const startupTimeout = 30 * time.Second

// StartTestServer builds the server binary (using the Go build cache) and starts it as a
// subprocess pointed at config_test.json, which configures port "0" — the OS assigns any free
// port, so a test server never collides with a fixed port already in use (e.g. a dev server).
// Packages still can't run in parallel: RunMigrations drops the shared test database, so run
// them with go test -p 1. It blocks until the subprocess reports its assigned port and the
// health endpoint responds, then mutates cfg.Host.Port/TestURL in place to that port so the
// caller's existing httprequest.GetRequester(cfg.Host.TestURL) calls reach the right place. The
// returned function stops the server and must be called before os.Exit in TestMain (defers do
// not run on os.Exit).
func StartTestServer(cfg *config.Config) (func(), error) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("StartTestServer: %w", err)
	}

	binPath := filepath.Join(projectRoot, "tmp", "testserver-bin")

	if err := os.MkdirAll(filepath.Dir(binPath), 0o750); err != nil {
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

	// Capture output so startup errors are visible in test failure messages, and so
	// waitForPort below can scan it for the subprocess's SERVER_LISTENING_PORT= marker line.
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

	deadline := time.Now().Add(startupTimeout)

	port, err := waitForPort(&mu, &buf, deadline, errCh)
	if err != nil {
		stop()

		return nil, err
	}

	cfg.Host.Port = strconv.Itoa(port)
	cfg.Host.TestURL = fmt.Sprintf("http://localhost:%d", port)

	if err := waitForServer(cfg.Host.TestURL+"/v1/health", time.Until(deadline), errCh); err != nil {
		stop()

		return nil, err
	}

	return stop, nil
}

// waitForPort polls buf (the subprocess's captured stdout/stderr) for the SERVER_LISTENING_PORT=
// marker line printed by cmd/server/run.go once it has bound its OS-assigned port.
func waitForPort(mu *sync.Mutex, buf *bytes.Buffer, deadline time.Time, errCh chan error) (int, error) {
	const marker = "SERVER_LISTENING_PORT="

	for time.Now().Before(deadline) {
		select {
		case err := <-errCh:
			return 0, fmt.Errorf("test server failed to start: %w", err)
		default:
		}

		mu.Lock()
		output := buf.String()
		mu.Unlock()

		if idx := strings.Index(output, marker); idx != -1 {
			rest := output[idx+len(marker):]
			if nl := strings.IndexByte(rest, '\n'); nl != -1 {
				rest = rest[:nl]
			}

			if port, err := strconv.Atoi(strings.TrimSpace(rest)); err == nil {
				return port, nil
			}
		}

		time.Sleep(50 * time.Millisecond)
	}

	return 0, fmt.Errorf("test server did not report a listening port within %s", startupTimeout)
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
