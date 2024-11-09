package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/vekkele/worddy/internal/store/mocks"
)

func newTestApplication() *application {
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour

	return &application{
		users:          &mocks.UserModel{},
		words:          &mocks.WordModel{},
		formDecoder:    form.NewDecoder(),
		sessionManager: sessionManager,
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewServer(h)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func (ts *testServer) get(t *testing.T, urlPath string) (int, http.Header, *goquery.Document) {
	rs, err := ts.Client().Get(ts.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	doc, err := goquery.NewDocumentFromReader(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, doc
}

func (ts *testServer) postForm(t *testing.T, urlPath string, form url.Values) (int, http.Header, *goquery.Document) {
	rs, err := ts.Client().PostForm(ts.URL+urlPath, form)
	if err != nil {
		t.Fatal(err)
	}

	defer rs.Body.Close()

	doc, err := goquery.NewDocumentFromReader(rs.Body)
	if err != nil {
		t.Fatal(err)
	}

	return rs.StatusCode, rs.Header, doc
}

func loggedInStubMiddleware(session *scs.SessionManager, h http.Handler) http.Handler {
	return loadAndSaveStub(session, "authenticatedUserID", int64(1))(h)
}

func loadAndSaveStub(session *scs.SessionManager, key string, value any) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return session.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session.Put(r.Context(), key, value)
			next.ServeHTTP(w, r)
		}))
	}
}

func extractCSRFToken(doc *goquery.Document) string {
	token, _ := doc.Find(`input[name="csrf_token"]`).First().Attr("value")

	return token
}
