package proxy

import (
	"net/http/httptest"
	"testing"

	"github.com/transientvariable/log-go"

	"github.com/stretchr/testify/assert"

	gohttp "net/http"
)

func TestBalancer(t *testing.T) {
	t.Skip("TODO: complete testing :)")

	h1 := httptest.NewServer(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		log.Info("[proxy:test] server 1", log.String("url", r.URL.String()))
	}))
	defer h1.Close()

	h2 := httptest.NewServer(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		log.Info("[proxy:test] server 2", log.String("url", r.URL.String()))
	}))
	defer h2.Close()

	h, err := prepareHosts(h1.URL, h2.URL)
	assert.NoError(t, err)

	_, err = NewBalancer(h)
	assert.NoError(t, err)

	//client := http.DefaultClient()
	//keys := "abc"
	//for _, key := range keys {
	//	req, err := gohttp.NewRequest("GET", "http://replace.me/foo/bar?key="+string(key), nil)
	//	assert.NoError(t, err)
	//
	//	resp, err := client.Do(req)
	//	assert.NoError(t, err)
	//
	//	_, err = io.Copy(os.Stdout, resp.Body)
	//	assert.NoError(t, err)
	//
	//	err := resp.Body.Close()
	//	assert.NoError(t, err)
	//}

	// Output:
	// A: /foo/bar?key=a
	// B: /foo/bar?key=b
	// A: /foo/bar?key=c
}

func prepareHosts(targets ...string) ([]*Host, error) {
	hosts := make([]*Host, len(targets))
	for i, t := range targets {
		h, err := NewHost(t)
		if err != nil {
			return hosts, err
		}
		hosts[i] = h
	}
	return hosts, nil
}
