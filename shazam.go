package shazam

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/exp/maps"
)

type wand struct {
	R *http.Request
	C []string
	Z string
}

// Boom take *http.Request and provides string giving a curl for in-usage request
func Boom(request *http.Request) (string, error) {
	var w wand
	if err := w.failSafe(request); err != nil {
		return "", err
	}
	w.startIntialCurl()
	err := w.addRequestFlag()
	if err != nil {
		return "", err
	}
	w.Z = strings.Join(w.C, " ")
	return w.Z, nil
}

func (w *wand) failSafe(request *http.Request) error {
	if request.Body != nil {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			if err != nil {
				return fmt.Errorf("shazam was unable to read body into bytes, error : %v", err)
			}
		}
		w.R = request.Clone(request.Context())
		request.Body = io.NopCloser(bytes.NewReader(body))
		w.R.Body = io.NopCloser(bytes.NewReader(body))
	} else {
		w.R = request.Clone(request.Context())

	}
	return nil
}

func (w *wand) startIntialCurl() {
	w.C = append(w.C, "curl --location")
}

func (w *wand) addRequestFlag() error {
	if w.R.URL == nil {
		return fmt.Errorf("URL of request is nil; couldn't proceed")
	}
	urlString := fmt.Sprintf("%s://%s%s", w.R.URL.Scheme, w.R.URL.Host, w.R.URL.Path)
	if len(maps.Keys(w.R.URL.Query())) != 0 {
		urlString = fmt.Sprintf("%s?%s", urlString, w.addSubstringForQueryParametersToURL())
	}
	w.C = append(w.C, fmt.Sprintf("--request %s '%s'", w.R.Method, urlString))

	for key, values := range w.R.Header {
		headerValue := strings.Join(values, ", ")
		w.C = append(w.C, fmt.Sprintf("--header '%s: %s'", key, headerValue))
	}

	if w.R.Body != nil {
		buf := new(strings.Builder)
		_, err := io.Copy(buf, w.R.Body)
		if err != nil {
			return fmt.Errorf("shazam was unable to read body into bytes, error : %v", err)
		}

		w.C = append(w.C, fmt.Sprintf("--data-raw '%s'", buf.String()))
	}
	return nil
}

func (w *wand) addSubstringForQueryParametersToURL() string {
	return w.R.URL.Query().Encode()
}
