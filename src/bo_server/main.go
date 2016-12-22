package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"boardInfo"
	"gameState"
)

var activeGame *gameState.Game

func init() {
	activeGame = gameState.NewGame([]string{
		"Player 1",
		"Player 2",
		"Player 3",
		"Player 4",
		"Player 5",
		"Player 6",
	})
}

func getboardInfo(w http.ResponseWriter, r *http.Request) {
	if buf, err := boardInfo.JsonMap(); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf)
	}
}

func getPlayers(w http.ResponseWriter, r *http.Request) {
	if buf, err := json.Marshal(activeGame.Players); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf)
	}
}

func getCompanies(w http.ResponseWriter, r *http.Request) {
	if buf, err := json.Marshal(activeGame.Companies); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf)
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/board_info", getboardInfo)
	r.HandleFunc("/companies", getCompanies)
	r.HandleFunc("/players", getPlayers)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(AssetFS{})))

	svr := http.Server{
		Handler: r,
		Addr:    ":8000",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	panic(svr.ListenAndServe())
}
