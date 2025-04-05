package main

import (
    "flag"
    "context"
    "log"

    "github.com/JeremiahVaughan/jobby/controllers" 
    "github.com/JeremiahVaughan/jobby/config" 
)


func main() {
    log.Println("jobby starting")
    ctx := context.Background()
    configPath := flag.String("c", "config.json", "location of config file")
    flag.Parse()
    
    theConfig, err := config.New(*configPath) 
    if err != nil {
        log.Fatalf("error, when creating config for main(). Error: %v", err)
    }

    theControllers := controllers.New(theConfig)
    for _, con := range theControllers {
        go func() {
            err = con.Start(ctx)
            if err != nil {
                log.Fatalf("error, when starting controllers for main(). Error: %v", err)
            }
        }()
    }

    log.Println("jobby running")
    select {
    case <- ctx.Done():
    }
}
