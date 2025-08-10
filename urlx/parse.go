package urlx

import "net/url"

// ParseOrPanic parses a url or panics.
func ParseOrPanic(in string) *url.URL {
	out, err := url.Parse(in)
	if err != nil {
		panic(err.Error())
	}
	return out
}
