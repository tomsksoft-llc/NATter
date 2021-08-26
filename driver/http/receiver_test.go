package http

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	m "NATter/mock"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestReceiverListen(t *testing.T) {
	receiver := &receiver{
		mux: chi.NewRouter(),
		wg:  &sync.WaitGroup{},
		uri: "/path",
	}

	sender := &m.DriverSender{}

	sender.
		On("Send", []byte("some-data")).
		Return(nil)

	err := receiver.Listen(sender)

	assert.Nil(t, err)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/path",
		strings.NewReader("some-data"),
	)

	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	receiver.mux.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestReceiverListenOnSendError(t *testing.T) {
	receiver := &receiver{
		mux: chi.NewRouter(),
		wg:  &sync.WaitGroup{},
		uri: "/path",
	}

	sender := &m.DriverSender{}

	sender.
		On("Send", []byte("some-data")).
		Return(errors.New("error"))

	err := receiver.Listen(sender)

	assert.Nil(t, err)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/path",
		strings.NewReader("some-data"),
	)

	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	receiver.mux.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusInternalServerError, resp.Code)
}

func TestReceiverListenOnIncorrectURI(t *testing.T) {
	receiver := &receiver{
		mux: chi.NewRouter(),
		wg:  &sync.WaitGroup{},
		uri: "incorrect-path",
	}

	err := receiver.Listen(&m.DriverSender{})

	assert.Error(t, err)
}

func TestReceiverListenRequest(t *testing.T) {
	receiver := &receiver{
		mux: chi.NewRouter(),
		wg:  &sync.WaitGroup{},
		uri: "/path",
	}

	sender := &m.DriverSender{}

	sender.
		On("Request", []byte("request-data")).
		Return([]byte("response-data"), nil)

	err := receiver.ListenRequest(sender)

	assert.Nil(t, err)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/path",
		strings.NewReader("request-data"),
	)

	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	receiver.mux.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "response-data", resp.Body.String())
}

func TestReceiverListenRequestOnAsync(t *testing.T) {
	srvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)

		reqb, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)
		assert.EqualValues(t, []byte("response-data"), reqb)
	}))

	receiver := &receiver{
		mux:      chi.NewRouter(),
		wg:       &sync.WaitGroup{},
		async:    true,
		uri:      "/path",
		endpoint: srvr.URL,
	}

	sender := &m.DriverSender{}

	sender.
		On("Request", []byte("request-data")).
		Return([]byte("response-data"), nil)

	err := receiver.ListenRequest(sender)

	assert.Nil(t, err)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/path",
		strings.NewReader("request-data"),
	)

	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	receiver.mux.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestReceiverListenRequestOnAsyncRequestError(t *testing.T) {
	receiver := &receiver{
		mux:      chi.NewRouter(),
		wg:       &sync.WaitGroup{},
		async:    true,
		uri:      "/path",
		endpoint: "endpoint.com",
	}

	sender := &m.DriverSender{}

	sender.
		On("Request", []byte("request-data")).
		Return([]byte(nil), errors.New("error"))

	err := receiver.ListenRequest(sender)

	assert.Nil(t, err)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/path",
		strings.NewReader("request-data"),
	)

	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	receiver.mux.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestReceiverListenRequestOnAsyncResponseError(t *testing.T) {
	srvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)

		reqb, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)
		assert.EqualValues(t, []byte("response-data"), reqb)
	}))

	receiver := &receiver{
		mux:      chi.NewRouter(),
		wg:       &sync.WaitGroup{},
		async:    true,
		uri:      "/path",
		endpoint: srvr.URL,
	}

	sender := &m.DriverSender{}

	sender.
		On("Request", []byte("request-data")).
		Return([]byte("response-data"), nil)

	err := receiver.ListenRequest(sender)

	assert.Nil(t, err)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/path",
		strings.NewReader("request-data"),
	)

	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	receiver.mux.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestReceiverListenRequestOnReservedURI(t *testing.T) {
	receiver := &receiver{
		mux: chi.NewRouter(),
		wg:  &sync.WaitGroup{},
		uri: "/i/reserved-pattern",
	}

	err := receiver.ListenRequest(&m.DriverSender{})

	assert.Error(t, err)
}
