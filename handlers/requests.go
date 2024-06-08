package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

type ApiRequest interface {
	R() *http.Request
	Values() (map[string]interface{}, error)
	Validate() error
}

type PathValueExtractor func(r *http.Request) (map[string]interface{}, error)

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

func (r *apiRequest) Values() (vs map[string]interface{}, err error) {
	vs, err = r.pathValueExtractor(r.R())
	if err != nil {
		return vs, err
	}

	err = r.R().ParseForm()
	if err != nil {
		return vs, err
	}
	for k := range r.R().Form {
		v := r.R().Form.Get(k)
		vs[k] = v
	}

	return vs, nil
}

func (r *apiRequest) Validate(rules map[string]interface{}) (err error) {
	_, err = r.Values()
	if err != nil {
		return err
	}

	return nil
}
