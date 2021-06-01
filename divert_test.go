package divert

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testValue = "test-divert-header-value"

func TestAddToContext(t *testing.T) {
	if expected, actual := testValue, FromContext(AddToContext(context.Background(), testValue)); actual != expected {
		t.Errorf("context value incorrect, expected %s, actual %s", expected, actual)
	}
}

func TestSetHeader(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Error(err)
	}

	SetHeader(AddToContext(context.Background(), testValue), r)

	if expected, actual := testValue, r.Header.Get(DivertHeaderName); actual != expected {
		t.Errorf("header not set, expected %s, actual %s", expected, actual)
	}
}

func TestFromContextNoValue(t *testing.T) {
	if expected, actual := "", FromContext(context.Background()); actual != expected {
		t.Errorf("missing value not handled, expected %s, actual %s", expected, actual)
	}
}

func TestFromHeaders(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Error(err)
	}
	r.Header.Set(DivertHeaderName, testValue)

	if expected, actual := testValue, FromHeaders(r); actual != expected {
		t.Errorf("header not retrieved, expected %s, actual %s", expected, actual)
	}
}

func TestInjectDivertHeader(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		t.Error(err)
	}
	r.Header.Set(DivertHeaderName, testValue)

	h := InjectDivertHeader()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if expected, actual := testValue, FromContext(r.Context()); actual != expected {
			t.Errorf("header not injected, expected %s, actual %s", expected, actual)
		}

	}))
	h.ServeHTTP(httptest.NewRecorder(), r)
}
