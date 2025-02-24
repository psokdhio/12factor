package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	nethttpmiddleware "github.com/oapi-codegen/nethttp-middleware"
	"github.com/psokdhio/12factor/api"
)

type pinger struct {
	api.StrictServerInterface
}

func (p *pinger) GetPing(ctx context.Context, request api.GetPingRequestObject) (api.GetPingResponseObject, error) {
	return api.GetPing200JSONResponse{Ping: "Hello? Yes, this is Pong"}, nil
}

type PongHandlerOptions struct {
	BaseURL string
}

func NewPongHandler(opt PongHandlerOptions) (http.Handler, error) {
	swagger, err := api.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("swagger: get spec: %w", err)
	}
	p := new(pinger)
	h := api.HandlerWithOptions(
		api.NewStrictHandlerWithOptions(p, nil, api.StrictHTTPServerOptions{}),
		api.ChiServerOptions{
			BaseURL: opt.BaseURL,
			Middlewares: []api.MiddlewareFunc{
				nethttpmiddleware.OapiRequestValidatorWithOptions(swagger,
					&nethttpmiddleware.Options{
						Options: openapi3filter.Options{
							AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
						},
					}),
			},
		},
	)
	return h, nil
}

func serve(ctx context.Context, hs *http.Server, opt ServeOptions) error {
	go func() {
		if wg, ok := ctx.Value("Closers").(*sync.WaitGroup); ok {
			wg.Add(1)
			defer wg.Done()
		}
		select {
		case <-ctx.Done():
			ctxG, cancel := context.WithTimeout(context.Background(), opt.ShutdownGracePeriod)
			defer cancel()
			err := hs.Shutdown(ctxG)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					err = fmt.Errorf("timed out after %v", opt.ShutdownGracePeriod)
				}
				log.Printf("serve: shutting down the server: %v\n", err)
			}
		}
		_ = hs.Close()
	}()
	go func() {
		if err := hs.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Printf("serve: listen and serve: %v\n", err)
		}
	}()
	return nil
}

type ServeOptions struct {
	Addr                string
	ReadHeaderTimeout   time.Duration
	ShutdownGracePeriod time.Duration
}

func Serve(ctx context.Context, h http.Handler, opt ServeOptions) error {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/", h)
	hs := &http.Server{
		Handler:           r,
		Addr:              opt.Addr,
		ReadHeaderTimeout: opt.ReadHeaderTimeout,
	}

	ctx, stop := signal.NotifyContext(
		context.WithValue(ctx, "Closers", new(sync.WaitGroup)),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()
	err := serve(ctx, hs, opt)
	if err != nil {
		return fmt.Errorf("serve: running listeners: %w", err)
	}
	select {
	case <-ctx.Done():
		stop()
		ctx.Value("Closers").(*sync.WaitGroup).Wait()

	}
	return nil
}
