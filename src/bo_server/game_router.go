package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"

	"gameState"
)

var activeGames = map[string]*gameRouter{}
var mapLock sync.RWMutex

type gameRouter struct {
	http.Handler
	game *gameState.Game
}

func (r gameRouter) getPlayers(writer http.ResponseWriter, request *http.Request) {
	writeJson(&jsonResponse{Result: r.game.Players}, writer)
}
func (r gameRouter) getCompanies(writer http.ResponseWriter, request *http.Request) {
	writeJson(&jsonResponse{Result: r.game.Companies}, writer)
}
func (r gameRouter) getGameState(writer http.ResponseWriter, request *http.Request) {
	writeJson(&jsonResponse{Result: r.game.GlobalState}, writer)
}

func addNewGame(gameId string, playerNames []string) error {
	mapLock.Lock()
	defer mapLock.Unlock()

	if _, exists := activeGames[gameId]; exists {
		return fmt.Errorf("game with id %q already exists", gameId)
	}

	game := gameState.NewGame(playerNames)
	router := mux.NewRouter()
	result := &gameRouter{
		Handler: http.StripPrefix("/"+gameId, router),
		game:    game,
	}
	activeGames[gameId] = result

	router.HandleFunc("/players", result.getPlayers)
	router.HandleFunc("/companies", result.getCompanies)
	router.HandleFunc("/state", result.getGameState)
	return nil
}

func initializeGameRoutes(router *mux.Router) {
	getter := router.Methods("GET").Subrouter()
	getter.HandleFunc("/{gameId}/players", serveGameContent)
	getter.HandleFunc("/{gameId}/companies", serveGameContent)
	getter.HandleFunc("/{gameId}/state", serveGameContent)
}

func serveGameContent(writer http.ResponseWriter, request *http.Request) {
	mapLock.RLock()
	defer mapLock.RUnlock()

	gameId := mux.Vars(request)["gameId"]
	if game, exists := activeGames[gameId]; !exists {
		data := jsonResponse{
			status: 404,
			Errors: []string{fmt.Sprintf("no active game with id %q", gameId)},
		}
		writeJson(&data, writer)
	} else {
		game.ServeHTTP(writer, request)
	}
}
