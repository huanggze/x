package decoderx

type (
	// HTTP decodes json and form-data from HTTP Request Bodies.
	HTTP struct{}
)

// NewHTTP creates a new HTTP decoder.
func NewHTTP() *HTTP {
	return new(HTTP)
}
