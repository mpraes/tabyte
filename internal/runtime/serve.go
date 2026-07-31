package runtime

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mpraes/tabyte/internal/application"
	"github.com/mpraes/tabyte/internal/httpapi"
	"github.com/mpraes/tabyte/internal/persistence/sqlite"
	"github.com/mpraes/tabyte/internal/platform"
)

type ServeOptions struct {
	Addr         string
	OpenBrowser  bool
	Persist      bool
	DBPath       string
}

func Serve(opts ServeOptions) error {
	var (
		store        *application.SessionStore
		settings     application.SettingsRepository
		persistence  bool
		db           *sqlite.DB
	)

	if opts.Persist {
		path := opts.DBPath
		if path == "" {
			var err error
			path, err = sqlite.DefaultPath()
			if err != nil {
				return fmt.Errorf("resolve db path: %w", err)
			}
		}
		opened, err := sqlite.Open(path)
		if err != nil {
			return fmt.Errorf("open sqlite: %w", err)
		}
		db = opened
		defer db.Close()

		store = application.NewSessionStore(db)
		if err := store.LoadFromRepo(); err != nil {
			return fmt.Errorf("load sessions: %w", err)
		}
		settings = db
		persistence = true
		fmt.Printf("persistence enabled: %s\n", path)
	} else {
		store = application.NewSessionStore(nil)
	}

	mux := httpapi.NewMux(store, settings, persistence)
	srv := &http.Server{
		Addr:    opts.Addr,
		Handler: mux,
	}

	url := "http://" + opts.Addr
	errCh := make(chan error, 1)
	go func() {
		fmt.Printf("Tabyte listening on %s\n", url)
		errCh <- srv.ListenAndServe()
	}()

	if opts.OpenBrowser {
		go func() {
			time.Sleep(200 * time.Millisecond)
			if err := platform.OpenBrowser(url); err != nil {
				fmt.Printf("could not open browser: %v\n", err)
			}
		}()
	}

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
