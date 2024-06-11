package responses

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

func ErrUnauthorized(userErr, err interface{}) render.Renderer {
	userErrStr, errStr := extractErrorStrings(userErr, err)
	if userErrStr == "" {
		userErrStr = "authorization is required to perform this action"
	}
	return &ErrResponse{
		Err:            errors.New(errStr),
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     http.StatusText(http.StatusUnauthorized),
		Error:          errStr,
		UserError:      userErrStr,
	}
}

func ErrForbidden(userErr, err interface{}) render.Renderer {
	userErrStr, errStr := extractErrorStrings(userErr, err)
	if userErrStr == "" {
		userErrStr = "you do not have permission to perform this action."
	}
	return &ErrResponse{
		Err:            errors.New(errStr),
		HTTPStatusCode: http.StatusForbidden,
		StatusText:     http.StatusText(http.StatusForbidden),
		Error:          errStr,
		UserError:      userErrStr,
	}
}

func ErrBadRequest(userErr, err interface{}) render.Renderer {
	userErrStr, errStr := extractErrorStrings(userErr, err)
	return &ErrResponse{
		Err:            errors.New(errStr),
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     http.StatusText(http.StatusBadRequest),
		Error:          errStr,
		UserError:      userErrStr,
	}
}

func ErrRender(userErr, err interface{}) render.Renderer {
	userErrStr, errStr := extractErrorStrings(userErr, err)
	if err == nil {
		err = errors.New("")
	}
	return &ErrResponse{
		Err:            errors.New(errStr),
		HTTPStatusCode: http.StatusUnprocessableEntity,
		Success:        false,
		StatusText:     "Error rendering response.",
		Error:          errStr,
		UserError:      userErrStr,
	}
}

func ErrNotFound(userErr, err interface{}) render.Renderer {
	userErrStr, errStr := extractErrorStrings(userErr, err)
	if userErrStr == "" {
		userErrStr = "resource not found"
	}
	return &ErrResponse{
		Err:            errors.New(errStr),
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     http.StatusText(http.StatusNotFound),
		Error:          errStr,
		UserError:      userErrStr,
	}
}

func ErrInternalServerError(userErr, err interface{}) render.Renderer {
	userErrStr, errStr := extractErrorStrings(userErr, err)
	return &ErrResponse{
		Err:            errors.New(errStr),
		HTTPStatusCode: http.StatusInternalServerError,
		Success:        false,
		StatusText:     http.StatusText(http.StatusInternalServerError),
		Error:          errStr,
		UserError:      userErrStr,
	}
}

func ErrValidationError(userErr, err interface{}) render.Renderer {
	userErrStr, errStr := extractErrorStrings(userErr, err)
	return &ErrResponse{
		Err:            errors.New(errStr),
		HTTPStatusCode: http.StatusInternalServerError,
		Success:        false,
		Error:          errStr,
		UserError:      userErrStr,
	}
}

func extractErrorStrings(userErr, err interface{}) (string, string) {
	userErrString := ""
	if value, ok := userErr.(string); ok {
		userErrString = value
	} else if value, ok := userErr.(error); ok {
		userErrString = value.Error()
	} else if userErr != nil {
		logrus.Errorf("extractErrorStrings: unknown user error type: %v", userErr)
	}

	errString := ""
	if value, ok := err.(string); ok {
		errString = value
	} else if value, ok := err.(error); ok {
		errString = value.Error()
	} else if err != nil {
		logrus.Errorf("extractErrorStrings: unknown error type: %v", err)
	}

	return userErrString, errString
}
