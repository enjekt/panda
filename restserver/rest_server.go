package restserver

import (
	"fmt"
	"github.com/enjekt/shannon-engine/models"
	"github.com/enjekt/shannon-engine/pipelines"
	"github.com/enjekt/shannon-engine/pools"
	"log"
	"panda/database"
	"time"

	"github.com/gorilla/mux"
	"net/http"
)

var db = database.NewDAO()

func NewRestServer() {
	log.Println("Starting REST Server...")
	encipherPool = constructBin6PlusLast4Pool(10)
	decipherPool = constructDecipherPipeline(10)
	db.Open()
	router := mux.NewRouter()
	router.HandleFunc("/storepanservice/v1/pan/{token}/{pad}", DecipherHandler)
	router.HandleFunc("/storepanservice/v1/identifiers/{pan}", EncipherHandler)
	http.Handle("/", router)
	srv := &http.Server{
		Handler: router,
		Addr:    "localhost:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Listen on ")
	log.Fatal(srv.ListenAndServe())

}

var encipherPool pools.PipeLinePool
var decipherPool pools.PipeLinePool

func EncipherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	p := models.NewPalette()
	p.GetPan().Set(vars["pan"])
	pipeline := encipherPool.CheckOut()
	pipeline.Execute(p)
	encipherPool.CheckIn(pipeline)
	db.Insert(p)
	fmt.Fprintf(w, "{\nToken: %v\nPad: %v\n}", p.GetToken().String(), p.GetPad().String())

}

func DecipherHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	p := models.NewPalette()
	p.GetPad().Set(vars["pad"])
	//This should come in on the request with the token...
	p.GetToken().Set(vars["token"])
	db.Get(p) //This will be looked up in the database.
	//TODO automate check out/check internally in the pool
	pipeline := decipherPool.CheckOut()
	pipeline.Execute(p)
	decipherPool.CheckIn(pipeline)

	fmt.Fprintf(w, "{\nPan: %v\n}", p.GetPan().String())
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
