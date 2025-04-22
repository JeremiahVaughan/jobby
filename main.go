package main

import (
    "context"
    "log"

    "github.com/JeremiahVaughan/jobby/config" 
    "github.com/JeremiahVaughan/jobby/clients" 
    "github.com/JeremiahVaughan/jobby/models" 
    "github.com/JeremiahVaughan/jobby/controllers" 
    "github.com/JeremiahVaughan/jobby/router" 
)


func main() {
    log.Println("jobby starting")
    ctx := context.Background()
    
    config, err := config.New(ctx) 
    if err != nil {
        log.Fatalf("error, when creating config for main(). Error: %v", err)
    }

    clients, err := clients.New(ctx, config.Clients)
    if err != nil {
        log.Fatalf("error, when creating clients for main(). Error: %v", err)
    }
    defer clients.Destroy()

    models := models.New(clients)

    theControllers := controllers.New(models)
    for _, con := range theControllers {
        go func() {
            err = con.Start(ctx)
            if err != nil {
                log.Fatalf("error, when starting controllers for main(). Error: %v", err)
            }
        }()
    }

    httpControllers := controllers.NewHttpControllers(models)
    router := router.New(httpControllers)

    log.Println("jobby running")
    err = router.Run()
    if err != nil {
        log.Fatalf("error, when starting router for main(). Error: %v", err)
    }
}
