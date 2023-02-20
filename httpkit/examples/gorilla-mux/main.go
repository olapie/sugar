package main

import (
	"log"
	"math/rand"
	"net"
	"net/http"

	"code.olapie.com/sugar/v2/httpkit"
	"code.olapie.com/sugar/v2/jsonutil"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Beer struct {
	ID    string `json:"id"`
	Price int    `json:"price"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/beers/{id}", getBeer).Methods("GET")
	r.HandleFunc("/beers", httpkit.JoinHandlerFuncs(authenticate, createBeer)).Methods("POST")
	r.HandleFunc("/beers", httpkit.JoinHandlerFuncs(authenticate, listBeers)).Methods("GET")
	l, err := net.Listen("tcp", "localhost:8001")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(l.Addr())
	err = http.Serve(l, r)
	log.Println(err)
}

func authenticate(w http.ResponseWriter, req *http.Request) {
	token := req.URL.Query().Get("token")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func getBeer(w http.ResponseWriter, req *http.Request) {
	beer := &Beer{
		ID:    mux.Vars(req)["id"],
		Price: int(100 + rand.Uint32()%100),
	}
	_, err := w.Write(jsonutil.ToBytes(beer))
	if err != nil {
		log.Println(err)
	}
}

func createBeer(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("created"))
}

func listBeers(w http.ResponseWriter, req *http.Request) {
	var beers []*Beer
	for i := 0; i < 5; i++ {
		beer := &Beer{
			ID:    uuid.NewString(),
			Price: int(100 + rand.Uint32()%100),
		}
		beers = append(beers, beer)
	}

	_, err := w.Write(jsonutil.ToBytes(beers))
	if err != nil {
		log.Println(err)
	}
}
