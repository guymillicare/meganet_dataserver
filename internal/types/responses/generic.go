package responses

import (
	"net/http"

	"github.com/go-chi/render"
)

type GenericResponse struct{}

func (rd *GenericResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type SuccessResponse struct {
	GenericResponse
	Success bool   `json:"success"`
	Id      string `json:"id"`
}
type SuccessResponseId32 struct {
	GenericResponse
	Success bool  `json:"success"`
	Id      int32 `json:"id"`
}
type SuccessResponseId struct {
	GenericResponse
	Success bool `json:"success"`
	Id      int  `json:"id"`
}

func (rd *SuccessResponse) Render(w http.ResponseWriter, r *http.Request) error {
	rd.Success = true
	return nil
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	Success    bool   `json:"success"`
	StatusText string `json:"status"`               // http-level status message
	AppCode    int64  `json:"code,omitempty"`       // application-specific error code
	Error      string `json:"error,omitempty"`      // technical-level error message, for debugging
	UserError  string `json:"user_error,omitempty"` // user-level error message
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
