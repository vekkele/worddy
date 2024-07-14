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
		name              string
		loggedIn          bool
		expectedCode      int
		expectedLocation  string
		expectedElemenets []string
	}{
		{
			name:              "No user logged in",
			expectedCode:      http.StatusOK,
			expectedElemenets: []string{`a[href="/user/signup"]`, `a[href="/user/login"]`},
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
				authedSessionMiddleware := LoadAndSaveMock(app.sessionManager, "authenticatedUserID", int64(1))
				h = authedSessionMiddleware(h)
			}

			ts := newTestServer(t, h)
			defer ts.Close()

			code, header, doc := ts.get(t, "/")

			assert.Equal(t, tc.expectedCode, code)

			if tc.expectedLocation != "" {
				assert.Equal(t, tc.expectedLocation, header.Get("Location"))
			}

			for _, selector := range tc.expectedElemenets {
				assert.True(t, doc.Find(selector).Length() > 0, "Page must contain element for selector: %s", selector)
			}
		})
	}
}
