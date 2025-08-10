package httpx

import "net/http"

// WrapTransportWithHeader wraps a http.Transport to always use the values from the given header.
func WrapTransportWithHeader(parent http.RoundTripper, h http.Header) *TransportWithHeader {
	return &TransportWithHeader{
		RoundTripper: parent,
		h:            h,
	}
}

// TransportWithHeader is an http.RoundTripper that always uses the values from the given header.
type TransportWithHeader struct {
	http.RoundTripper
	h http.Header
}
