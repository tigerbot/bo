package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"boardInfo"
)

type jsonResponse struct {
	status int         `json:"-"`
	Errors []string    `json:"errors"`
	Result interface{} `json:"result"`
}

func writeJson(data *jsonResponse, writer http.ResponseWriter) {
	if buf, err := json.Marshal(data); err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
	} else {
		writer.Header().Set("Content-Type", "application/json")
		if data.status != 0 {
			writer.WriteHeader(data.status)
		}
		writer.Write(buf)
	}
}

func init() {
	addNewGame("game", []string{
		"Player 1",
		"Player 2",
		"Player 3",
		"Player 4",
		"Player 5",
		"Player 6",
	})
}

func getboardInfo(writer http.ResponseWriter, request *http.Request) {
	data := jsonResponse{}
	defer writeJson(&data, writer)

	if buf, err := boardInfo.JsonMap(); err != nil {
		data.status = 500
		data.Errors = []string{err.Error()}
	} else {
		data.Result = (*json.RawMessage)(&buf)
	}
}

func main() {
	router := mux.NewRouter()
	initializeGameRoutes(router)

	static := router.Methods("GET").Subrouter()
	static.HandleFunc("/board_info", getboardInfo)
	static.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(AssetFS{})))

	svr := http.Server{
		Handler: router,
		Addr:    ":8000",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	panic(svr.ListenAndServe())
}
