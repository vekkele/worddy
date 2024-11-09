package main

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vekkele/worddy/internal/store/mocks"
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

func TestLoginPost(t *testing.T) {
	app := newTestApplication()
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, doc := ts.get(t, "/user/login")

	validCSRFToken := extractCSRFToken(doc)

	var (
		validEmail    = mocks.User.Email
		validPassword = mocks.UserPassword
	)

	tests := []struct {
		name                string
		email               string
		password            string
		csrfToken           string
		expectedCode        int
		expectedLocation    string
		expectedFormAction  string
		expectedFieldErrors map[string]string
		expectedFormError   string
	}{
		{
			name:             "Valid submission",
			email:            validEmail,
			password:         validPassword,
			csrfToken:        validCSRFToken,
			expectedCode:     http.StatusSeeOther,
			expectedLocation: "/dashboard",
		},
		{
			name:                "Empty email",
			email:               "",
			password:            validPassword,
			csrfToken:           validCSRFToken,
			expectedCode:        http.StatusUnprocessableEntity,
			expectedFormAction:  "/user/login",
			expectedFieldErrors: map[string]string{"email": "blank"},
		},
		{
			name:                "Invalid email",
			email:               "invalidEmail",
			password:            validPassword,
			csrfToken:           validCSRFToken,
			expectedCode:        http.StatusUnprocessableEntity,
			expectedFormAction:  "/user/login",
			expectedFieldErrors: map[string]string{"email": "Invalid"},
		},
		{
			name:                "Empty password",
			email:               validEmail,
			password:            "",
			csrfToken:           validCSRFToken,
			expectedCode:        http.StatusUnprocessableEntity,
			expectedFormAction:  "/user/login",
			expectedFieldErrors: map[string]string{"password": "blank"},
		},
		{
			name:               "Wrong email",
			email:              "wrong@email.test",
			password:           validEmail,
			csrfToken:          validCSRFToken,
			expectedCode:       http.StatusUnprocessableEntity,
			expectedFormAction: "/user/login",
			expectedFormError:  "email or password",
		},
		{
			name:               "Wrong password",
			email:              validEmail,
			password:           "wrongpassword",
			csrfToken:          validCSRFToken,
			expectedCode:       http.StatusUnprocessableEntity,
			expectedFormAction: "/user/login",
			expectedFormError:  "email or password",
		},
		{
			name:         "Invalid CSRF token",
			email:        validEmail,
			password:     validPassword,
			csrfToken:    "wrong token",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tc.email)
			form.Add("password", tc.password)
			form.Add("csrf_token", tc.csrfToken)

			code, header, doc := ts.postForm(t, "/user/login", form)

			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expectedLocation, header.Get("Location"))

			if tc.expectedFormAction != "" {
				formAction, _ := doc.Find("form").Attr("action")
				assert.Equal(t, tc.expectedFormAction, formAction)
			}

			if tc.expectedFormError != "" {
				formErr := doc.Find("form div[data-form-error]")
				assert.Contains(t, formErr.Text(), tc.expectedFormError)
			}

			if tc.expectedFieldErrors != nil {
				for name, message := range tc.expectedFieldErrors {
					fieldErr := doc.Find(fmt.Sprintf(`form div[data-field-error="%s"]`, name))
					assert.Contains(t, fieldErr.Text(), message)
				}
			}
		})
	}
}

