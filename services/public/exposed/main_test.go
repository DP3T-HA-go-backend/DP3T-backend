package main

import (
	"log"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestExposed(t *testing.T) {
	var req *http.Request
	var err error
	var rr *httptest.ResponseRecorder

	err = initConfig("config.ini")
	if err != nil {
		t.Fatal(err)
	}

	// Disable logging
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	t.Log("GET without date")
	req, err = http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = response(req, "GET", "/:date", exposed)
	if rr.Code != 404 {
		t.Error("Expected response code to be 404")
	}

	t.Log("GET with date")
	req, err = http.NewRequest("GET", "/1234", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr = response(req, "GET", "/:date", exposed)
	if rr.Code != 200 {
		t.Error("Expected response code to be 200")
	}
}

type handler = func(w http.ResponseWriter, r *http.Request, param httprouter.Params)

// Mocks a handler and returns a httptest.ResponseRecorder
func response(req *http.Request, method string, path string, h handler) *httptest.ResponseRecorder {
	router := httprouter.New()
	router.Handle(method, path, h)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	return rr
}
