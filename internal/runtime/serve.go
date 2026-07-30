package runtime

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mpraes/tabyte/internal/httpapi"
	"github.com/mpraes/tabyte/internal/application"
)

func Serve(addr string) error {
	store := application.NewSessionStore()
	mux := httpapi.NewMux(store)
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	errCh := make(chan error, 1)
	go func() {
		fmt.Printf("Tabyte listening on http://%s\n", addr)
		errCh <- srv.ListenAndServe()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	case sig := <-sigCh:
		fmt.Printf("shutting down: %s\n", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Printf("shutdown error: %v\n", err)
			return err
		}
		return nil
	}
}