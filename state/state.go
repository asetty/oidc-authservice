// Copyright © 2019 Arrikto Inc.  All Rights Reserved.

package state

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
)

const oidcLoginSessionCookie = "non_existent_cookie"

var nonceChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

type state struct {
	OrigURL string
}

type StateFunc func(*http.Request) *state

func RelativeURL(r *http.Request) *state {
	return &state{
		OrigURL: r.URL.String(),
	}
}

// scheme could be
// 1. configured by user
// 2. gotten from header e.g. X-Envoy-Forwarded-Proto
func SchemeAndHost(r *http.Request) *state {
	// TODO get scheme from header or config
	return &state{
		OrigURL: "https://" + r.Host + r.URL.String(),
	}
}

// load retrieves a state from the store given its id.
func Load(store sessions.Store, id string) (*state, error) {
	// Make a fake request so that the store will find the cookie
	r := &http.Request{Header: make(http.Header)}
	r.AddCookie(&http.Cookie{Name: oidcLoginSessionCookie, Value: id, MaxAge: 10})

	session, err := store.Get(r, oidcLoginSessionCookie)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if session.IsNew {
		return nil, errors.New("session does not exist")
	}

	return &state{
		OrigURL: session.Values["origURL"].(string),
	}, nil
}

// save persists a state to the store and returns the entry's id.
func (s *state) Save(store sessions.Store) (string, error) {
	session := sessions.NewSession(store, oidcLoginSessionCookie)
	session.ID = createNonce(16)
	session.Options.MaxAge = int(time.Hour)
	session.Values["origURL"] = s.OrigURL

	// The current gorilla/sessions Store interface doesn't allow us
	// to set the session ID.
	// Because of that, we have to retrieve it from the cookie value.
	w := httptest.NewRecorder()
	err := session.Save(&http.Request{}, w)
	if err != nil {
		return "", errors.Wrap(err, "error trying to save session")
	}
	// Cookie is persisted in ResponseWriter, make a request to parse it.
	r := &http.Request{Header: make(http.Header)}
	r.Header.Set("Cookie", w.Header().Get("Set-Cookie"))
	c, err := r.Cookie(oidcLoginSessionCookie)
	if err != nil {
		return "", errors.Wrap(err, "error trying to save session")
	}
	return c.Value, nil
}

func createNonce(length int) string {

	var nonce = make([]rune, length)
	for i := range nonce {
		nonce[i] = nonceChars[rand.Intn(len(nonceChars))]
	}

	return string(nonce)
}
