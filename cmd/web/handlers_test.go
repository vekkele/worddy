package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	app := newTestApplication()
	h := app.routes()

	tests := []struct {
		name             string
		loggedIn         bool
		expectedCode     int
		expectedLocation string
		expectedElements []string
	}{
		{
			name:             "No user logged in",
			expectedCode:     http.StatusOK,
			expectedElements: []string{`a[href="/user/signup"]`, `a[href="/user/login"]`},
		},
		{
			name:             "User logged in",
			loggedIn:         true,
			expectedCode:     http.StatusSeeOther,
			expectedLocation: "/dashboard",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.loggedIn {
				h = loggedInStubMiddleware(app.sessionManager, h)
			}

			ts := newTestServer(t, h)
			defer ts.Close()

			code, header, doc := ts.get(t, "/")

			assert.Equal(t, tc.expectedCode, code)

			if tc.expectedLocation != "" {
				assert.Equal(t, tc.expectedLocation, header.Get("Location"))
			}

			for _, selector := range tc.expectedElements {
				assert.True(t, doc.Find(selector).Length() > 0, "Page must contain element for selector: %s", selector)
			}
		})
	}
}
