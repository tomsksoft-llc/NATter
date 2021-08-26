package http

import (
	"net/url"
	"regexp"
	"sync"

	"NATter/driver"
	"NATter/log"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

type receiver struct {
	mux *chi.Mux
	wg  *sync.WaitGroup

	async    bool
	uri      string
	endpoint string
}

func (r *receiver) Listen(sender driver.Sender) error {
	if err := r.validateURI(); err != nil {
		return err
	}

	r.mux.Post(r.uri, routeHTTP(func(payload []byte) ([]byte, error) {
		log.Debugf("received request from uri: %s", r.uri)

		return nil, sender.Send(payload)
	}))

	return nil
}

func (r *receiver) ListenRequest(sender driver.Sender) error {
	if err := r.validateURI(); err != nil {
		return err
	}

	r.mux.Post(r.uri, routeHTTP(func(payload []byte) ([]byte, error) {
		log.Debugf("received request from uri: %s", r.uri)

		if !r.async {
			return sender.Request(payload)
		}

		r.asyncRequest(sender, payload)

		return nil, nil
	}))

	return nil
}

func (r *receiver) validateURI() error {
	if _, err := url.ParseRequestURI(r.uri); err != nil {
		return err
	}

	if regexp.MustCompile(reservedURIPattern).MatchString(r.uri) {
		return errors.Errorf("use of reserved uri pattern: %s", r.uri)
	}

	return nil
}

func (r *receiver) asyncRequest(sender driver.Sender, payload []byte) {
	r.wg.Add(1)

	go func() {
		defer r.wg.Done()

		respb, err := sender.Request(payload)

		if err != nil {
			log.Error(err)

			return
		}

		if _, err := request(r.endpoint, respb); err != nil {
			log.Error(err)

			return
		}

		log.Debugf("responded to endpoint: %s", r.endpoint)
	}()
}
