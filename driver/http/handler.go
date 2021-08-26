package http

import (
	"io/ioutil"
	"net/http"

	"NATter/driver/http/response"
	"NATter/entity"
)

func routeHTTP(handler func([]byte) ([]byte, error)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		reqb, err := ioutil.ReadAll(r.Body)

		if err != nil {
			response.RenderError(w, r, err)

			return
		}

		respb, err := handler(reqb)

		if err != nil {
			response.RenderError(w, r, err)

			return
		}

		if _, err := w.Write(respb); err != nil {
			response.RenderError(w, r, err)

			return
		}
	}
}

func responseRoutes(routes []*entity.Route) []*entity.Route {
	if routes == nil {
		routes = []*entity.Route{}
	}

	return routes
}

func routes(handler func() []*entity.Route) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		routes := handler()

		response.Render(responseRoutes(routes), w)
	}
}
