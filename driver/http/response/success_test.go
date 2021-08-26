package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderWithStatus(t *testing.T) {
	resp := httptest.NewRecorder()

	RenderWithStatus(struct {
		Data string `json:"data"`
	}{
		Data: "text",
	}, resp, http.StatusCreated)

	assert.Equal(t, "application/json; charset=utf-8", resp.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, `{"data":"text"}`+"\n", resp.Body.String())
}

func TestRenderWithStatusOnMarshalError(t *testing.T) {
	resp := httptest.NewRecorder()

	RenderWithStatus(struct {
		Data chan bool `json:"data"`
	}{
		Data: make(chan bool),
	}, resp, http.StatusCreated)

	assert.Equal(t, "application/json; charset=utf-8", resp.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusCreated, resp.Code)
	assert.Equal(t, "{}\n", resp.Body.String())
}

func TestRender(t *testing.T) {
	resp := httptest.NewRecorder()

	Render(struct {
		Data string `json:"data"`
	}{
		Data: "text",
	}, resp)

	assert.Equal(t, "application/json; charset=utf-8", resp.Header().Get("Content-Type"))
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, `{"data":"text"}`+"\n", resp.Body.String())
}
