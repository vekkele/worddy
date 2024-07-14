package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	t.Run("Not logged in", func(t *testing.T) {
		app := newTestApplication()

		ts := newTestServer(t, app.routes())
		defer ts.Close()

		code, _, doc := ts.get(t, "/")

		assert.Equal(t, http.StatusOK, code)
		assert.True(t, doc.Find(`a[href="/user/signup"]`).Length() > 0, `Page must contain link to "/user/signup"`)
		assert.True(t, doc.Find(`a[href="/user/login"]`).Length() > 0, `Page must contain link to "/user/login"`)
	})

	t.Run("Logged in", func(t *testing.T) {
		app := newTestApplication()

		authedSessionMiddleware := LoadAndSaveMock(app.sessionManager, "authenticatedUserID", int64(1))

		ts := newTestServer(t, authedSessionMiddleware(app.routes()))
		defer ts.Close()

		code, header, _ := ts.get(t, "/")

		assert.Equal(t, http.StatusSeeOther, code)
		assert.Equal(t, "/dashboard", header.Get("Location"))
	})
}
