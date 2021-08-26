package response

import (
	"io/ioutil"
	"net/http"

	"NATter/errtpl"
	"NATter/log"

	"github.com/pkg/errors"
)

func prepareError(err error) int {
	if errors.Is(errors.Cause(err), errtpl.ErrNotFound) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	httperr := prepareError(err)

	http.Error(w, http.StatusText(httperr), httperr)

	if httperr != http.StatusInternalServerError {
		return
	}

	body := []byte("empty body")

	if r.Body != nil {
		var errread error

		body, errread = ioutil.ReadAll(r.Body)

		if errread != nil {
			body = []byte("broken")
		}
	}

	log.WithFields(log.Fields{
		"method": r.Method,
		"uri":    r.RequestURI,
		"status": httperr,
		"body":   string(body),
	}).Error(err)
}
