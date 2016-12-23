package main

import (
	"encoding/json"
	"io/ioutil"
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

func readBody(data interface{}, request *http.Request) error {
	if body, err := ioutil.ReadAll(request.Body); err != nil {
		return err
	} else if err = json.Unmarshal(body, data); err != nil {
		return err
	}
	return nil
}

func writeJson(resp *jsonResponse, writer http.ResponseWriter) {
	if buf, err := json.Marshal(resp); err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
	} else {
		writer.Header().Set("Content-Type", "application/json")
		if resp.status != 0 {
			writer.WriteHeader(resp.status)
		}
		writer.Write(buf)
	}
}

func convertErrors(errs []error) []string {
	if len(errs) == 0 {
		return nil
	}
	result := make([]string, 0, len(errs))
	for _, err := range errs {
		result = append(result, err.Error())
	}
	return result
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
	resp := jsonResponse{}
	defer writeJson(&resp, writer)

	if buf, err := boardInfo.JsonMap(); err != nil {
		resp.status = 500
		resp.Errors = []string{err.Error()}
	} else {
		resp.Result = (*json.RawMessage)(&buf)
	}
}

func getTrainCosts(writer http.ResponseWriter, request *http.Request) {
	writeJson(&jsonResponse{Result: boardInfo.AllTrainCosts()}, writer)
}

func main() {
	router := mux.NewRouter()
	initializeGameRoutes(router)

	static := router.Methods("GET").Subrouter()
	static.HandleFunc("/board_info", getboardInfo)
	static.HandleFunc("/train_costs", getTrainCosts)
	static.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(AssetFS{})))

	svr := http.Server{
		Handler: router,
		Addr:    ":8000",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	panic(svr.ListenAndServe())
}
