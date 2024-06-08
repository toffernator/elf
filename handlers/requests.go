package handlers

import (
	"net/http"
	"net/url"
)

type ApiRequest interface {
	R() *http.Request
	Values() (url.Values, error)
}

type PathValueExtractor func(r *http.Request) (url.Values, error)

type apiRequest struct {
	r                  *http.Request
	pathValueExtractor PathValueExtractor
}

func NewApiRequest(r *http.Request, p PathValueExtractor) *apiRequest {
	return &apiRequest{r: r, pathValueExtractor: p}

}

func (r *apiRequest) R() *http.Request {
	return r.r
}

func (r *apiRequest) Values() (vs url.Values, err error) {
	vs, err = r.pathValueExtractor(r.R())
	if err != nil {
		return vs, err
	}

	r.R().ParseForm()
	if err != nil {
		return vs, err
	}
	for k, v := range r.R().Form {
		vs[k] = v
	}

	return vs, nil
}
