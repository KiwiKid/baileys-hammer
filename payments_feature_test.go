package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPaymentsRouteUsesPaymentsFeatureFlag(t *testing.T) {
	t.Setenv("GOOGLE_CLIENT_ID", "")
	db := testAuthDB(t)
	season := Season{Title: "2026", StartDate: time.Now().AddDate(0, -1, 0), IsActive: true}
	require.NoError(t, db.Create(&season).Error)
	team := Team{TeamName: "AFC", EnablePlayersModule: true, EnablePaymentsModule: false}
	require.NoError(t, db.Create(&team).Error)
	router := setupRouter(db)

	req := httptest.NewRequest(http.MethodGet, "/season/"+S(season.ID)+"/payments?displayType=button", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusForbidden, rec.Code)
	require.Contains(t, rec.Body.String(), "Payments module is disabled")

	team.EnablePlayersModule = false
	team.EnablePaymentsModule = true
	require.NoError(t, db.Save(&team).Error)
	req = httptest.NewRequest(http.MethodGet, "/season/"+S(season.ID)+"/payments?displayType=button", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "Open Player Payments")
}
