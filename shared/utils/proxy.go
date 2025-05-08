package utils

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func ProxyHandler(target string) http.HandlerFunc {
	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)

	return func(w http.ResponseWriter, r *http.Request) {
		r.Host = url.Host
		proxy.ServeHTTP(w, r)
	}
}
