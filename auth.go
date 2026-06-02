package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

const (
	adminSessionCookieName = "admin-session"
	adminRoleSuperAdmin    = "super-admin"
	adminRoleTeamAdmin     = "team-admin"
	adminRoleLineupAccess  = "lineup-access"
)

type googleUserInfo struct {
	Subject     string
	Email       string
	DisplayName string
}

var googleCredentialVerifier = verifyGoogleCredentialWithTokenInfo

func googleAuthEnabled() bool {
	return strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")) != ""
}

func googleAdminSessionSecret() string {
	secret := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET"))
	if secret != "" {
		return secret
	}
	return strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
}

func verifyGoogleCredentialWithTokenInfo(credential string) (*googleUserInfo, error) {
	clientID := strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID"))
	if clientID == "" {
		return nil, errors.New("google auth is not configured")
	}
	if strings.TrimSpace(credential) == "" {
		return nil, errors.New("missing google credential")
	}

	endpoint := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(credential)
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google token verification failed: %s", resp.Status)
	}

	var payload struct {
		Audience      string `json:"aud"`
		Subject       string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		Name          string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.Audience != clientID {
		return nil, errors.New("google credential audience mismatch")
	}
	if payload.Subject == "" || payload.Email == "" {
		return nil, errors.New("google credential is missing subject or email")
	}
	if payload.EmailVerified != "true" {
		return nil, errors.New("google email is not verified")
	}
	return &googleUserInfo{
		Subject:     payload.Subject,
		Email:       strings.ToLower(strings.TrimSpace(payload.Email)),
		DisplayName: payload.Name,
	}, nil
}

func makeAdminSessionToken(userID uint, teamID uint, secret string) string {
	body := fmt.Sprintf("%d:%d", userID, teamID)
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(body))
	return body + ":" + hex.EncodeToString(mac.Sum(nil))
}

func parseAdminSessionToken(token string, secret string) (uint, uint, bool) {
	parts := strings.Split(token, ":")
	if len(parts) != 3 {
		return 0, 0, false
	}
	body := parts[0] + ":" + parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(body))
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(parts[2]), []byte(expected)) {
		return 0, 0, false
	}
	userID64, userErr := strconv.ParseUint(parts[0], 10, 64)
	teamID64, teamErr := strconv.ParseUint(parts[1], 10, 64)
	if userErr != nil || teamErr != nil || userID64 == 0 {
		return 0, 0, false
	}
	return uint(userID64), uint(teamID64), true
}

func setAdminSessionCookie(w http.ResponseWriter, userID uint, teamID uint) {
	http.SetCookie(w, &http.Cookie{
		Name:     adminSessionCookieName,
		Value:    makeAdminSessionToken(userID, teamID, googleAdminSessionSecret()),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   60 * 60 * 24 * 30,
	})
}

func currentAdminUser(r *http.Request, db *gorm.DB) (*AdminUser, uint, bool) {
	c, err := r.Cookie(adminSessionCookieName)
	if err != nil {
		return nil, 0, false
	}
	userID, teamID, ok := parseAdminSessionToken(c.Value, googleAdminSessionSecret())
	if !ok {
		return nil, 0, false
	}
	user, err := GetAdminUser(db, userID)
	if err != nil {
		return nil, 0, false
	}
	return user, teamID, true
}

func isSuperAdmin(db *gorm.DB, userID uint) bool {
	ok, err := AdminUserHasRole(db, userID, 0, adminRoleSuperAdmin)
	return err == nil && ok
}

func canAdminAccessTeam(db *gorm.DB, userID uint, teamID uint) bool {
	if userID == 0 {
		return false
	}
	if isSuperAdmin(db, userID) {
		return true
	}
	if teamID == 0 {
		return false
	}
	ok, err := AdminUserHasRole(db, userID, teamID, adminRoleTeamAdmin)
	return err == nil && ok
}

func canGoogleUserAccessLineups(db *gorm.DB, userID uint, teamID uint) bool {
	if userID == 0 || teamID == 0 {
		return false
	}
	if canAdminAccessTeam(db, userID, teamID) {
		return true
	}
	ok, err := AdminUserHasRole(db, userID, teamID, adminRoleLineupAccess)
	return err == nil && ok
}

func selectedAdminTeamID(r *http.Request, db *gorm.DB) uint {
	if user, teamID, ok := currentAdminUser(r, db); ok {
		if teamID > 0 && canAdminAccessTeam(db, user.ID, teamID) {
			return teamID
		}
		if team, err := FirstTeamForAdminUser(db, user.ID); err == nil && team != nil {
			return team.ID
		}
	}
	return 0
}

func requireGoogleAdmin(db *gorm.DB, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !googleAuthEnabled() {
			next(w, r)
			return
		}
		user, teamID, ok := currentAdminUser(r, db)
		if !ok {
			w.Header().Set("HX-Redirect", "/finemaster/auth")
			http.Redirect(w, r, "/finemaster/auth", http.StatusSeeOther)
			return
		}
		if teamID > 0 && !canAdminAccessTeam(db, user.ID, teamID) {
			http.Error(w, "Admin user cannot access this team", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func requireGoogleAdminForTeam(db *gorm.DB, r *http.Request, teamID uint) (*AdminUser, bool) {
	if !googleAuthEnabled() {
		return nil, true
	}
	user, _, ok := currentAdminUser(r, db)
	if !ok {
		return nil, false
	}
	return user, canAdminAccessTeam(db, user.ID, teamID)
}

func applyStartupSuperAdminGrant(db *gorm.DB) error {
	email := strings.ToLower(strings.TrimSpace(os.Getenv("ADD_SUPER_ADMIN_TO_EMAIL_ON_STARTUP")))
	if email == "" {
		return nil
	}
	var users []AdminUser
	if err := db.Where("LOWER(email) = ?", email).Find(&users).Error; err != nil {
		return err
	}
	for _, user := range users {
		if err := GrantAdminRole(db, user.ID, 0, adminRoleSuperAdmin); err != nil {
			return err
		}
	}
	return nil
}

func FilterTeamsForAdminUser(db *gorm.DB, userID uint, teams []Team) []Team {
	if userID == 0 || isSuperAdmin(db, userID) {
		return teams
	}
	filtered := make([]Team, 0, len(teams))
	for _, team := range teams {
		if canAdminAccessTeam(db, userID, team.ID) {
			filtered = append(filtered, team)
		}
	}
	return filtered
}
