package foo

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// standard library => http test server
func TestHandleGetFoo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleGetFoo))
	res, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected 200 but got %d", res.StatusCode)
	}
	defer res.Body.Close()

	expected := "Foo"
	b, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
	if string(b) != expected {
		t.Errorf("Expected %s but got %s", expected, string(b))

	}
}

// Response Record
func TestHandleGetFooRR(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "?userId=123", nil)
	if err != nil {
		t.Error(err)
	}

	handleGetFoo(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected 200 but got %d", rr.Result().StatusCode)
	}
	defer rr.Result().Body.Close()

	expected := "Foo"
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Error(err)
	}
	if string(b) != expected {
		t.Errorf("Expected %s but got %s", expected, string(b))

	}
}
