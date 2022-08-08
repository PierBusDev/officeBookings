package main

import (
	"net/http"
	"os"
	"testing"
)

//type just for test purposes
type myHandler struct{}

func (m *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
