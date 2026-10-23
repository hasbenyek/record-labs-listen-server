package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
	"time"
)

// ClientAuthHeader carries the official apps' proof of origin:
// "v1.<unix seconds>.<hex HMAC-SHA256(key, "recordlabs-listen.v1.<unix seconds>")>".
// The key is a server secret (CLIENT_KEY) baked into official builds at build
// time; it is never part of the public source code.
const ClientAuthHeader = "X-RecordLabs-Auth"

// ClientAuthMaxSkew is how far a token's time may be from the server clock.
const ClientAuthMaxSkew = 12 * time.Hour

// UnofficialClientMessage is what a client without a valid key is told when
// the key is required.
const UnofficialClientMessage = "This Listen Together server only works with the official Record Labs app."

type clientAuth struct {
	key      []byte
	required bool
}

// loadClientAuth reads CLIENT_KEY and REQUIRE_CLIENT_KEY. Without a key every
// client is handled by the User-Agent policy alone, as before.
func loadClientAuth() *clientAuth {
	key := strings.TrimSpace(os.Getenv("CLIENT_KEY"))
	if key == "" {
		return nil
	}
	required, _ := strconv.ParseBool(strings.TrimSpace(os.Getenv("REQUIRE_CLIENT_KEY")))
	return &clientAuth{key: []byte(key), required: required}
}

func clientAuthToken(key []byte, now time.Time) string {
	ts := strconv.FormatInt(now.Unix(), 10)
	return "v1." + ts + "." + clientAuthSignature(key, ts)
}

func clientAuthSignature(key []byte, ts string) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte("recordlabs-listen.v1." + ts))
	return hex.EncodeToString(mac.Sum(nil))
}

// valid reports whether token was made with this server's key, recently.
func (a *clientAuth) valid(token string, now time.Time) bool {
	if a == nil || len(a.key) == 0 {
		return false
	}
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 3 || parts[0] != "v1" {
		return false
	}
	seconds, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return false
	}
	skew := now.Sub(time.Unix(seconds, 0))
	if skew < -ClientAuthMaxSkew || skew > ClientAuthMaxSkew {
		return false
	}
	return hmac.Equal([]byte(parts[2]), []byte(clientAuthSignature(a.key, parts[1])))
}

// tier decides a connection's tier from its token and its User-Agent tier.
// rejected is true when the connection is refused for lacking a valid key.
func (a *clientAuth) tier(token string, uaTierFromPolicy uaTier, now time.Time) (tier uaTier, rejected bool) {
	if a == nil {
		return uaTierFromPolicy, false
	}
	if a.valid(token, now) {
		return uaAllow, false
	}
	if a.required {
		return uaBlock, true
	}
	return uaTierFromPolicy, false
}
