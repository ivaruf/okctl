package fetch

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Fetcher provides an interface for downloading a file from a URL
type Fetcher interface {
	Fetch(io.Writer) (int64, error)
}

type staticFetcher struct {
	data []byte
}

// Fetch returns the static data
func (f *staticFetcher) Fetch(writer io.Writer) (int64, error) {
	return io.Copy(writer, bytes.NewReader(f.data))
}

// NewStaticFetcher creates a fetcher that returns the
// provided data, good for use with testing
func NewStaticFetcher(data []byte) Fetcher {
	return &staticFetcher{
		data: data,
	}
}

type httpFetcher struct {
	url string
}

// NewHTTPFetcher creates a fetcher for downloading a package.
func NewHTTPFetcher(url string) Fetcher {
	return &httpFetcher{
		url: url,
	}
}

// Fetch the content of the provided URL.
func (f *httpFetcher) Fetch(w io.Writer) (int64, error) {
	u, err := url.Parse(f.url)
	if err != nil {
		return 0, err
	}

	if u.Scheme != "https" {
		return 0, fmt.Errorf("a valid pkgURL must begin https://, got: %s", f.url)
	}

	resp, err := http.Get(f.url)
	if err != nil {
		return 0, err
	}

	defer func() {
		err = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad status: %s", resp.Status)
	}

	return io.Copy(w, resp.Body)
}
