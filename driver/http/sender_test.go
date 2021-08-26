package http

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSenderSend(t *testing.T) {
	srvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)

		reqb, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)
		assert.EqualValues(t, []byte("request-data"), reqb)
	}))

	sender := &sender{
		endpoint: srvr.URL,
	}

	err := sender.Send([]byte("request-data"))

	assert.Nil(t, err)
}

func TestSenderRequest(t *testing.T) {
	srvr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)

		n, err := w.Write([]byte("response-data"))

		assert.Equal(t, 13, n)
		assert.Nil(t, err)

		reqb, err := ioutil.ReadAll(r.Body)

		assert.Nil(t, err)
		assert.EqualValues(t, []byte("request-data"), reqb)
	}))

	sender := &sender{
		endpoint: srvr.URL,
	}

	respb, err := sender.Request([]byte("request-data"))

	assert.Nil(t, err)
	assert.Equal(t, []byte("response-data"), respb)
}

func TestSenderRequestOnError(t *testing.T) {
	sender := &sender{
		endpoint: "incorrect-address",
	}

	respb, err := sender.Request([]byte("request-data"))

	assert.Error(t, err)
	assert.Nil(t, respb)
}
