// Copyright © 2019 Arrikto Inc.  All Rights Reserved.

package state

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/quasoft/memstore"
)

func TestSaveLoad(t *testing.T) {
	for _, fn := range []StateFunc{RelativeURL, SchemeAndHost} {
		store := memstore.NewMemStore([]byte("randomstring"))
		r := httptest.NewRequest("GET", "https://example.com/", nil)
		s := fn(r)

		// Check that save works with no errors
		id, err := s.Save(store)
		if err != nil {
			t.Fatalf("Unexpected error while saving: %+v", err)
		}

		// Check that load works with no errors
		loadedState, err := Load(store, id)
		if err != nil {
			t.Fatalf("Unexpected error while loading: %+v", err)
		}

		if !reflect.DeepEqual(loadedState, s) {
			t.Fatalf("Saved state and Loaded state and not equal. Got: '%v' ; Want: '%v'",
				loadedState, s)
		}
	}
}
