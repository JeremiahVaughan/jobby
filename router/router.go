package router

import (
    "fmt"
    "net/http"

    "github.com/JeremiahVaughan/jobby/controllers" 
)

type Router struct {
    mux *http.ServeMux
    controllers *controllers.HttpControllers
}

func New(controllers *controllers.HttpControllers) *Router {
    mux := http.NewServeMux()
    mux.HandleFunc("/health", controllers.Health.Check)
    return &Router{
        mux: mux,
        controllers: controllers,
    }
}

func (r *Router) Run() error {
    r.controllers.AcmeChallenger.StartChallengeProvider("5001")
    err := http.ListenAndServe(":6666", r.mux)
    if err != nil {
        return fmt.Errorf("error, when starting http server. Error: %v", err)
    }
    return nil
}
