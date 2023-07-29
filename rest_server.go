package main

import (
	"fmt"
	"github.com/enjekt/shannon-engine/models"
	"github.com/enjekt/shannon-engine/pipelines"
	"github.com/enjekt/shannon-engine/pools"

	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var encipherPool pools.PipeLinePool
var decipherPool pools.PipeLinePool

// Use this until we get the database wired in...
var current models.Palette

func EncipherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	p := models.NewPalette()
	p.GetPan().Set(vars["pan"])
	pipeline := encipherPool.CheckOut()
	current = pipeline.Execute(p)
	encipherPool.CheckIn(pipeline)

	fmt.Fprintf(w, "{\n%v\n%v\n}", p.GetToken().String(), p.GetPad())

}

func DecipherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	p := models.NewPalette()
	if current.GetToken().String() == vars["token"] {
		pipeline := decipherPool.CheckOut()
		p.GetToken().Set(vars["token"])
		p.GetPad().Set(current.GetPad().String())             //This should come in on the request with the token...
		p.GetPaddedPan().Set(current.GetPaddedPan().String()) //This will be looked up in the database.
		pipeline.Execute(p)
		decipherPool.CheckIn(pipeline)
	}
	fmt.Fprintf(w, "{\n%v\n}", p.GetPan().String())
}
func main() {
	encipherPool = constructBin6PlusLast4Pool(10)
	decipherPool = constructDecipherPipeline(10)

	router := mux.NewRouter()
	router.HandleFunc("/decipher/{token}", DecipherHandler)
	router.HandleFunc("/encipher/{pan}", EncipherHandler)
	http.Handle("/", router)
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())

}
func constructBin6PlusLast4Pool(number int) pools.PipeLinePool {
	pool := pools.NewPool(number)
	for i := 0; i < number-1; i++ {
		pool.CheckIn(pipelines.NewPipeline().Add(pipelines.CompactAndStripPanFunc).Add(pipelines.CreatePadFunc).Add(pipelines.EncipherFunc).Add(pipelines.TokenFunc(6, 4)))
	}
	return pool
}

func constructDecipherPipeline(number int) pools.PipeLinePool {
	pool := pools.NewPool(number)
	for i := 0; i < number-1; i++ {
		pool.CheckIn(pipelines.NewPipeline().Add(pipelines.DecipherFunc))
	}
	return pool
}
