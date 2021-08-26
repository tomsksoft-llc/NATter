package response

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"NATter/errtpl"
	m "NATter/mock"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestProcessError(t *testing.T) {
	inputs := []struct {
		err        error
		statusText string
		statusCode int
	}{
		{
			err:        errtpl.ErrNotFound,
			statusText: http.StatusText(http.StatusNotFound),
			statusCode: http.StatusNotFound,
		},
		{
			err:        errors.New("unknown error"),
			statusText: http.StatusText(http.StatusInternalServerError),
			statusCode: http.StatusInternalServerError,
		},
	}

	for i, input := range inputs {
		w := httptest.NewRecorder()
		RenderError(w, &http.Request{}, input.err)

		assert.Equalf(t, fmt.Sprintf("%s\n", input.statusText), w.Body.String(), "case %d", i+1)
		assert.Equalf(t, input.statusCode, w.Code, "case %d", i+1)
	}
}

func TestProcessErrorOnBrokenBody(t *testing.T) {
	w := httptest.NewRecorder()
	RenderError(w, &http.Request{Body: &m.ReadCloserErr{}}, errors.New("error"))

	assert.Equal(t, fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)), w.Body.String())
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
