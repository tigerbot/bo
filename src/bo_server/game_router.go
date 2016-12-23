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

func (r gameRouter) getGameState(writer http.ResponseWriter, request *http.Request) {
	writeJson(&jsonResponse{Result: r.game.GlobalState}, writer)
}
func (r gameRouter) getPlayers(writer http.ResponseWriter, request *http.Request) {
	writeJson(&jsonResponse{Result: r.game.Players}, writer)
}
func (r gameRouter) getCompanies(writer http.ResponseWriter, request *http.Request) {
	writeJson(&jsonResponse{Result: r.game.Companies}, writer)
}

func (r gameRouter) takeMarketTurn(writer http.ResponseWriter, request *http.Request) {
	resp := jsonResponse{}
	defer writeJson(&resp, writer)

	var body struct {
		Player string `json:"player_name"`
		gameState.MarketTurn
	}
	if err := readBody(&body, request); err != nil {
		resp.status = 400
		resp.Errors = []string{fmt.Sprintf("invalid request: %v", err)}
		return
	}

	if errs := r.game.PerformMarketTurn(body.Player, body.MarketTurn); len(errs) > 0 {
		resp.status = 400
		resp.Errors = convertErrors(errs)
	}
}

func (r gameRouter) takeBussinessTurnOne(writer http.ResponseWriter, request *http.Request) {
	resp := jsonResponse{}
	defer writeJson(&resp, writer)

	var body struct {
		Player string `json:"player_name"`
		gameState.CompanyInventory
	}
	if err := readBody(&body, request); err != nil {
		resp.status = 400
		resp.Errors = []string{fmt.Sprintf("invalid request: %v", err)}
		return
	}

	if errs := r.game.UpdateCompanyInventory(body.Player, body.CompanyInventory); len(errs) > 0 {
		resp.status = 400
		resp.Errors = convertErrors(errs)
	}
}

func (r gameRouter) takeBussinessTurnTwo(writer http.ResponseWriter, request *http.Request) {
	resp := jsonResponse{}
	defer writeJson(&resp, writer)

	var body struct {
		Player string `json:"player_name"`
		gameState.CompanyEarnings
	}
	if err := readBody(&body, request); err != nil {
		resp.status = 400
		resp.Errors = []string{fmt.Sprintf("invalid request: %v", err)}
		return
	}

	if errs := r.game.HandleCompanyEarnings(body.Player, body.CompanyEarnings); len(errs) > 0 {
		resp.status = 400
		resp.Errors = convertErrors(errs)
	}
}

func addNewGame(gameId string, playerNames []string) error {
	mapLock.Lock()
	defer mapLock.Unlock()

	if _, exists := activeGames[gameId]; exists {
		return fmt.Errorf("game with id %q already exists", gameId)
	}

	router := mux.NewRouter()
	result := &gameRouter{
		Handler: http.StripPrefix("/"+gameId, router),
		game:    gameState.NewGame(playerNames),
	}
	activeGames[gameId] = result

	router.HandleFunc("/state", result.getGameState)
	router.HandleFunc("/players", result.getPlayers)
	router.HandleFunc("/companies", result.getCompanies)

	router.HandleFunc("/market_turn", result.takeMarketTurn)
	router.HandleFunc("/business_turn_one", result.takeBussinessTurnOne)
	router.HandleFunc("/business_turn_two", result.takeBussinessTurnTwo)
	return nil
}

func initializeGameRoutes(router *mux.Router) {
	getter := router.Methods("GET").Subrouter()
	getter.HandleFunc("/{gameId}/state", serveGameContent)
	getter.HandleFunc("/{gameId}/players", serveGameContent)
	getter.HandleFunc("/{gameId}/companies", serveGameContent)

	poster := router.Methods("POST").Subrouter()
	poster.HandleFunc("/{gameId}/market_turn", serveGameContent)
	poster.HandleFunc("/{gameId}/business_turn_one", serveGameContent)
	poster.HandleFunc("/{gameId}/business_turn_two", serveGameContent)
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
