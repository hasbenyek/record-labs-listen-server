package server

import (
	"strings"
	"testing"
	"time"
)

func TestClientAuthToken(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	auth := &clientAuth{key: []byte("secret-key")}
	good := clientAuthToken(auth.key, now)
	if !auth.valid(good, now) {
		t.Fatalf("fresh token %q rejected", good)
	}
	if !auth.valid(good, now.Add(ClientAuthMaxSkew-time.Minute)) || !auth.valid(good, now.Add(-ClientAuthMaxSkew+time.Minute)) {
		t.Fatal("token inside the allowed clock skew rejected")
	}
	if auth.valid(good, now.Add(ClientAuthMaxSkew+time.Minute)) {
		t.Fatal("old token accepted")
	}
	for _, bad := range []string{
		"",
		"v1",
		clientAuthToken([]byte("other-key"), now),
		strings.Replace(good, "v1.", "v2.", 1),
		good[:len(good)-1] + "0",
		"v1.abc." + strings.Split(good, ".")[2],
	} {
		if bad != good && auth.valid(bad, now) {
			t.Errorf("token %q accepted", bad)
		}
	}
}

func TestClientAuthTier(t *testing.T) {
	now := time.Unix(1_800_000_000, 0)
	key := []byte("secret-key")
	good := clientAuthToken(key, now)

	var none *clientAuth
	if tier, rejected := none.tier("", uaAdvert, now); tier != uaAdvert || rejected {
		t.Fatalf("no key configured: got %v %v", tier, rejected)
	}

	optional := &clientAuth{key: key}
	if tier, rejected := optional.tier(good, uaAdvert, now); tier != uaAllow || rejected {
		t.Fatalf("optional, valid token: got %v %v", tier, rejected)
	}
	if tier, rejected := optional.tier("", uaAllow, now); tier != uaAllow || rejected {
		t.Fatalf("optional, no token keeps the User-Agent tier: got %v %v", tier, rejected)
	}

	required := &clientAuth{key: key, required: true}
	if tier, rejected := required.tier(good, uaAdvert, now); tier != uaAllow || rejected {
		t.Fatalf("required, valid token: got %v %v", tier, rejected)
	}
	if tier, rejected := required.tier("", uaAllow, now); tier != uaBlock || !rejected {
		t.Fatalf("required, no token: got %v %v", tier, rejected)
	}
	if tier, rejected := required.tier(clientAuthToken([]byte("fork"), now), uaAllow, now); tier != uaBlock || !rejected {
		t.Fatalf("required, wrong key: got %v %v", tier, rejected)
	}
}

func TestLoadClientAuth(t *testing.T) {
	t.Setenv("CLIENT_KEY", "")
	if loadClientAuth() != nil {
		t.Fatal("empty CLIENT_KEY should disable the check")
	}
	t.Setenv("CLIENT_KEY", " k ")
	t.Setenv("REQUIRE_CLIENT_KEY", "true")
	auth := loadClientAuth()
	if auth == nil || string(auth.key) != "k" || !auth.required {
		t.Fatalf("got %+v", auth)
	}
	t.Setenv("REQUIRE_CLIENT_KEY", "")
	if auth := loadClientAuth(); auth == nil || auth.required {
		t.Fatalf("REQUIRE_CLIENT_KEY unset should not require the key: %+v", auth)
	}
}
