package main

import (
	"net/http"

	"gorm.io/gorm"
)

func playerNamesHandler(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
			case "GET":
				pwfs, err := GetPlayers(db, 0, 999)
				if err != nil {
					http.Error(w, "Could not get matches", http.StatusNotFound)
					return
				}
				matchComp := playerNames(pwfs)
				matchComp.Render(GetContext(r), w)
			}
		}
	}
	