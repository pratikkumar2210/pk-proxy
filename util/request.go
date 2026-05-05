package util

import "net/http"

func RequestURI(r *http.Request) string {
	return r.Host + r.URL.RequestURI()
}