func TestSignupPost(t *testing.T) {
	app := newTestApplication()
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	_, _, doc := ts.get(t, "/user/signup")

	validCSRFToken := extractCSRFToken(doc)

	const (
		validEmail    = "email@example.com"
		validPassword = "pass123"
	)

	tests := []struct {
		name               string
		email              string
		password           string
		passwordConfirm    string
		csrfToken          string
		expectedCode       int
		expectedLocation   string
		expectedFormAction string
		expectedError      map[string]string
	}{
		{
			name:             "Valid submission",
			email:            validEmail,
			password:         validPassword,
			passwordConfirm:  validPassword,
			csrfToken:        validCSRFToken,
			expectedCode:     http.StatusSeeOther,
			expectedLocation: "/user/login",
		},
		{
			name:               "Empty email",
			email:              "",
			password:           validPassword,
			passwordConfirm:    validPassword,
			csrfToken:          validCSRFToken,
			expectedCode:       http.StatusUnprocessableEntity,
			expectedFormAction: "/user/signup",
			expectedError:      map[string]string{"email": "blank"},
		},
		{
			name:               "Invalid email",
			email:              "invalidEmail",
			password:           validPassword,
			passwordConfirm:    validPassword,
			csrfToken:          validCSRFToken,
			expectedCode:       http.StatusUnprocessableEntity,
			expectedFormAction: "/user/signup",
			expectedError:      map[string]string{"email": "Invalid"},
		},
		{
			name:               "Email duplication",
			email:              mocks.User.Email,
			password:           validPassword,
			passwordConfirm:    validPassword,
			csrfToken:          validCSRFToken,
			expectedCode:       http.StatusUnprocessableEntity,
			expectedFormAction: "/user/signup",
			expectedError:      map[string]string{"email": "already in use"},
		},
		{
			name:               "Empty password",
			email:              validEmail,
			password:           "",
			passwordConfirm:    "",
			csrfToken:          validCSRFToken,
			expectedCode:       http.StatusUnprocessableEntity,
			expectedFormAction: "/user/signup",
			expectedError:      map[string]string{"password": "blank"},
		},
		{
			name:               "Short password",
			email:              validEmail,
			password:           "123",
			passwordConfirm:    "123",
			csrfToken:          validCSRFToken,
			expectedCode:       http.StatusUnprocessableEntity,
			expectedFormAction: "/user/signup",
			expectedError:      map[string]string{"password": "6"},
		},
		{
			name:               "Password confirm does not match",
			email:              validEmail,
			password:           validPassword,
			passwordConfirm:    "another password",
			csrfToken:          validCSRFToken,
			expectedCode:       http.StatusUnprocessableEntity,
			expectedFormAction: "/user/signup",
			expectedError:      map[string]string{"password-confirm": "match"},
		},
		{
			name:            "Invalid CSRF token",
			email:           validEmail,
			password:        validPassword,
			passwordConfirm: validPassword,
			csrfToken:       "wrong token",
			expectedCode:    http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("email", tc.email)
			form.Add("password", tc.password)
			form.Add("password-confirm", tc.passwordConfirm)
			form.Add("csrf_token", tc.csrfToken)

			code, header, doc := ts.postForm(t, "/user/signup", form)

			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expectedLocation, header.Get("Location"))

			if tc.expectedFormAction != "" {
				formAction, _ := doc.Find("form").Attr("action")
				assert.Equal(t, tc.expectedFormAction, formAction)
			}

			if tc.expectedError != nil {
				for name, message := range tc.expectedError {
					fieldErr := doc.Find(fmt.Sprintf(`form div[data-field-error="%s"]`, name))
					assert.Contains(t, fieldErr.Text(), message)
				}
			}

		})
	}
}

func TestAddWordPost(t *testing.T) {
	app := newTestApplication()
	h := loggedInStubMiddleware(app.sessionManager, app.routes())
	ts := newTestServer(t, h)
	defer ts.Close()

	_, _, doc := ts.get(t, "/word/add")
	validCSRFToken := extractCSRFToken(doc)

	const (
		validWord         = "Word"
		validTranslations = "translation1, translation2"
	)

	tests := []struct {
		name               string
		word               string
		translations       string
		csrfToken          string
		expectedCode       int
		expectedLocation   string
		expectedFormAction string
		expectedErrors     map[string]string
	}{
		{
			name:             "Valid submission",
			word:             validWord,
			translations:     validTranslations,
			csrfToken:        validCSRFToken,
			expectedCode:     http.StatusSeeOther,
			expectedLocation: "/dashboard",
		},
		{
			name:               "Empty word",
			word:               "",
			translations:       validTranslations,
			expectedCode:       http.StatusUnprocessableEntity,
			csrfToken:          validCSRFToken,
			expectedFormAction: "/word/add",
			expectedErrors:     map[string]string{"word": "blank"},
		},
		{
			name:               "Empty translations",
			word:               validWord,
			translations:       "",
			expectedCode:       http.StatusUnprocessableEntity,
			csrfToken:          validCSRFToken,
			expectedFormAction: "/word/add",
			expectedErrors:     map[string]string{"translations": "empty"},
		},
		{
			name:         "Invalid csrf token",
			word:         validWord,
			translations: validTranslations,
			csrfToken:    "wrong token",
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("word", tc.word)
			form.Add("translations", tc.translations)
			form.Add("csrf_token", tc.csrfToken)

			code, header, doc := ts.postForm(t, "/word/add", form)

			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expectedLocation, header.Get("Location"))

			if tc.expectedFormAction != "" {
				formAction, _ := doc.Find("form").Attr("action")
				assert.Equal(t, tc.expectedFormAction, formAction)
			}

			if tc.expectedErrors != nil {
				for name, message := range tc.expectedErrors {
					fieldErr := doc.Find(fmt.Sprintf(`form div[data-field-error="%s"]`, name))
					assert.Contains(t, fieldErr.Text(), message)
				}
			}

		})
	}
}
