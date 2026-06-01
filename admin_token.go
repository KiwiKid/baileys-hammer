package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strconv"
	"strings"
)

const adminTokenCookieName = "admin-token"

func adminTokenSecret() string {
	// Keep this consistent with the existing DEFAULT_PASS fallback in adminHandler.
	pass := os.Getenv("PASS")
	if pass == "" {
		pass = "pass"
	}
	return pass
}

func makeAdminToken(teamID uint, secret string) string {
	idStr := strconv.FormatUint(uint64(teamID), 10)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(idStr))
	sig := hex.EncodeToString(mac.Sum(nil))
	return idStr + ":" + sig
}

func parseAdminToken(token string, secret string) (uint, bool) {
	parts := strings.Split(token, ":")
	if len(parts) != 2 {
		return 0, false
	}
	idStr := parts[0]
	sig := parts[1]

	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(idStr))
	expectedSig := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(sig), []byte(expectedSig)) {
		return 0, false
	}

	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(id64), true
}

