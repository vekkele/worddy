package main

import (
	"net/http"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	t.Run("Not logged in", func(t *testing.T) {
		app := newTestApplication()

		ts := newTestServer(t, app.routes())

		resp, err := ts.Client().Get(ts.URL + "/")
		if err != nil {
			t.Fatal(err)
		}

		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		assert.True(t, doc.Find(`a[href="/user/signup"]`).Length() > 0, `Page must contain link to "/user/signup"`)
		assert.True(t, doc.Find(`a[href="/user/login"]`).Length() > 0, `Page must contain link to "/user/login"`)
	})

	t.Run("Logged in", func(t *testing.T) {
		app := newTestApplication()

		authedSessionMiddleware := LoadAndSaveMock(app.sessionManager, "authenticatedUserID", int64(1))

		ts := newTestServer(t, authedSessionMiddleware(app.routes()))

		resp, err := ts.Client().Get(ts.URL + "/")
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, http.StatusSeeOther, resp.StatusCode)
		assert.Equal(t, "/dashboard", resp.Header.Get("Location"))
	})
}
