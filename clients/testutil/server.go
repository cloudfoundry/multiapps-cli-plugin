package testutil

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/go-openapi/runtime"
)

// NewGetXMLOKServer creates a new server that responds to GET requests on the specified path
// and accepting application/xml with 200 OK and the specified payload, and to all other
// requests with 400 Bad Request.
func NewGetXMLOKServer(path string, payload []byte) *httptest.Server {
	return newServer(newStandardHandlerFunc(http.MethodGet, path, runtime.XMLMime, http.StatusOK, payload, nil, nil))
}

// NewGetXMLOAuthOKServer creates a new server that responds to GET requests on the specified path,
// accepting application/xml, and having the specified OAuth token with 200 OK and the specified
// payload, and to all other requests with 400 Bad Request.
func NewGetXMLOAuthOKServer(path string, token string, payload []byte) *httptest.Server {
	return newServer(newStandardHandlerFunc(http.MethodGet, path, runtime.XMLMime, http.StatusOK, payload,
		[]requestChecker{oauthChecker{token}}, nil))
}

// NewPostXMLOKServer creates a new server that responds to POST requests on the specified path,
// accepting application/xml, and having content type also application/xml and the specified
// body with 200 OK and the specified payload, and to all other requests with 400 Bad Request.
func NewPostXMLOKServer(path string, body []byte, payload []byte) *httptest.Server {
	return newServer(newStandardHandlerFunc(http.MethodPost, path, runtime.XMLMime, http.StatusOK, payload,
		[]requestChecker{bodyChecker{runtime.XMLMime, body}}, nil))
}

// NewPostFileXMLOKServer creates a new server that responds to POST requests on the specified path,
// accepting application/xml, and having content type multipart/form-data and the specified
// file key and content with 200 OK and the specified payload, and to all other requests with
// 400 Bad Request.
func NewPostFileXMLOKServer(path string, key string, content []byte, payload []byte) *httptest.Server {
	return newServer(newStandardHandlerFunc(http.MethodPost, path, runtime.XMLMime, http.StatusOK, payload,
		[]requestChecker{fileChecker{key, content}}, nil))
}

// NewDeleteNoContentServer creates a new server that responds to DELETE requests on the
// specified path with 204 No Content, and to all other requests with 400 Bad Request.
func NewDeleteNoContentServer(path string) *httptest.Server {
	return newServer(newStandardHandlerFunc(http.MethodDelete, path, "", http.StatusNoContent, nil, nil, nil))
}

// NewGetXMLOKDeleteNoContentCsrfServer creates a new composite server that responds to
// GET and DELETE requests on the specified path and implements a CSRF protection mechanism.
func NewGetXMLOKDeleteNoContentCsrfServer(path string, payload []byte) *httptest.Server {
	ch := csrfCheckerHandler{"dummy"}
	return newServer(
		newStandardHandlerFunc(http.MethodGet, path, runtime.XMLMime, http.StatusOK, payload, nil, []responseHandler{ch}),
		newStandardHandlerFunc(http.MethodDelete, path, "", http.StatusNoContent, nil, []requestChecker{ch}, nil),
	)
}

// NewStatusServer creates a new server that responds to all requests with the specified status and
// no payload.
func NewStatusServer(status int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(status)
	}))
}

type handlerFunc func(rw http.ResponseWriter, req *http.Request) bool

type requestChecker interface {
	check(req *http.Request) (bool, error)
}

type responseHandler interface {
	handle(rw http.ResponseWriter) error
}

func newServer(funcs ...handlerFunc) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		for _, f := range funcs {
			if f(rw, req) {
				return
			}
		}
		rw.WriteHeader(http.StatusBadRequest)
	}))
}

func newStandardHandlerFunc(method, path, accept string, status int, payload []byte,
	checkers []requestChecker, handlers []responseHandler) handlerFunc {
	return newHandlerFunc(
		append([]requestChecker{standardChecker{method, path, accept}}, checkers...),
		append(handlers, standardHandler{status, accept, payload}),
	)
}

func newHandlerFunc(checkers []requestChecker, handlers []responseHandler) handlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) bool {
		ok, err := checkRequest(req, checkers)
		if err != nil {
			handleError(rw, err)
			return true
		}
		if !ok {
			return false
		}
		err = handleResponse(rw, handlers)
		if err != nil {
			handleError(rw, err)
			return true
		}
		return true
	}
}

func handleError(rw http.ResponseWriter, err error) {
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write([]byte(err.Error()))
}

func checkRequest(req *http.Request, checkers []requestChecker) (bool, error) {
	if checkers != nil {
		for _, c := range checkers {
			ok, err := c.check(req)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
		}
	}
	return true, nil
}

func handleResponse(rw http.ResponseWriter, handlers []responseHandler) error {
	if handlers != nil {
		for _, h := range handlers {
			err := h.handle(rw)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type standardChecker struct {
	method string
	path   string
	accept string
}

func (c standardChecker) check(req *http.Request) (bool, error) {
	return req.Method == c.method &&
		req.URL.Path == c.path &&
		(c.accept == "" || req.Header.Get(runtime.HeaderAccept) == c.accept), nil
}

type standardHandler struct {
	status      int
	contentType string
	payload     []byte
}

func (h standardHandler) handle(rw http.ResponseWriter) error {
	if h.contentType != "" {
		rw.Header().Add(runtime.HeaderContentType, h.contentType)
	}
	rw.WriteHeader(h.status)
	if h.payload != nil {
		_, err := rw.Write(h.payload)
		if err != nil {
			return err
		}
	}
	return nil
}

type oauthChecker struct {
	token string
}

func (c oauthChecker) check(req *http.Request) (bool, error) {
	return req.Header.Get("Authorization") == "Bearer "+c.token, nil
}

type bodyChecker struct {
	contentType string
	body        []byte
}

func (c bodyChecker) check(req *http.Request) (bool, error) {
	if req.Header.Get(runtime.HeaderContentType) != c.contentType {
		return false, nil
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return false, err
	}
	return bytes.Equal(body, c.body), nil
}

type fileChecker struct {
	key     string
	content []byte
}

func (c fileChecker) check(req *http.Request) (bool, error) {
	if !strings.HasPrefix(req.Header.Get(runtime.HeaderContentType), runtime.MultipartFormMime) {
		return false, nil
	}
	file, _, err := req.FormFile(c.key)
	if err != nil {
		return false, err
	}
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return false, err
	}
	return bytes.Equal(buf.Bytes(), c.content), nil
}

const csrfHeader = "X-CSRF-TOKEN"

type csrfCheckerHandler struct {
	token string
}

func (ch csrfCheckerHandler) check(req *http.Request) (bool, error) {
	return req.Header.Get(csrfHeader) == ch.token, nil
}

func (ch csrfCheckerHandler) handle(rw http.ResponseWriter) error {
	rw.Header().Add("X-Csrf-Header", csrfHeader)
	rw.Header().Add("X-Csrf-Token", ch.token)
	return nil
}
