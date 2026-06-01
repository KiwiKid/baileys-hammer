package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestParseUnavailablePlayerIDsAcceptsMultiSelectNames(t *testing.T) {
	form := url.Values{
		"unavailablePlayerIds[]": {"1", "2"},
		"unavailablePlayerIds":   {"3", ""},
	}
	request := httptest.NewRequest(http.MethodPost, "/match/1", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err := request.ParseForm(); err != nil {
		t.Fatalf("parse form: %v", err)
	}

	playerIDs, err := parseUnavailablePlayerIDs(request)
	if err != nil {
		t.Fatalf("parse unavailable player IDs: %v", err)
	}

	want := []uint{3, 1, 2}
	if len(playerIDs) != len(want) {
		t.Fatalf("expected %v, got %v", want, playerIDs)
	}
	for i := range want {
		if playerIDs[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, playerIDs)
		}
	}
}
