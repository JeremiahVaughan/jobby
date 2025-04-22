package router

import (
    "fmt"
    "net/http"

    "github.com/JeremiahVaughan/jobby/controllers" 
)

type Router struct {
    mux *http.ServeMux
}

func New(controllers *controllers.HttpControllers) *Router {
    mux := http.NewServeMux()
    mux.HandleFunc("/well-known/acme-challenge/", controllers.AcmeChallenger.Challenge)
    return &Router{
        mux: mux,
    }
}

func (r *Router) Run() error {
    err := http.ListenAndServe(":6666", r.mux)
    if err != nil {
        return fmt.Errorf("error, when starting http server. Error: %v", err)
    }
    return nil
}
