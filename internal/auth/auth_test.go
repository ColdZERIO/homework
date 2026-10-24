package auth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTokenTypes(t *testing.T) {
	t.Setenv("ACCESS_KEY", "test-secret")
	t.Setenv("REFRESH_KEY", "test-secret")
	access, err := GenerateAccessToken("user-1", "user")
	if err != nil {
		t.Fatal(err)
	}
	refresh, err := GenerateRefreshToken("user-1", "user")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		token string
		kind  jwtType
	}{{access, AccessToken}, {refresh, RefreshToken}} {
		claims, err := ParseToken(tc.token, tc.kind, []byte("test-secret"))
		if err != nil {
			t.Fatal(err)
		}
		if claims.UserID != "user-1" {
			t.Fatal("wrong user ID")
		}
	}
	if _, err := ParseToken(refresh, AccessToken, []byte("test-secret")); err == nil {
		t.Fatal("refresh accepted as access")
	}
	if _, err := ParseToken(access, AccessToken, []byte("wrong-secret")); err == nil {
		t.Fatal("wrong secret accepted")
	}
}

func TestAccessMiddleware(t *testing.T) {
	t.Setenv("ACCESS_KEY", "test-secret")
	token, err := GenerateAccessToken("user-1", "user")
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		header string
		status int
	}{{"Bearer " + token, 204}, {"", 401}, {"Bearer invalid", 401}} {
		r := httptest.NewRequest(http.MethodPut, "/update", strings.NewReader(`{"id":"user-1"}`))
		r.Header.Set("Authorization", tc.header)
		w := httptest.NewRecorder()
		CheckAccessTokenMiddleware(func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			if string(body) != `{"id":"user-1"}` {
				t.Error("middleware consumed body")
			}
			if !CanAccessUser(r, "user-1") || CanAccessUser(r, "other") {
				t.Error("incorrect user access")
			}
			w.WriteHeader(http.StatusNoContent)
		})(w, r)
		if w.Code != tc.status {
			t.Errorf("status %d, want %d", w.Code, tc.status)
		}
	}
}
