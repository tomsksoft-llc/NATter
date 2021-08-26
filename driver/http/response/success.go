package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"NATter/errtpl"
	"NATter/log"
)

func Render(payload interface{}, w http.ResponseWriter) {
	RenderWithStatus(payload, w, http.StatusOK)
}

func RenderWithStatus(payload interface{}, w http.ResponseWriter, status int) {
	response, err := json.Marshal(payload)

	if err != nil {
		log.Error(errtpl.ErrMarshal(err, response))

		response = []byte("{}")
	}

	RenderPlainWithStatus(string(response), w, status)
}

func RenderPlainWithStatus(payload string, w http.ResponseWriter, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintln(w, payload)
}
