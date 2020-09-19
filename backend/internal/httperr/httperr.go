package httperr

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type StatusCoder interface {
	StatusCode() int
}

func ErrCode(err error) int {
	var sc StatusCoder

	if errors.As(err, &sc) {
		return sc.StatusCode()
	}

	return 500
}

type JSONError struct {
	Error string `json:"error"`
}

// WriteErr writes the error code into the status and body.
func WriteErr(w http.ResponseWriter, err error) error {
	w.WriteHeader(ErrCode(err))
	return json.NewEncoder(w).Encode(JSONError{err.Error()})
}

type basicError struct {
	code int
	msg  string
}

var (
	_ error       = (*basicError)(nil)
	_ StatusCoder = (*basicError)(nil)
)

func New(code int, msg string) error {
	return basicError{code, msg}
}

func (e basicError) Error() string {
	return e.msg
}

func (e basicError) StatusCode() int {
	return e.code
}

type wrapError struct {
	code int
	wrap error
}

var (
	_ error       = (*wrapError)(nil)
	_ StatusCoder = (*wrapError)(nil)
)

func Wrap(err error, code int, msg string) error {
	if err == nil {
		return nil
	}
	return wrapError{code, errors.Wrap(err, msg)}
}

func Wrapf(err error, code int, f string, v ...interface{}) error {
	if err == nil {
		return nil
	}
	return wrapError{code, errors.Wrapf(err, f, v...)}
}

func (e wrapError) Error() string {
	return e.wrap.Error()
}

func (e wrapError) StatusCode() int {
	return e.code
}
