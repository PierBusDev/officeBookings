package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler
	httpHandler := NoSurf(&myH)

	switch v := httpHandler.(type) {
	case http.Handler:
		//do nothing, it's the correct one
	default:
		t.Errorf("type is not http.Handler but %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	httpHandler := SessionLoad(&myH)

	switch v := httpHandler.(type) {
	case http.Handler:
		//do nothing, it's the correct one
	default:
		t.Errorf("type is not http.Handler but %T", v)
	}
}
