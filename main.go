package main

import (
    "context"
    "log"

    "github.com/JeremiahVaughan/jobby/config" 
    "github.com/JeremiahVaughan/jobby/clients" 
    "github.com/JeremiahVaughan/jobby/clients/healthy" 
    "github.com/JeremiahVaughan/jobby/models" 
    "github.com/JeremiahVaughan/jobby/controllers" 
    "github.com/JeremiahVaughan/jobby/router" 
)

func main() {
    log.Println("jobby starting")
    ctx := context.Background()
    
    serviceName := "jobby"
    config, err := config.New(ctx, serviceName) 
    if err != nil {
        log.Fatalf("error, when creating config for main(). Error: %v", err)
    }

    clients, err := clients.New(
        ctx,
        config.Clients,
        serviceName,
    )
    if err != nil {
        log.Fatalf("error, when creating clients for main(). Error: %v", err)
    }
    defer clients.Destroy()

    models, err := models.New(ctx, clients, config)
    if err != nil {
        log.Fatalf("error, when creating models for main(). Error: %v", err)
    }

    theControllers, httpControllers := controllers.New(models)
    healthyControllers := make(map[string]healthy.Controller)
    for k, v := range theControllers {
        healthyControllers[k] = v
    }
    err = clients.Healthy.Run(healthyControllers)
    if err != nil {
        log.Fatalf("error, when running healthy client for main(). Error: %v", err)
    }

    router := router.New(httpControllers)
    log.Println("jobby running")
    go func() {
        err2 := router.Run()
        if err2 != nil {
            log.Fatalf("error, when starting router for main(). Error: %v", err2)
        }
    }()

    // must run these last since the acme challenge controller depends on the router hosting an endpoint
    for _, con := range theControllers {
        go con.Start(ctx)
    }

    select {
    case <- ctx.Done():
        return
    }
}
