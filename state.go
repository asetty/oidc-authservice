// Copyright © 2019 Arrikto Inc.  All Rights Reserved.
//
// Utils related to handling the OIDC state parameter
// for CSRF.

package main

const (
	oauthStateCookie = "oidc_state_csrf"
)

type State struct {
	// FirstVisitedURL is the URL that the user visited when we redirected them
	// to login.
	FirstVisitedURL string
}

func newState(firstVisitedURL string) *State {
	return &State{
		FirstVisitedURL: firstVisitedURL,
	}
}
