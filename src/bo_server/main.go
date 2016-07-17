package main

import (
	"fmt"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gorilla/mux"
)

func main() {
	staticFS := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo}
	r := mux.NewRouter()
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(staticFS)))

	fmt.Println("Hello World!")
	http.ListenAndServe(":8000", r)
}
