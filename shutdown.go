package bttp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ShutdownConfig holds all variables needed in the shutdown process.
type ShutdownConfig struct {
	// the timeout for performing a graceful shutdown.
	Timeout time.Duration
}

func newShutdownConfig() *ShutdownConfig {
	return &ShutdownConfig{
		Timeout: time.Second * 15,
	}
}

// ShutdownOpt is a modifier function. They can be used to modify a ShutdownConfig.
type ShutdownOpt func(c *ShutdownConfig)

// SetTimeout sets the timeout for performing a graceful shutdown.
func SetTimeout(t time.Duration) ShutdownOpt {
	return func(c *ShutdownConfig) {
		c.Timeout = t
	}
}

// ListenGracefully starts a http.Server and handles the shutdown process,
// when a syscall.SIGINT or syscall.SIGTERM are sent to the process.
func ListenGracefully(srv *http.Server, opts ...ShutdownOpt) error {
	// Creating a default config and override it with given opts
	c := newShutdownConfig()
	for _, applier := range opts {
		applier(c)
	}

	done := make(chan os.Signal, 1)
	listenErr := make(chan error, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Starting the blocking ListenAndServe call in a separate goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			// Ignoring this error because this is caused by a Shutdown or Close command.
			// So, this is not really an error in this case.
			if err != http.ErrServerClosed {
				listenErr <- err
			}
		}
	}()

	select {
	case err := <-listenErr:
		return fmt.Errorf("unexpected error while listening: %w", err)
	case <-done:
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			return fmt.Errorf("error performing shutdown: %w", err)
		}
		return nil
	}
}
