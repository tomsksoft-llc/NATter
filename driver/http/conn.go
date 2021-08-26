package http

import (
	"context"
	"net"
	"net/http"
	"sync"

	"NATter/driver"
	"NATter/entity"
	"NATter/errtpl"
	"NATter/log"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

const (
	DriverName = "http"

	reservedURIPattern = `^/i(/.*)?$` // /i/* or /i
)

type ConnConfig struct {
	Host string
	Port string
}

type conn struct {
	host string
	port string

	mux    *chi.Mux
	routes []*entity.Route
	wg     *sync.WaitGroup
}

func NewConn(cfg *ConnConfig) driver.Conn {
	return &conn{
		host:   cfg.Host,
		port:   cfg.Port,
		mux:    chi.NewRouter(),
		routes: []*entity.Route{},
		wg:     &sync.WaitGroup{},
	}
}

func (c *conn) Serve(ctx context.Context) error {
	c.registerInternalRoutes()

	srv := http.Server{
		Addr:    net.JoinHostPort(c.host, c.port),
		Handler: c.mux,
	}

	go func() {
		<-ctx.Done()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error(errtpl.ErrClose(err, "http server"))
		}
	}()

	err := srv.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errtpl.ErrClose(err, "http server")
	}

	log.Debug("closing HTTP Server")

	c.wg.Wait()

	return nil
}

func (c *conn) registerInternalRoutes() {
	c.mux.Get("/i/routes", routes(func() []*entity.Route {
		return c.routes
	}))
}

func (c *conn) Close() error {
	return nil
}

func (c *conn) Receiver(route *entity.Route) driver.Receiver {
	c.routes = append(c.routes, route)

	return &receiver{
		mux:      c.mux,
		wg:       c.wg,
		async:    route.Async,
		uri:      route.URI,
		endpoint: route.Endpoint,
	}
}

func (c *conn) Sender(route *entity.Route) driver.Sender {
	c.routes = append(c.routes, route)

	return &sender{endpoint: route.Endpoint}
}
