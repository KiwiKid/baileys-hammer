package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testAuthDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&Team{}, &AdminUser{}, &AdminUserTeamRole{}, &AdminAccessRequest{}, &Season{}, &Player{}, &Fine{}, &PresetFine{}, &Match{}, &MatchEvent{}, &FineImage{}, &PlayerPayment{}, &LineupUser{}, &Formation{}, &FormationPosition{}, &Lineup{}, &LineupPlayer{}, &MLNote{}, &PlayerMatchUnavailability{}))
	return db
}

func stubGoogleVerifier(t *testing.T, info googleUserInfo) {
	t.Helper()
	previous := googleCredentialVerifier
	googleCredentialVerifier = func(_ string) (*googleUserInfo, error) {
		return &info, nil
	}
	t.Cleanup(func() {
		googleCredentialVerifier = previous
	})
}

func adminSessionRequest(method string, target string, userID uint, teamID uint, body url.Values) *http.Request {
	var reader *strings.Reader
	if body == nil {
		reader = strings.NewReader("")
	} else {
		reader = strings.NewReader(body.Encode())
	}
	req := httptest.NewRequest(method, target, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	req.AddCookie(&http.Cookie{Name: adminSessionCookieName, Value: makeAdminSessionToken(userID, teamID, googleAdminSessionSecret())})
	return req
}

func TestFirstGoogleRegistrationCreatesSuperAdmin(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC", AllowAdminRegistration: true}
	require.NoError(t, db.Create(&team).Error)
	stubGoogleVerifier(t, googleUserInfo{Subject: "sub-1", Email: "first@example.com", DisplayName: "First Admin"})

	body := url.Values{"credential": {"token"}, "teamId": {S(team.ID)}}
	req := httptest.NewRequest(http.MethodPost, "/finemaster/auth", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	finemasterAuthHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	user, err := FindAdminUserByEmail(db, "first@example.com")
	require.NoError(t, err)
	ok, err := AdminUserHasRole(db, user.ID, 0, adminRoleSuperAdmin)
	require.NoError(t, err)
	require.True(t, ok)
	require.NotEmpty(t, rec.Result().Cookies())
}

func TestStartupSuperAdminEnvPromotesExistingUser(t *testing.T) {
	t.Setenv("ADD_SUPER_ADMIN_TO_EMAIL_ON_STARTUP", "owner@example.com")
	db := testAuthDB(t)
	user := AdminUser{Email: "owner@example.com", GoogleSubjectID: "sub-owner"}
	require.NoError(t, db.Create(&user).Error)

	require.NoError(t, applyStartupSuperAdminGrant(db))

	ok, err := AdminUserHasRole(db, user.ID, 0, adminRoleSuperAdmin)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestTeamCreateAssignsTeamAdminToCreator(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	user := AdminUser{Email: "creator@example.com", GoogleSubjectID: "sub-creator"}
	require.NoError(t, db.Create(&user).Error)
	body := url.Values{
		"teamName":                {"New Team"},
		"teamKey":                 {"new"},
		"teamAdminPass":           {"secret"},
		"lineupPlayerCount":       {"11"},
		"enableFinesModule":       {"on"},
		"enableLineupsModule":     {"on"},
		"enableMatchesModule":     {"on"},
		"enablePlayersModule":     {"on"},
		"enablePaymentsModule":    {"on"},
		"enableCourtModule":       {"on"},
		"enableLeaderboardModule": {"on"},
	}
	req := adminSessionRequest(http.MethodPost, "/teams", user.ID, 0, body)
	rec := httptest.NewRecorder()

	teamHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var team Team
	require.NoError(t, db.Where("team_name = ?", "New Team").First(&team).Error)
	ok, err := AdminUserHasRole(db, user.ID, team.ID, adminRoleTeamAdmin)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestTeamAdminCannotActivateAnotherTeam(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	teamOne := Team{TeamName: "One"}
	teamTwo := Team{TeamName: "Two"}
	user := AdminUser{Email: "team-admin@example.com", GoogleSubjectID: "sub-team"}
	require.NoError(t, db.Create(&teamOne).Error)
	require.NoError(t, db.Create(&teamTwo).Error)
	require.NoError(t, db.Create(&user).Error)
	require.NoError(t, GrantAdminRole(db, user.ID, teamOne.ID, adminRoleTeamAdmin))

	req := adminSessionRequest(http.MethodPost, F("/teams/activate?teamId=%d", teamTwo.ID), user.ID, teamOne.ID, nil)
	rec := httptest.NewRecorder()

	teamActivateHandler(db)(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
}

func TestAllowAdminRegistrationControlsSignupButton(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "Closed", AllowAdminRegistration: false}
	require.NoError(t, db.Create(&team).Error)

	req := httptest.NewRequest(http.MethodGet, "/finemaster/auth", nil)
	rec := httptest.NewRecorder()
	finemasterAuthHandler(db)(rec, req)
	require.NotContains(t, rec.Body.String(), "Sign up")

	team.AllowAdminRegistration = true
	require.NoError(t, db.Save(&team).Error)
	rec = httptest.NewRecorder()
	finemasterAuthHandler(db)(rec, req)
	require.Contains(t, rec.Body.String(), "Sign up")
}

func TestPasswordAdminRejectedWhenGoogleEnabled(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC"}
	require.NoError(t, db.Create(&team).Error)
	body := url.Values{"password": {"pass"}, "teamId": {S(team.ID)}}
	req := httptest.NewRequest(http.MethodPost, "/admin", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	adminHandler(db)(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestPasswordAdminWorksWhenGoogleUnset(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "")
	t.Setenv("PASS", "pass")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC"}
	require.NoError(t, db.Create(&team).Error)
	body := url.Values{"password": {"pass"}, "teamId": {S(team.ID)}}
	req := httptest.NewRequest(http.MethodPost, "/admin", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	adminHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "/finemaster", rec.Header().Get("HX-Redirect"))
}

func TestSuperAdminCanAssignTeamRole(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC"}
	superAdmin := AdminUser{Email: "owner@example.com", GoogleSubjectID: "sub-owner"}
	target := AdminUser{Email: "target@example.com", GoogleSubjectID: "sub-target"}
	require.NoError(t, db.Create(&team).Error)
	require.NoError(t, db.Create(&superAdmin).Error)
	require.NoError(t, db.Create(&target).Error)
	require.NoError(t, GrantAdminRole(db, superAdmin.ID, 0, adminRoleSuperAdmin))
	body := url.Values{"adminUserId": {S(target.ID)}, "role": {adminRoleTeamAdmin}, "teamId": {S(team.ID)}}
	req := adminSessionRequest(http.MethodPost, "/admin/users/roles", superAdmin.ID, team.ID, body)
	rec := httptest.NewRecorder()

	adminUserRoleHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	ok, err := AdminUserHasRole(db, target.ID, team.ID, adminRoleTeamAdmin)
	require.NoError(t, err)
	require.True(t, ok)
	require.Contains(t, rec.Body.String(), "target@example.com")
}

func TestTeamAdminCannotAssignRolesOutsideOwnTeam(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC"}
	otherTeam := Team{TeamName: "Other"}
	teamAdmin := AdminUser{Email: "team@example.com", GoogleSubjectID: "sub-team"}
	target := AdminUser{Email: "target@example.com", GoogleSubjectID: "sub-target"}
	require.NoError(t, db.Create(&team).Error)
	require.NoError(t, db.Create(&otherTeam).Error)
	require.NoError(t, db.Create(&teamAdmin).Error)
	require.NoError(t, db.Create(&target).Error)
	require.NoError(t, GrantAdminRole(db, teamAdmin.ID, team.ID, adminRoleTeamAdmin))
	body := url.Values{"adminUserId": {S(target.ID)}, "role": {adminRoleTeamAdmin}, "teamId": {S(otherTeam.ID)}}
	req := adminSessionRequest(http.MethodPost, "/admin/users/roles", teamAdmin.ID, team.ID, body)
	rec := httptest.NewRecorder()

	adminUserRoleHandler(db)(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
	ok, err := AdminUserHasRole(db, target.ID, otherTeam.ID, adminRoleTeamAdmin)
	require.NoError(t, err)
	require.False(t, ok)
}

func TestGoogleUserCanRequestTeamAdminAccess(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC"}
	existing := AdminUser{Email: "owner@example.com", GoogleSubjectID: "sub-owner"}
	require.NoError(t, db.Create(&team).Error)
	require.NoError(t, db.Create(&existing).Error)
	require.NoError(t, GrantAdminRole(db, existing.ID, 0, adminRoleSuperAdmin))
	stubGoogleVerifier(t, googleUserInfo{Subject: "sub-request", Email: "request@example.com", DisplayName: "Requester"})
	body := url.Values{"credential": {"token"}, "mode": {"request"}, "teamId": {S(team.ID)}}
	req := httptest.NewRequest(http.MethodPost, "/finemaster/auth", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	finemasterAuthHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	user, err := FindAdminUserByEmail(db, "request@example.com")
	require.NoError(t, err)
	ok, err := AdminUserHasRole(db, user.ID, team.ID, adminRoleTeamAdmin)
	require.NoError(t, err)
	require.False(t, ok)
	var requests []AdminAccessRequest
	require.NoError(t, db.Where("admin_user_id = ? AND team_id = ? AND status = ?", user.ID, team.ID, "pending").Find(&requests).Error)
	require.Len(t, requests, 1)
}

func TestGoogleUserCanRequestLineupAccess(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC"}
	require.NoError(t, db.Create(&team).Error)
	stubGoogleVerifier(t, googleUserInfo{Subject: "sub-lineup-request", Email: "lineup@example.com", DisplayName: "Lineup User"})
	body := url.Values{"credential": {"token"}, "teamId": {S(team.ID)}}
	req := httptest.NewRequest(http.MethodPost, "/lineups/login", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	lineupLoginHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	user, err := FindAdminUserByEmail(db, "lineup@example.com")
	require.NoError(t, err)
	ok, err := AdminUserHasRole(db, user.ID, team.ID, adminRoleLineupAccess)
	require.NoError(t, err)
	require.False(t, ok)
	var requests []AdminAccessRequest
	require.NoError(t, db.Where("admin_user_id = ? AND team_id = ? AND role = ? AND status = ?", user.ID, team.ID, adminRoleLineupAccess, "pending").Find(&requests).Error)
	require.Len(t, requests, 1)
}

func TestApprovedGoogleLineupUserGetsLineupCookie(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC"}
	user := AdminUser{Email: "lineup@example.com", GoogleSubjectID: "sub-lineup", DisplayName: "Lineup User"}
	require.NoError(t, db.Create(&team).Error)
	require.NoError(t, db.Create(&user).Error)
	require.NoError(t, GrantAdminRole(db, user.ID, team.ID, adminRoleLineupAccess))
	stubGoogleVerifier(t, googleUserInfo{Subject: "sub-lineup", Email: "lineup@example.com", DisplayName: "Lineup User"})
	body := url.Values{"credential": {"token"}, "teamId": {S(team.ID)}}
	req := httptest.NewRequest(http.MethodPost, "/lineups/login", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()

	lineupLoginHandler(db)(rec, req)

	require.Equal(t, http.StatusSeeOther, rec.Code)
	require.Equal(t, "/lineups", rec.Header().Get("Location"))
	var foundLineupCookie bool
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name != lineupUserCookieName {
			continue
		}
		teamID, userID, ok := parseLineupUserToken(cookie.Value, adminTokenSecret())
		require.True(t, ok)
		require.Equal(t, team.ID, teamID)
		require.NotZero(t, userID)
		foundLineupCookie = true
	}
	require.True(t, foundLineupCookie)
}

func TestHomeTeamKeyActivatesAccessibleTeam(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC", TeamKey: "afc"}
	user := AdminUser{Email: "team@example.com", GoogleSubjectID: "sub-team"}
	require.NoError(t, db.Create(&team).Error)
	require.NoError(t, db.Create(&user).Error)
	require.NoError(t, GrantAdminRole(db, user.ID, team.ID, adminRoleTeamAdmin))
	req := adminSessionRequest(http.MethodGet, "/?teamKey=afc", user.ID, 0, nil)
	rec := httptest.NewRecorder()

	homeHandler(db)(rec, req)

	require.Equal(t, http.StatusSeeOther, rec.Code)
	require.Equal(t, "/", rec.Header().Get("Location"))
	var foundAdminSession bool
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == adminSessionCookieName {
			_, teamID, ok := parseAdminSessionToken(cookie.Value, googleAdminSessionSecret())
			require.True(t, ok)
			require.Equal(t, team.ID, teamID)
			foundAdminSession = true
		}
	}
	require.True(t, foundAdminSession)
}

func TestHomeTeamKeyShowsRequestAccessForLoggedInUserWithoutAccess(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC", TeamKey: "afc"}
	user := AdminUser{Email: "viewer@example.com", GoogleSubjectID: "sub-viewer"}
	require.NoError(t, db.Create(&team).Error)
	require.NoError(t, db.Create(&user).Error)
	req := adminSessionRequest(http.MethodGet, "/?teamKey=afc", user.ID, 0, nil)
	rec := httptest.NewRecorder()

	homeHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Request access to AFC")
	require.Contains(t, rec.Body.String(), `action="/admin/team-access-requests"`)
}

func TestLoggedInUserCanRequestHomeTeamAccess(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC", TeamKey: "afc"}
	user := AdminUser{Email: "viewer@example.com", GoogleSubjectID: "sub-viewer"}
	require.NoError(t, db.Create(&team).Error)
	require.NoError(t, db.Create(&user).Error)
	body := url.Values{"teamId": {S(team.ID)}}
	req := adminSessionRequest(http.MethodPost, "/admin/team-access-requests", user.ID, 0, body)
	rec := httptest.NewRecorder()

	adminTeamAccessRequestHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	var requests []AdminAccessRequest
	require.NoError(t, db.Where("admin_user_id = ? AND team_id = ? AND status = ?", user.ID, team.ID, "pending").Find(&requests).Error)
	require.Len(t, requests, 1)
}

func TestAdminTeamSwitchActivatesAccessibleTeam(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	teamOne := Team{TeamName: "One", TeamKey: "one"}
	teamTwo := Team{TeamName: "Two", TeamKey: "two"}
	user := AdminUser{Email: "team@example.com", GoogleSubjectID: "sub-team"}
	require.NoError(t, db.Create(&teamOne).Error)
	require.NoError(t, db.Create(&teamTwo).Error)
	require.NoError(t, db.Create(&user).Error)
	require.NoError(t, GrantAdminRole(db, user.ID, teamOne.ID, adminRoleTeamAdmin))
	require.NoError(t, GrantAdminRole(db, user.ID, teamTwo.ID, adminRoleTeamAdmin))
	body := url.Values{"teamId": {S(teamTwo.ID)}}
	req := adminSessionRequest(http.MethodPost, "/admin/team-switch", user.ID, teamOne.ID, body)
	rec := httptest.NewRecorder()

	adminTeamSwitchHandler(db)(rec, req)

	require.Equal(t, http.StatusSeeOther, rec.Code)
	require.Equal(t, "/", rec.Header().Get("Location"))
	var foundAdminSession bool
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == adminSessionCookieName {
			_, teamID, ok := parseAdminSessionToken(cookie.Value, googleAdminSessionSecret())
			require.True(t, ok)
			require.Equal(t, teamTwo.ID, teamID)
			foundAdminSession = true
		}
	}
	require.True(t, foundAdminSession)
}

func TestTeamAdminCanApproveTeamAccessRequest(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC"}
	teamAdmin := AdminUser{Email: "team@example.com", GoogleSubjectID: "sub-team"}
	target := AdminUser{Email: "target@example.com", GoogleSubjectID: "sub-target"}
	require.NoError(t, db.Create(&team).Error)
	require.NoError(t, db.Create(&teamAdmin).Error)
	require.NoError(t, db.Create(&target).Error)
	require.NoError(t, GrantAdminRole(db, teamAdmin.ID, team.ID, adminRoleTeamAdmin))
	require.NoError(t, RequestAdminAccess(db, target.ID, team.ID, adminRoleTeamAdmin))
	var accessRequest AdminAccessRequest
	require.NoError(t, db.Where("admin_user_id = ? AND team_id = ?", target.ID, team.ID).First(&accessRequest).Error)
	body := url.Values{"requestId": {S(accessRequest.ID)}, "action": {"approve"}}
	req := adminSessionRequest(http.MethodPost, "/admin/access-requests", teamAdmin.ID, team.ID, body)
	rec := httptest.NewRecorder()

	adminAccessRequestHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	ok, err := AdminUserHasRole(db, target.ID, team.ID, adminRoleTeamAdmin)
	require.NoError(t, err)
	require.True(t, ok)
	require.NoError(t, db.First(&accessRequest, accessRequest.ID).Error)
	require.Equal(t, "approved", accessRequest.Status)
}

func TestTeamAdminCanApproveLineupAccessRequest(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC"}
	teamAdmin := AdminUser{Email: "team@example.com", GoogleSubjectID: "sub-team"}
	target := AdminUser{Email: "target@example.com", GoogleSubjectID: "sub-target"}
	require.NoError(t, db.Create(&team).Error)
	require.NoError(t, db.Create(&teamAdmin).Error)
	require.NoError(t, db.Create(&target).Error)
	require.NoError(t, GrantAdminRole(db, teamAdmin.ID, team.ID, adminRoleTeamAdmin))
	require.NoError(t, RequestAdminAccess(db, target.ID, team.ID, adminRoleLineupAccess))
	var accessRequest AdminAccessRequest
	require.NoError(t, db.Where("admin_user_id = ? AND team_id = ?", target.ID, team.ID).First(&accessRequest).Error)
	body := url.Values{"requestId": {S(accessRequest.ID)}, "action": {"approve"}}
	req := adminSessionRequest(http.MethodPost, "/admin/access-requests", teamAdmin.ID, team.ID, body)
	rec := httptest.NewRecorder()

	adminAccessRequestHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	ok, err := AdminUserHasRole(db, target.ID, team.ID, adminRoleLineupAccess)
	require.NoError(t, err)
	require.True(t, ok)
}

func TestTeamAdminCanGrantAndRevokeTeamAdminForOwnTeam(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	team := Team{TeamName: "AFC"}
	teamAdmin := AdminUser{Email: "team@example.com", GoogleSubjectID: "sub-team"}
	target := AdminUser{Email: "target@example.com", GoogleSubjectID: "sub-target"}
	require.NoError(t, db.Create(&team).Error)
	require.NoError(t, db.Create(&teamAdmin).Error)
	require.NoError(t, db.Create(&target).Error)
	require.NoError(t, GrantAdminRole(db, teamAdmin.ID, team.ID, adminRoleTeamAdmin))
	body := url.Values{"adminUserId": {S(target.ID)}, "role": {adminRoleTeamAdmin}, "teamId": {S(team.ID)}, "action": {"grant"}}
	req := adminSessionRequest(http.MethodPost, "/admin/users/roles", teamAdmin.ID, team.ID, body)
	rec := httptest.NewRecorder()

	adminUserRoleHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	ok, err := AdminUserHasRole(db, target.ID, team.ID, adminRoleTeamAdmin)
	require.NoError(t, err)
	require.True(t, ok)

	body.Set("action", "revoke")
	req = adminSessionRequest(http.MethodPost, "/admin/users/roles", teamAdmin.ID, team.ID, body)
	rec = httptest.NewRecorder()
	adminUserRoleHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	ok, err = AdminUserHasRole(db, target.ID, team.ID, adminRoleTeamAdmin)
	require.NoError(t, err)
	require.False(t, ok)
}

func TestSuperAdminCanGrantAndRevokeSuperAdmin(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "client-id")
	t.Setenv("GOOGLE_CLIENT_SECRET", "client-secret")
	db := testAuthDB(t)
	superAdmin := AdminUser{Email: "owner@example.com", GoogleSubjectID: "sub-owner"}
	target := AdminUser{Email: "target@example.com", GoogleSubjectID: "sub-target"}
	require.NoError(t, db.Create(&superAdmin).Error)
	require.NoError(t, db.Create(&target).Error)
	require.NoError(t, GrantAdminRole(db, superAdmin.ID, 0, adminRoleSuperAdmin))
	body := url.Values{"adminUserId": {S(target.ID)}, "role": {adminRoleSuperAdmin}, "action": {"grant"}}
	req := adminSessionRequest(http.MethodPost, "/admin/users/roles", superAdmin.ID, 0, body)
	rec := httptest.NewRecorder()

	adminUserRoleHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	ok, err := AdminUserHasRole(db, target.ID, 0, adminRoleSuperAdmin)
	require.NoError(t, err)
	require.True(t, ok)

	body.Set("action", "revoke")
	req = adminSessionRequest(http.MethodPost, "/admin/users/roles", superAdmin.ID, 0, body)
	rec = httptest.NewRecorder()
	adminUserRoleHandler(db)(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	ok, err = AdminUserHasRole(db, target.ID, 0, adminRoleSuperAdmin)
	require.NoError(t, err)
	require.False(t, ok)
}
