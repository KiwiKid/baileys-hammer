package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func stubGoogleConnectivity(t *testing.T, ok bool, message string, err error) {
	t.Helper()
	previous := googleConnectivityChecker
	googleConnectivityChecker = func() (bool, string, error) {
		return ok, message, err
	}
	t.Cleanup(func() {
		googleConnectivityChecker = previous
	})
}

func TestGoogleAuthStatusShowsLegacyAdminStatus(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "")
	t.Setenv("GOOGLE_CLIENT_SECRET", "")
	stubGoogleConnectivity(t, true, "Google OAuth discovery is reachable", nil)
	db := testAuthDB(t)

	req := httptest.NewRequest(http.MethodGet, "/admin/google-auth-status", nil)
	rec := httptest.NewRecorder()
	googleAuthStatusHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Google admin auth")
	require.Contains(t, rec.Body.String(), "Google OAuth discovery is reachable")
}

func TestGoogleAuthStatusRequiresSuperAdminWhenGoogleAuthEnabled(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "session-secret")
	stubGoogleConnectivity(t, false, "Could not reach Google OAuth discovery", errors.New("dial timeout"))
	db := testAuthDB(t)
	team := Team{TeamName: "AFC"}
	teamAdmin := AdminUser{Email: "team@example.com", GoogleSubjectID: "sub-team"}
	superAdmin := AdminUser{Email: "owner@example.com", GoogleSubjectID: "sub-owner"}
	require.NoError(t, db.Create(&team).Error)
	require.NoError(t, db.Create(&teamAdmin).Error)
	require.NoError(t, db.Create(&superAdmin).Error)
	require.NoError(t, GrantAdminRole(db, teamAdmin.ID, team.ID, adminRoleTeamAdmin))
	require.NoError(t, GrantAdminRole(db, superAdmin.ID, 0, adminRoleSuperAdmin))

	req := adminSessionRequest(http.MethodGet, "/admin/google-auth-status", teamAdmin.ID, team.ID, nil)
	rec := httptest.NewRecorder()
	googleAuthStatusHandler(db)(rec, req)
	require.Equal(t, http.StatusForbidden, rec.Code)

	req = adminSessionRequest(http.MethodGet, "/admin/google-auth-status", superAdmin.ID, team.ID, nil)
	rec = httptest.NewRecorder()
	googleAuthStatusHandler(db)(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Client ID")
	require.Contains(t, rec.Body.String(), "dial timeout")
}
