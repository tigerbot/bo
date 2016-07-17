package main

import (
	"fmt"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"

	"boardInfo"
)

func getboardInfo(w http.ResponseWriter, r *http.Request) {
	if buf, err := boardInfo.JsonMap(); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf)
	}
}

func main() {
	staticFS := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo}
	r := mux.NewRouter()
	r.HandleFunc("/board_info", getboardInfo)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(staticFS)))

	fmt.Println("Hello World!")
	http.ListenAndServe(":8000", r)
}
