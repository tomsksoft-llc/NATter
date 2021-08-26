package http

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"

	"NATter/log"

	"github.com/pkg/errors"
)

func request(endpoint string, payload []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(payload),
	)

	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	log.Debugf("requested to endpoint: %s", endpoint)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected response code %d from endpoint", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}
