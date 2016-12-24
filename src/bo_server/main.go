package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
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
	port := flag.Int("port", 8000, "the port the web server will listen on")
	flag.Parse()

	rand.Seed(int64(time.Now().Nanosecond()))
	if flag.NArg() > 0 {
		addNewGame("game", flag.Args())
	} else {
		addNewGame("game", []string{"1st", "2nd", "3rd", "4th"})
	}

	router := mux.NewRouter()
	initializeGameRoutes(router)

	static := router.Methods("GET").Subrouter()
	static.HandleFunc("/board_info", getboardInfo)
	static.HandleFunc("/train_costs", getTrainCosts)
	static.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(AssetFS{})))

	svr := http.Server{
		Handler: router,
		Addr:    fmt.Sprintf(":%d", *port),

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	panic(svr.ListenAndServe())
}
