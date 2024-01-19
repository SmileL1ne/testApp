package handlers

import (
	"net/http"
	"runtime/debug"
)

func (r *routes) serverError(w http.ResponseWriter, req *http.Request, err error) {
	var (
		method = req.Method
		uri    = req.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	r.logger.Error(err.Error(), "method", method, "uri", uri, "stack", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (r *routes) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (r *routes) identifyStatus(w http.ResponseWriter, req *http.Request, status int, err error) {
	switch {
	case status >= 500:
		r.serverError(w, req, err)
	case status >= 400:
		r.clientError(w, status)
	default:
		r.logger.Error("Unknown status")
		r.clientError(w, status)
	}
}
