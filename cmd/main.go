package main

import (
	"log"
	"net/http"

	"github.com/XIAOR1Ck/go-redis-voting/internal/voting"
		"github.com/XIAOR1Ck/go-redis-voting/internal/adapters/redis"
)

func main(){
	mux := http.NewServeMux() // http multiplexer
	store := voting.NewRedisStore(redis.NewClient("103.180.241.216:6379"))
	svc := voting.NewService(store)
	handler := voting.NewHandler(svc)

	//getprojects endpoint
	mux.HandleFunc("GET /getprojects", handler.GetProjects)
	mux.HandleFunc(("POST /addproject"), handler.CreateProject)
	if err := http.ListenAndServe(":9090", mux); err != nil {
		log.Fatal(err)
	}
}


