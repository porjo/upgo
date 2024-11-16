package upgo

import (
	"log"
	"net/http"
)

// credit: https://stackoverflow.com/a/51326483/202311

type withHeader struct {
	http.Header
	rt http.RoundTripper
}

func getRTWithHeader(rt http.RoundTripper) withHeader {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return withHeader{Header: make(http.Header), rt: rt}
}

func (h withHeader) RoundTrip(req *http.Request) (*http.Response, error) {

	// FIXME debug
	log.Printf("%s %s", req.Method, req.URL)

	if len(h.Header) == 0 {
		return h.rt.RoundTrip(req)
	}

	req = req.Clone(req.Context())
	for k, v := range h.Header {
		req.Header[k] = v
	}

	return h.rt.RoundTrip(req)
}
