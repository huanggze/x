package tlsx

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"

	"github.com/huanggze/x/watcherx"
)

// ErrNoCertificatesConfigured is returned when no TLS configuration was found.
var ErrNoCertificatesConfigured = errors.New("no tls configuration was found")

// CertificateFromBase64 loads a TLS certificate from a base64-encoded string of
// the PEM representations of the cert and key.
func CertificateFromBase64(certBase64, keyBase64 string) (tls.Certificate, error) {
	certPEM, err := base64.StdEncoding.DecodeString(certBase64)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to base64 decode the TLS certificate: %v", err)
	}
	keyPEM, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to base64 decode the TLS private key: %v", err)
	}
	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("unable to load X509 key pair: %v", err)
	}
	return cert, nil
}

// GetCertificate returns a function for use with
// "net/tls".Config.GetCertificate.
//
// The certificate and private key are read from the specified filesystem paths.
// The certificate file is watched for changes, upon which the cert+key are
// reloaded in the background. Errors during reloading are deduplicated and
// reported through the errs channel if it is not nil. When the provided context
// is canceled, background reloading stops and the errs channel is closed.
//
// The returned function always yields the latest successfully loaded
// certificate; ClientHelloInfo is unused.
func GetCertificate(
	ctx context.Context,
	certPath, keyPath string,
	errs chan<- error,
) (func(*tls.ClientHelloInfo) (*tls.Certificate, error), error) {
	if certPath == "" || keyPath == "" {
		return nil, errors.WithStack(ErrNoCertificatesConfigured)
	}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, errors.WithStack(fmt.Errorf("unable to load X509 key pair from files: %v", err))
	}
	var store atomic.Value
	store.Store(&cert)

	events := make(chan watcherx.Event)
	// The cert could change without the key changing, but not the other way around.
	// Hence, we only watch the cert.
	_, err = watcherx.WatchFile(ctx, certPath, events)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	go func() {
		if errs != nil {
			defer close(errs)
		}
		var lastReportedError string
		for {
			select {
			case <-ctx.Done():
				return

			case event := <-events:
				var err error
				switch event := event.(type) {
				case *watcherx.ChangeEvent:
					var cert tls.Certificate
					cert, err = tls.LoadX509KeyPair(certPath, keyPath)
					if err == nil {
						store.Store(&cert)
						lastReportedError = ""
						continue
					}
					err = fmt.Errorf("unable to load X509 key pair from files: %v", err)

				case *watcherx.ErrorEvent:
					err = fmt.Errorf("file watch: %v", event)
				default:
					continue
				}

				if err.Error() == lastReportedError { // same message as before: don't spam the error channel
					continue
				}
				// fresh error
				select {
				case errs <- errors.WithStack(err):
					lastReportedError = err.Error()
				case <-time.After(500 * time.Millisecond):
				}
			}
		}
	}()

	return func(_ *tls.ClientHelloInfo) (*tls.Certificate, error) {
		if cert, ok := store.Load().(*tls.Certificate); ok {
			return cert, nil
		}
		return nil, errors.WithStack(ErrNoCertificatesConfigured)
	}, nil
}
