package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"NATter/entity"
	m "NATter/mock"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestConnServeOnInternalRoutes(t *testing.T) {
	conn := &conn{
		mux: chi.NewRouter(),
		wg:  &sync.WaitGroup{},
		routes: []*entity.Route{
			{
				Mode:  "receiver-sender-oneway",
				Topic: "topic1",
				URI:   "/path1",
			},
			{
				Mode:  "http-sender-twoway",
				Topic: "topic2",
				URI:   "/path2",
			},
		},
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/i/routes",
		strings.NewReader(""),
	)

	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)

	defer cancel()

	err = conn.Serve(ctx)

	conn.mux.ServeHTTP(resp, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t,
		`[{"mode":"receiver-sender-oneway","topic":"topic1","uri":"/path1"},`+
			`{"mode":"http-sender-twoway","topic":"topic2","uri":"/path2"}]`+"\n", resp.Body.String())
}

func TestConnServeOnInternalRoutesEmpty(t *testing.T) {
	conn := &conn{
		mux: chi.NewRouter(),
		wg:  &sync.WaitGroup{},
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/i/routes",
		strings.NewReader(""),
	)

	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)

	defer cancel()

	err = conn.Serve(ctx)

	conn.mux.ServeHTTP(resp, req)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, `[]`+"\n", resp.Body.String())
}

func TestConnClose(t *testing.T) {
	conn := NewConn(&ConnConfig{})

	err := conn.Close()

	assert.Nil(t, err)
}

func TestConnReceiver(t *testing.T) {
	conn := &conn{
		mux:    chi.NewRouter(),
		wg:     &sync.WaitGroup{},
		routes: []*entity.Route{},
	}

	receiver := conn.Receiver(&entity.Route{
		URI: "/path",
	})

	assert.Equal(t, []*entity.Route{
		{URI: "/path"},
	}, conn.routes)

	err := receiver.Listen(&m.DriverSender{})

	assert.Nil(t, err)

	routes := conn.mux.Routes()

	assert.Equal(t, 1, len(routes))
	assert.Equal(t, "/path", routes[0].Pattern)
}

func TestConnSender(t *testing.T) {
	conn := &conn{
		mux:    chi.NewRouter(),
		wg:     &sync.WaitGroup{},
		routes: []*entity.Route{},
	}

	srvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)

		reqb, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)
		assert.EqualValues(t, []byte("request-data"), reqb)
	}))

	sender := conn.Sender(&entity.Route{
		Endpoint: srvr.URL,
	})

	assert.Equal(t, []*entity.Route{
		{Endpoint: srvr.URL},
	}, conn.routes)

	err := sender.Send([]byte("request-data"))

	assert.Nil(t, err)
}
