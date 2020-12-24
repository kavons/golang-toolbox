package errors

import (
	"testing"

	"fmt"
	"io"
	"reflect"

	"github.com/pkg/errors"
)

func TestWrapNil(t *testing.T) {
	got := errors.Wrap(nil, "no error")
	if got != nil {
		t.Errorf("Wrap(nil, \"no error\"): got %#v, expected nil", got)
	}
}

func TestWrap(t *testing.T) {
	tests := []struct {
		err     error
		message string
		want    string
	}{
		{io.EOF, "read error", "read error: EOF"},
		{errors.Wrap(io.EOF, "read error"), "client error", "client error: read error: EOF"},
	}

	for _, tt := range tests {
		got := errors.Wrap(tt.err, tt.message).Error()
		if got != tt.want {
			t.Errorf("Wrap(%v, %q): got: %v, want %v", tt.err, tt.message, got, tt.want)
		}
	}
}

type nilError struct{}

func (nilError) Error() string { return "nil error" }

func TestCause(t *testing.T) {
	x := errors.New("error")
	tests := []struct {
		err  error
		want error
	}{{
		// nil error is nil
		err:  nil,
		want: nil,
	}, {
		// explicit nil error is nil
		err:  (error)(nil),
		want: nil,
	}, {
		// typed nil is nil
		err:  (*nilError)(nil),
		want: (*nilError)(nil),
	}, {
		// uncaused error is unaffected
		err:  io.EOF,
		want: io.EOF,
	}, {
		// caused error returns cause
		err:  errors.Wrap(io.EOF, "ignored"),
		want: io.EOF,
	}, {
		err:  x, // return from errors.New
		want: x,
	}, {
		errors.WithMessage(nil, "whoops"),
		nil,
	}, {
		errors.WithMessage(io.EOF, "whoops"),
		io.EOF,
	}, {
		errors.WithStack(nil),
		nil,
	}, {
		errors.WithStack(io.EOF),
		io.EOF,
	}}

	for i, tt := range tests {
		got := errors.Cause(tt.err)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("test %d: got %#v, want %#v", i+1, got, tt.want)
		}
	}
}

func returnError() (status int, err error) {
	return 0, errors.New("bad")
}

func ReturnError() (err error) {
	s, err := returnError()
	if err != nil {
		return
	}
	fmt.Println(s)

	return
}

func TestErrorShading(t *testing.T) {
	t.Log(ReturnError())
}
