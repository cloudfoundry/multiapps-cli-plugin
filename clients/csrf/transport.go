package csrf

import (
	"net/http"

	"github.com/jinzhu/copier"
)

type Csrf struct {
	Header string
	Token  string
}

type Transport struct {
	Transport http.RoundTripper
	Csrf      *Csrf
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := http.Request{}
	copier.Copy(&req2, req)
	if t.Csrf.Header != "" && t.Csrf.Token != "" {
		req2.Header.Set(t.Csrf.Header, t.Csrf.Token)
	}
	res, err := t.Transport.RoundTrip(&req2)
	if err != nil {
		return res, err
	}
	csrfHeader, csrfToken := res.Header.Get("X-Csrf-Header"), res.Header.Get("X-Csrf-Token")
	if csrfHeader != "" && csrfToken != "" {
		t.Csrf.Header, t.Csrf.Token = csrfHeader, csrfToken
	}
	return res, err
}
