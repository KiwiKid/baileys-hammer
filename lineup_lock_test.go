package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestLockedLineupRequiresUnlockBeforeEdit(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &LineupUser{}, &Formation{}, &FormationPosition{}, &Match{}, &Lineup{}, &LineupPlayer{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Test team"}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	lineup := Lineup{
		TeamID: team.ID,
		Name:   "Locked line-up",
		Status: lineupStatusDraft,
		Locked: true,
	}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create line-up: %v", err)
	}

	router := chi.NewRouter()
	router.HandleFunc("/lineups/{lineupId}", lineupDetailHandler(db))

	form := url.Values{}
	form.Set("name", "Edited line-up")
	form.Set("status", lineupStatusDraft)
	editRequest := httptest.NewRequest(http.MethodPost, "/lineups/"+S(lineup.ID), strings.NewReader(form.Encode()))
	editRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	editRequest.AddCookie(&http.Cookie{Name: adminTokenCookieName, Value: makeAdminToken(team.ID, adminTokenSecret())})
	editResponse := httptest.NewRecorder()
	router.ServeHTTP(editResponse, editRequest)

	if editResponse.Code != http.StatusForbidden {
		t.Fatalf("expected locked edit to be forbidden, got %d", editResponse.Code)
	}
	if err := db.First(&lineup, lineup.ID).Error; err != nil {
		t.Fatalf("reload line-up: %v", err)
	}
	if lineup.Name != "Locked line-up" {
		t.Fatalf("expected locked line-up name unchanged, got %q", lineup.Name)
	}

	unlockForm := url.Values{}
	unlockForm.Set("lineupAction", "unlock")
	unlockRequest := httptest.NewRequest(http.MethodPost, "/lineups/"+S(lineup.ID), strings.NewReader(unlockForm.Encode()))
	unlockRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	unlockRequest.AddCookie(&http.Cookie{Name: adminTokenCookieName, Value: makeAdminToken(team.ID, adminTokenSecret())})
	unlockResponse := httptest.NewRecorder()
	router.ServeHTTP(unlockResponse, unlockRequest)

	if unlockResponse.Code != http.StatusSeeOther {
		t.Fatalf("expected unlock redirect, got %d", unlockResponse.Code)
	}
	if err := db.First(&lineup, lineup.ID).Error; err != nil {
		t.Fatalf("reload unlocked line-up: %v", err)
	}
	if lineup.Locked {
		t.Fatalf("expected line-up to be unlocked")
	}
	if lineup.LockedByName != "" || lineup.LockedByUserID != 0 {
		t.Fatalf("expected lock attribution cleared, got name=%q user=%d", lineup.LockedByName, lineup.LockedByUserID)
	}
}

func TestLockLineupRecordsLocker(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &LineupUser{}, &Formation{}, &FormationPosition{}, &Match{}, &Lineup{}, &LineupPlayer{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Test team"}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	user := LineupUser{TeamID: team.ID, DisplayName: "Sam"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	lineup := Lineup{
		TeamID:          team.ID,
		CreatedByUserID: user.ID,
		Name:            "Draft line-up",
		Status:          lineupStatusDraft,
	}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create line-up: %v", err)
	}

	router := chi.NewRouter()
	router.HandleFunc("/lineups/{lineupId}", lineupDetailHandler(db))

	form := url.Values{}
	form.Set("lineupAction", "lock")
	request := httptest.NewRequest(http.MethodPost, "/lineups/"+S(lineup.ID), strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.AddCookie(&http.Cookie{Name: lineupUserCookieName, Value: makeLineupUserToken(team.ID, user.ID, adminTokenSecret())})
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusSeeOther {
		t.Fatalf("expected lock redirect, got %d", response.Code)
	}
	if err := db.First(&lineup, lineup.ID).Error; err != nil {
		t.Fatalf("reload line-up: %v", err)
	}
	if !lineup.Locked {
		t.Fatalf("expected line-up to be locked")
	}
	if lineup.LockedByUserID != user.ID || lineup.LockedByName != "Sam" {
		t.Fatalf("expected lock attribution to Sam/%d, got %q/%d", user.ID, lineup.LockedByName, lineup.LockedByUserID)
	}
}

func TestUsedLineupCannotBeEditedButCanBeCopied(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&Team{}, &LineupUser{}, &Formation{}, &FormationPosition{}, &Match{}, &Lineup{}, &LineupPlayer{}); err != nil {
		t.Fatalf("migrate db: %v", err)
	}
	team := Team{TeamName: "Test team"}
	if err := db.Create(&team).Error; err != nil {
		t.Fatalf("create team: %v", err)
	}
	lineup := Lineup{
		TeamID:  team.ID,
		Name:    "Used line-up",
		Details: "Original details",
		Status:  lineupStatusUsed,
	}
	if err := db.Create(&lineup).Error; err != nil {
		t.Fatalf("create line-up: %v", err)
	}

	router := chi.NewRouter()
	router.HandleFunc("/lineups/{lineupId}", lineupDetailHandler(db))
	router.HandleFunc("/lineups/{lineupId}/copy", lineupCopyHandler(db))

	adminCookie := &http.Cookie{Name: adminTokenCookieName, Value: makeAdminToken(team.ID, adminTokenSecret())}

	form := url.Values{}
	form.Set("name", "Edited line-up")
	form.Set("status", lineupStatusUsed)
	editRequest := httptest.NewRequest(http.MethodPost, "/lineups/"+S(lineup.ID), strings.NewReader(form.Encode()))
	editRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	editRequest.AddCookie(adminCookie)
	editResponse := httptest.NewRecorder()
	router.ServeHTTP(editResponse, editRequest)

	if editResponse.Code != http.StatusForbidden {
		t.Fatalf("expected used edit to be forbidden, got %d", editResponse.Code)
	}

	copyRequest := httptest.NewRequest(http.MethodPost, "/lineups/"+S(lineup.ID)+"/copy", strings.NewReader(""))
	copyRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	copyRequest.AddCookie(adminCookie)
	copyResponse := httptest.NewRecorder()
	router.ServeHTTP(copyResponse, copyRequest)

	if copyResponse.Code != http.StatusSeeOther {
		t.Fatalf("expected copy redirect, got %d", copyResponse.Code)
	}
	var copied Lineup
	if err := db.Where("id <> ?", lineup.ID).First(&copied).Error; err != nil {
		t.Fatalf("find copied line-up: %v", err)
	}
	if copied.Status != lineupStatusDraft {
		t.Fatalf("expected copied line-up to be draft, got %q", copied.Status)
	}
	if copied.Details != lineup.Details {
		t.Fatalf("expected copied details %q, got %q", lineup.Details, copied.Details)
	}
}
